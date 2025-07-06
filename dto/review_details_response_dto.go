package dto

type ReviewDetails struct {
	Id       string  `json:"id"`
	Title    string  `json:"title"`
	Username string  `json:"userNickname"`
	Rating   float32 `json:"rating"`
	Desc     string  `json:"text"`
	Date     string  `json:"submissionTime"`
}

type ReviewDetailsResponse struct {
	TotalResults      int              `json:"totalResults"`
	ReviewDetailsList []*ReviewDetails `json:"reviews"`
}
