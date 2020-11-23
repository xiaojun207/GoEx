package jbex

import (
	"errors"
	"fmt"
	. "github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/internal/logger"
	"log"
	"net/http"
	"sort"
)

var (
	marketBaseUrl = "https://api.jbex.com"
)

type Jbex struct {
	client *http.Client
	accesskey,
	secretkey string
}

func New(client *http.Client, accesskey, secretkey string) *Jbex {
	return &Jbex{client: client, accesskey: accesskey, secretkey: secretkey}
}

func (g *Jbex) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	panic("not implement")
}
func (g *Jbex) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	panic("not implement")
}
func (g *Jbex) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (g *Jbex) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (g *Jbex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implement")
}
func (g *Jbex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (g *Jbex) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implement")
}
func (g *Jbex) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implement")
}
func (g *Jbex) GetAccount() (*Account, error) {
	panic("not implement")
}

/**
{
  "time": 1538725500422,
  "symbol": "ETHBTC",
  "bestBidPrice": "4.00000200",
  "bestAskPrice": "4.00000200",
  "lastPrice": "4.00000200",
  "openPrice": "99.00000000",
  "highPrice": "100.00000000",
  "lowPrice": "0.10000000",
  "volume": "8913.30000000"
}
*/
func (g *Jbex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	uri := fmt.Sprintf("%s/openapi/quote/v1/ticker/24hr?symbol=%s", marketBaseUrl, currency.ToSymbol(""))

	log.Println("uri:", uri)
	resp, err := HttpGet(g.client, uri)
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}

	return &Ticker{
		Pair: currency,
		Date: uint64(ToInt(resp["time"]) / 1000),
		Last: ToFloat64(resp["lastPrice"]),
		Sell: ToFloat64(resp["bestAskPrice"]),
		Buy:  ToFloat64(resp["bestBidPrice"]),
		High: ToFloat64(resp["highPrice"]),
		Low:  ToFloat64(resp["lowPrice"]),
		Vol:  ToFloat64(resp["volume"]),
	}, nil
}

func (g *Jbex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	resp, err := HttpGet(g.client, fmt.Sprintf("%s/openapi/quote/v1/depth?symbol=%s", marketBaseUrl, currency.ToSymbol("")))
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}

	bids, _ := resp["bids"].([]interface{})
	asks, _ := resp["asks"].([]interface{})

	dep := new(Depth)

	for _, v := range bids {
		r := v.([]interface{})
		dep.BidList = append(dep.BidList, DepthRecord{ToFloat64(r[0]), ToFloat64(r[1])})
	}

	for _, v := range asks {
		r := v.([]interface{})
		dep.AskList = append(dep.AskList, DepthRecord{ToFloat64(r[0]), ToFloat64(r[1])})
	}

	sort.Sort(sort.Reverse(dep.AskList))

	return dep, nil
}

/**
[
  [
    1499040000000,      // 开盘时间
    "0.01634790",       // 开盘价
    "0.80000000",       // 最高价
    "0.01575800",       // 最低价
    "0.01577100",       // 收盘价
    "148976.11427815",  // 交易量
    1499644799999,      // 收盘时间
    "2434.19055334",    // Quote asset数量
    308,                // 交易次数
    "1756.87402397",    // Taker buy base asset数量
    "28.46694368"       // Taker buy quote asset数量
  ]
]
*/
func (g *Jbex) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	granularity := "SECOND"
	switch period {
	case KLINE_PERIOD_1MIN:
		granularity = "1m"
	case KLINE_PERIOD_5MIN:
		granularity = "5m"
	case KLINE_PERIOD_15MIN:
		granularity = "15m"
	case KLINE_PERIOD_1H, KLINE_PERIOD_60MIN:
		granularity = "60m"
	case KLINE_PERIOD_6H:
		granularity = "6h"
	case KLINE_PERIOD_1DAY:
		granularity = "1d"
	default:
		return nil, errors.New("unsupport the kline period")
	}

	uri := fmt.Sprintf("%s/openapi/quote/v1/klines?symbol=%s&interval=%s", marketBaseUrl, currency.ToSymbol(""), granularity)
	log.Println("uri:", uri)

	resp, err := HttpGet3(g.client, uri, nil)
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}

	var klines []Kline
	for i := 0; i < len(resp); i++ {
		k, ok := resp[i].([]interface{})
		if !ok {
			logger.Error("data format err data =", resp[i])
			continue
		}
		klines = append(klines, Kline{
			Pair:      currency,
			Timestamp: ToInt64(k[0]),
			Open:      ToFloat64(k[1]),
			High:      ToFloat64(k[2]),
			Low:       ToFloat64(k[3]),
			Close:     ToFloat64(k[4]),
			Vol:       ToFloat64(k[5]),
		})
	}

	return klines, nil
}

/**
[
	{
		price: "0.19101",
		time: 1606115698968,
		qty: "21",
		isBuyerMaker: false
	},
]
*/
//非个人，整个交易所的交易记录
func (g *Jbex) GetTrades(currency CurrencyPair, since int64) ([]Trade, error) {
	if since == 0 {
		since = 30
	}
	uri := fmt.Sprintf("%s/openapi/quote/v1/trades?symbol=%s&limit=%d", marketBaseUrl, currency.ToSymbol(""), since)
	log.Println("uri:", uri)
	resp, err := HttpGet3(g.client, uri, map[string]string{})
	if err != nil {
		return nil, err
	}

	var trades []Trade
	for _, v := range resp {
		m := v.(map[string]interface{})
		ty := SELL
		if m["isBuyerMaker"].(bool) {
			ty = BUY
		}
		trades = append(trades, Trade{
			Tid:    ToInt64(m["id"]),
			Type:   ty,
			Amount: ToFloat64(m["qty"]),
			Price:  ToFloat64(m["price"]),
			Date:   ToInt64(m["time"]),
			Pair:   currency,
		})
	}

	return trades, nil
}

func (g *Jbex) GetExchangeName() string {
	return JBEX
}
