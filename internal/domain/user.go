package domain

type User struct {
	Id       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	AppRole  string `json:"app_role"`
	Enabled  bool   `json:"enabled"`
}
