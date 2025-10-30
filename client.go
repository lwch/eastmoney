package eastmoney

import (
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/lwch/logging"
)

type Client struct {
	cli     *http.Client
	ua      string
	preCall func(*Client) error
}

func New(opts ...configFn) *Client {
	var cfg config
	cfg.timeout = 10 * time.Second
	for _, opt := range opts {
		opt(&cfg)
	}
	var tr http.Transport
	if cfg.disableKeepAlive {
		tr.DisableKeepAlives = true
	}
	if len(cfg.proxy) > 0 {
		proxyURL, err := url.Parse(cfg.proxy)
		if err != nil {
			logging.Error("invalid proxy URL: %v", err)
			return nil
		}
		tr.Proxy = http.ProxyURL(proxyURL)
	}
	return &Client{
		cli: &http.Client{
			Transport: &tr,
			Timeout:   cfg.timeout,
		},
		ua:      cfg.ua,
		preCall: cfg.preCall,
	}
}

func (cli *Client) call(api string, args url.Values, data any) error {
	if cli.preCall != nil {
		if err := cli.preCall(cli); err != nil {
			return err
		}
	}
	url := "https://push2his.eastmoney.com/api" + api + "?" + args.Encode()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if len(cli.ua) > 0 {
		req.Header.Set("User-Agent", cli.ua)
	}
	resp, err := cli.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return errors.New(string(data))
	}
	str, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var ret struct {
		Code int `json:"rc"`
		Data any `json:"data"`
	}
	ret.Data = data
	err = json.Unmarshal(str, &ret)
	if err != nil {
		logging.Error("eastmoney error: %v\n%s", err, string(str))
		return err
	}
	if ret.Code != 0 {
		return errors.New("eastmoney error: " + string(str))
	}
	return nil
}

func (cli *Client) callPaged(url string, args url.Values, append func(any) int) error {
	req, err := http.NewRequest(http.MethodGet, url+"?"+args.Encode(), nil)
	if err != nil {
		return err
	}
	if len(cli.ua) > 0 {
		req.Header.Set("User-Agent", cli.ua)
	}
	resp, err := cli.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return errors.New(string(data))
	}
	str, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var ret struct {
		Code int `json:"rc"`
		Data struct {
			Total int `json:"total"`
			List  any `json:"diff"`
		} `json:"data"`
	}
	err = json.Unmarshal(str, &ret)
	if err != nil {
		logging.Error("eastmoney error: %v\n%s", err, string(str))
		return err
	}
	if ret.Code != 0 {
		return errors.New("eastmoney error: " + string(str))
	}
	perPage := append(ret.Data.List)
	pages := int(math.Ceil(float64(ret.Data.Total) / float64(perPage)))
	for i := 2; i <= pages; i++ {
		args.Set("pn", strconv.Itoa(i))
		req, err := http.NewRequest(http.MethodGet, url+"?"+args.Encode(), nil)
		if err != nil {
			return err
		}
		resp, err := cli.cli.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			data, _ := io.ReadAll(resp.Body)
			return errors.New(string(data))
		}
		str, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var ret struct {
			Code int `json:"rc"`
			Data struct {
				List any `json:"diff"`
			} `json:"data"`
		}
		err = json.Unmarshal(str, &ret)
		if err != nil {
			logging.Error("eastmoney error: %v\n%s", err, string(str))
			return err
		}
		if ret.Code != 0 {
			return errors.New("eastmoney error: " + string(str))
		}
		append(ret.Data.List)
	}
	return nil
}

func (cli *Client) SetProxy(proxy string) error {
	if cli.cli.Transport == nil {
		cli.cli.Transport = &http.Transport{}
	}
	if len(proxy) == 0 {
		cli.cli.Transport.(*http.Transport).Proxy = nil
		return nil
	}
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return err
	}
	cli.cli.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
	return nil
}
