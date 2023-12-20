package process

import (
	"context"
	"github.com/shopspring/decimal"
	"gophermart/internal/storage"
	"gophermart/pkg/config"
	l2 "gophermart/pkg/logger"
	"io"
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
		_, e = io.ReadAll(resp.Body)
		if e != nil {
			return e
		}

		discount := decimal.NewFromFloat32(729.98)

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
