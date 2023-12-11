package service

import (
	"context"
	"github.com/shopspring/decimal"
	"gophermart/internal/accural/storage"
	"gophermart/internal/logger"
	"log"
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
		log.Println(err)
		return result, err
	}

	return result, nil
}
func (om OrderManager) AcceptOrder(ctx context.Context, order Order) error {
	discounts, err := om.GetAllDiscounts(ctx)
	if err != nil {
		log.Println(err)
		return err
	}
	dc := order.CalculateDiscount(discounts)

	_, err = om.db.Conn.Exec(context.Background(), "INSERT INTO orders(order_id, discount_id) VALUES ($1, $2)", order.Order, dc)
	if err != nil {
		om.db.Logger.Error("unable insert into orders table", logger.LogMap{"error": err, "AcceptOrder": "discounts"})
		log.Println(err)
		return err
	}

	return nil
}
func (om OrderManager) AcceptDiscount(ctx context.Context, discount storage.Discount) error {

	_, err := om.db.Conn.Exec(ctx, "INSERT INTO discounts(match, reward, reward_type) VALUES ($1, $2, $3)", discount.Match, discount.Reward, discount.RewardType)
	if err != nil {
		om.db.Logger.Error("unable insert into discounts table", logger.LogMap{"error": err, "AcceptDiscounts": "discount"})
		log.Println(err)
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

func (om OrderManager) GetAllDiscounts(ctx context.Context) ([]storage.Discount, error) {
	var discounts []storage.Discount

	rows, err := om.db.Conn.Query(ctx, "SELECT match, reward, reward_type FROM discounts")
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return nil, err
	}

	return discounts, nil
}
