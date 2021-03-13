package shelf

import (
	"fmt"
	"github.com/Propertyfinder/ck-order-delivery-system/pkg/order"
	"time"
)

func GetAllShelves() []*Shelf {
	 return []*Shelf{
		&HotShelf,
		&ColdShelf,
		&FrozenShelf,
	}
}


func isShelfAvailable(s *Shelf) bool {
	oCount := len(s.Orders)
	return oCount < s.Capacity
}

func CalculateOrderShelfLife(ord order.Order, s *Shelf)  float64 {
	var sdm float64
	sdm = 1
	if s.Name == "Overflow shelf" {
		sdm = 2
	}
	//current time - created time covert to seconds
	orderAge := ((time.Now()).Sub(ord.CreateTime)).Seconds() //calculated in seconds
	return (ord.ShelfLife - orderAge - ord.DecayRate * orderAge * sdm)/ord.ShelfLife
}

// return order index to delete
func findShelfOrderToDelete(s *Shelf)  *order.Order {

	for _,o := range s.Orders {
		v := CalculateOrderShelfLife(o, s)
		if v == 0{
			return &o
		}
	}

	// can improve here to get next decading shelf order
	return nil
}

func GetOrderByTemperature(s *Shelf, t string)  (*order.Order, bool){
	for _,o := range s.Orders {
		if o.Temp == t {
			return &o, true
		}
	}
	return nil, false
}

// can improve by using shelf id
func GetShelfByName(n string)  *Shelf {
	switch n {
	case "Hot Shelf":
		return &HotShelf
	case "Cold Shelf":
		return &ColdShelf
	case "Frozen Shelf":
		return &FrozenShelf
	case "Overflow Shelf":
		return &FrozenShelf
	}
	return nil
}

func GetOrderShelf(ord order.Order) *Shelf {

	fmt.Println("*******",ord)
s := GetShelfByName(ord.AssignedShelfName)
fmt.Println(s.Name)
	for _,o := range s.Orders {
		if o.ID == ord.ID {
			return  s
		}
	}
	return nil
}