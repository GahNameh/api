package entity

type Event struct {
	Id          int
	Year        int
	Month       int
	Day         int
	Type        int
	IsHoliday   bool
	Description string
}
