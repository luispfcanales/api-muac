package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/luispfcanales/api-muac/internal/core/domain"
)

// roleRepository implementa la interfaz RoleRepository para PostgreSQL
type roleRepository struct {
	db *sql.DB
}

// NewRoleRepository crea una nueva instancia de RoleRepository
func NewRoleRepository(db *sql.DB) *roleRepository {
	return &roleRepository{
		db: db,
	}
}

// Create inserta un nuevo rol en la base de datos
func (r *roleRepository) Create(ctx context.Context, role *domain.Role) error {
	query := `INSERT INTO ROLE (ID, NAME, DESCRIPTION) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, role.ID, role.Name, role.Description)
	if err != nil {
		return fmt.Errorf("error al crear rol: %w", err)
	}

	return nil
}

// GetByID obtiene un rol por su ID
func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	query := `SELECT ID, NAME, DESCRIPTION FROM ROLE WHERE ID = $1`

	var role domain.Role
	err := r.db.QueryRowContext(ctx, query, id).Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrRoleNotFound
		}
		return nil, fmt.Errorf("error al obtener rol: %w", err)
	}

	return &role, nil
}

// GetAll obtiene todos los roles
func (r *roleRepository) GetAll(ctx context.Context) ([]*domain.Role, error) {
	query := `SELECT ID, NAME, DESCRIPTION FROM ROLE`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al obtener roles: %w", err)
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			return nil, fmt.Errorf("error al escanear rol: %w", err)
		}
		roles = append(roles, &role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar roles: %w", err)
	}

	return roles, nil
}

// Update actualiza un rol existente
func (r *roleRepository) Update(ctx context.Context, role *domain.Role) error {
	query := `UPDATE ROLE SET NAME = $1, DESCRIPTION = $2 WHERE ID = $3`

	result, err := r.db.ExecContext(ctx, query, role.Name, role.Description, role.ID)
	if err != nil {
		return fmt.Errorf("error al actualizar rol: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrRoleNotFound
	}

	return nil
}

// Delete elimina un rol por su ID
func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM ROLE WHERE ID = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar rol: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al obtener filas afectadas: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrRoleNotFound
	}

	return nil
}
