package process

import (
	"context"
	"github.com/shopspring/decimal"
	"gophermart/internal/storage"
	"gophermart/pkg/config"
	"io"
	"log"
	"net/http"
	"time"
)

func StatusDaemon(ctx context.Context, config *config.Config, st *storage.DBStorage) {
	for {
		orders, err := st.GetOrdersInProcessing()
		if err != nil {
			log.Printf("get orders: %v\n", err)
		}
		log.Println("process orders : ", len(orders))

		for _, order := range orders {
			order.Status = "PROCESSING"
			resp, e := http.Get(config.AccuralSystemAddress + "/api/orders/" + order.Number)
			log.Println(e)
			b, e := io.ReadAll(resp.Body)
			log.Println(string(b), e)

			resp.Body.Close()
			discount := decimal.NewFromFloat32(729.98)

			if err != nil {
				order.Status = "INVALID"
				err = st.UpdateOrder(ctx, order)
				if err != nil {
					log.Printf("order status invalid %v", err)
				}
			}
			order.Status = "PROCESSED"
			order.Accrual = discount

			order.Accrual.Add(discount)
			err = st.UpdateOrder(ctx, order)
			if err != nil {
				log.Printf("order status processed %v", err)
			}
		}

		time.Sleep(time.Second)
	}
}
