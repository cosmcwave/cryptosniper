package signal

import (
	"cryptosniper/statistic"
	"math"
)

type VolatilityConstant int

const (
	NormalVolatility VolatilityConstant = iota
	VeryHighVolatility
	HighVolatility
)

const (
	VeryHighVolatilityConstant = 1.49
	HighVolatilityConstant = 0.99
)

func VolatilityToText(volatility VolatilityConstant) string {
	vmap := map[VolatilityConstant]string{
		VeryHighVolatility: "very high volatility",
		HighVolatility:     "high volatility",
	}
	return vmap[volatility]
}

func PriceVolatility(period int, series *statistic.TimeSeries) (float64, VolatilityConstant) {
	stdev := statistic.NewStdDev(series)
	adjMean := statistic.NewAdjMean(series)

	volatility := math.Abs(1-(adjMean+stdev)/(adjMean-stdev)) * 100
	if volatility > VeryHighVolatilityConstant {
		return volatility, VeryHighVolatility
	}
	if volatility > HighVolatilityConstant && volatility < VeryHighVolatilityConstant {
		return volatility, HighVolatility
	}
	return 0, NormalVolatility
}
