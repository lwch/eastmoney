package eastmoney

type right byte // 复权

const (
	PreRight  right = iota // 前复权
	PostRight              // 后复权
	NoRight                // 不复权
)

func (r right) arg() string {
	switch r {
	case PreRight:
		return "1"
	case PostRight:
		return "2"
	default:
		return "0"
	}
}

func (r right) String() string {
	switch r {
	case PreRight:
		return "前复权"
	case PostRight:
		return "后复权"
	case NoRight:
		return "不复权"
	default:
		return "未知复权类型"
	}
}
