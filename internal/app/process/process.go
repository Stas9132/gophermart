package process

import (
	"context"
	"gophermart/internal/accural/service"
	"gophermart/internal/storage"
	"log"
	"time"
)

func StatusDaemon(ctx context.Context) {
	st := storage.New()
	accural := service.New()
	for {
		orders, err := st.GetOrders(ctx)
		if err != nil {
			log.Printf("get orders: %v\n", err)
		}

		for _, order := range orders {
			if order.Status == "NEW" {
				order.Status = "PROCESSING"
				discount, err := accural.GetCalculatedDiscountByOrderID(order.Number)
				if err != nil {
					order.Status = "INVALID"
					continue
				}
				order.Accrual.Add(discount)
				err = st.UpdateOrder(ctx, order)
				if err != nil {
					log.Println(err)
				}

			}
		}

		time.Sleep(time.Second)
	}
}
