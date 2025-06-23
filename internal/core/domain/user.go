package domain

import (
	"time"

	"github.com/google/uuid"
)

// User representa la entidad de usuario en el dominio
type User struct {
	// ID           uuid.UUID `json:"id" gorm:"type:char(36);primaryKey;default:uuid_generate_v4()"`
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name         string    `json:"name" gorm:"column:NAME;type:varchar(100);not null"`
	LastName     string    `json:"lastname" gorm:"column:LASTNAME;type:varchar(100);not null"`
	Username     string    `json:"username" gorm:"column:USERNAME;type:varchar(100);not null;unique"`
	Email        string    `json:"email" gorm:"column:EMAIL;type:varchar(255);not null;unique"`
	DNI          string    `json:"dni" gorm:"column:DNI;type:varchar(20);unique"`
	Phone        string    `json:"phone" gorm:"column:PHONE;type:varchar(20)"`
	PasswordHash string    `json:"-" gorm:"column:PASSWORD_HASH;type:varchar(255);not null"`
	Active       bool      `json:"active" gorm:"column:ACTIVE;default:true"`

	// Relaciones (FKs)
	RoleID uuid.UUID `json:"role_id" gorm:"column:ROLE_ID;type:uuid;not null"`
	Role   Role      `json:"role" gorm:"foreignKey:RoleID"`

	LocalityID *uuid.UUID `json:"locality_id" gorm:"column:LOCALITY_ID;type:uuid"`
	Locality   *Locality  `json:"locality" gorm:"foreignKey:LocalityID"`

	Patients []Patient `json:"patients" gorm:"foreignKey:UserID"`

	CreatedAt time.Time  `json:"created_at,omitempty" gorm:"column:CREATE_AT;autoCreateTime"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:UPDATE_AT;autoUpdateTime"`
}

// TableName especifica el nombre de la tabla para GORM
func (User) TableName() string {
	return "users"
}

// NewUser crea una nueva instancia de User
func NewUser(
	name, lastName, username, dni, phone, email, passwordHash string,
	// patients []Patient,
	roleID uuid.UUID,
	localityID *uuid.UUID,
) *User {
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
		LocalityID:   localityID,
		// Patients:     patients,
		CreatedAt: time.Now(),
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

	now := time.Now()
	u.UpdatedAt = &now
}

// UpdatePassword actualiza la contraseña del usuario
func (u *User) UpdatePassword(passwordHash string) {
	u.PasswordHash = passwordHash

	now := time.Now()
	u.UpdatedAt = &now
}

// UpdateRole actualiza el rol del usuario
func (u *User) UpdateRole(roleID uuid.UUID) {
	u.RoleID = roleID

	now := time.Now()
	u.UpdatedAt = &now
}
