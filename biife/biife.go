package biife

import (
	"fmt"
	. "github.com/nntaoli-project/goex"
	"log"
	"net/http"
)

var (
	marketBaseUrl = "https://api.biife.com"
)

type Biife struct {
	client *http.Client
	accesskey,
	secretkey string
}

func New(client *http.Client, accesskey, secretkey string) *Biife {
	return &Biife{client: client, accesskey: accesskey, secretkey: secretkey}
}

func (g *Biife) GetTicker(currency CurrencyPair) (*Ticker, error) {
	uri := fmt.Sprintf("%s/openapi/quote/v1/ticker/24hr?symbol=%s", marketBaseUrl, currency.ToSymbol(""))

	log.Println("uri:", uri)
	resp, err := HttpGet(g.client, uri)
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}

	/*
		{
		time: 1614920299438,
		symbol: "NT3USDT",
		bestBidPrice: "6.0339",
		bestAskPrice: "6.13",
		volume: "10509.67",
		quoteVolume: "63217.798477",
		lastPrice: "6.0754",
		highPrice: "6.7474",
		lowPrice: "5.908",
		openPrice: "6.6328"
		}
	*/
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

/*
	获取Symbol信息
*/
func (client *Biife) Symbols(args map[string]interface{}) (Response, http.Header, error) {
	return Response{}, nil, nil
}
