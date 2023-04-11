package transport

type Word struct {
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Status        bool    `json:"status"`
	Id            int     `json:"id"`
	Speech        string  `json:"speech"`
	TodayActivity float64 `json:"today_activity"`
	Coverage      float64 `json:"coverage"`
}
type Feedback struct {
	Fail    bool `json:"fail"`
	Success bool `json:"success"`
	Status  bool `json:"status"`
	Id      int  `json:"id"`
}
