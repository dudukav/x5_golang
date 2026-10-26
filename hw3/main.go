package main

import "fmt"

func main() {
	notifiers := []Notifier{
		&EmailNotifier{
			BaseNotifier: BaseNotifier{Channel: "email"},
			Address:      "student@example.com",
		},
		&PushNotifier{
			BaseNotifier: BaseNotifier{Channel: "push"},
			DeviceID:     "device-42",
		},
	}

	if err := Dispatch(notifiers, "Заказ готов"); err != nil {
		fmt.Println("ошибка отправки:", err)
	}
}
