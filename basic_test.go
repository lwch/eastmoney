package eastmoney

import "testing"

func TestBasicInfo(t *testing.T) {
	cli := New()
	info, err := cli.Info("600519")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", info)
}
