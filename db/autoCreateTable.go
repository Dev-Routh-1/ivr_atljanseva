package db

import "log"

func AutoMigrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS citizen_ivr (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		phone_number VARCHAR(20) UNIQUE NOT NULL,
		language VARCHAR(10),
		pincode VARCHAR(10),
		ward VARCHAR(50),
		nagarsevak_id UUID,
		created_at TIMESTAMP DEFAULT NOW()
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("citizen_ivr table ready")
	return nil
}