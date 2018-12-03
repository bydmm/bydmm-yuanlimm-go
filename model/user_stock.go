package model

import "time"

// UserStock 用户持股
type UserStock struct {
	ID        uint `gorm:"primary_key"`
	Balance   int64
	User      User   `gorm:"foreignkey:UserID"`
	UserID    uint   `sql:"index"`
	Stock     Stock  `gorm:"foreignkey:StockCode"`
	StockCode string `sql:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
