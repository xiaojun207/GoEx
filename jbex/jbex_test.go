package jbex

import (
	"github.com/nntaoli-project/goex"
	"net/http"
	"testing"
)

var jbex = New(http.DefaultClient, "", "")

func Test_GetTicker(t *testing.T) {
	ticker, err := jbex.GetTicker(goex.BTC_USDT)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}

func Test_GetDepth(t *testing.T) {
	dep, err := jbex.GetDepth(1, goex.BTC_USDT)

	t.Log("err=>", err)
	t.Log("asks=>", dep.AskList)
	t.Log("bids=>", dep.BidList)
}

func TestJbex_GetTrades(t *testing.T) {
	res, err := jbex.GetTrades(goex.BTC_USDT, 10)
	t.Log("err=>", err)
	t.Log("res=>", res)
}
