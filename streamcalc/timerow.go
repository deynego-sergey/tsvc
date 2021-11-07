package streamcalc

import (
	"time"
)

type (
	TimeRow []DataItem
)

//
func (row *TimeRow) Add(v DataItem) {
	*row = append(*row, v)
}

// UpdateDataFrame - удаляет данные которые вышли за пределы временного окна
//
func (row *TimeRow) UpdateDataFrame(frameWidth int64) {
	lowTreshold := time.Now().UnixNano() - frameWidth
	var _tmp TimeRow
	for _, v := range *row {
		if v.Time >= lowTreshold {
			_tmp = append(_tmp, v)
		}
	}
	if len(_tmp) > 0 {
		*row = _tmp
		_tmp = nil //runtime.GC()
	}

}

// Average - возвращает скользящее среднее
//
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
