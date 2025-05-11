package eastmoney

import "net/url"

type RealTimeItem struct {
	Code   string  `json:"f12"` // 代码
	Name   string  `json:"f14"` // 名称
	Value  float64 `json:"f2"`  // 最新价
	Delta  float64 `json:"f3"`  // 涨跌幅
	Diff   float64 `json:"f4"`  // 涨跌额
	Amount float64 `json:"f5"`  // 成交量
	Trun   float64 `json:"f6"`  // 成交额
	Swing  float64 `json:"f7"`  // 振幅
	High   float64 `json:"f15"` // 最高价
	Low    float64 `json:"f16"` // 最低价
	Open   float64 `json:"f17"` // 开盘价
	Close  float64 `json:"f18"` // 昨收盘
}

// KC 获取科创板实时数据
func (cli *Client) KC() ([]RealTimeItem, error) {
	args := make(url.Values)
	args.Set("pn", "1")
	args.Set("pz", "100")
	args.Set("po", "1")
	args.Set("np", "1")
	args.Set("fltt", "2")
	args.Set("invt", "2")
	args.Set("fid", "f12")
	args.Set("fs", "m:0 t:80")
	args.Set("fields", "f2,f3,f4,f5,f6,f7,f12,f14,f15,f16,f17,f18")
	getFloat := func(a any) float64 {
		if a == nil {
			return 0
		}
		switch v := a.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		default:
			return 0
		}
	}
	var list []RealTimeItem
	err := cli.callPaged("https://88.push2.eastmoney.com/api/qt/clist/get", args, func(a any) int {
		arr := a.([]any)
		for _, item := range arr {
			if _, ok := item.(map[string]any)["f2"].(string); ok {
				continue
			}
			var fund RealTimeItem
			fund.Code = item.(map[string]any)["f12"].(string)
			fund.Name = item.(map[string]any)["f14"].(string)
			fund.Value = item.(map[string]any)["f2"].(float64)
			fund.Delta = getFloat(item.(map[string]any)["f3"])
			fund.Diff = getFloat(item.(map[string]any)["f4"])
			fund.Amount = getFloat(item.(map[string]any)["f5"])
			fund.Trun = getFloat(item.(map[string]any)["f6"])
			fund.Swing = getFloat(item.(map[string]any)["f7"])
			fund.High = getFloat(item.(map[string]any)["f15"])
			fund.Low = getFloat(item.(map[string]any)["f16"])
			fund.Open = getFloat(item.(map[string]any)["f17"])
			fund.Close = getFloat(item.(map[string]any)["f18"])
			list = append(list, fund)
		}
		return len(arr)
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
