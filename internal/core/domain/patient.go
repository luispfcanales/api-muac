package domain

import (
	"time"

	"github.com/google/uuid"
)

// Patient representa la entidad de paciente en el dominio
type Patient struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`
	Lastname     string    `json:"lastname" gorm:"type:varchar(100);not null"`
	Gender       string    `json:"gender" gorm:"type:varchar(50)"`
	Age          int       `json:"age" gorm:"type:int"`
	DNI          string    `json:"dni" gorm:"column:DNI;type:varchar(20);unique"`
	UrlDNI       string    `json:"url_dni" gorm:"type:text"`
	BirthDate    string    `json:"birth_date" gorm:"type:varchar(20)"`
	ArmSize      string    `json:"arm_size" gorm:"type:varchar(50)"`
	Weight       string    `json:"weight" gorm:"type:varchar(50)"`
	Size         string    `json:"size" gorm:"type:varchar(50)"`
	ConsentGiven bool      `json:"consent_given" gorm:"type:boolean;default:true"`
	ConsentDate  time.Time `json:"consent_date,omitempty" gorm:"type:date"`
	Description  string    `json:"description" gorm:"type:text"`
	CreatedAt    time.Time `json:"created_at,omitempty" gorm:"column:CREATE_AT;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" gorm:"column:UPDATE_AT"`

	Measurements []Measurement `json:"measurements" gorm:"foreignKey:PatientID"`
	UserID       *uuid.UUID    `json:"user_id" gorm:"column:USER_ID;type:uuid"`
	User         *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName especifica el nombre de la tabla para GORM
func (Patient) TableName() string {
	return "patients"
}

// NewPatient crea una nueva instancia de Patient
func NewPatient(
	name, lastname, gender, birthDate, armSize, weight, size, description string,
	age int,
	dni string,
	consentGiven bool,
	createdBy *uuid.UUID,
) *Patient {

	return &Patient{
		ID:           uuid.New(),
		Name:         name,
		Lastname:     lastname,
		Gender:       gender,
		Age:          age,
		DNI:          dni,
		BirthDate:    birthDate,
		ArmSize:      armSize,
		Weight:       weight,
		Size:         size,
		ConsentGiven: consentGiven,
		Description:  description,
		UserID:       createdBy,
		ConsentDate:  time.Now(),
		CreatedAt:    time.Now(),
	}
}

// Validate valida que el paciente tenga los campos requeridos
func (p *Patient) Validate() error {
	if p.Name == "" {
		return ErrEmptyPatientName
	}
	if p.Lastname == "" {
		return ErrEmptyPatientLastName
	}
	return nil
}

// Update actualiza los campos del paciente
func (p *Patient) Update(name, lastname, gender, birthDate, armSize, weight, size, description string, age int, consentGiven bool) {
	p.Name = name
	p.Lastname = lastname
	p.Gender = gender
	p.Age = age
	p.BirthDate = birthDate
	p.ArmSize = armSize
	p.Weight = weight
	p.Size = size
	p.ConsentGiven = consentGiven
	p.Description = description
	p.UpdatedAt = time.Now()
}
