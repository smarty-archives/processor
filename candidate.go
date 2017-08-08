package processor

type Candidate struct {
	DeliveryLine1 string `json:"delivery_line_1"`
	LastLine      string `json:"last_line"`
	Components    struct {
		City    string `json:"city_name"`
		State   string `json:"state_abbreviation"`
		ZIPCode string `json:"zipcode"`
	} `json:"components"`
	Analysis struct {
		Match  string `json:"dpv_match_code"`
		Vacant string `json:"dpv_vacant"`
		Active string `json:"active"`
	} `json:"analysis"`
}
