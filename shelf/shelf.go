package shelf

import (
	"fmt"
	"github.com/Propertyfinder/ck-order-delivery-system/orders"
)

type Shelf struct {
	Name      string
	Temp      string
	Capacity  int
	Available int
	Orders    []orders.Order
}

// help me with interface
type Service interface {
	AddToShelf(orders.Order) error
	RemoveFromShelf(orders.Order) error
}

func NewShelf() Shelf {
	return Shelf{}
}

var HotShelf Shelf
var ColdShelf Shelf
var FrozenShelf Shelf
var OverflowShelf Shelf

func (s Shelf) AddOrderToShelf(ord orders.Order) error {
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
		fmt.Println("Order do not have temp")
	}
	return nil
}

func addToShelfHandler(ord orders.Order, s *Shelf) error {

	// if hot shelf slot is avaialble
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

func deleteOrderFromShelf(ordIndex int, s *Shelf) error {
	so := s.Orders
	fmt.Println("Deleting order with id %s  from %s shelf ", so[ordIndex].ID, s.Temp)
	mo := append(so[:ordIndex-1], so[ordIndex+1:]...)
	s.Orders = mo
	sl := calculateOrderShelfLife(so[ordIndex], *s)
	fmt.Println("Deleted order with id %s  from %s shelf, shelf life: %f", so[ordIndex].ID, s.Temp, sl)

	// inform delivery about deletion
	return nil
}

func moveOrder() error {
	for _, s := range GetAllShelves() {
		if isShelfAvailable(s) {
			oIndex, b := GetOrderByTemperature(&OverflowShelf, s.Temp)
			if !b  {
				continue
			}
			// add to new shelf
			err := addToShelf(s.Orders[oIndex], s)
			if err != nil {
				continue
			}
			// delete from overflow
			err = deleteOrderFromShelf(oIndex, &OverflowShelf)
			if err != nil {
				return err
			}
			fmt.Println("Moved order with id %s from Overflow shelf to %s shelf", s.Orders[oIndex], s.Temp)
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

func overflowShelfHandler(ord orders.Order, s Shelf) error {

	fmt.Println("Adding order with id %s and tempareture %s to overflow shelf", ord.ID, ord.Temp)
	// 1. add to overflow if slot is available
	if isShelfAvailable(&s) {
		err := addToShelf(ord, &s)
		if err != nil {
			return err
		}
		fmt.Println("Added to overflow shelf %s ", ord.ID)
		return nil
	}

	// 2. move order from overflow to available shelf
	fmt.Println("Overflow shelf is not available for %s ", ord.ID)
	err, mo := moveOrderHandler()
	if err != nil {
		return err
	}
	if mo {
		err = addToShelf(ord, &s)
		if err != nil {
			return err
		}
		fmt.Println("Added to overflow shelf %s ", ord.ID)
		return nil
	}

	// 3. discard random order(can be with less life) from overflow
	fmt.Println("Overflow shelf is not available for %s ", ord.ID)
	do := findShelfOrderToDelete(s)
	err = deleteOrderFromShelf(do, &s)
	if err != nil {
		return err
	}

	err = addToShelf(ord, &s)
	if err != nil {
		return err
	}
	fmt.Println("Added to overflow shelf %s ", ord.ID)
	return nil
}

func addToShelf(ord orders.Order, s *Shelf) error {

	fmt.Println("Adding order with id %s  to %s shelf ", ord.ID, s.Temp)
	if isShelfAvailable(s) { // double check if shelf is not occupied by concurrency run
		uOrds := append(s.Orders, ord)

		s.Orders = uOrds
		sl := calculateOrderShelfLife(ord, *s)
		fmt.Println("Adding order with id %s  to %s shelf, Shelf life: %f", ord.ID, s.Temp, sl)

		return nil
	}

	return fmt.Errorf("Shelf is occupied for order with id : %s", ord.ID)
}
