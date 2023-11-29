package models

import (
	"context"
	"github.com/shopspring/decimal"
	"gophermart/internal/accural/storage"
	"log"
	"log/slog"
	"strconv"
	"strings"
)

type OrderManager struct {
	db *storage.DBStorage
}

func NewOrderManager(db *storage.DBStorage) *OrderManager {
	return &OrderManager{db: db}
}

type Order struct {
	Order string `json:"Order"`
	Goods []Good `json:"goods"`
}
type Good struct {
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
}

func (om OrderManager) GetCalculatedDiscountByOrderID(orderID string) (decimal.Decimal, error) {
	var result decimal.Decimal
	id, err := strconv.Atoi(orderID)
	if err != nil {
		log.Println(err)
	}

	err = om.db.Conn.QueryRow(context.Background(), "SELECT SUM(discounts.reward) FROM discounts JOIN orders ON discounts.id = orders.discount_id WHERE orders.order_id = $1", id).Scan(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
func (om OrderManager) AcceptOrder(order Order) error {
	discounts, err := om.GetAllDiscounts()
	if err != nil {
		log.Println(err)
		return err
	}
	dc := order.CalculateDiscount(discounts)

	o, _ := strconv.Atoi(order.Order)
	_, err = om.db.Conn.Exec(context.Background(), "INSERT INTO discounts (order_id, discount_id) VALUES ($1, $2)", o, dc)
	if err != nil {
		om.db.Logger.Error("unable insert into orders table", slog.String("error", err.Error()), slog.String("AcceptOrder", "discounts"))
		return err
	}

	return nil
}
func (om OrderManager) AcceptDiscount(discount storage.Discount) error {

	_, err := om.db.Conn.Exec(context.Background(), "INSERT INTO discounts (match, reward, reward_type) VALUES ($1, $2, $3)", discount.Match, discount.Reward, discount.RewardType)
	if err != nil {
		om.db.Logger.Error("unable insert into discounts table", slog.String("error", err.Error()), slog.String("AcceptDiscounts", "discount"))
		return err
	}

	om.db.Logger.Info("Data inserted successfully")
	return nil
}

func (o Order) CalculateDiscount(discounts []storage.Discount) decimal.Decimal {
	var result decimal.Decimal

	for _, g := range o.Goods {
		for _, d := range discounts {
			if !strings.Contains(g.Description, d.Match) {
				continue
			}
			switch d.RewardType {
			case "%":
				result = result.Add(g.Price.Mul(d.Reward).Div(decimal.NewFromInt(100)))
			case "pt":
				result = result.Add(d.Reward)
			}
		}
	}

	return result
}

func (om OrderManager) GetAllDiscounts() ([]storage.Discount, error) {
	var discounts []storage.Discount

	rows, err := om.db.Conn.Query(context.Background(), "SELECT match, reward, reward_type FROM discounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var discount storage.Discount
		if err := rows.Scan(&discount.Match, &discount.Reward, &discount.RewardType); err != nil {
			return nil, err
		}
		discounts = append(discounts, discount)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return discounts, nil
}
