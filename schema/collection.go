package schemas

import "time"

type Collection struct {
	ID        string    `json:"id,omitempty"`
	AdminID   string    `json:"admin_id"`
	MemberID  string    `json:"member_id"`
	Amount    float64   `json:"amount"`
	Reason    string    `json:"reason,omitempty"`
	ForMonth  string    `json:"for_month,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Notes     string    `json:"notes,omitempty"`
}
