package model

type City struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	StateID    int64   `json:"state_id"`
	StateCode  string  `json:"state_code"`
	CountryID  int64   `json:"country_id"`
	CountryCode string  `json:"country_code"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
	Flag       int16   `json:"flag"`
	WikiDataID *string `json:"wiki_data_id,omitempty"`
}
