package model

type Country struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Iso3           string  `json:"iso3"`
	NumericCode    string  `json:"numeric_code"`
	Iso2           string  `json:"iso2"`
	PhoneCode      string  `json:"phonecode"`
	Capital        string  `json:"capital"`
	Currency       string  `json:"currency"`
	CurrencyName   string  `json:"currency_name"`
	CurrencySymbol string  `json:"currency_symbol"`
	Tld            string  `json:"tld"`
	Native         *string `json:"native,omitempty"`
	Region         string  `json:"region"`
	RegionID       *int64  `json:"region_id,omitempty"`
	Subregion      string  `json:"subregion"`
	SubregionID    *int64  `json:"subregion_id,omitempty"`
	Nationality    string  `json:"nationality"`
	Timezones      string  `json:"timezones"`
	Translations   string  `json:"translations"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Emoji          string  `json:"emoji"`
	EmojiU         string  `json:"emoji_u"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	Flag           int16   `json:"flag"`
	WikiDataID     *string `json:"wiki_data_id,omitempty"`
}
