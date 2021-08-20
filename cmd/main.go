package main

import (
	"fmt"

	"github.com/user0608/kcheck"
)

func main() {
	k := kcheck.New()

	edad := "33"
	costo := "32"
	total := "12"
	v := k.Target("num max=8 min=2", edad, costo, total)
	if err := v.Ok(); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Ok")
	}
}
