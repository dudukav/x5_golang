package sumstatistic

func SumStatistic(kopeyki []int) (min, max, sum int, average float64) {
	if len(kopeyki) == 0 {
		return 0, 0, 0, 0
	}

	min = kopeyki[0]
	max = kopeyki[0]

	for _, price := range kopeyki {
		if price < min {
			min = price
		}
		if price > max {
			max = price
		}
		sum += price
	}

	average = float64(sum) / float64(len(kopeyki))
	return min, max, sum, average
}
