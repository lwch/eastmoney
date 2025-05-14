package eastmoney

import "encoding/json"

type naFloat64 float64

func (n *naFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == `"-"` {
		*n = 0
		return nil
	}
	var v float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*n = naFloat64(v)
	return nil
}

func (n naFloat64) Value() float64 {
	return float64(n)
}
