package model

type DateType int64

const (
	Year DateType = iota
	Month
	Day
)

func (s DateType) String() string {
	switch s {
	case Year:
		return "year"
	case Month:
		return "month"
	case Day:
		return "day"
	}
	return ""
}

func NewDateType(s string) DateType {
	switch s {
	case "year":
		return Year
	case "month":
		return Month
	case "day":
		return Day
	}
	return Month
}
