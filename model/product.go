package model

type Product struct {
	Name                string          // F
	TitleOfDesc         string          // F
	Description         string          // F
	ItemizedDescription []string        // F
	SizeChart           []string        // F // "28(Waist=76cm, inseam=34cm, Rise=31cm, Hip=94cm, Hem Circumference=63cm)" Need to construct
	Category            string          // F
	Price               string          // F
	AvailableSizes      []string        // F
	ImagesURLs          []string        // F
	DetailsURL          string          // F // will construct through func and productId
	AllReviews          *Review         // F // Need to construct
	Coordinates         *AllCoordinates // F // Need to construct
	// Keywords []string
	// SenseOfSize
	// SpecialFunctionalities
}

// 1 product has multiple sets of coordinates, where each coordinates has multiple products.
type AllCoordinates struct {
	TotalCoordinates int
	Coordinates      []*CoordinateSet //  as multiple sets of coordinates for each product exists
}

type CoordinateSet struct {
	CoordinateProducts []*CoordinateProduct
}
type CoordinateProduct struct {
	Name          string
	Price         string
	ProductNumber string
	ImageURL      string
	DetailsURL    string
}

type Review struct {
	AverageRating   string
	NumberOfReviews string
	RecommendRating string
	GeneralReviews  []*GeneralReview
	// SenseOfOverallProperties
}

type GeneralReview struct {
	ReviewId    string
	Username    string
	Rating      string
	ReviewTitle string
	ReviewDesc  string
	ReviewDate  string
}

/*
SKIPPED DATA DUE TO TIME CONSTRAINT:
- Review Rating of each Item
- Special Function
- Sense Of Overall Properties
- Keywords
*/
