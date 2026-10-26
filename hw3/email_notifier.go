package main

import "fmt"

type EmailNotifier struct {
	BaseNotifier
	Address string
}

var _ Notifier = (*EmailNotifier)(nil)

func (n *EmailNotifier) Send(message string) error {
	if len(message) == 0 {
		return ErrEmptyMessage
	}

	n.Log(message)
	fmt.Printf("Email отправлен на %s: %s\n", n.Address, message)

	return nil
}
