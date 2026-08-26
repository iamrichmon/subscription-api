package model

import "time"

type SubscriptionStatus string

const (
	FreeSubscription    SubscriptionStatus = "free"
	ProSubscription     SubscriptionStatus = "pro"
	PremiumSubscription SubscriptionStatus = "premium"
)

type User struct {
	ID        uint               `gorm:"primaryKey" json:"id"`
	Name      string             `gorm:"not null" json:"name" binding:"required"`
	Email     string             `gorm:"unique; not null" json:"email" binding:"required`
	Password  string             `gorm:"not null" json:"-" binding:"required,min=12`
	Plan      SubscriptionStatus `gorm:"type:varchar(20);not null;default:free" json:"plan"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
