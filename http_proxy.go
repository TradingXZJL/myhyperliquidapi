package myhyperliquidapi

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"github.com/robfig/cron/v3"
)

type RestProxy struct {
	ProxyUrl     string // 代理的协议IP端口URL
	InfoWeight   ProxyWeight
	ActionWeight ActionProxyWeight
}

type ProxyWeight struct {
	RemainWeight int64 //剩余可用权重
	IsLimited    bool  //是否已被限制
}

func (w *ProxyWeight) restore() {
	w.RemainWeight = 1200
	w.IsLimited = false
}

type ActionProxyWeight struct {
	MaxWeight  int64
	UsedWeight int64
	IsLimited  bool
}

var proxyList = []*RestProxy{}

func GetCurrentProxyList() []*RestProxy {
	return proxyList
}

var UseProxy = false
var WsUseProxy = false

func SetUseProxy(useProxy bool, proxyUrls ...string) {
	UseProxy = useProxy
	var newProxyList []*RestProxy
	for _, proxyUrl := range proxyUrls {
		newProxyList = append(newProxyList, &RestProxy{
			ProxyUrl: proxyUrl,
			InfoWeight: ProxyWeight{
				RemainWeight: 1200,
				IsLimited:    false,
			},
			ActionWeight: ActionProxyWeight{},
		})
	}
	proxyList = newProxyList
}

func SetWsUseProxy(useProxy bool) error {
	if !UseProxy {
		return errors.New("please set UseProxy first")
	}
	WsUseProxy = useProxy
	return nil
}

func isUseProxy() bool {
	return UseProxy
}

func init() {
	c := cron.New(cron.WithSeconds())
	//每1分钟权重清零，状态恢复
	_, err := c.AddFunc("0 */1 * * * *", func() {
		for _, proxy := range proxyList {
			proxy.InfoWeight.restore()
		}
	})
	if err != nil {
		log.Error(err)
		return
	}
	c.Start()
}

func refreshActionWeight(wallet *Wallet) error {
	h := &MyHyperliquid{}
	infoClient := h.NewInfoRestClient()

	// 立即获取一次
	res, err := infoClient.NewBasicUserRateLimit().User(wallet.WalletAddress).Do()
	if err != nil {
		log.Errorf("get info client rate limit failed: %s", err)
		return err
	}
	maxWeight := res.NRequestsCap
	usedWeight := res.NRequestsUsed
	isLimited := maxWeight-usedWeight <= 0
	for _, p := range proxyList {
		p.ActionWeight = ActionProxyWeight{
			MaxWeight:  maxWeight,
			UsedWeight: usedWeight,
			IsLimited:  isLimited,
		}
	}

	c := cron.New(cron.WithSeconds())
	// 每5分钟获取一次Action权重
	_, err = c.AddFunc("0 */5 * * * *", func() {
		res, err := infoClient.NewBasicUserRateLimit().User(wallet.WalletAddress).Do()
		if err != nil {
			log.Errorf("get info client rate limit failed: %s", err)
			return
		}
		maxWeight := res.NRequestsCap
		usedWeight := res.NRequestsUsed
		isLimited := maxWeight-usedWeight <= 0
		for _, p := range proxyList {
			p.ActionWeight = ActionProxyWeight{
				MaxWeight:  maxWeight,
				UsedWeight: usedWeight,
				IsLimited:  isLimited,
			}
		}
	})
	if err != nil {
		c.Stop()
		return err
	}
	c.Start()

	return nil
}

func getBestProxyAndWeight() (*RestProxy, *ProxyWeight) {
	var maxWeightProxy *RestProxy
	var maxWeight *ProxyWeight

	for _, proxy := range proxyList {
		proxyWeight := &proxy.InfoWeight
		if proxyWeight.IsLimited {
			continue
		}
		if maxWeightProxy == nil {
			maxWeightProxy = proxy
			maxWeight = proxyWeight
			continue
		}
		if proxyWeight.RemainWeight > maxWeight.RemainWeight {
			maxWeightProxy = proxy
			maxWeight = proxyWeight
		}
	}

	return maxWeightProxy, maxWeight
}

// 获取随机代理
func getRandomProxy() (*RestProxy, error) {
	length := len(proxyList)
	if length == 0 {
		return nil, errors.New("proxyList is empty")
	}

	return proxyList[rand.Intn(length)], nil
}

var currentProxy *RestProxy
var currentProxyWeight *ProxyWeight

func Request(url string, reqBody []byte, method RequestType, isGzip bool, ipWeight int64, actionWeight int64) ([]byte, error) {
	return RequestWithHeader(url, reqBody, method, map[string]string{}, isGzip, ipWeight, actionWeight)
}

func RequestWithHeader(urlStr string, reqBody []byte, method RequestType, headerMap map[string]string, isGzip bool, ipWeight int64, actionWeight int64) ([]byte, error) {
	req, err := http.NewRequest(method.String(), urlStr, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headerMap {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	if isGzip {
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
	}

	log.Debug("reqURL: ", req.URL.String())
	if len(reqBody) > 0 {
		log.Debug("reqBody: ", string(reqBody))
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
	}

	if UseProxy {
		currentProxy, currentProxyWeight = getBestProxyAndWeight()
		if currentProxy == nil || currentProxyWeight.RemainWeight <= 0 {
			currentProxyWeight.IsLimited = true
			return nil, errors.New("all proxy ip weight limit reached")
		}

		url_i := url.URL{}
		bestProxy, _ := url_i.Parse(currentProxy.ProxyUrl)

		reqProxy := &http.Transport{}
		reqProxy.Proxy = http.ProxyURL(bestProxy)                        // set proxy
		reqProxy.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // set ssl

		client.Transport = reqProxy
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body := resp.Body

	// 判断返回值是否为json
	isJson := resp.Header.Get("Content-Type") == "application/json"
	if !isJson {
		respStr, err := io.ReadAll(body)
		if err != nil {
			return nil, err
		}
		log.Errorf("rest接口错误: %s", strings.TrimSpace(string(respStr)))
		return nil, errors.New("rest接口错误: " + strings.TrimSpace(string(respStr)))
	}

	log.Debugf("respHeader: %s", resp.Header)
	if resp.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	log.Debug("respBody: ", string(data))

	// 基础权重回填
	if UseProxy {
		currentProxy.InfoWeight.RemainWeight -= ipWeight
		currentProxy.ActionWeight.UsedWeight += actionWeight
		log.Debugf("Proxy: %s, InfoWeight: %+v, ActionWeight: %+v", currentProxy.ProxyUrl, currentProxy.InfoWeight, currentProxy.ActionWeight)
	}

	return data, nil
}

func (w *ProxyWeight) CalculateAdditionalWeight(weight int64) {
	if currentProxy == nil {
		return
	}
	currentProxy.InfoWeight.RemainWeight -= weight
	if currentProxy.InfoWeight.RemainWeight <= 0 {
		currentProxy.InfoWeight.IsLimited = true
	}
	log.Warnf("CalculateAdditionalWeight: %d", weight)
	log.Warnf("Proxy: %s, InfoWeight: %+v, ActionWeight: %+v", currentProxy.ProxyUrl, currentProxy.InfoWeight, currentProxy.ActionWeight)
}
