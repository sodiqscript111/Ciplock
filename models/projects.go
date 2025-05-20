package models

import (
	"Ciplock/cliplockdb"
	"fmt"
	"time"
)

type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
	AdminID   string    `json:"admin_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

func AddProjects(project Project) error {
	query := `INSERT INTO projects (id, name, api_key, admin_id, created_at, updated_at, is_active)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := cliplockdb.DB.Exec(query,
		project.ID,
		project.Name,
		project.APIKey,
		project.AdminID,
		project.CreatedAt,
		project.UpdatedAt,
		project.IsActive,
	)
	return err
}

func GetProjects(adminID string) ([]Project, error) {
	projects := []Project{}
	query := `SELECT id, name, api_key, admin_id, created_at, updated_at, is_active FROM projects WHERE admin_id = $1`

	rows, err := cliplockdb.DB.Query(query, adminID)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var project Project
		err := rows.Scan(&project.ID, &project.Name, &project.APIKey, &project.AdminID,
			&project.CreatedAt, &project.UpdatedAt, &project.IsActive)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project row: %w", err)
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return projects, nil
}
