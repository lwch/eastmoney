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
	cli *http.Client
}

func New() *Client {
	return &Client{
		cli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (cli *Client) call(api string, args url.Values, data any) error {
	url := "https://push2his.eastmoney.com/api" + api + "?" + args.Encode()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36")
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
