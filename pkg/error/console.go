package error

import (
	"fmt"
)

type Console struct {
}

func NewConsole() Console {
	return Console{}
}


func (c Console) Print(out interface{}) {
	fmt.Print("\n")
	fmt.Println(out)
}
