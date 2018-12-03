package model

import (
	"time"
)

// TrendType 走势类型
type TrendType uint

// StockTrend 股票走势记录
type StockTrend struct {
	ID        uint `gorm:"primary_key"`
	TrendType TrendType
	Price     int64
	Datetime  time.Time
	Stock     Stock `gorm:"foreignkey:StockCode"`
	StockCode string
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	// DayClose 日收盘价
	DayClose TrendType = iota
	// DayAvg 日均价
	DayAvg
	// HourClose 时收盘价
	HourClose
	// HourAvg 时均价
	HourAvg
)
