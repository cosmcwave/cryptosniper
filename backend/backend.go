package backend

import (
	"context"
	"cryptosniper/cache"
	"cryptosniper/cache/memorycache"
	"cryptosniper/extension"
	"cryptosniper/signal"
	"cryptosniper/statistic"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	ApiKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
}

type Backend struct {
	client     *binance.Client
	symbols    []string
	cachePool  map[string]cache.Cache
	extensions map[string]string
	wg         sync.WaitGroup
	m          sync.Mutex
	Out        zerolog.Logger
	Error      zerolog.Logger
}

func New(cfg *Config) *Backend {
	out := zerolog.ConsoleWriter{Out: colorable.NewColorableStdout(), TimeFormat: time.RFC3339}
	err := zerolog.ConsoleWriter{Out: colorable.NewColorableStderr(), TimeFormat: time.RFC3339}

	b := &Backend{
		client:     binance.NewClient(cfg.ApiKey, cfg.SecretKey),
		cachePool:  make(map[string]cache.Cache),
		extensions: make(map[string]string),
		Out:        zerolog.New(out).With().Timestamp().Logger(),
		Error:      zerolog.New(err).With().Timestamp().Logger(),
	}

	b.cachePool["volume"] = memorycache.New()
	b.cachePool["volatility"] = memorycache.New()

	return b
}

func NewConfig(filename string) (*Config, error) {
	var cfg *Config

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewDefaultConfig() (*Config, error) {
	return NewConfig("./.cryptosniper/config.json")
}

func (b *Backend) AddSymbols(symbols ...string) {
	for _, symbol := range symbols {
		symbol = strings.TrimSpace(symbol)
		b.symbols = append(b.symbols, symbol)
		b.Out.Info().Msgf("%v added to the sniping list", symbol)
	}
	fmt.Println("")
}

func (b *Backend) AddExtensions(extensions ...string) {
	for _, e := range extensions {
		e = strings.TrimSpace(e)
		k, v := extension.Parse(e)
		b.extensions[k] = v
		b.Out.Info().Msgf("extension '%v' set to '%v'", k, v)
	}
	fmt.Println("")
}

func (b *Backend) Snipe(ctx context.Context, symbol string) {
	defer b.wg.Done()

	for {
		select {
		case <-ctx.Done():
			b.Out.Warn().Msgf("%v sniper shutting down..", symbol)
			return
		default:
		}

		klines, err := b.client.NewKlinesService().Symbol(symbol).Interval(b.extensions["interval"]).Do(ctx)
		if err != nil {
			b.Error.Error().Msgf("%v", err)
			return
		}

		series := statistic.NewTimeSeries(14, klines)

		b.m.Lock()

		var score int
		volumeThreshold, _ := strconv.ParseFloat(b.extensions["volume_threshold"], 64)
		if v, ok := signal.Volume(series, volumeThreshold); ok {
			cacheVal := b.cachePool["volume"].Get(symbol)
			var t int64

			if cacheVal == nil {
				t = 0
			} else {
				t = cacheVal.(int64)
			}

			if t != (*series)[len(*series)-1].Timestamp {
				b.cachePool["volume"].Set(symbol, (*series)[len(*series)-1].Timestamp)
				b.Out.Info().Msgf("%v high volume (%f)", symbol, math.Abs(1-v)*100)
				score++
			}
		}

		volatilityThreshold, _ := strconv.ParseFloat(b.extensions["volatility_threshold"], 64)
		if volatility := signal.PriceVolatility(14, series); volatility < volatilityThreshold {
			cacheVal := b.cachePool["volatility"].Get(symbol)
			var v int64

			if cacheVal == nil {
				v = 0
			} else {
				v = cacheVal.(int64)
			}

			if v != (*series)[len(*series)-1].Timestamp {
				b.cachePool["volatility"].Set(symbol, (*series)[len(*series)-1].Timestamp)
				b.Out.Info().Msgf("%v high volatility (%f)", symbol, volatility)
				score++
			}
		}
		b.m.Unlock()

		if score == 2 {
			b.Out.Warn().Msgf("%v high volume and volatility", symbol)
		}

		time.Sleep(5 * time.Second)
	}
}

func (b *Backend) Start(ctx context.Context) {
	for _, symbol := range b.symbols {
		b.wg.Add(1)
		go b.Snipe(ctx, symbol)
	}

	b.Out.Info().Msg("cryptosniper started\n\n\n")
	b.wg.Wait()

	fmt.Println()
	b.Out.Info().Msg("cryptosniper terminated")
}
