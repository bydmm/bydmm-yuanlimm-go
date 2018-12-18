package model

import "time"

// UserWallet 用户钱包
type UserWallet struct {
	ID        uint `gorm:"primary_key"`
	Balance   int64
	User      User `gorm:"foreignkey:UserID"`
	UserID    uint `sql:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
