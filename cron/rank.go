package cron

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"time"
	"yuanlimm-worker/cache"
	"yuanlimm-worker/model"
)

// FreshPrice 刷新股票的最后交易价格
func FreshPrice() {
	for _, stock := range model.Stocks() {
		key := fmt.Sprintf("stock:%s:last_price", stock.Code)
		var trade model.Transaction
		err := model.DB.Preload("BuyOrder").
			Where("stock_code = ?", stock.Code).
			Where("pay_type = ?", model.Trade).
			Order("id DESC").
			First(&trade).Error
		if err == nil {
			cache.RedisClient.Set(key, trade.BuyOrder.Price, 0)
		} else {
			cache.RedisClient.Set(key, 0, 0)
		}
	}
}

// BuyPriceRank 刷新股票的最新买价
func BuyPriceRank() {
	for _, stock := range model.Stocks() {
		var order model.StockOrder
		err := model.DB.
			Where("status = ?", model.Padding).
			Where("stock_code = ?", stock.Code).
			Where("amount > 0").
			Order("price DESC").
			First(&order).Error
		if err == nil {
			model.DB.Model(&stock).Update("buy_price", order.Price)
		} else {
			model.DB.Model(&stock).Update("buy_price", 0)
		}
	}
}

// SalePriceRank 刷新股票的最新买价
func SalePriceRank() {
	for _, stock := range model.Stocks() {
		var order model.StockOrder
		err := model.DB.
			Where("status = ?", model.Padding).
			Where("stock_code = ?", stock.Code).
			Where("amount < 0").
			Order("price").
			First(&order).Error
		if err == nil {
			model.DB.Model(&stock).Update("sale_price", order.Price)
		} else {
			model.DB.Model(&stock).Update("sale_price", 0)
		}
	}
}

// HotRankResult 单个股票属性
type HotRankResult struct {
	Count     int64
	Hot       float64
	StockCode string
}

// ByHot 股票排名
type ByHot []HotRankResult

func (a ByHot) Len() int           { return len(a) }
func (a ByHot) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHot) Less(i, j int) bool { return a[i].Hot > a[j].Hot }

// HotRank 热门排名
func HotRank() {
	results := []HotRankResult{}
	err := model.DB.Table("transactions").
		Select("COUNT(*) as count, stock_code").
		Where("created_at > ?", time.Now().Add(time.Duration(-30)*time.Minute)).
		Group("stock_code").
		Scan(&results)
	if err != nil && len(results) == 0 {
		return
	}

	for index := 0; index < len(results); index++ {
		result := results[index]
		var stock model.Stock
		err := model.DB.Where("code = ?", result.StockCode).First(&stock).Error
		if err != nil {
			return
		}
		time := (time.Now().Second() - stock.CreatedAt.Second()) / (24 * 60 * 60)
		hot := math.Pow(float64(result.Count)/float64(time+2), 1.8)
		results[index].Hot = hot
	}

	sort.Sort(ByHot(results))

	for index := 0; index < len(results); index++ {
		model.DB.Table("stocks").Where("code = ?", results[index].StockCode).Update("order", index)
	}
}

// MarketValueRank 市值排序
func MarketValueRank() {
	for _, stock := range model.Stocks() {
		price := stock.Price()
		share := stock.TotalShare()
		marketValue := price * share
		model.DB.Table("stocks").Where("code = ?", stock.Code).Update("market_value", marketValue)
	}
}

// CapitalRank 财富结构
type CapitalRank struct {
	ID            uint   `json:"id"`
	Capital       int64  `json:"capital"`
	TrueLoveStock string `json:"true_love_stock"`
}

// ByCapital 财富排名
type ByCapital []CapitalRank

func (a ByCapital) Len() int           { return len(a) }
func (a ByCapital) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCapital) Less(i, j int) bool { return a[i].Capital > a[j].Capital }

// UserRank 用户排序
func UserRank() {
	var ranks []CapitalRank
	var users []model.User
	model.DB.Where("id != ?", 1).Find(&users)
	for _, user := range users {
		trueLoveStock := ""
		capital := int64(0)
		var holdings []model.UserStock
		model.DB.Preload("Stock").Where("user_id =? AND balance > 0", user.ID).Order("balance DESC").Find(&holdings)
		for index, hold := range holdings {
			if index == 0 {
				trueLoveStock = hold.Stock.Code
			}
			capital += hold.Stock.Price() * hold.Balance
		}
		rank := CapitalRank{
			ID:            user.ID,
			Capital:       capital,
			TrueLoveStock: trueLoveStock,
		}
		ranks = append(ranks, rank)
	}
	sort.Sort(ByCapital(ranks))
	if rankJSON, err := json.Marshal(ranks[:50]); err == nil {
		cache.RedisClient.Set("user:capital:rank", rankJSON, 0)
	}
}
