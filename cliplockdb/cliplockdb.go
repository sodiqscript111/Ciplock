package cliplockdb

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open("postgres", "user=postgres password=password dbname=ciplock sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)

	log.Println("Database connection established successfully.")

	err = enablePGCrypto()
	if err != nil {
		log.Fatalf("Failed to enable pgcrypto extension: %v", err)
	}

	err = CreateUserTable()
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	err = CreateProjectTbale()
	if err != nil {
		log.Fatalf("Failed to create project tbale: %v", err)
	}
	err = CreateCustomers()
	if err != nil {
		log.Fatalf("Failed to create customers: %v", err)
	}
}

func enablePGCrypto() error {
	_, err := DB.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`)
	return err
}

func CreateUserTable() error {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email TEXT UNIQUE NOT NULL,
		hashed_password TEXT NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		is_verified BOOLEAN DEFAULT FALSE,
		role TEXT DEFAULT 'user',
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW()
	);
	`

	_, err := DB.Exec(createUserTable)
	return err
}

func CreateProjectTbale() error {
	createProjectTbale := `
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    api_key TEXT UNIQUE NOT NULL,
    admin_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (admin_id) REFERENCES users(id)
);
`
	_, err := DB.Exec(createProjectTbale)
	if err != nil {
		log.Println("Table 'projects' already exists, skipping creation.")
	}
	return nil
}

func CreateCustomers() error {
	createCustomers := `CREATE TABLE IF NOT EXISTS customers (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		project_id UUID REFERENCES projects(id),
		full_name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	);`

	_, err := DB.Exec(createCustomers)
	if err != nil {
		log.Println("Error creating 'customers' table:", err)
		return err
	}

	log.Println("Table 'customers' checked/created successfully.")
	return nil
}
