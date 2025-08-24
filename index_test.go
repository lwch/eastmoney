package eastmoney

import (
	"fmt"
	"testing"
	"time"
)

func TestIndexList(t *testing.T) {
	cli := New()
	list, err := cli.IndexSH()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range list {
		fmt.Println(item)
	}
	list, err = cli.IndexSZ()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range list {
		fmt.Println(item)
	}
}

func TestIndex(t *testing.T) {
	cli := New()
	ticks, err := cli.IndexDaily("sh000001", time.Now().Add(-time.Hour*24*30), time.Now(), PreRight)
	if err != nil {
		t.Fatal(err)
	}
	for _, tick := range ticks {
		fmt.Println(tick)
	}
}
