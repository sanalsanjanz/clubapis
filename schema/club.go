package schemas

type Club struct {
	ClubName   string  `json:"club_name"`
	Contact    string  `json:"contact"`
	Location   string  `json:"location"`
	MonthlyFee float64 `json:"monthly_fee"`
	RegNo      string  `json:"reg_no"`
	Password   string  `json:"password"`
}
