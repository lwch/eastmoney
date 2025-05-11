package eastmoney

import (
	"net/url"
)

type FundBasicItem struct {
	Code string `json:"f12"` // 代码
	Name string `json:"f14"` // 名称
}

// FundBasic 获取基金基本信息
func (cli *Client) FundBasic() ([]FundBasicItem, error) {
	args := make(url.Values)
	args.Set("pn", "1")
	args.Set("pz", "100")
	args.Set("po", "1")
	args.Set("np", "1")
	args.Set("fltt", "2")
	args.Set("invt", "2")
	args.Set("fid", "f3")
	args.Set("fs", "b:MK0021,b:MK0022,b:MK0023,b:MK0024")
	args.Set("fields", "f12,f14")
	var list []FundBasicItem
	err := cli.callPaged("https://88.push2.eastmoney.com/api/qt/clist/get", args, func(a any) int {
		arr := a.([]any)
		for _, item := range arr {
			var fund FundBasicItem
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
