package main

import "fmt"

func Dispatch(notifiers []Notifier, message string) error {
	for _, notifier := range notifiers {
		if err := notifier.Send(message); err != nil {
			return fmt.Errorf("send message: %w", err)
		}
	}

	return nil
}
