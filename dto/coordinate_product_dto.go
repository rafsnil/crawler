package dto

type CoordinateProductsSetDto struct {
	Products []*Products `json:"products"`
}

type Products struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"imgNormal"`
	Price    string `json:"currentPrice"`
}
