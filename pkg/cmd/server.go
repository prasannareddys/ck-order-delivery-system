package cmd

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Prasannareddys/ck-order-delivery-system/pkg/order"
	"github.com/Prasannareddys/ck-order-delivery-system/pkg/shelf"
	"github.com/spf13/cobra"
)

func NewServerCommand() *cobra.Command {
	var ordersPerSecond int
	var filePath string

	processOrdersCmd := &cobra.Command{
		Use:   "server [--ops][--order-file-path]",
		Short: "Order delivery system",
		RunE: func(cmd *cobra.Command, args []string) error {
			var wg sync.WaitGroup
			wg.Add(1)

			// this channel will receive orders from loop at the bottom of file
			orderChannel := make(chan order.Order, 10)

			// delivery Channel,  this will be courier channel
			deliveryChannel := make(chan order.Order, 10)

			// Shelf receiving orders go routine
			go func() {
				defer close(deliveryChannel)
				for ord := range orderChannel {
					// Logic to select shelf goes here
					fmt.Print("\n")
					fmt.Print("\n")
					fmt.Println("Order received to be placed on shelf ", ord)
					s := shelf.NewShelf()
					aOrd, err := s.AddOrderToShelf(ord)
					if err != nil {
						fmt.Errorf("failed to add order to the shelf error : %w",  err)
					} else {
						// After placing on shelf forward it to courier channel
						deliveryChannel <- aOrd
					}
				}
			}()

			// delivery receiving orders go routine
			go func() {
				for ord := range deliveryChannel {
					s := shelf.GetOrderShelf(ord)
					if  s != nil {
						err := shelf.DeleteOrderFromShelf(ord, s)
						if err != nil {
							fmt.Errorf("failed to deviver order to shelf error : %w",  err)
						}
						fmt.Print("\n")
						fmt.Printf("Order with id %s from %s is delivered ", ord.ID, ord.AssignedShelfName)
					}

				}
				wg.Done()
			}()

			newOrd := order.NewOrder()
			ord, err := newOrd.ReadOrders(filePath)
			if err != nil {
				log.Fatal(err)
			}

			var process int

			for process < len(ord) {
				for i := 0; i < ordersPerSecond; i++ {
					o := ord[process]
					o.CreateTime = time.Now() // adding creation time
					orderChannel <- o
					process++
					if process >= len(ord) {
						break
					}
				}
				time.Sleep(time.Second)
			}

			close(orderChannel)
			wg.Wait()
			return nil
		},
	}

	processOrdersCmd.Flags().StringVar(&filePath, "order-file-path", "./data/orders.json", "File path for orders")
	processOrdersCmd.Flags().IntVar(&ordersPerSecond, "ops", 2, "Orders throughput per second")
	return processOrdersCmd
}
