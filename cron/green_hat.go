package cron

import "yuanlimm-worker/model"

// GreenHatChecker 定时检查老公
func GreenHatChecker() {
	for _, stock := range model.Stocks() {
		var hold model.UserStock
		err := model.DB.Where("stock_code = ?", stock.Code).Order("balance DESC").First(&hold).Error
		if err != nil {
			continue
		}
		var oldHusband model.StockDannan
		err = model.DB.Where("stock_code = ?", stock.Code).Order("id DESC").First(&oldHusband).Error
		if err != nil || hold.UserID == oldHusband.UserID {
			continue
		}
		husband := model.StockDannan{
			StockCode: stock.Code,
			UserID:    hold.UserID,
		}
		model.DB.Create(&husband)
	}
}
