package models

import (
	"github.com/shopspring/decimal"
	"strings"
)

type order struct {
	Order string `json:"order"`
	Goods []good `json:"goods"`
}
type good struct {
	Description string          `json:"description"`
	Price       decimal.Decimal `json:"price"`
}
type discount struct {
	match       string
	reward      decimal.Decimal
	reward_type string
}

// todo ручка которая принемает дискаунты
// todo ручка которая принмает ордер , при создании идем за скидками, получаем дискаунт и результат в базу
// todo ручка по ид заказа достает расчитанный дискаунт

//ид, матч, реворд и реворд тайп
//таблица заказов -ордер ид и дискаунт

func acceptOrder(o order) {

}

func (c order) acceptDiscount(ds []discount) {
	for _, v := range ds {
		//пишем в базу
	}
}

func (c order) calculateDiscount(ds []discount) decimal.Decimal {
	var result decimal.Decimal
	for _, g := range c.Goods {
		for _, d := range ds {
			if !strings.Contains(g.Description, d.match) {
				continue
			}
			switch d.reward_type {
			case "%":
				result = result.Add(g.Price.Mul(d.reward).Div(decimal.NewFromInt(100)))
			case "pt":
				result = result.Add(d.reward)
			}
		}
	}

	return result
}
