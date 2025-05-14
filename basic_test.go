package eastmoney

import "testing"

func TestBasicInfo(t *testing.T) {
	cli := New()
	info, err := cli.Info("000976")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", info)
}
