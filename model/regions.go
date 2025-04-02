package model

type Region struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Translations string  `json:"translations"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Flag        int16   `json:"flag"`
	WikiDataID  *string `json:"wiki_data_id,omitempty"`
}