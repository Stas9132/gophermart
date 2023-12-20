package process

import (
	"context"
	"encoding/json"
	"github.com/shopspring/decimal"
	"gophermart/internal/storage"
	"gophermart/pkg/config"
	l2 "gophermart/pkg/logger"
	"net/http"
	"time"
)

func StatusDaemon(ctx context.Context, config *config.Config, st *storage.DBStorage, logger l2.Logger) {

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			orders, err := st.GetOrdersInProcessing()
			if err != nil {
				logger.Error("Get Orders in processing error", l2.LogMap{"error": err})
			}
			logger.Info("process orders : ", l2.LogMap{"len": len(orders)})

			if err := process(ctx, config, st, orders); err != nil {
				logger.Error("Status daemon canceled", l2.LogMap{"error": ctx.Err()})
			}
		case <-ctx.Done():

			logger.Warn("Status daemon canceled", l2.LogMap{"error": ctx.Err()})
			return
		}
	}
}

func process(ctx context.Context, config *config.Config, st *storage.DBStorage, orders []storage.Order) (err error) {

	for _, order := range orders {
		order.Status = "PROCESSING"
		resp, e := http.Get(config.AccuralSystemAddress + "/api/orders/" + order.Number)
		if e != nil {
			return e
		}
		defer resp.Body.Close()

		type AccrualRespT struct {
			Order   string          `json:"order"`
			Status  string          `json:"status"`
			Accrual decimal.Decimal `json:"accrual"`
		}
		var accrResp AccrualRespT

		if e = json.NewDecoder(resp.Body).Decode(&accrResp); e != nil {
			return e
		}

		discount := accrResp.Accrual

		if err != nil {
			order.Status = "INVALID"
			err = st.UpdateOrder(ctx, order)
			if err != nil {
				return err
			}
		}
		order.Status = "PROCESSED"
		order.Accrual = discount

		order.Accrual.Add(discount)
		err = st.UpdateOrder(ctx, order)
		if err != nil {
			return err
		}
	}

	return nil
}
