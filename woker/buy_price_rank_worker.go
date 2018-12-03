package woker

// BuyPriceRankWorker 买价排序
type BuyPriceRankWorker struct {
	name string
}

// Perform 执行函数
func (w BuyPriceRankWorker) Perform() error {
	// var stocks []model.Stock
	// err := model.DB.Find(&stocks).Error
	// if err != nil {
	// 	return err
	// }
	// for _, stock := range stocks {
	// 	var order model.StockOrder
	// 	err := model.DB.
	// 		Where(&model.StockOrder{Status: model.Padding}).
	// 		Where("stock_code = ?", stock.Code).
	// 		Where("amount > 0").
	// 		Order("price").
	// 		First(&order).Error
	// 	if err == nil {
	// 		model.DB.Model(&stock).Update("buy_price", order.Price)
	// 	}
	// }
	return nil
}
