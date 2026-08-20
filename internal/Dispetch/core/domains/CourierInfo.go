package dispetch_domains

type CourierInfo struct {
	ID      int
	Version int
	City    string
}

func NewCourierInfo(
	id int,
	version int,
	city string,
) CourierInfo {
	return CourierInfo{
		ID:      id,
		Version: version,
		City:    city,
	}
}
