package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
	"tinderutf/db"
	"tinderutf/db/models"
	"tinderutf/domain"
)

func GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model models.User

	err := db.Db.QueryRowContext(ctx, `SELECT name, email, birth_date, password, is_active, external_id, sex
		FROM users where email = $1`, email).Scan(&model.Name, &model.Email, &model.BirthDate, &model.Password,
		&model.IsActive, &model.ExternalId, &model.Sex)
	if err != nil {
		return nil, err
	}

	return domain.NewUserFromDb(model.ExternalId, model.Name, model.BirthDate,
		model.Email, model.Password, model.IsActive, model.Sex), nil
}

func CreateUser(ctx context.Context, user *domain.User) (string, error) {
	tx, err := db.Db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var externalId string

	if err := tx.QueryRowContext(ctx, `INSERT INTO users (name, email, password, birth_date, is_active, sex) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING external_id`,
		user.Name(), user.Email(), user.Password(), user.BirthDate(), true, user.Sex()).Scan(&externalId); err != nil {

		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return externalId, nil
}

func SetUserLocation(ctx context.Context, lat, lng decimal.Decimal, userId string) error {
	tx, err := db.Db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `UPDATE users 
		SET geolocation = ST_GeomFromText($1) WHERE external_id = $2`,
		convertToPoint(lat, lng), userId); err != nil {

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func SetUserCustomization(ctx context.Context, instagram, about string, distance int, sexPref domain.Sex, userId string) error {
	tx, err := db.Db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `UPDATE users 
		SET instagram = $1, about = $2, find_distance = $3, sex_preference = $4 WHERE external_id = $5`,
		instagram, about, distance, sexPref, userId); err != nil {

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func GetUserById(ctx context.Context, id string) (*domain.User, error) {
	var model models.User

	err := db.Db.QueryRowContext(ctx, `SELECT name, email, instagram, about, sex_preference,  find_distance, birth_date, password, is_active, external_id, st_x(geolocation) as lat, st_y(geolocation) as lng, sex
		FROM users where external_id = $1`, id).Scan(&model.Name, &model.Email, &model.Instagram, &model.About,
		&model.SexPreference, &model.FindDistance, &model.BirthDate, &model.Password, &model.IsActive, &model.ExternalId, &model.Latitude, &model.Longitude, &model.Sex)
	if err != nil {
		return nil, err
	}

	user := domain.NewUserFromDb(model.ExternalId, model.Name, model.BirthDate,
		model.Email, model.Password, model.IsActive, model.Sex)

	user.SetOptions(domain.NewUserOptions(model.SexPreference, model.FindDistance))
	user.SetInfo(domain.NewUserInfo(model.About, model.Instagram, domain.NewGeo(model.Latitude.Decimal, model.Longitude.Decimal)))

	return user, nil
}

func convertToPoint(lat, lng decimal.Decimal) string {
	return fmt.Sprintf("Point(%s %s)", lat.String(), lng.String())
}
