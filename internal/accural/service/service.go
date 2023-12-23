package service

import (
	"context"
	"github.com/shopspring/decimal"
	"gophermart/internal/accural/storage"
	"gophermart/pkg/logger"
	"log"
	"strings"
)

type OrderManager struct {
	db *storage.DBStorage
	logger.Logger
}

func NewOrderManager(db *storage.DBStorage, logger logger.Logger) *OrderManager {
	return &OrderManager{db: db, Logger: logger}
}

type AccuralService interface {
	GetCalculatedDiscountByOrderID(orderID string) (decimal.Decimal, error)
	AcceptOrder(ctx context.Context, order Order) error
	AcceptDiscount(ctx context.Context, discount storage.Discount) error
	GetAllDiscounts(ctx context.Context) ([]storage.Discount, error)
}

type Order struct {
	Order string `json:"Order"`
	Goods []Good `json:"goods"`
}
type Good struct {
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
}

func New(st *storage.DBStorage, logger logger.Logger) AccuralService {
	return NewOrderManager(st, logger)
}

func (om OrderManager) GetCalculatedDiscountByOrderID(orderID string) (decimal.Decimal, error) {
	var result decimal.Decimal

	rows, err := om.db.Conn.Query(context.Background(), "select * from aorders")
	if err != nil {
		om.db.Error("om.db.Conn.Query(orderID) error", logger.LogMap{"error": err})
		return result, err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Err(); err != nil {
			om.db.Error("rows.Scan(&result) error", logger.LogMap{"error": err})
			return result, err
		}
		var a, b, c any
		log.Println(rows.Scan(&a, &b, &c), a, b, c, orderID)
	}

	return result, nil
}
func (om OrderManager) AcceptOrder(ctx context.Context, order Order) error {
	discounts, err := om.GetAllDiscounts(ctx)
	if err != nil {
		om.Error("om.GetAllDiscounts(ctx) error", logger.LogMap{"error": err})
		return err
	}
	dc, err := order.CalculateDiscount(discounts)
	if err != nil {
		dc = decimal.NewFromInt(-1)
	}

	_, err = om.db.Conn.Exec(context.Background(), "INSERT INTO aorders(order_id, discount) VALUES ($1, $2)", order.Order, dc)
	if err != nil {
		om.db.Logger.Error("unable insert into orders table", logger.LogMap{"error": err, "AcceptOrder": "discounts"})
		return err
	}

	return nil
}
func (om OrderManager) AcceptDiscount(ctx context.Context, discount storage.Discount) error {

	_, err := om.db.Conn.Exec(ctx, "INSERT INTO discounts(match, reward, reward_type) VALUES ($1, $2, $3)", discount.Match, discount.Reward, discount.RewardType)
	if err != nil {
		om.db.Logger.Error("unable insert into discounts table", logger.LogMap{"error": err, "AcceptDiscounts": "discount"})
		return err
	}

	om.db.Logger.Info("Data inserted successfully")
	return nil
}

func (o Order) CalculateDiscount(discounts []storage.Discount) (decimal.Decimal, error) {
	var result decimal.Decimal

	for _, g := range o.Goods {
		for _, d := range discounts {
			if !strings.Contains(g.Description, d.Match) {
				continue
			}
			switch d.RewardType {
			case "%", "":
				result = result.Add(g.Price.Mul(d.Reward).Div(decimal.NewFromInt(100)))
			case "pt":
				result = result.Add(d.Reward)
			}
		}
	}

	return result, nil
}

func (om OrderManager) GetAllDiscounts(ctx context.Context) ([]storage.Discount, error) {
	var discounts []storage.Discount

	rows, err := om.db.Conn.Query(ctx, "SELECT match, reward, reward_type FROM discounts")
	if err != nil {
		om.db.Logger.Error("om.db.Conn.Query() error", logger.LogMap{"error": err})
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
		om.db.Logger.Error("rows.Err() error", logger.LogMap{"error": err})
		return nil, err
	}

	return discounts, nil
}
