package app

type Order struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	JoinedAt string `json:"joined_at,omitempty"`
	OrderID  string `json:"order_id,omitempty"`
	Amount   string `json:"amount,omitempty"`
	OrderAt  string `json:"ordered_at,omitempty"`
}
