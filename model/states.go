package model

type State struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	CountryID   int64   `json:"country_id"`
	CountryCode string  `json:"country_code"`
	FipsCode    string `json:"fips_code"`
	Iso2        string `json:"iso2"`
	Type        string `json:"type"`
	Level       int64  `json:"level"`
	ParentID    int64   `json:"parent_id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Flag        int16   `json:"flag"`
	WikiDataID *string `json:"wiki_data_id,omitempty"`
}
