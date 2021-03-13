package main

import (
	"fmt"
	"github.com/Propertyfinder/ck-order-delivery-system/shelf"
	"log"
	"sync"
	"time"

	"github.com/Propertyfinder/ck-order-delivery-system/orders"
)

var orderFilePath = "./data/orders.json"
func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	// this channel will receive orders from loop at the bottom of file
	orderChannel := make(chan orders.Order, 10)

	// delivery Channel,  this will be courier channel
	deliveryChannel := make(chan orders.Order, 10)

	// Shelf receiving orders go routine
	go func() {
		defer close(deliveryChannel)
		for ord := range orderChannel {
			// Logic to select shelf goes here
			fmt.Println("Order received to be placed on shelf ", ord)
			s := shelf.NewShelf()
			err := s.AddOrderToShelf(ord)
			if err != nil {
				fmt.Println("failed to add order %s to shelf error : %w", ord.ID, err)
			} else {
				// After placing on shelf forward it to courier channel
				deliveryChannel <- ord
			}
		}
	}()

	// delivery receiving orders go routine
	go func() {
		for ord := range deliveryChannel {
			// Logic to deliver order goes here
			time.Sleep(3 * time.Second) // this should be random time b/w 1 to 6 to sleep and then deliver

			// Calculate shelf life here
			fmt.Println("Order received Now we will be delivering and removing from shelf ", ord)
		}
		wg.Done()
	}()

	order := orders.NewOrder()
	ord, err := order.ReadOrders(orderFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var process int

	//change,refactor this read logic to control speed
	var orderThroughPutPerSecond = 2
	for process < len(ord) {
		for i := 0; i < orderThroughPutPerSecond; i++ {
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
}
