package generator

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/Jhnvlglmlbrt/wb-order/internal/models"
)

func GenerateOrder() *models.Order {
	orderId := hash32bit()
	var orderCount = 1 + rand.Intn(2)
	items := make([]models.Item, orderCount)
	for i := 0; i < orderCount; i += 1 {
		items[i] = models.Item{
			ChrtId:      0,
			TrackNumber: "",
			Price:       0,
			Rid:         "",
			Name:        "",
			Sale:        0,
			Size:        "",
			TotalPrice:  0,
			NmId:        0,
			Brand:       "",
			Status:      0}
	}
	order := models.Order{
		OrderUid:    orderId[:len(orderId)-15],
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery:    generateOrderDelivery(),
		Payment: models.Payment{
			Transaction:  "",
			RequestId:    "",
			Currency:     "",
			Provider:     "",
			Amount:       0,
			PaymentDt:    0,
			Bank:         "",
			DeliveryCost: 0,
			GoodsTotal:   0,
			CustomFee:    0},
		Items:             items,
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "text",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmId:              99,
		DateCreated:       time.Now().Format(time.RFC3339),
		OofShard:          "1",
	}
	generateOrderItems(&order)
	generateOrderPayment(&order)
	return &order
}

func generateOrderPayment(order *models.Order) {
	currency := []string{"USD", "RUB", "EUR"}
	banks := []string{"sber", "alpha", "tinkoff"}
	var amount float32 = 0
	for i := range order.Items {
		amount += order.Items[i].TotalPrice
	}
	deliveryCost := float32(rand.Intn(1500))
	order.Payment.Transaction = order.OrderUid + order.CustomerId
	order.Payment.RequestId = ""
	order.Payment.Currency = currency[rand.Intn(len(currency))]
	order.Payment.Provider = "wbpay"
	order.Payment.Amount = amount + deliveryCost
	order.Payment.PaymentDt = uint32(1000000000 + rand.Intn(1000000000))
	order.Payment.Bank = banks[rand.Intn(len(banks))]
	order.Payment.DeliveryCost = uint32(deliveryCost)
	order.Payment.GoodsTotal = amount
	order.Payment.CustomFee = 0
}

func generateOrderDelivery() models.Delivery {
	names := []string{"Rykov Maxim", "Ivanov Ivan", "Random Random", "Maximov Maxim"}
	addresses := []string{"Ploshad Mira 15", "Pokryshkina 8", "Orshanskaya 3", "Red Square"}
	newDelivery := models.Delivery{
		Name:    names[rand.Intn(len(names))],
		Phone:   "+" + strconv.Itoa(1000000000+rand.Intn(8000000000)),
		Zip:     strconv.Itoa(100000 + rand.Intn(150000)),
		City:    "Moscow",
		Address: addresses[rand.Intn(len(addresses))],
		Region:  "Moscow",
		Email:   "test@gmail.com",
	}
	return newDelivery
}

func generateOrderItems(order *models.Order) {
	for i := 0; i < len(order.Items); i += 1 {
		amount := float32(1500 + rand.Intn(10000))
		sale := float32(rand.Intn(50))
		total := ((100 - sale) / 100.0) * amount
		totalPrice := math.Round(float64(total)*10) / 10
		order.Items[i] = models.Item{
			ChrtId:      uint32(rand.Intn(1000000)),
			TrackNumber: order.TrackNumber,
			Price:       uint16(amount),
			Rid:         hash32bit()[:len(hash32bit())-15] + order.CustomerId,
			Name:        "Mascaras",
			Sale:        uint16(sale),
			Size:        "0",
			TotalPrice:  float32(totalPrice),
			NmId:        uint32(rand.Intn(1000000)),
			Brand:       "Vivienne Sabo",
			Status:      202,
		}
	}
}

func hash32bit() string {
	sum := md5.Sum([]byte(strconv.Itoa(rand.Intn(150000))))
	return hex.EncodeToString(sum[:])
}
