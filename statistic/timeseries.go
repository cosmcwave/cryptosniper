package statistic

import (
	"github.com/adshao/go-binance/v2"
	"strconv"
)

type Candle struct {
	Timestamp int64
	High      float64
	Low       float64
	Open      float64
	Close     float64
	Volume    float64
}

type TimeSeries []*Candle

func NewTimeSeries(period int, klines []*binance.Kline) *TimeSeries {
	var series TimeSeries
	index := len(klines) - (period + 1)

	for i := index; i < len(klines); i++ {
		highPrice, _ := strconv.ParseFloat(klines[i].High, 64)
		lowPrice, _ := strconv.ParseFloat(klines[i].Low, 64)
		openPrice, _ := strconv.ParseFloat(klines[i].Open, 64)
		closePrice, _ := strconv.ParseFloat(klines[i].Close, 64)
		volume, _ := strconv.ParseFloat(klines[i].Volume, 64)

		series = append(series, &Candle{Timestamp: klines[i].OpenTime, High: highPrice, Low: lowPrice, Open: openPrice, Close: closePrice, Volume: volume})
	}

	return &series
}
