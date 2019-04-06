package main

import (
	"Calla/core"
	"fmt"
)

func main() {
	if err := core.Run(); err != nil {
		fmt.Println(err)
	}
}
