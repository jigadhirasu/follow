package types

type Time int64

func (ts Time) Int64() int64 {
	return int64(ts)
}
