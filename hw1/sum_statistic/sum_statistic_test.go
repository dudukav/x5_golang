package sumstatistic_test

import (
	"testing"

	sumstatistic "sum_statistic"
)

func TestSumStatistic(t *testing.T) {
	tests := []struct {
		name        string
		orders      []int
		wantMin     int
		wantMax     int
		wantSum     int
		wantAverage float64
	}{
		{
			name:        "orders from the task",
			orders:      []int{12000, 7500, 31000, 4500},
			wantMin:     4500,
			wantMax:     31000,
			wantSum:     55000,
			wantAverage: 13750,
		},
		{
			name:        "one order",
			orders:      []int{999},
			wantMin:     999,
			wantMax:     999,
			wantSum:     999,
			wantAverage: 999,
		},
		{
			name:        "fractional average",
			orders:      []int{1, 2},
			wantMin:     1,
			wantMax:     2,
			wantSum:     3,
			wantAverage: 1.5,
		},
		{
			name:        "zero-valued orders",
			orders:      []int{0, 0, 0},
			wantMin:     0,
			wantMax:     0,
			wantSum:     0,
			wantAverage: 0,
		},
		{
			name:        "empty slice returns zeros",
			orders:      []int{},
			wantMin:     0,
			wantMax:     0,
			wantSum:     0,
			wantAverage: 0,
		},
		{
			name:        "nil slice returns zeros",
			orders:      nil,
			wantMin:     0,
			wantMax:     0,
			wantSum:     0,
			wantAverage: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax, gotSum, gotAverage := sumstatistic.SumStatistic(tt.orders)

			if gotMin != tt.wantMin {
				t.Errorf("minimum = %d, want %d", gotMin, tt.wantMin)
			}
			if gotMax != tt.wantMax {
				t.Errorf("maximum = %d, want %d", gotMax, tt.wantMax)
			}
			if gotSum != tt.wantSum {
				t.Errorf("sum = %d, want %d", gotSum, tt.wantSum)
			}
			if gotAverage != tt.wantAverage {
				t.Errorf("average = %v, want %v", gotAverage, tt.wantAverage)
			}
		})
	}
}
