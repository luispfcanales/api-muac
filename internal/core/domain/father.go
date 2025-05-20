package domain

import (
	"time"

	"github.com/google/uuid"
)

// Father representa la entidad de padre/administrador en el dominio
type Father struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name         string    `json:"name" gorm:"column:NAME;type:varchar(100);not null"`
	LastName     string    `json:"lastname" gorm:"column:LASTNAME;type:varchar(100);not null"`
	Email        string    `json:"email" gorm:"column:EMAIL;type:varchar(255);not null"`
	DNI          int       `json:"dni" gorm:"column:DNI"`
	Phone        string    `json:"phone" gorm:"column:PHONE;type:varchar(20)"`
	PasswordHash string    `json:"password_hash" gorm:"column:PASSWORD_HASH;type:varchar(255);not null"`
	Active       bool      `json:"active" gorm:"column:ACTIVE;default:true"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:CREATE_AT;autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:UPDATE_AT;autoUpdateTime"`
	RoleID       uuid.UUID `json:"role_id" gorm:"column:ROLE_ID;type:uuid"`
	LocalityID   uuid.UUID `json:"locality_id" gorm:"column:LOCALITY_ID;type:uuid"`
	PatientID    uuid.UUID `json:"patient_id" gorm:"column:PATIENT_ID;type:uuid"`
	Role         *Role     `json:"role" gorm:"foreignKey:RoleID"`
	Locality     *Locality `json:"locality" gorm:"foreignKey:LocalityID"`
	Patient      *Patient  `json:"patient" gorm:"foreignKey:PatientID"`
}

// TableName especifica el nombre de la tabla para GORM
func (Father) TableName() string {
	return "FATHER"
}

// NewFather crea una nueva instancia de Father
func NewFather(name, lastName, email, passwordHash string, roleID, localityID, patientID uuid.UUID) *Father {
	return &Father{
		ID:           uuid.New(),
		Name:         name,
		LastName:     lastName,
		Email:        email,
		PasswordHash: passwordHash,
		Active:       true,
		RoleID:       roleID,
		LocalityID:   localityID,
		PatientID:    patientID,
		CreatedAt:    time.Now(),
	}
}

// Validate valida que el padre tenga los campos requeridos
func (f *Father) Validate() error {
	if f.Name == "" {
		return ErrEmptyFatherName
	}
	if f.LastName == "" {
		return ErrEmptyFatherLastName
	}
	if f.Email == "" {
		return ErrEmptyFatherEmail
	}
	if f.PasswordHash == "" {
		return ErrEmptyFatherPassword
	}
	return nil
}

// Update actualiza los campos del padre
func (f *Father) Update(name, lastname, email, phone string, dni int, active bool, roleID, localityID, patientID uuid.UUID) {
	f.Name = name
	f.LastName = lastname
	f.Email = email
	f.DNI = dni
	f.Phone = phone
	f.Active = active
	f.RoleID = roleID
	f.LocalityID = localityID
	f.PatientID = patientID
	f.UpdatedAt = time.Now()
}

// UpdatePassword actualiza la contraseña del padre
func (f *Father) UpdatePassword(passwordHash string) {
	f.PasswordHash = passwordHash
	f.UpdatedAt = time.Now()
}

// UpdateActive actualiza el estado activo del padre
func (f *Father) UpdateActive(active bool) {
	f.Active = active
	f.UpdatedAt = time.Now()
}
