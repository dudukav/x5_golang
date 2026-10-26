package main

import (
	"errors"
	"testing"
)

type stubNotifier struct {
	called  bool
	message string
	err     error
}

func (s *stubNotifier) Send(message string) error {
	s.called = true
	s.message = message
	return s.err
}

func TestDispatchSendsMessageToNotifier(t *testing.T) {
	first := &stubNotifier{}
	second := &stubNotifier{}
	message := "Заказ готов"

	err := Dispatch([]Notifier{first, second}, message)
	if err != nil {
		t.Fatalf("Dispatch() returned an unexpected error: %v", err)
	}

	for i, stub := range []*stubNotifier{first, second} {
		if !stub.called {
			t.Errorf("notifier %d: Send() was not called", i)
		}
		if stub.message != message {
			t.Errorf("notifier %d: message = %q, want %q", i, stub.message, message)
		}
	}
}

func TestNotifierRejectsEmptyMessage(t *testing.T) {
	tests := []struct {
		name     string
		notifier Notifier
	}{
		{
			name: "email",
			notifier: &EmailNotifier{
				BaseNotifier: BaseNotifier{Channel: "email"},
				Address:      "student@example.com",
			},
		},
		{
			name: "push",
			notifier: &PushNotifier{
				BaseNotifier: BaseNotifier{Channel: "push"},
				DeviceID:     "device-42",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.notifier.Send("")
			if !errors.Is(err, ErrEmptyMessage) {
				t.Errorf("Send(\"\") error = %v, want %v", err, ErrEmptyMessage)
			}
		})
	}
}

func TestDispatchReturnsNotifierError(t *testing.T) {
	wantErr := errors.New("channel is unavailable")
	failing := &stubNotifier{err: wantErr}

	err := Dispatch([]Notifier{failing}, "Заказ готов")

	if !errors.Is(err, wantErr) {
		t.Fatalf("Dispatch() error = %v, want wrapped error %v", err, wantErr)
	}
}
