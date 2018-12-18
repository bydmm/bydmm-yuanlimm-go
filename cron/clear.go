package cron

import (
	"time"
	"yuanlimm-worker/model"
)

// LoveClear 清理过期的许愿
func LoveClear() {
	model.DB.
		Where("pay_type = ?", model.Love).
		Where("created_at < ?", time.Now().AddDate(0, 0, -2)).
		Delete(model.Transaction{})
}

// OrderClear 清理过期的订单
func OrderClear() {
	model.DB.
		Table("stock_orders").
		Where("created_at < ?", time.Now().AddDate(0, -1, 0)).
		Where("status = ?", model.Padding).
		Update("status", model.Cancel)
}
