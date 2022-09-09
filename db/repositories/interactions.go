package repositories

import (
	"context"
	"fmt"
	"tinderutf/db"
	"tinderutf/db/models"
	"tinderutf/domain"
)

func FindPeople(ctx context.Context, userId string, distance int, geo *domain.Geo, sexPref domain.Sex, sex domain.Sex) ([]*domain.User, error) {
	var sexFilter string
	if sexPref == domain.Male {
		sexFilter = " AND sex in ('MALE', 'ALL')"
	}

	if sexPref == domain.Female {
		sexFilter = " AND sex in ('FEMALE', 'ALL')"
	}

	var finderSexFilter = fmt.Sprintf(" AND sex_preference in ('%s', 'ALL')", sex)
	if sex == domain.All {
		finderSexFilter = " AND sex_preference in ('MALE', 'FEMALE', 'ALL')"
	}

	var result []*domain.User
	query := ` SELECT name, email, instagram, about, sex_preference,  find_distance, birth_date, password, is_active, external_id, st_x(geolocation) as lat, st_y(geolocation) as lng, sex
		FROM users where external_id <> $1 and ST_Contains(ST_Buffer(ST_GeomFromText($2), $3),
			ST_Buffer(geolocation, 0.001)) = true AND id NOT IN (
			select target_user_id from interactions WHERE user_id = (select id from users where external_id = $1)
			)` + sexFilter + finderSexFilter

	rows, err := db.Db.QueryContext(ctx, query, userId, convertToPoint(geo.Latitude(), geo.Longitude()), float64(distance)/100)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var model models.User
		if err := rows.Scan(&model.Name, &model.Email, &model.Instagram, &model.About,
			&model.SexPreference, &model.FindDistance, &model.BirthDate, &model.Password,
			&model.IsActive, &model.ExternalId, &model.Latitude, &model.Longitude, &model.Sex); err != nil {
			return result, err
		}

		user := domain.NewUserFromDb(model.ExternalId, model.Name, model.BirthDate,
			model.Email, model.Password, model.IsActive, model.Sex)
		user.SetOptions(domain.NewUserOptions(model.SexPreference, model.FindDistance))
		user.SetInfo(domain.NewUserInfo(model.About, model.Instagram, domain.NewGeo(model.Latitude.Decimal, model.Longitude.Decimal)))

		result = append(result, user)
	}

	return result, nil
}

func ShowLiked(ctx context.Context, userId string) ([]*domain.User, error) {
	var result []*domain.User
	query := `SELECT u.name, u.email, u.instagram, u.about, u.sex_preference,  u.find_distance, u.birth_date,
	u.password, u.is_active, u.external_id, st_x(u.geolocation) as lat, st_y(u.geolocation) as lng, u.sex
			FROM interactions i
			INNER JOIN users u ON i.target_user_id = u.id
			where i.user_id = (select id from users where external_id = $1) AND i.liked = true`

	rows, err := db.Db.QueryContext(ctx, query, userId)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var model models.User
		if err := rows.Scan(&model.Name, &model.Email, &model.Instagram, &model.About,
			&model.SexPreference, &model.FindDistance, &model.BirthDate, &model.Password,
			&model.IsActive, &model.ExternalId, &model.Latitude, &model.Longitude, &model.Sex); err != nil {
			return result, err
		}

		user := domain.NewUserFromDb(model.ExternalId, model.Name, model.BirthDate,
			model.Email, model.Password, model.IsActive, model.Sex)
		user.SetOptions(domain.NewUserOptions(model.SexPreference, model.FindDistance))
		user.SetInfo(domain.NewUserInfo(model.About, model.Instagram, domain.NewGeo(model.Latitude.Decimal, model.Longitude.Decimal)))

		result = append(result, user)
	}

	return result, nil
}

func ShowMatches(ctx context.Context, userId string) ([]*domain.User, error) {
	var result []*domain.User
	query := ` SELECT u.name, u.email, u.instagram, u.about, u.sex_preference,  u.find_distance, u.birth_date, 
        u.password, u.is_active, u.external_id, st_x(u.geolocation) as lat, st_y(u.geolocation) as lng, u.sex
		FROM users u
		INNER JOIN interactions i ON i.user_id = u.id
		where u.id in (
			select user_id from interactions where target_user_id = (select id from users where external_id = $1) and liked = true
			) and u.id in (
			select target_user_id from interactions where user_id = (select id from users where external_id = $1) and liked = true
			)`

	rows, err := db.Db.QueryContext(ctx, query, userId)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var model models.User
		if err := rows.Scan(&model.Name, &model.Email, &model.Instagram, &model.About,
			&model.SexPreference, &model.FindDistance, &model.BirthDate, &model.Password,
			&model.IsActive, &model.ExternalId, &model.Latitude, &model.Longitude, &model.Sex); err != nil {
			return result, err
		}

		user := domain.NewUserFromDb(model.ExternalId, model.Name, model.BirthDate,
			model.Email, model.Password, model.IsActive, model.Sex)
		user.SetOptions(domain.NewUserOptions(model.SexPreference, model.FindDistance))
		user.SetInfo(domain.NewUserInfo(model.About, model.Instagram, domain.NewGeo(model.Latitude.Decimal, model.Longitude.Decimal)))

		result = append(result, user)
	}

	return result, nil
}

func CreateInteraction(ctx context.Context, userId, targetId string, like bool) (bool, error) {
	query := `INSERT INTO interactions (user_id, target_user_id, liked)
	VALUES ((select id from users where external_id = $1), 
			(select id from users where external_id = $2), $3)`
	
	_, err := db.Db.ExecContext(ctx, query, userId, targetId, like)
	if err != nil {
		return false, err
	}

	matchQuery := `SELECT liked FROM interactions 
	WHERE user_id = (select id from users where external_id = $1) 
	AND target_user_id = (select id from users where external_id = $2)`

	var isMatch bool
	_ = db.Db.QueryRowContext(ctx, matchQuery, targetId, userId).Scan(&isMatch)

	return isMatch, nil
}

func CancelInteraction(ctx context.Context, userId, targetId string) error {
	query := `UPDATE interactions SET liked = false
	WHERE user_id = (select id from users where external_id = $1) 
	AND target_user_id = (select id from users where external_id = $2)`

	_, err := db.Db.ExecContext(ctx, query, userId, targetId)
	if err != nil {
		return err
	}

	return nil
}