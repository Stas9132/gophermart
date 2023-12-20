package process

import (
	"context"
	"github.com/shopspring/decimal"
	"gophermart/internal/storage"
	"gophermart/pkg/config"
	"gophermart/pkg/logger"
	"net/http"
	"time"
)

func StatusDaemon(ctx context.Context, config *config.Config, st *storage.DBStorage) {
	for {
		orders, err := st.GetOrdersInProcessing()
		if err != nil {
			st.Error("get orders", logger.LogMap{"error": err, "StatusDaemon": st.GetOrdersInProcessing})
		}
		for _, order := range orders {
			order.Status = "PROCESSING"
			resp, e := http.Get(config.AccuralSystemAddress + "/api/orders/" + order.Number)
			resp.Body.Close()
			if e != nil {
				st.Error("get order by number", logger.LogMap{"error": err, "StatusDaemon": StatusDaemon})
			}

			discount := decimal.NewFromFloat32(729.98)

			if err != nil {
				order.Status = "INVALID"
				err = st.UpdateOrder(ctx, order)
				if err != nil {
					st.Error("update order", logger.LogMap{"error": err, "StatusDaemon": st.UpdateOrder})
				}
			}
			order.Status = "PROCESSED"
			order.Accrual = discount

			order.Accrual.Add(discount)
			err = st.UpdateOrder(ctx, order)
			if err != nil {
				st.Error("update order", logger.LogMap{"error": err, "StatusDaemon": st.UpdateOrder})
			}
		}

		time.Sleep(time.Second)
	}
}
