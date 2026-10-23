package hw2

import (
	"errors"
	"testing"
)

func TestNewOrder(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{name: "positive ID", id: 1},
		{name: "zero ID", id: 0, wantErr: ErrInvalidID},
		{name: "negative ID", id: -1, wantErr: ErrInvalidID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order, err := NewOrder(tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewOrder(%d) error = %v, want %v", tt.id, err, tt.wantErr)
			}

			if tt.wantErr != nil {
				if order != nil {
					t.Fatalf("NewOrder(%d) returned order on error: %#v", tt.id, order)
				}
				return
			}

			if order == nil {
				t.Fatal("NewOrder returned nil order without an error")
			}
			if order.ID != tt.id {
				t.Errorf("order ID = %d, want %d", order.ID, tt.id)
			}
			if order.Status != StatusNames["new"] {
				t.Errorf("order status = %q, want %q", order.Status, StatusNames["new"])
			}
		})
	}
}

func TestEmptyOrderTotalIsZero(t *testing.T) {
	order, err := NewOrder(1)
	if err != nil {
		t.Fatalf("NewOrder() error = %v", err)
	}

	if got := order.Total(); got != 0 {
		t.Errorf("Total() = %v, want 0", got)
	}
}

func TestTotalAfterAddingTwoItems(t *testing.T) {
	order, err := NewOrder(1)
	if err != nil {
		t.Fatalf("NewOrder() error = %v", err)
	}

	items := []Item{
		{Name: "Coffee", Price: 250.5},
		{Name: "Cake", Price: 149.5},
	}
	for _, item := range items {
		if err := order.AddItem(item); err != nil {
			t.Fatalf("AddItem(%#v) error = %v", item, err)
		}
	}

	if got, want := order.Total(), 400.0; got != want {
		t.Errorf("Total() = %v, want %v", got, want)
	}
	if got, want := len(order.Items), 2; got != want {
		t.Errorf("len(Items) = %d, want %d", got, want)
	}
}

func TestAddItemRejectsInvalidData(t *testing.T) {
	tests := []struct {
		name    string
		order   *Order
		item    Item
		wantErr error
	}{
		{
			name:    "nil order",
			item:    Item{Name: "Coffee", Price: 250},
			wantErr: ErrNilOrder,
		},
		{
			name:    "empty item name",
			order:   &Order{ID: 1},
			item:    Item{Price: 250},
			wantErr: ErrInvalidName,
		},
		{
			name:    "negative price",
			order:   &Order{ID: 1},
			item:    Item{Name: "Coffee", Price: -1},
			wantErr: ErrInvalidPrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := 0
			if tt.order != nil {
				before = len(tt.order.Items)
			}

			if err := tt.order.AddItem(tt.item); !errors.Is(err, tt.wantErr) {
				t.Errorf("AddItem() error = %v, want %v", err, tt.wantErr)
			}

			if tt.order != nil && len(tt.order.Items) != before {
				t.Error("AddItem() changed order after invalid input")
			}
		})
	}
}

func TestChangingOrderCopyDoesNotChangeOriginal(t *testing.T) {
	original := Order{ID: 1, Status: StatusNames["new"]}
	copyOrder := original

	copyOrder.ID = 2
	copyOrder.Status = StatusNames["paid"]

	if original.ID != 1 {
		t.Errorf("original ID = %d, want 1", original.ID)
	}
	if original.Status != StatusNames["new"] {
		t.Errorf("original status = %q, want %q", original.Status, StatusNames["new"])
	}
}

func TestStatusNamesMapIsInitialized(t *testing.T) {
	if StatusNames == nil {
		t.Fatal("StatusNames is nil; initialize it with make")
	}

	wantNames := map[string]string{
		"new":       "Новый",
		"paid":      "Оплачен",
		"cancelled": "Отменён",
	}
	for status, want := range wantNames {
		if got := StatusNames[status]; got != want {
			t.Errorf("StatusNames[%q] = %q, want %q", status, got, want)
		}
	}

	const status = "returned"
	StatusNames[status] = "Возвращён"
	t.Cleanup(func() {
		delete(StatusNames, status)
	})

	if got := StatusNames[status]; got != "Возвращён" {
		t.Errorf("StatusNames[%q] = %q, want %q", status, got, "Возвращён")
	}
}
