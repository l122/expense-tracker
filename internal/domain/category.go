package domain

type Category struct {
	Id        int    `json:"id,omitempty"`
	ProfileId int    `json:"profile_id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Icon      string `json:"icon"`
}
