package eastmoney

import (
	"encoding/csv"
	"errors"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// KLine15Min 获取15分钟K线数据
func (cli *Client) KLine15Min(code string, from, to time.Time, right right) ([]Tick, error) {
	args := make(url.Values)
	// f1（代码）、f2（市场）、f3（名称）、f4（保留小数点位数）、f5（上市到如今交易日总数）、f6（K 线的前 1 日收盘价，即开始日期的前 1日收盘价）、f7（昨日收盘价，即结束日期的前 1 日收盘价）、f8（杂项，比如 4 是指 ETF 基金，7 是 A 股）
	args.Set("fields1", "f1,f2,f3,f4,f5,f6,f7")
	// f51（日期时间）、f52（开盘价）、f53（收盘价）、f54（最高价）、f55（最低价）、f56（成交量）、f57（成交额）、f58（振幅%）、f59（涨跌幅%）、f60（涨跌额）、f61（换手率%）
	args.Set("fields2", "f51,f52,f53,f54,f55,f56,f57")
	// 15分钟K线
	args.Set("klt", "15")
	// 0（不复权）、1（前复权）、2（后复权）
	args.Set("fqt", right.arg())
	if code[0] == '6' {
		args.Set("secid", "1."+code)
	} else {
		args.Set("secid", "0."+code)
	}
	args.Set("beg", from.Format("20060102"))
	args.Set("end", to.Format("20060102"))
	var data struct {
		Code    string   `json:"code"`      // 代码
		Market  int      `json:"market"`    // 市场
		Name    string   `json:"name"`      // 名称
		Decimal int      `json:"decimal"`   // 保留小数点位数
		Total   int      `json:"dktotal"`   // 上市到如今交易日总数
		Day1    float32  `json:"prePrice"`  // 前一交易日收盘价
		Day2    float32  `json:"preKPrice"` // 前二交易日收盘价
		KLines  []string `json:"klines"`
	}
	err := cli.call("/qt/stock/kline/get", args, &data)
	if err != nil {
		return nil, err
	}
	var lines string
	for _, line := range data.KLines {
		lines += line + "\n"
	}
	var ret []Tick
	r := csv.NewReader(strings.NewReader(lines))
	for {
		record, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return ret, nil
			}
			return nil, err
		}
		if len(record) < 7 {
			return nil, errors.New("invalid data")
		}
		t, err := time.ParseInLocation("2006-01-02 15:04", record[0], time.Local)
		if err != nil {
			return nil, err
		}
		open, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		close, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}
		high, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}
		low, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}
		volume, err := strconv.ParseInt(record[5], 10, 64)
		if err != nil {
			return nil, err
		}
		turn, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			return nil, err
		}
		ret = append(ret, Tick{
			Open:   open,
			Close:  close,
			High:   high,
			Low:    low,
			Volume: volume,
			Turn:   turn,
			Time:   dateTime(t),
		})
	}
}
