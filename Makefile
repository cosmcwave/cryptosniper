.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o ./exec/cryptosniper main.go

.PHONY: run
run:
	./exec/cryptosniper -s "BTCUSDT, ETHUSDT, BNBUSDT, ADAUSDT, DOTUSDT, XRPUSDT, UNIUSDT, THETAUSDT, LTCUSDT, LINKUSDT, BCHUSDT, XLMUSDT, LUNAUSDT, DOGEUSDT, FILUSDT, VETUSDT, TRXUSDT, IOTAUSDT, XMRUSDT, AVAXUSDT, ANKRUSDT, NPXSUSDT, HOTUSDT, ONEUSDT" -e "interval:1m, volume_threshold:0.2, volatility_threshold:0.995"
