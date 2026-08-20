package dispetch_domains

type FreeCourierInfo struct {
	ID      int
	Version int
}

func NewFreeCourierInfo(
	id int,
	version int,
) FreeCourierInfo {
	return FreeCourierInfo{
		ID:      id,
		Version: version,
	}
}
