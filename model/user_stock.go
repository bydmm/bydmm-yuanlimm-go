package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

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

// AddStock 加股
func AddStock(tx *gorm.DB, userID uint, code string, amount int64) (err error) {
	if userID == 0 {
		return
	}
	var wallet UserStock
	err = tx.FirstOrCreate(&wallet, UserStock{UserID: userID, StockCode: code}).Error
	if err != nil {
		return err
	}
	return tx.Model(&wallet).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
}

// MinusStock 减股
func MinusStock(tx *gorm.DB, userID uint, code string, amount int64) (err error) {
	if userID == 0 {
		return
	}
	var wallet UserStock
	err = tx.FirstOrCreate(&wallet, UserStock{UserID: userID, StockCode: code}).Error
	if err != nil {
		return err
	}
	return tx.Model(&wallet).UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error
}
