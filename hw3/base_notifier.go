package main

import "fmt"

type BaseNotifier struct {
	Channel string
}

func (b BaseNotifier) Log(message string) {
	fmt.Printf("[%s] %s\n", b.Channel, message)
}
