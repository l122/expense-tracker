package domain

type Category struct {
	Id        int    `json:"id"`
	ProfileId int    `json:"profile_id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Icon      string `json:"icon"`
}
