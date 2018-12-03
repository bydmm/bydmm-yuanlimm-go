package cron

import (
	"time"
	"yuanlimm-worker/model"

	"github.com/jinzhu/now"
)

func trend(beginning time.Time, trendType model.TrendType) {
	for _, stock := range model.Stocks() {
		count := 0
		model.DB.Model(&model.StockTrend{}).
			Where("stock_code = ?", stock.Code).
			Where("trend_type = ?", trendType).
			Where("created_at > ?", beginning).
			Count(&count)
		if count > 0 {
			continue
		}
		var trade model.Transaction
		err := model.DB.Preload("BuyOrder").
			Where("stock_code = ?", stock.Code).
			Where("pay_type = ?", model.Trade).
			Where("created_at < ?", beginning).
			Order("id DESC").
			First(&trade).Error
		if err != nil {
			continue
		}
		trend := model.StockTrend{
			TrendType: trendType,
			StockCode: stock.Code,
			Price:     trade.BuyOrder.Price,
			Datetime:  beginning.Add(time.Duration(-1) * time.Second),
		}
		model.DB.Create(&trend)
	}
}

// TrendHour 按小时记录价格
func TrendHour() {
	trend(now.BeginningOfHour(), model.HourClose)
}

// TrendDayWorker 记录日价
func TrendDayWorker() {
	trend(now.BeginningOfDay(), model.DayClose)
}
