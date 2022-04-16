package model

type Day struct {
	Title     string
	Value     string
	Weekday   string
	IsToday   bool
	IsHoliday bool
	Events    []Event
	IsWeekend bool
	Row       int
	Column    int
}
