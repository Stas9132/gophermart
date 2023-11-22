package models

import (
	"context"
	"github.com/shopspring/decimal"
	"gophermart/internal/accural/storage"
	"log/slog"
	"strconv"
	"strings"
)

type OrderManager struct {
	db storage.DBStorage
}
type Order struct {
	Order string `json:"Order"`
	Goods []Good `json:"goods"`
}
type Good struct {
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
}

func (om OrderManager) GetCalculatedDiscountByOrderID(orderId string) (decimal.Decimal, error) {
	var result decimal.Decimal
	id, err := strconv.Atoi(orderId)

	err = om.db.Conn.QueryRow(context.Background(), "SELECT SUM(discounts.reward) FROM discounts JOIN orders ON discounts.id = orders.discount_id WHERE orders.order_id = $1", id).Scan(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
func (om OrderManager) AcceptOrder(order Order) error {
	var orderID int

	err := om.db.Conn.QueryRow(context.Background(), "INSERT INTO orders DEFAULT VALUES RETURNING order_id").Scan(&orderID)
	if err != nil {
		om.db.Logger.Error("unable insert into orders table", slog.String("error", err.Error()), slog.String("AcceptOrder", "discounts"))
		return err
	}

	for _, good := range order.Goods {
		_, err = om.db.Conn.Exec(context.Background(), "INSERT INTO discounts (id) VALUES ($1)", orderID, good.Description)
		if err != nil {
			om.db.Logger.Error("unable insert into orders table", slog.String("error", err.Error()), slog.String("AcceptOrder", "discounts"))
			return err
		}
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

func (om OrderManager) CalculateDiscount(discounts []storage.Discount) (decimal.Decimal, error) {
	var result decimal.Decimal
	orders, err := om.GetAllOrders()
	if err != nil {
		return result, err
	}
	for _, o := range orders {
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

	}

	return result, nil
}
func (om OrderManager) GetAllOrders() ([]Order, error) {
	var orders []Order

	rows, err := om.db.Conn.Query(
		context.Background(),
		"SELECT order FROM orders",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		err := rows.Scan(&order.Order)
		if err != nil {
			return nil, err
		}

		goodsRows, err := om.db.Conn.Query(
			context.Background(),
			"SELECT description, price FROM goods WHERE order = $1",
			order.Order,
		)
		if err != nil {
			return nil, err
		}

		var goods []Good
		for goodsRows.Next() {
			var good Good
			err := goodsRows.Scan(&good.Description, &good.Price)
			if err != nil {
				return nil, err
			}
			goods = append(goods, good)
		}

		if err := goodsRows.Err(); err != nil {
			return nil, err
		}

		order.Goods = goods
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
