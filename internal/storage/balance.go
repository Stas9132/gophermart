package storage

import (
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

var Balance = struct {
	Current   decimal.Decimal `json:"current"`
	Withdrawn decimal.Decimal `json:"withdrawn"`
}{
	Current:   decimal.NewFromFloat32(729.98),
	Withdrawn: decimal.Zero,
}

var lock sync.Mutex

type HistT struct {
	Order       string          `json:"order"`
	Sum         decimal.Decimal `json:"sum"`
	ProcessedAt time.Time       `json:"processed_at"`
}

var Hist []HistT

func AddBalance(value decimal.Decimal) {
	Balance.Current = Balance.Current.Add(value)
}

func SubBalance(value decimal.Decimal) {
	Balance.Current = Balance.Current.Sub(value)
	Balance.Withdrawn = Balance.Withdrawn.Add(value)
	lock.Lock()
	Hist = append(Hist, HistT{Sum: value})
	lock.Unlock()
}
