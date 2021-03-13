package shelf

import (
	"fmt"

	"github.com/Propertyfinder/ck-order-delivery-system/pkg/order"
)

type Shelf struct {
	Name      string
	Temp      string
	Capacity  int
	Available int
	Orders    []order.Order
}

type Service interface {
	AddOrderToShelf(order.Order) error
	DeleteOrderFromShelf(order.Order) error
}

func NewShelf() Shelf{
	return Shelf{}
}

var HotShelf = Shelf {
	Name      :"Hot Shelf",
	Temp      :"Hot",
	Capacity  :10,
}
var ColdShelf = Shelf {
	Name      :"Cold Shelf",
	Temp      :"Cold",
	Capacity  :10,
}
var FrozenShelf = Shelf {
	Name      :"Frozen Shelf",
	Temp      :"Frozen",
	Capacity  :10,
}
var OverflowShelf = Shelf {
	Name      :"Overflow Shelf",
	Temp      :"any",
	Capacity  :15,
}

func (s Shelf) AddOrderToShelf(ord order.Order) error {
	switch oTmp := ord.Temp; oTmp {
	case "hot":
		err := addToShelfHandler(ord, &HotShelf)
		if err != nil {
			return err
		}
		break
	case "cold":
		err := addToShelfHandler(ord, &ColdShelf)
		if err != nil {
			return err
		}
		break
	case "frozen":
		err := addToShelfHandler(ord, &FrozenShelf)
		if err != nil {
			return err
		}
		break
	default:
		fmt.Printf("Order do not have temp")
	}
	return nil
}

func addToShelfHandler(ord order.Order, s *Shelf) error {

	// ifshelf slot is avaialble
	if isShelfAvailable(s) {
		err := addToShelf(ord, s)
		if err != nil {
			return err
		}
		return nil
	}
	// else add to overflow
	err := overflowShelfHandler(ord, OverflowShelf)
	if err != nil {
		return err
	}

	return nil
}

func DeleteOrderFromShelf(o order.Order, s *Shelf) error {
	indexToRemove := -1
	for i, orders := range s.Orders {
		if orders.ID == o.ID {
			indexToRemove = i
		}
	}
	so := s.Orders
	if indexToRemove > -1 {
		fmt.Print("\n")
		fmt.Printf("Deleting order with id %s  from %s shelf ", o.ID, s.Temp)
		mo := append(so[:indexToRemove], so[indexToRemove+1:]...)
		s.Orders = mo
		sl := CalculateOrderShelfLife(o, *s)
		fmt.Print("\n")
		fmt.Printf("Deleted order with id %s  from %s shelf, shelf life: %f", o.ID, s.Temp, sl)
	}
	
	// inform delivery about deletion
	return nil
}

func moveOrder() error {
	for _, s := range GetAllShelves() {
		if isShelfAvailable(s) {
			o, b := GetOrderByTemperature(&OverflowShelf, s.Temp)
			if !b {
				continue
			}
			// add to new shelf
			err := addToShelf(*o, s)
			if err != nil {
				continue
			}
			// delete from overflow
			err = DeleteOrderFromShelf(*o, &OverflowShelf)
			if err != nil {
				return err
			}
			fmt.Printf("Moved order with id %s from Overflow shelf to %s shelf", o.ID, s.Temp)
			return nil
		}
	}
	return fmt.Errorf("Order is not moved")
}

func moveOrderHandler() (error, bool) {
	err := moveOrder()
	if err != nil {
		return err, false
	}

	return nil, true
}

func overflowShelfHandler(ord order.Order, s Shelf) error {
	fmt.Print("\n")
	fmt.Printf("Adding order with id %s and tempareture %s to overflow shelf", ord.ID, ord.Temp)
	// 1. add to overflow if slot is available
	if isShelfAvailable(&s) {
		err := addToShelf(ord, &s)
		if err != nil {
			return err
		}
		fmt.Print("\n")
		fmt.Printf("Added to overflow shelf %s ", ord.ID)
		return nil
	}

	// 2. move order from overflow to available shelf
	fmt.Print("\n")
	fmt.Printf("Overflow shelf is not available for %s ", ord.ID)
	err, mo := moveOrderHandler()
	if err != nil {
		return err
	}
	if mo {
		err = addToShelf(ord, &s)
		if err != nil {
			return err
		}
		fmt.Print("\n")
		fmt.Printf("Added to overflow shelf %s ", ord.ID)
		return nil
	}

	// 3. discard random order(can be with less life) from overflow
	fmt.Printf("Overflow shelf is not available for %s ", ord.ID)
	do := findShelfOrderToDelete(s)
	err = DeleteOrderFromShelf(*do, &s)
	if err != nil {
		return err
	}

	err = addToShelf(ord, &s)
	if err != nil {
		return err
	}
	fmt.Print("\n")
	fmt.Printf("Added to overflow shelf %s ", ord.ID)
	return nil
}

func addToShelf(ord order.Order, s *Shelf) error {
	fmt.Print("\n")
	fmt.Printf("Adding order with id %s  to %s shelf ", ord.ID, s.Temp)
	if isShelfAvailable(s) { // double check if shelf is not occupied by concurrency run
		ord.AssignedShelfName = s.Name
		uOrds := append(s.Orders, ord)

		s.Orders = uOrds
		sl := CalculateOrderShelfLife(ord, *s)
		fmt.Print("\n")
		fmt.Printf("Adding order with id %s  to %s shelf, Shelf life: %f", ord.ID, s.Temp, sl)

		return nil
	}

	return fmt.Errorf("Shelf is occupied for order with id : %s", ord.ID)
}
