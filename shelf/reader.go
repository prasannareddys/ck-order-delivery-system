package shelf

import (
	"github.com/Propertyfinder/ck-order-delivery-system/orders"
	"time"
)

func GetAllShelves() []*Shelf {
	 return []*Shelf {
		&HotShelf,
		&ColdShelf,
		&FrozenShelf,
	}
}


func isShelfAvailable(s *Shelf) bool {
	oCount := len(s.Orders)
	return oCount < s.Capacity
}

func calculateOrderShelfLife(ord orders.Order, s Shelf)  float64 {
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
func findShelfOrderToDelete(s Shelf)  int {
	ov := make(map[string]float64)

	for i,o := range s.Orders {
		v := calculateOrderShelfLife(o, s)
		if v == 0{
			return i
		}
		ov[o.ID] = v
	}

	// can improve here to get next decading shelf order
	return 0
}

func GetOrderByTemperature(s *Shelf, t string)  (int, bool){
	for i,o := range s.Orders {
		if o.Temp == t {
			return i, true
		}
	}
	return 0, false
}