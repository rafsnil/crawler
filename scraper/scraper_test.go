package scraper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"simple-go-crawler/dto"
	"testing"
)

// Recursive function to extract all reviews
// Retrieved all 236 reviews for this model, kinda proud of myself T-T
func Test_extractAllReviews(t *testing.T) {
	model := "BSZ08"
	ans, err := getReviewDetails(model)
	if err != nil {
		t.Errorf("Error getting review details: %v", err)
	}
	t.Logf("Response body: %s", func() string { b, _ := json.Marshal(ans); return string(b) }())
}

func Test_extractSizeChart(t *testing.T) {
	productId := "JC6719"
	sizeChartId := "size-chart-size-m_bottoms"
	ans, err := getSizeChartDetails(productId, sizeChartId)
	if err != nil {
		t.Errorf("Error getting review details: %v", err)
	}
	t.Logf("Response body: %s", func() string { b, _ := json.Marshal(ans); return string(b) }())
}

func Test_extractCoordinates(t *testing.T) {
	productId := "B75806"
	modelNumber := "BSZ08"
	ans, err := getCoordinateDetails(productId, modelNumber)
	if err != nil {
		t.Errorf("Error getting review details: %v", err)
	}
	t.Logf("Response body: %s", func() string { b, _ := json.Marshal(ans); return string(b) }())
}

func Test_GetCoordinatesSet(t *testing.T) {
	item1 := &dto.CoordinateProductResp{
		Id: "25342141",
	}

	item2 := &dto.CoordinateProductResp{
		Id: "25308915",
	}

	item3 := &dto.CoordinateProductResp{
		Id: "25252042",
	}

	list := []*dto.CoordinateProductResp{item1, item2, item3}

	param := &dto.CoordinatesResponse{
		ProductList: list,
	}
	ans, err := extractLookBookData(param)
	if err != nil {
		t.Errorf("Error getting review details: %v", err)
	}

	t.Logf("Response body: %s", func() string { b, _ := json.Marshal(ans); return string(b) }())

}

func TestExtractSizeChartStrings(t *testing.T) {
	resp := dto.SizeChartDetailsResponse{
		Data: []*dto.Section{
			{
				Name: "garment-measurement",
				Components: &dto.Component{
					Table: [][]string{
						{"", "28", "30", "32"},
						{"ウエスト", "76cm", "78cm", "83cm"},
						{"股下", "34cm", "34cm", "34cm"},
						{"股上", "31cm", "31cm", "32cm"},
						{"ヒップ", "94cm", "97cm", "102cm"},
						{"裾回り", "63cm", "63cm", "65cm"},
					},
				},
			},
			// Other sections like size-guide or statement can be ignored
		},
	}

	//ans := extractSizeChartStrings(&resp)
	//for _, v := range ans {
	//	t.Logf("Size: %s", v)
	//}

	expected := []string{
		"28 (ウエスト=76cm, 股下=34cm, 股上=31cm, ヒップ=94cm, 裾回り=63cm)",
		"30 (ウエスト=78cm, 股下=34cm, 股上=31cm, ヒップ=97cm, 裾回り=63cm)",
		"32 (ウエスト=83cm, 股下=34cm, 股上=32cm, ヒップ=102cm, 裾回り=65cm)",
	}

	result := extractSizeChartStrings(&resp)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected:\n%v\n\nGot:\n%v", expected, result)
	}
}

func TestMakeReviewDataModel(t *testing.T) {
	details := &dto.ReviewDetailsResponse{
		TotalResults: 2,
		ReviewDetailsList: []*dto.ReviewDetails{
			{
				Id:       "rev1",
				Title:    "Awesome shoes",
				Username: "John",
				Rating:   4.5,
				Desc:     "Very comfortable and stylish.",
				Date:     "2024-01-01",
			},
			{
				Id:       "rev2",
				Title:    "Good fit",
				Username: "Alice",
				Rating:   3.5,
				Desc:     "Fits well but a bit pricey.",
				Date:     "2024-02-01",
			},
		},
	}

	result := makeReviewDataModel(details)

	// print the whoel result
	fmt.Printf("%+v\n", result)

	for _, v := range result.GeneralReviews {
		fmt.Printf("%+v\n", v)
	}
	//if !reflect.DeepEqual(result, expected) {
	//	t.Errorf("Expected:\n%+v\n\nGot:\n%+v", expected, result)
	//}
}
