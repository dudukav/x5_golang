package main

import "fmt"

type PushNotifier struct {
	BaseNotifier
	DeviceID string
}

var _ Notifier = (*PushNotifier)(nil)

func (n *PushNotifier) Send(message string) error {
	if len(message) == 0 {
		return ErrEmptyMessage
	}

	n.Log(message)
	fmt.Printf("Push отправлен на устройство %s: %s\n", n.DeviceID, message)
	return nil
}
