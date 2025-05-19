package subscription

import (
	"weatherApi/internal/common/constants"
	"weatherApi/internal/repository/user"
	"time"

	"gorm.io/gorm"
)

type SubscriptionModel struct {
	gorm.Model

	City      string              `gorm:"size:32;not null"`
	Frequency constants.Frequency `gorm:"type:VARCHAR(10);not null;default:'daily'"`

	UserID uint
	User   user.UserModel `gorm:"foreignKey:UserID"`

	IsConfirmed   bool      `gorm:"default:false"`
	ConfirmToken  string    `gorm:"uniqueIndex;size:64"`
	TokenExpires  time.Time `gorm:"not null"`
	ConfirmedAt   *time.Time
}

func (SubscriptionModel) TableName() string {
	return "subscriptions"
}
