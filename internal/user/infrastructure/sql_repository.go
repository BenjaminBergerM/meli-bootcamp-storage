package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"meli-bootcamp-storage/internal/models"
	"meli-bootcamp-storage/internal/user"
)

const (
	insertQuery    = "INSERT INTO users (uuid, firstname, lastname, username, password, email, ip, macAddress, website, image) values (?,?,?,?,?,?,?,?,?,?);"
	selectOneQuery = "SELECT * FROM users WHERE uuid = ?"
	updateQuery    = "UPDATE users SET firstname = ?, lastname = ?, username = ?, password = ?, email = ?, ip = ?, macAddress = ?, website = ?, image = ? WHERE uuid = ?"
	selectQuery    = "SELECT * FROM users"
	deleteQuery    = "DELETE from users where uuid = ?"
)

var _ user.Repository = (*sqlRepository)(nil)

type sqlRepository struct {
	db *sql.DB
}

func NewSqlRepository(db *sql.DB) *sqlRepository {
	return &sqlRepository{db: db}
}

func (receiver *sqlRepository) Store(ctx context.Context, model *models.User) error {
	_, err := receiver.db.ExecContext(ctx, insertQuery, model.UUID, model.Firstname, model.Lastname, model.Username, model.Password, model.Email, model.IP, model.MacAddress, model.Website, model.Image)
	return err
}

func (receiver *sqlRepository) GetOne(ctx context.Context, id uuid.UUID) (*models.User, error) {
	result := new(models.User)
	row := receiver.db.QueryRowContext(ctx, selectOneQuery, id)

	err := row.Err()

	if err != nil {
		return nil, err
	}

	err = row.Scan(
		&result.UUID,
		&result.Firstname,
		&result.Lastname,
		&result.Username,
		&result.Password,
		&result.Email,
		&result.IP,
		&result.MacAddress,
		&result.Website,
		&result.Image,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (receiver *sqlRepository) Update(ctx context.Context, model *models.User) error {
	result, err := receiver.db.ExecContext(
		ctx,
		updateQuery,
		model.Firstname,
		model.Lastname,
		model.Username,
		model.Password,
		model.Email,
		model.IP,
		model.MacAddress,
		model.Website,
		model.Image,
		model.UUID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return fmt.Errorf("%d users updated", rowsAffected)
	}

	return nil
}

func (receiver *sqlRepository) GetAll(ctx context.Context) ([]models.User, error) {
	result := make([]models.User, 0)
	rows, err := receiver.db.QueryContext(ctx, selectQuery)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := models.User{}
		err = rows.Scan(
			&user.UUID,
			&user.Firstname,
			&user.Lastname,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.IP,
			&user.MacAddress,
			&user.Website,
			&user.Image,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return result, nil
}

func (receiver *sqlRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := receiver.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return fmt.Errorf("%d users deleted", rowsAffected)
	}

	return nil
}
