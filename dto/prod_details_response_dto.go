package dto

type ViewList struct {
	ImageUrl string `json:"image_url"`
}

type BreadCrumb struct {
	Url string `json:"link"`
}

type AttributeList struct {
	Category    string `json:"category"`
	SizeChartId string `json:"size_chart_id"` // Needed for fetching size chart
}

type Metadata struct {
	Keywords []string `json:"keywords"`
}

type PriceInformation struct {
	CurrentPrice int `json:"currentPrice"`
}

type ProductDescription struct {
	Name                string   `json:"title"`
	TitleOfDesc         string   `json:"subtitle"`
	Description         string   `json:"text"`
	ItemizedDescription []string `json:"usps"`
}

type VariationList struct {
	Size string `json:"size"`
}

type ProductDetailsResponse struct {
	Id          string `json:"id"`
	ModelNumber string `json:"model_number"` // Needed for Reviews
	//MetaData *Metadata `json:"meta_data"` // Keywords.
	ViewList           []*ViewList         `json:"view_list"` // images_url list
	AttributeList      *AttributeList      `json:"attribute_list"`
	BreadCrumbs        []*BreadCrumb       `json:"breadcrumb_list"`
	PriceInformation   *PriceInformation   `json:"pricing_information"`
	ProductDescription *ProductDescription `json:"product_description"`
	VariationList      []*VariationList    `json:"variation_list"`
}
