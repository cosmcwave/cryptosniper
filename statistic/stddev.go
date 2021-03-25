package statistic

import (
	"math"
)

func NewStdDev(series *TimeSeries) float64 {
	var ss float64
	adjMean := NewAdjMean(series)

	for _, candle := range *series {
		ss += math.Pow((candle.Close - adjMean),2)
	}

	// Sample standard deviation with N-1 denominator
	return math.Sqrt(ss/float64(len(*series)-1))
}
