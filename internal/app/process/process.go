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
		orders, err := st.GetOrdersInProcessing()
		if err != nil {
			log.Printf("get orders: %v\n", err)
		}
		log.Println("process orders : ", len(orders))

		for _, order := range orders {
			if order.Status == "NEW" {
				order.Status = "PROCESSING"
				log.Printf("order status processing created")
				discount, err := accural.GetCalculatedDiscountByOrderID(order.Number)
				if err != nil {
					order.Status = "INVALID"
					err = st.UpdateOrder(ctx, order)
					if err != nil {
						log.Printf("order status invalid %v", err)
					}
				}
				order.Accrual.Add(discount)
				err = st.UpdateOrder(ctx, order)
				if err != nil {
					log.Printf("order status processed %v", err)
				}

			}
		}

		time.Sleep(time.Second)
	}
}
