package service

import (
	"bytes"
	"crypto/sha512"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
	"yuanlimm-worker/cache"
	"yuanlimm-worker/model"
	"yuanlimm-worker/serializer"

	"github.com/jinzhu/now"
	yaml "gopkg.in/yaml.v2"
)

// Wish 许愿词期待的JSON
type Wish struct {
	CheerWord string `yaml:"cheer_word"`
	LovePower string `yaml:"love_power"`
}

// SuperWishsService 获取WPS下发配置的服务
type SuperWishsService struct {
	CheerWord string `form:"cheer_word" json:"cheer_word"`
	Address   string `form:"address" json:"address" binding:"required"`
	Code      string `form:"code" json:"code" binding:"required"`
	LovePower string `form:"love_power" json:"love_power" binding:"required"`
}

// RawOre 根据条件拼接字符串用来产生哈希值
func (service *SuperWishsService) RawOre() []byte {
	ore := bytes.Join([][]byte{
		[]byte(service.CheerWord),
		[]byte(service.Address),
		[]byte(service.LovePower),
		[]byte(strconv.FormatInt(now.BeginningOfMinute().Unix(), 10)),
		[]byte(service.Code),
	}, []byte{})
	return ore
}

// Hash 生成哈希值
func (service *SuperWishsService) Hash(ore []byte) [64]byte {
	return sha512.Sum512(ore)
}

// MatchWish 匹配是否满足标准
func (service *SuperWishsService) MatchWish() bool {
	ore := service.RawOre()
	hard := int(service.Hard())
	bin := service.Hash(ore)
	zero := (hard / 8)
	for index := 1; index <= zero; index++ {
		if bin[len(bin)-index] != 0 {
			return false
		}
	}

	residual := (hard % 8)

	if residual > 0 {
		last := bin[len(bin)-(hard/8)-1]
		head := fmt.Sprintf("%08b", last)

		if len(head) < residual {
			return false
		}

		headZero := ""
		for index := 0; index < residual; index++ {
			headZero += "0"
		}

		return head[len(head)-residual:] == headZero
	}
	return true
}

// ReplayAttack 保护重放攻击
func (service *SuperWishsService) ReplayAttack() bool {
	key := fmt.Sprintf("super_wish:%s:%s", service.Code, service.LovePower)
	set, _ := cache.RedisClient.SetNX(key, "true", 2*time.Minute).Result()
	return !set
}

// HardAddition 难度倍数
func (service *SuperWishsService) HardAddition() int64 {
	addition := (service.Hard() - 13) / 2
	if addition > 0 {
		return 1
	}
	return addition
}

// PayStockAmount 支付股票的算法
func (service *SuperWishsService) PayStockAmount() int64 {
	randNumber := rand.Intn(1000)
	if randNumber >= 995 {
		return 233 * service.HardAddition()
	}
	return (rand.Int63n(4) + 1) * service.HardAddition()
}

// PayCoinAmount 支付硬币的算法
func (service *SuperWishsService) PayCoinAmount() int64 {
	randNumber := rand.Intn(1000)
	if randNumber >= 995 {
		return 233000 * service.HardAddition()
	}
	return 1000 * (rand.Int63n(9) + 1) * service.HardAddition()
}

// Hard 难度
func (service *SuperWishsService) Hard() int64 {
	return 24
}

// PayRandStock 给随机股票
func (service *SuperWishsService) PayRandStock(user model.User) (*model.Transaction, error) {
	stock, err := model.RandStock()
	if err != nil {
		return &model.Transaction{}, err
	}
	return service.PayStock(user, stock)
}

// WishDetail 保存许愿词
func (service *SuperWishsService) WishDetail(stock model.Stock) string {
	jsonWish := map[string]string{}
	jsonWish["love_power"] = service.LovePower
	if stock.Code == service.Code {
		jsonWish["cheer_word"] = service.CheerWord
	}
	d, err := yaml.Marshal(&jsonWish)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return fmt.Sprintf("---\n%s\n", string(d))
}

// PayStock 给股票
func (service *SuperWishsService) PayStock(user model.User, stock model.Stock) (*model.Transaction, error) {
	transaction := model.Transaction{
		Type:      model.StockTransaction,
		PayerID:   0,
		PayeeID:   user.ID,
		Amount:    service.PayStockAmount(),
		StockCode: stock.Code,
		PayType:   model.Love,
		Detail:    service.WishDetail(stock),
	}
	err := model.DB.Create(&transaction).Error
	return &transaction, err
}

// PayCoin 给钱
func (service *SuperWishsService) PayCoin(user model.User, stock model.Stock) (*model.Transaction, error) {
	transaction := model.Transaction{
		Type:      model.CoinTransaction,
		PayerID:   0,
		PayeeID:   user.ID,
		Amount:    service.PayCoinAmount(),
		StockCode: stock.Code,
		PayType:   model.Love,
		Detail:    service.WishDetail(stock),
	}
	err := model.DB.Create(&transaction).Error
	return &transaction, err
}

// PayWish 奖励
func (service *SuperWishsService) PayWish(user model.User, stock model.Stock) (*model.Transaction, error) {
	randNumber := rand.Intn(100)
	switch {
	case randNumber >= 90:
		return service.PayStock(user, stock)
	case randNumber >= 60 && randNumber < 90:
		return service.PayRandStock(user)
	default:
		return service.PayCoin(user, stock)
	}
}

// User 根据钱包地址找用户
func (service *SuperWishsService) User() (model.User, error) {
	var user model.User
	err := model.DB.Where("address = ?", service.Address).Find(&user).Error
	return user, err
}

// Wish 许愿函数
func (service *SuperWishsService) Wish() interface{} {
	rand.Seed(time.Now().UnixNano())
	if service.ReplayAttack() {
		return serializer.WishResponse{
			Success: false,
			Hard:    service.Hard(),
			Msg:     "许愿失败，重复提交",
		}
	}
	if !service.MatchWish() {
		return serializer.WishResponse{
			Success: false,
			Hard:    service.Hard(),
			Msg:     "CheckLovePower failed",
		}
	}
	user, err := service.User()
	if err != nil {
		return serializer.WishResponse{
			Success: false,
			Hard:    service.Hard(),
			Msg:     "许愿失败，钱包地址不存在",
		}
	}
	stock, err := model.FindStock(service.Code)
	if err != nil {
		return serializer.WishResponse{
			Success: false,
			Hard:    service.Hard(),
			Msg:     "许愿失败，股票不存在",
		}
	}
	transaction, err := service.PayWish(user, stock)
	if err != nil {
		return serializer.WishResponse{
			Success: false,
			Hard:    service.Hard(),
			Msg:     "许愿失败， 内部错误",
		}
	}
	payType := ""
	if transaction.Type == model.CoinTransaction {
		payType = "coin"
	} else if transaction.Type == model.StockTransaction {
		payType = "stock"
	}
	return serializer.WishResponse{
		Success: true,
		Hard:    service.Hard(),
		Type:    payType,
		Amount:  transaction.Amount,
		Stock:   transaction.StockCode,
	}
}
