package models

import (
	"Ciplock/cliplockdb"
	"Ciplock/utils"
	"errors"
	"log"
	"time"
)

type Admin struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func AddAdmin(admin Admin) error {
	query := `
		INSERT INTO users (email, hashed_password, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := cliplockdb.DB.Exec(
		query,
		admin.Email,
		admin.PasswordHash,
		admin.CreatedAt,
		admin.UpdatedAt,
	)

	if err != nil {
		log.Printf("Failed to add admin: %v", err)
		return err
	}

	return nil
}

func (admin *Admin) ValidateUser(plainPassword string) (bool, error) {
	query := `SELECT id, hashed_password FROM users WHERE email = $1`
	row := cliplockdb.DB.QueryRow(query, admin.Email)

	var storedHash string
	err := row.Scan(&admin.ID, &storedHash)
	if err != nil {
		log.Printf("User not found or DB error: %v", err)
		return false, err
	}

	err = utils.CheckPasswordHash(plainPassword, storedHash)
	if err != nil {
		log.Println("Invalid password attempt")
		return false, errors.New("invalid credentials")
	}

	return true, nil
}
