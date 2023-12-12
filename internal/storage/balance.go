package storage

import "github.com/shopspring/decimal"

var Balance = struct {
	Current   decimal.Decimal `json:"current"`
	Withdrawn decimal.Decimal `json:"withdrawn"`
}{
	Current:   decimal.NewFromFloat32(729.98),
	Withdrawn: decimal.Zero,
}

func AddBalance(value decimal.Decimal) {
	Balance.Current = Balance.Current.Add(value)
}

func SubBalance(value decimal.Decimal) {
	Balance.Current = Balance.Current.Sub(value)
	Balance.Withdrawn = Balance.Withdrawn.Add(value)
}
