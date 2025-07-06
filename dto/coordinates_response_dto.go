package dto

type CoordinateProductResp struct {
	Id           string `json:"id"`
	ProductTotal int    `json:"productTotal"`
}

type CoordinatesResponse struct {
	Total       int                      `json:"total"`
	ProductList []*CoordinateProductResp `json:"styles"`
}
