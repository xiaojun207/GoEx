package biife

import (
	"github.com/nntaoli-project/goex"
	"log"
	"net/http"
	"testing"
)

var client *Biife
var userId string = ""

func TestBiife_Time(t *testing.T) {
	apiKey := ""
	secret := ""

	client = New(http.DefaultClient, apiKey, secret)
	resp, err := client.GetTicker(goex.BTC_USDT)
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)
}
