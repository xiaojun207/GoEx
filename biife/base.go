package biife

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const API_VERSION string = "v1"

type Request struct {
	entryPoint string
	apiKey     string
	secret     string
	httpClient *http.Client
}
type Response struct {
	Status int
	Data   interface{}
	Msg    string
}

func (request *Request) Init(entryPoint string, apiKey string, secret string, proxyUrl string) {
	if !strings.HasSuffix(entryPoint, "/") {
		entryPoint = entryPoint + "/"
	}
	request.entryPoint = entryPoint
	request.apiKey = apiKey
	request.secret = secret
	if len(proxyUrl) != 0 {
		proxy := func(_ *http.Request) (*url.URL, error) {
			return url.Parse(proxyUrl)
		}
		httpTransport := &http.Transport{Proxy: proxy}
		request.httpClient = &http.Client{Transport: httpTransport}
	} else {
		request.httpClient = &http.Client{}
	}
}

func (request *Request) GenerateSignature(data map[string]string) string {
	var params []string
	for key, value := range data {
		params = append(params, key+"="+value)
	}
	sort.Strings(params)
	paramsStr := strings.Join(params, "&")
	println("Params:%v", paramsStr)
	key := []byte(request.secret)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(paramsStr))
	return hex.EncodeToString(mac.Sum(nil))
}

func (request *Request) _request(method string, uri string, args map[string]interface{}) (Response, http.Header, error) {
	if _, ok := args["timeout"]; ok {
		args["timeout"] = 10
	}
	//dataKey := "data"
	//
	//if _, ok := args[dataKey]; !ok {
	//	args[dataKey] = make(map[string]string)
	//	_, ok := args["params"]
	//	if  ok {
	//		args[dataKey] = args["params"]
	//	}
	//}
	data := make(map[string]string)
	for key, value := range args {
		data[key] = value.(string)
	}
	//data := args[dataKey].(map[string]string)
	data["timestamp"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	data["signature"] = request.GenerateSignature(data)

	println(data["timestamp"])

	var httpRequest *http.Request
	if method == "GET" {
		urlParams := url.Values{}
		Url, err := url.Parse(uri)
		if err != nil {
		}
		for key, value := range data {
			urlParams.Set(key, value)
		}
		Url.RawQuery = urlParams.Encode()
		uri = Url.String()
		httpRequest, err = http.NewRequest(method, uri, nil)
		if err != nil {
			log.Println("http failed error", err.Error())
		}
	} else {
		var err error
		var params []string
		for key, value := range data {
			params = append(params, key+"="+value)
		}
		httpRequest, err = http.NewRequest(method, uri, strings.NewReader(strings.Join(params, "&")))
		httpRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		if err != nil {
			log.Println("http failed error", err.Error())
		}
	}

	httpRequest.Header.Add("X-ACCESS-KEY", request.apiKey)
	httpRequest.Header.Add("X-BH-APIKEY", request.apiKey)

	headers, ok := args["headers"].(map[string]string)
	if ok {
		for key, value := range headers {
			httpRequest.Header.Add(key, value)
		}
	}

	log.Printf("Request method:[%v], url:[%v], headers:[%v], params: [%v]", method, uri, httpRequest.Header, data)

	beginTime := time.Now()
	resp, err := request.httpClient.Do(httpRequest)
	if err != nil {
		log.Println("Request Failed ", err.Error())
		return Response{100, nil, ""}, nil, err
	}
	defer resp.Body.Close()
	endTime := time.Now()

	log.Printf("Response code:[%v], cost:[%v], uri: [%v], body: %v", resp.StatusCode, endTime.Sub(beginTime), uri, resp.Body)

	return request.handleResponse(resp), resp.Header, nil
}

func (request *Request) handleResponse(resp *http.Response) Response {
	ret := Response{}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error body", resp)
		return Response{999, nil, ""}
	}
	if resp.StatusCode != 200 {
		return Response{resp.StatusCode, string(bodyBytes), ""}
	}
	println(string(bodyBytes))
	json.Unmarshal(bodyBytes, &ret)
	return ret
}

func (request *Request) Get(path string, args map[string]interface{}) (Response, http.Header, error) {
	uri := request.entryPoint + API_VERSION + path
	return request._request("GET", uri, args)
}

func (request *Request) Post(path string, args map[string]interface{}) (Response, http.Header, error) {
	uri := request.entryPoint + API_VERSION + path
	return request._request("POST", uri, args)
}
