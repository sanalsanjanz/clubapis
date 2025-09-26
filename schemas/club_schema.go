package schemas

type Club struct {
	ClubName   string  `json:"club_name"`
	Contact    string  `json:"contact"`
	Location   string  `json:"location"`
	MonthlyFee float64 `json:"monthly_fee"`
	RegNo      string  `json:"reg_no"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
