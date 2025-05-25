package eastmoney

import (
	"testing"
)

func TestDisableKeepAlive(t *testing.T) {
	cli := New(DisableKeepAlive())
	cli.call("/api/test", nil, nil)
	cli.call("/api/test", nil, nil)
}
