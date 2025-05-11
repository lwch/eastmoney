package eastmoney

import (
	"net/url"
	"strings"
)

type BasicItem struct {
	Code string
	Name string
}

// Basic 获取个股基本信息
func (cli *Client) Basic() ([]BasicItem, error) {
	args := make(url.Values)
	args.Set("pn", "1")
	args.Set("pz", "100")
	args.Set("po", "1")
	args.Set("np", "1")
	args.Set("fltt", "2")
	args.Set("invt", "2")
	args.Set("fid", "f12")
	args.Set("fs", "m:0 t:6,m:0 t:80,m:1 t:2,m:1 t:23,m:0 t:81 s:2048")
	args.Set("fields", "f12,f14")
	var list []BasicItem
	err := cli.callPaged("https://88.push2.eastmoney.com/api/qt/clist/get", args, func(a any) int {
		arr := a.([]any)
		for _, item := range arr {
			var fund BasicItem
			fund.Code = item.(map[string]any)["f12"].(string)
			fund.Name = item.(map[string]any)["f14"].(string)
			list = append(list, fund)
		}
		return len(arr)
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

type Info struct {
	Code     string
	Name     string
	Industry string
	Area     string
	Sectors  []string
}

// Info 获取某只个股的基本信息
func (cli *Client) Info(code string) (Info, error) {
	args := make(url.Values)
	// f57(代码)、f58(名称)、f127(行业)、f128(地区)、f129(概念板块)
	args.Set("fields", "f57,f58,f127,f128,f129")
	args.Set("fltt", "2")
	args.Set("invt", "2")
	if code[0] == '6' {
		args.Set("secid", "1."+code)
	} else {
		args.Set("secid", "0."+code)
	}
	var data struct {
		Code     string `json:"f57"`  // 代码
		Name     string `json:"f58"`  // 名称
		Industry string `json:"f127"` // 行业
		Area     string `json:"f128"` // 地区
		Sectors  string `json:"f129"` // 概念板块
	}
	var info Info
	err := cli.call("/qt/stock/get", args, &data)
	if err != nil {
		return info, err
	}
	data.Area = strings.TrimSuffix(data.Area, "板块")
	info.Code = data.Code
	info.Name = data.Name
	info.Industry = data.Industry
	info.Area = data.Area
	if len(data.Sectors) > 0 {
		info.Sectors = strings.Split(data.Sectors, ",")
	}
	return info, nil
}
