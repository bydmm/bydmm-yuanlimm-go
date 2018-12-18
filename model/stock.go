package model

import (
	"fmt"
	"strconv"
	"time"
	"yuanlimm-worker/cache"
)

// Stock 股票
type Stock struct {
	Code        string `gorm:"primary_key"`
	name        string
	Avatar      string `sql:"type:text"`
	MusicLink   string `sql:"type:text"`
	VideoLink   string
	Order       uint
	Tags        string `sql:"type:text"`
	MarketValue int64
	BuyPrice    int64
	SalePrice   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Stocks 返回所有股票
func Stocks() []Stock {
	var stocks []Stock
	DB.Find(&stocks)
	return stocks
}

// RandStock 随机股票
func RandStock() (Stock, error) {
	var stock Stock
	err := DB.Order("RAND()").Find(&stock).Error
	return stock, err
}

// FindStock 根据代号找股票
func FindStock(code string) (Stock, error) {
	var stock Stock
	err := DB.Where("code = ?", code).Find(&stock).Error
	return stock, err
}

// Price 当前价格
func (stock *Stock) Price() int64 {
	key := fmt.Sprintf("go:stock:%s:price", stock.Code)
	cachePrice, _ := cache.Fetch(key, 5*time.Minute, func() string {
		var trade Transaction
		DB.Preload("BuyOrder").
			Where("pay_type = ? AND stock_code = ?", Trade, stock.Code).
			Order("id DESC").First(&trade)
		return strconv.FormatInt(trade.BuyOrder.Price, 10)
	})
	intCache, _ := strconv.ParseInt(cachePrice, 10, 64)
	return intCache
}

// TotalShare 总股份
func (stock *Stock) TotalShare() int64 {
	type SumResult struct {
		Sum int64
	}
	result := SumResult{}
	DB.Table("user_stocks").
		Select("sum(balance) as sum").
		Where("stock_code = ?", stock.Code).
		Scan(&result)
	return result.Sum
}
