package eastmoney

import (
	"net/url"
	"strings"
	"time"
)

type IndexItem struct {
	Code string `json:"f12"` // 代码
	Name string `json:"f14"` // 名称
}

// IndexSH 获取上交所指数列表
func (cli *Client) IndexSH() ([]IndexItem, error) {
	args := make(url.Values)
	args.Set("pn", "1")
	args.Set("pz", "100")
	args.Set("po", "1")
	args.Set("np", "1")
	args.Set("fltt", "2")
	args.Set("invt", "2")
	args.Set("fid", "f3")
	args.Set("fs", "m:1+t:1")
	args.Set("fields", "f12,f14")
	var list []IndexItem
	err := cli.callPaged("https://88.push2.eastmoney.com/api/qt/clist/get", args, func(a any) int {
		arr := a.([]any)
		for _, item := range arr {
			var index IndexItem
			index.Code = "sh" + item.(map[string]any)["f12"].(string)
			index.Name = item.(map[string]any)["f14"].(string)
			list = append(list, index)
		}
		return len(arr)
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// IndexSZ 获取深交所指数列表
func (cli *Client) IndexSZ() ([]IndexItem, error) {
	args := make(url.Values)
	args.Set("pn", "1")
	args.Set("pz", "100")
	args.Set("po", "1")
	args.Set("np", "1")
	args.Set("fltt", "2")
	args.Set("invt", "2")
	args.Set("fid", "f3")
	args.Set("fs", "m:0+t:5")
	args.Set("fields", "f12,f14")
	var list []IndexItem
	err := cli.callPaged("https://88.push2.eastmoney.com/api/qt/clist/get", args, func(a any) int {
		arr := a.([]any)
		for _, item := range arr {
			var index IndexItem
			index.Code = "sz" + item.(map[string]any)["f12"].(string)
			index.Name = item.(map[string]any)["f14"].(string)
			list = append(list, index)
		}
		return len(arr)
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// IndexDaily 获取指数在指定范围内的数据
// sh开头表示上交所指数
// sz开头表示深交所指数
// csi开头表示中信指数
func (cli *Client) IndexDaily(code string, begin, end time.Time, right right) ([]Tick, error) {
	args := make(url.Values)
	// f1（代码）、f2（市场）、f3（名称）、f4（保留小数点位数）、f5（上市到如今交易日总数）、f6（K 线的前 1 日收盘价，即开始日期的前 1日收盘价）、f7（昨日收盘价，即结束日期的前 1 日收盘价）、f8（杂项，比如 4 是指 ETF 基金，7 是 A 股）
	args.Set("fields1", "f1,f2,f3,f4,f5,f6")
	// f51（日期时间）、f52（开盘价）、f53（收盘价）、f54（最高价）、f55（最低价）、f56（成交量）、f57（成交额）、f58（振幅%）、f59（涨跌幅%）、f60（涨跌额）、f61（换手率%）
	args.Set("fields2", "f51,f52,f53,f54,f55,f56,f57")
	// 日线
	args.Set("klt", "101")
	// 0（不复权）、1（前复权）、2（后复权）
	args.Set("fqt", right.arg())
	if strings.HasPrefix(code, "sz") {
		args.Set("secid", "0."+code[2:])
	} else if strings.HasPrefix(code, "sh") {
		args.Set("secid", "1."+code[2:])
	} else if strings.HasPrefix(code, "csi") {
		args.Set("secid", "2."+code[3:])
	}
	args.Set("beg", begin.Format("20060102"))
	args.Set("end", end.Format("20060102"))
	var data struct {
		Code    string   `json:"code"`     // 代码
		Market  int      `json:"market"`   // 市场
		Name    string   `json:"name"`     // 名称
		Decimal int      `json:"decimal"`  // 保留小数点位数
		Day1    float32  `json:"prePrice"` // 前一交易日收盘价
		KLines  []string `json:"klines"`
	}
	err := cli.call("/qt/stock/kline/get", args, &data)
	if err != nil {
		return nil, err
	}
	return parseKLineData(strings.Join(data.KLines, "\n"))
}
