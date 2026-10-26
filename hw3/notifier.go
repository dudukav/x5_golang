package main

type Notifier interface {
	Send(message string) error
}
