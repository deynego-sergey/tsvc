package main

import "time"

type (
	TimeRow []DataItem
)

func (row *TimeRow) Add(v DataItem) {
	*row = append(*row, v)
}

func (row *TimeRow) UpdateWindow(width int64) {
	tc := time.Now().UnixNano()
	lowTreshold := tc - width

	for k, v := range *row {
		if v.Time >= lowTreshold {
			*row = (*row)[k:]
			break
		}
	}
}

func (row *TimeRow) Average() float64 {

	if len(*row) > 0 {
		var s float64
		for _, v := range *row {
			s += v.Value
		}
		return s / float64(len(*row))
	}
	return 0
}
