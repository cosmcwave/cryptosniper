package signal

import (
	"context"
	"cryptosniper/statistic"
	"sort"
)

func Volume(ctx context.Context, series *statistic.TimeSeries, threshold float64) (float64, bool) {
	var volume []float64

	for _, candle := range *series {
		volume = append(volume, candle.Volume)
	}

	sort.Float64s(volume[:len(volume)-2])

	var maVol float64
	for i:=0; i<len(volume)-2; i++ {
		maVol += volume[i]
	}

	if volume != nil {
		q := volume[len(volume)-1] / maVol
		if (q) > threshold {
			return volume[len(volume)-1], true
		}
	}
	return 0, false
}