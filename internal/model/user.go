package model

import "time"

type SubscriptionStatus string

const (
	FreeSubscription    SubscriptionStatus = "free"
	ProSubscription     SubscriptionStatus = "pro"
	PremiumSubscription SubscriptionStatus = "premium"
)

type User struct {
	ID           uint               `gorm:"primaryKey" json: "id"`
	Name         string             `gorm:"unique; not null" json: "name"`
	Email        string             `gorm:"unique; not null" json: "email"`
	Password     string             `gorm:"unique; not null" json: "password"`
	IsSubscriber SubscriptionStatus `gorm:"type: varchar(20); not null;" json:"is_subscriber"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}
