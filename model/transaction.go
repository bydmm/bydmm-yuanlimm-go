package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// PayType 交易类型
type PayType uint

// Transaction 交易
type Transaction struct {
	ID               uint       `gorm:"primary_key"`
	Type             string     `sql:"index"`
	PayType          PayType    `sql:"index"`
	PayerID          uint       `sql:"index"`
	PayeeID          uint       `sql:"index"`
	BuyOrder         StockOrder `gorm:"foreignkey:BuyStockOrderID"`
	BuyStockOrderID  uint       `sql:"index"`
	SaleOrder        StockOrder `gorm:"foreignkey:SaleStockOrderID"`
	SaleStockOrderID uint       `sql:"index"`
	Stock            Stock      `gorm:"foreignkey:StockCode"`
	StockCode        string     `sql:"index"`
	Amount           int64
	Detail           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

const (
	// Love 挖矿
	Love PayType = iota
	// Trade 交易
	Trade
	// Give 转账
	Give
)

const (
	// CoinTransaction 金币转账
	CoinTransaction = "CoinTransaction"
	// StockTransaction 股票转账
	StockTransaction = "StockTransaction"
)

// TransactionPreload 必要预加载
func TransactionPreload() *gorm.DB {
	return DB.Preload("BuyOrder").Preload("SaleOrder").Preload("Stock")
}

// AfterCreate 创建之后的钩子
func (item *Transaction) AfterCreate(tx *gorm.DB) (err error) {
	switch item.Type {
	case "CoinTransaction":
		err = AddCoin(tx, item.PayeeID, item.Amount)
		if err != nil {
			return err
		}
		err = MinusCoin(tx, item.PayerID, item.Amount)
		if err != nil {
			return err
		}
	case "StockTransaction":
		err = AddStock(tx, item.PayeeID, item.StockCode, item.Amount)
		if err != nil {
			return err
		}
		err = MinusStock(tx, item.PayerID, item.StockCode, item.Amount)
		if err != nil {
			return err
		}
	}
	return
}
