package domain

import (
	"time"

	"github.com/google/uuid"
)

// User representa la entidad de usuario en el dominio
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name         string    `json:"name" gorm:"column:NAME;type:varchar(100);not null"`
	LastName     string    `json:"lastname" gorm:"column:LASTNAME;type:varchar(100);not null"`
	Username     string    `json:"username" gorm:"column:USER;type:varchar(50);not null"`
	Email        string    `json:"email" gorm:"column:EMAIL;type:varchar(255);not null"`
	DNI          string    `json:"dni" gorm:"column:DNI;type:varchar(8);not null"`
	Phone        string    `json:"phone" gorm:"column:PHONE;type:varchar(20)"`
	PasswordHash string    `json:"password_hash" gorm:"column:PASSWORD_HASH;type:varchar(255);not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:CREATE_AT;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:UPDATE_AT;autoUpdateTime"`
	RoleID       uuid.UUID `json:"role_id" gorm:"column:ROLE_ID;type:uuid"`
	Role         *Role     `json:"role" gorm:"foreignKey:RoleID"`
}

// TableName especifica el nombre de la tabla para GORM
func (User) TableName() string {
	return "USER"
}

// NewUser crea una nueva instancia de User
func NewUser(name, lastName, username, dni, phone, email, passwordHash string, roleID uuid.UUID) *User {
	return &User{
		ID:           uuid.New(),
		Name:         name,
		LastName:     lastName,
		Username:     username,
		Email:        email,
		DNI:          dni,
		Phone:        phone,
		PasswordHash: passwordHash,
		RoleID:       roleID,
		CreatedAt:    time.Now(),
	}
}

// Validate valida que el usuario tenga los campos requeridos
func (u *User) Validate() error {
	if u.Name == "" {
		return ErrEmptyUserName
	}
	if u.LastName == "" {
		return ErrEmptyUserLastName
	}
	if u.Username == "" {
		return ErrEmptyUsername
	}
	if u.Email == "" {
		return ErrEmptyUserEmail
	}
	if u.PasswordHash == "" {
		return ErrEmptyUserPassword
	}
	return nil
}

// Update actualiza los campos del usuario
func (u *User) Update(name, lastname, user, email, phone, dni string, roleID uuid.UUID) {
	u.Name = name
	u.LastName = lastname
	u.Username = user
	u.Email = email
	u.DNI = dni
	u.Phone = phone
	u.RoleID = roleID
	u.UpdatedAt = time.Now()
}

// UpdatePassword actualiza la contraseña del usuario
func (u *User) UpdatePassword(passwordHash string) {
	u.PasswordHash = passwordHash
	u.UpdatedAt = time.Now()
}

// UpdateRole actualiza el rol del usuario
func (u *User) UpdateRole(roleID uuid.UUID) {
	u.RoleID = roleID
	u.UpdatedAt = time.Now()
}
