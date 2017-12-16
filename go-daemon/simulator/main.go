package main

import (
	"fmt"
)

func main() {
	fmt.Println("this is simulator begin")
	var str string
	for {
		fmt.Print("$>")
		fmt.Scan(&str)
		if str == "e" {
			break
		}
		fmt.Println(str)
	}
	fmt.Println("this is simulator end")

}
