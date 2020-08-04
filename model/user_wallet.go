package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// UserWallet 用户钱包
type UserWallet struct {
	ID        uint `gorm:"primary_key"`
	Balance   int64
	User      User `gorm:"foreignkey:UserID"`
	UserID    uint `sql:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AddCoin 加钱
func AddCoin(tx *gorm.DB, userID uint, amount int64) (err error) {
	if userID == 0 {
		return
	}
	var wallet UserWallet
	err = tx.FirstOrCreate(&wallet, UserWallet{UserID: userID}).Error
	if err != nil {
		return err
	}
	return tx.Model(&wallet).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
}

// MinusCoin 减钱
func MinusCoin(tx *gorm.DB, userID uint, amount int64) (err error) {
	if userID == 0 {
		return
	}
	var wallet UserWallet
	err = tx.FirstOrCreate(&wallet, UserWallet{UserID: userID}).Error
	if err != nil {
		return err
	}
	return tx.Model(&wallet).UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error
}
