package cron

import (
	"fmt"
	"yuanlimm-worker/model"

	"github.com/jinzhu/gorm"
)

func findMatchingSaleOrders(buyOrder model.StockOrder) []model.StockOrder {
	var orders []model.StockOrder
	model.DB.
		Where("status = ?", model.Padding).
		Where("stock_code = ?", buyOrder.StockCode).
		Where("price <= ?", buyOrder.Price).
		Where("user_id != ?", buyOrder.UserID).
		Where("amount < 0").
		Order("id ASC").
		Find(&orders)
	return orders
}

func findBuyOrders(stock model.Stock) []model.StockOrder {
	var orders []model.StockOrder
	model.DB.
		Where("status = ?", model.Padding).
		Where("stock_code = ?", stock.Code).
		Where("amount > 0").
		Order("price DESC, id ASC").
		Find(&orders)
	return orders
}

// int64Abs 取绝对值
func int64Abs(value int64) int64 {
	if value < 0 {
		return -value
	}
	return value
}

// MatchingTrade 撮合交易
func MatchingTrade() {
	for _, stock := range model.Stocks() {
		buyOrders := findBuyOrders(stock)
		for _, buyOrder := range buyOrders {
			matchedSaleOrders := findMatchingSaleOrders(buyOrder)
			for _, saleOrder := range matchedSaleOrders {
				helper := TradeOrderHelper{
					BuyOrderID:  buyOrder.ID,
					SaleOrderID: saleOrder.ID,
				}
				helper.Init()
				helper.Trade()
			}
		}
	}
}

// TradeOrderHelper 交易用的帮助类
type TradeOrderHelper struct {
	BuyOrderID      uint
	SaleOrderID     uint
	dealStockAmount int64
	buyOrder        *model.StockOrder
	saleOrder       *model.StockOrder
}

// Init 初始化
func (helper *TradeOrderHelper) Init() {
	var buyOrder model.StockOrder
	model.StockOrderPreload().First(&buyOrder, helper.BuyOrderID)
	helper.buyOrder = &buyOrder
	var saleOrder model.StockOrder
	model.StockOrderPreload().First(&saleOrder, helper.SaleOrderID)
	helper.saleOrder = &saleOrder
}

// Buyer 买家
func (helper *TradeOrderHelper) Buyer() model.User {
	return helper.buyOrder.User
}

// Seller 卖家
func (helper *TradeOrderHelper) Seller() model.User {
	return helper.saleOrder.User
}

// TotalFund 总资金
func (helper *TradeOrderHelper) TotalFund() int64 {
	var wallet model.UserWallet
	model.DB.Where("user_id = ?", helper.buyOrder.UserID).First(&wallet)
	return wallet.Balance
}

// TotalStock 卖方总股份
func (helper *TradeOrderHelper) TotalStock() int64 {
	var hold model.UserStock
	model.DB.Where("user_id = ?", helper.saleOrder.UserID).First(&hold)
	return hold.Balance
}

// DealPrice 成交价
func (helper *TradeOrderHelper) DealPrice() int64 {
	return helper.buyOrder.Price
}

// SaleAmount 可卖数量
func (helper *TradeOrderHelper) SaleAmount() int64 {
	wantSale := int64Abs(int64Abs(helper.saleOrder.Amount) - int64Abs(helper.saleOrder.FinishedAmount))
	if helper.TotalStock() <= wantSale {
		return helper.TotalStock()
	}
	return wantSale
}

// BuyAmount 最大可买数量
func (helper *TradeOrderHelper) BuyAmount() int64 {
	wantBuy := helper.buyOrder.Amount - helper.buyOrder.FinishedAmount
	maxBuyAmount := helper.TotalFund() / helper.DealPrice()
	if maxBuyAmount <= wantBuy {
		return maxBuyAmount
	}
	return wantBuy
}

// DealStockAmount 成交数量
func (helper *TradeOrderHelper) DealStockAmount() int64 {
	if helper.dealStockAmount > 0 {
		return helper.dealStockAmount
	}
	if helper.BuyAmount() >= helper.SaleAmount() {
		helper.dealStockAmount = helper.SaleAmount()
	} else {
		helper.dealStockAmount = helper.BuyAmount()
	}
	return helper.dealStockAmount
}

// Pay 付款总数
func (helper *TradeOrderHelper) Pay() int64 {
	return helper.DealStockAmount() * helper.DealPrice()
}

// TransferCoin 一手交钱
func (helper *TradeOrderHelper) TransferCoin(db *gorm.DB) error {
	return db.Create(&model.Transaction{
		Type:             model.CoinTransaction,
		PayerID:          helper.Buyer().ID,
		PayeeID:          helper.Seller().ID,
		Amount:           helper.Pay(),
		StockCode:        helper.buyOrder.StockCode,
		BuyStockOrderID:  helper.buyOrder.ID,
		SaleStockOrderID: helper.saleOrder.ID,
		PayType:          model.Trade,
	}).Error
}

// TransferStock 一手交货
func (helper *TradeOrderHelper) TransferStock(db *gorm.DB) error {
	return db.Create(&model.Transaction{
		Type:             model.StockTransaction,
		PayerID:          helper.Seller().ID,
		PayeeID:          helper.Buyer().ID,
		Amount:           helper.DealStockAmount(),
		StockCode:        helper.buyOrder.StockCode,
		BuyStockOrderID:  helper.buyOrder.ID,
		SaleStockOrderID: helper.saleOrder.ID,
		PayType:          model.Trade,
	}).Error
}

// FinishBuyOrder 处理买单
func (helper *TradeOrderHelper) FinishBuyOrder(db *gorm.DB) error {
	helper.buyOrder.FinishedAmount = helper.buyOrder.FinishedAmount + helper.DealStockAmount()
	if int64Abs(helper.buyOrder.FinishedAmount) >= int64Abs(helper.buyOrder.Amount) {
		helper.buyOrder.Status = model.Success
	}
	return db.Save(helper.buyOrder).Error
}

// FinishSaleOrder 处理卖单
func (helper *TradeOrderHelper) FinishSaleOrder(db *gorm.DB) error {
	helper.saleOrder.FinishedAmount = helper.saleOrder.FinishedAmount - helper.DealStockAmount()
	if int64Abs(helper.saleOrder.FinishedAmount) >= int64Abs(helper.saleOrder.Amount) {
		helper.saleOrder.Status = model.Success
	}
	return db.Save(helper.saleOrder).Error
}

// Transaction 交易事务
func (helper *TradeOrderHelper) Transaction() error {
	transaction := model.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()
	if err := helper.TransferCoin(transaction); err != nil {
		transaction.Rollback()
		return fmt.Errorf("TransferCoin: %s", err.Error())
	}
	if err := helper.TransferStock(transaction); err != nil {
		transaction.Rollback()
		return fmt.Errorf("TransferStock: %s", err.Error())
	}
	if err := helper.FinishBuyOrder(transaction); err != nil {
		transaction.Rollback()
		return fmt.Errorf("FinishBuyOrder: %s", err.Error())
	}
	if err := helper.FinishSaleOrder(transaction); err != nil {
		transaction.Rollback()
		return fmt.Errorf("FinishSaleOrder: %s", err.Error())
	}
	return transaction.Commit().Error
}

// AutoCancel 自动取消
func (helper *TradeOrderHelper) AutoCancel() {
	if helper.DealPrice() > helper.TotalFund() {
		model.DB.Model(helper.buyOrder).Update("status", model.Cancel)
	}
	if helper.TotalStock() < 1 {
		model.DB.Model(helper.saleOrder).Update("status", model.Cancel)
	}
}

// Trade 当前价格
func (helper *TradeOrderHelper) Trade() {
	helper.AutoCancel()
	if !helper.Valid() {
		return
	}
	err := helper.Transaction()
	if err != nil {
		panic(err.Error())
	}
}

// Valid 验证
func (helper *TradeOrderHelper) Valid() bool {
	return helper.BuyOrderValid() && helper.SaleOrderValid() && helper.DealStockAmount() > 0
}

// BuyOrderValid 验证买单
func (helper *TradeOrderHelper) BuyOrderValid() bool {
	if helper.buyOrder.Status != model.Padding {
		return false
	}
	if helper.buyOrder.Amount <= 0 || helper.buyOrder.Price == 0 {
		return false
	}
	return true
}

// SaleOrderValid 验证卖单
func (helper *TradeOrderHelper) SaleOrderValid() bool {
	if helper.saleOrder.Status != model.Padding {
		return false
	}
	if helper.saleOrder.Amount >= 0 || helper.saleOrder.Price == 0 {
		return false
	}
	return true
}
