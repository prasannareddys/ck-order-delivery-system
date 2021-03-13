package order

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func NewOrder() Order {
	return Order{}
}

// Reads Orders from orders file
// map each order to struct Order
// returns []Orders

func (ord Order) ReadOrders(path string) ([]Order, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var orders []Order
	err = json.Unmarshal(data, &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
