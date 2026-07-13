package domain

type User struct {
	Id        int    `json:"id"`
	AuthId    string `json:"auth_id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	AppRole   string `json:"app_role"`
	Enabled   bool   `json:"enabled"`
	AvatarUrl string `json:"avatar_url"`
}
