package signal

import (
	"cryptosniper/statistic"
)

func PriceVolatility(period int, series *statistic.TimeSeries) float64 {
	stdev := statistic.NewStdDev(series)
	adjMean := statistic.NewAdjMean(series)

	return adjMean / (adjMean + stdev)
}
