package statistic

func NewAdjMean(series *TimeSeries) float64 {
	var avg float64
	for _, candle := range *series {
		avg += ((candle.Close + candle.High + candle.Low)/3)
	}
	return avg/float64(len(*series))
}
