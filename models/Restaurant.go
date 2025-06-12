package models

import "time"

type Restaurant struct {
	ID               uint      `gorm:"primaryKey"`
	OwnerID          uint      `gorm:"not null"`                   // FK to users.id
	Name             string    `gorm:"not null"`
	Address          string
	TrialStartedAt   time.Time `gorm:"not null;default:current_timestamp"`
	TrialExpiresAt   time.Time
	ContractSignedAt *time.Time
	ContractActive   bool      `gorm:"not null;default:false"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}