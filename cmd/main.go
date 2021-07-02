package main

import (
	"fmt"

	"github.com/user0608/kcheck"
)

func main() {
	k := kcheck.New()
	edad := "67"
	num := "40"
	if err := k.Traget("num plus", edad, num).Ok(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Valido")
}
