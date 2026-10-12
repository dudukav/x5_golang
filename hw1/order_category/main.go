package main

import "fmt"

func main() {
	var orderPrice int
	fmt.Scan(&orderPrice)

	if orderPrice < 0 {
		fmt.Println("Невалидная сумма заказа")
		return
	}

	if orderPrice < 10000 {
		fmt.Println("small")
	} else if orderPrice >= 10000 && orderPrice < 50000 {
		fmt.Println("regular")
	} else {
		fmt.Println("large")
	}
}