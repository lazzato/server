package models

import (
	"time"
)

type User struct {
	ID            uint           `gorm:"primaryKey"`
	Email         string         `gorm:"type:varchar(100);unique;not null"`
	Name          string         `gorm:"type:varchar(100)"`
	Phone         string         `gorm:"type:varchar(30)"`
	Role          string         `gorm:"type:varchar(20);not null;default:employee"` // 'owner' or 'employee'
	RestaurantID  *uint          `gorm:"index"`              // Nullable for owners
	Restaurant    *Restaurant    `gorm:"foreignKey:RestaurantID"`
	GoogleID      *string        `gorm:"type:text;unique"`   // Nullable for employees
	HireDate      *time.Time
	SalaryType    string         `gorm:"type:varchar(20)"`
	SalaryAmount  float64        `gorm:"type:numeric(12,2)"`
	IsActive      bool           `gorm:"not null;default:true"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
