package repository

import (
	"context"
	"database/sql"

	"ivr_ataljanseva/models"
	"github.com/google/uuid"

)

type CitizenRepository struct {
	db *sql.DB
}

func NewCitizenRepository(db *sql.DB) *CitizenRepository {
	return &CitizenRepository{db: db}
}

func (r *CitizenRepository) FindByPhone(
	ctx context.Context,
	phone string,
) (*models.CitizenIVR, error) {

	query := `
	SELECT
		language,
		pincode,
		ward,
		nagarsevak_id
	FROM citizen_ivr
	WHERE phone_number = $1
	LIMIT 1;
	`

	var citizen models.CitizenIVR

	err := r.db.QueryRowContext(ctx, query, phone).Scan(
		&citizen.Language,
		&citizen.Pincode,
		&citizen.Ward,
		&citizen.NagarsevakID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &citizen, nil
}

func (r *CitizenRepository) Create(
	ctx context.Context,
	req *models.RegisterCitizenRequest,
) error {

	query := `
	INSERT INTO citizen_ivr (
		phone_number,
		language,
		pincode,
		ward,
		nagarsevak_id
	)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (phone_number)
	DO UPDATE SET
		language = EXCLUDED.language,
		pincode = EXCLUDED.pincode,
		ward = EXCLUDED.ward,
		nagarsevak_id = EXCLUDED.nagarsevak_id;
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		req.PhoneNumber,
		req.Language,
		req.Pincode,
		req.Ward,
		req.NagarsevakID,
	)

	return err
}


func (r *CitizenRepository) UpsertCitizen(
	ctx context.Context,
	phone string,
	language string,
	pincode string,
	ward string,
	nagarsevakID uuid.UUID,
) error {

	query := `
	INSERT INTO citizen_ivr (
		phone_number,
		language,
		pincode,
		ward,
		nagarsevak_id
	)
	VALUES ($1,$2,$3,$4,$5)

	ON CONFLICT (phone_number)
	DO UPDATE SET
		language = EXCLUDED.language,
		pincode = EXCLUDED.pincode,
		ward = EXCLUDED.ward,
		nagarsevak_id = EXCLUDED.nagarsevak_id
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		phone,
		language,
		pincode,
		ward,
		nagarsevakID,
	)

	return err
}