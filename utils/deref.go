package goutils

func DerefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func DerefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}
