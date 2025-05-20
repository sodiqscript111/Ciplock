package models

import (
	"Ciplock/cliplockdb"
	"Ciplock/utils"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	ProjectID uuid.UUID `json:"project_id"`
	CreatedAt time.Time `json:"created_at"`
}

func CustomerExists(email string, projectID uuid.UUID) (bool, error) {
	query := `SELECT COUNT(*) FROM customers WHERE email = $1 AND project_id = $2`
	var count int
	err := cliplockdb.DB.QueryRow(query, email, projectID).Scan(&count)
	return count > 0, err
}

func CreateCustomer(customer Customer) error {
	query := `INSERT INTO customers (id, full_name, email, password_hash, project_id, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := cliplockdb.DB.Exec(query,
		customer.ID,
		customer.FullName,
		customer.Email,
		customer.Password,
		customer.ProjectID,
		time.Now(),
	)

	return err
}

func (customer *Customer) ValidateCustomers(plainPassword, projectID string) (bool, error) {
	query := `SELECT id, password_hash FROM customers WHERE email = $1 AND project_id = $2`

	row := cliplockdb.DB.QueryRow(query, customer.Email, projectID)
	var storedHash string
	err := row.Scan(&customer.ID, &storedHash)
	if err != nil {
		return false, err
	}

	if err := utils.CheckPasswordHash(plainPassword, storedHash); err != nil {
		return false, fmt.Errorf("invalid credentials")
	}
	return true, nil
}
