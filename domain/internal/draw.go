package internal

import "time"

type Word struct {
	Id          int
	Word        string
	Description string
	DrawNo      int
	DrawSuccess int
	DrawFail    int
	UpdatedAt   time.Time
}
