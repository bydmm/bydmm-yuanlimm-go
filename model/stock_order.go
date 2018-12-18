package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// OrderStatus 订单状态
type OrderStatus uint

// StockOrder 股票订单
type StockOrder struct {
	ID             uint `gorm:"primary_key"`
	Status         OrderStatus
	Price          int64
	Amount         int64
	FinishedAmount int64
	Detail         string `sql:"type:text"`
	User           User   `gorm:"foreignkey:UserID"`
	UserID         uint
	Stock          Stock `gorm:"foreignkey:StockCode"`
	StockCode      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

const (
	// Padding 交易中
	Padding OrderStatus = iota
	// Success 交易成功
	Success
	// Cancel 交易取消
	Cancel
)

// StockOrderPreload 必要预加载
func StockOrderPreload() *gorm.DB {
	return DB.Preload("User").Preload("Stock")
}
