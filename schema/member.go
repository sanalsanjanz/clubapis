package schemas

type Member struct {
	ID        string `json:"id,omitempty"`
	AdminID   string `json:"admin_id"` // club reg_no
	Name      string `json:"name"`
	Mobile    string `json:"mobile"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	// ClubName string `json:"updated_at,omitempty"`
	// Contact string  `json:"mobile"`
	// Location string `json:"updated_at,omitempty"`
}
