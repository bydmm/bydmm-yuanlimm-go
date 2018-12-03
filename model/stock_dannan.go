package model

import "time"

// StockDannan 老公
type StockDannan struct {
	ID        uint  `gorm:"primary_key"`
	User      User  `gorm:"foreignkey:UserID"`
	UserID    uint  `sql:"index"`
	Stock     Stock `gorm:"foreignkey:StockCode"`
	StockCode string
	CreatedAt time.Time
	UpdatedAt time.Time
}
