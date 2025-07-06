package common

import (
	"encoding/json"
	"fmt"
	"simple-go-crawler/dto"
	"testing"
)

func TestGetWebPage(t *testing.T) {
	url := "https://shop.adidas.jp/men"
	_, _, err := GetWebPage(url)
	if err != nil {
		t.Errorf("Error getting web page: %v", err)
	}
	//t.Logf("Response body: %s", node)
}

func TestMakeAPICall(t *testing.T) {
	url := "https://www.adidas.jp/api/products/B75806"
	body, err := MakeAPICall(url)
	if err != nil {
		t.Errorf("Error making API call: %v", err)
	}
	t.Logf("Response body: %s", body)
}

func TestProductDetailsResponseParsing(t *testing.T) {
	url := "https://www.adidas.jp/api/products/B75806"
	body, err := MakeAPICall(url)
	if err != nil {
		t.Errorf("Error making API call: %v", err)
	}
	var resp dto.ProductDetailsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}

	t.Logf("Response body: %s", func() string { b, _ := json.Marshal(resp); return string(b) }())
}

func TestReviewDetailsResponse(t *testing.T) {
	model := "BSZ08"
	url := GetReviewAPIURL(model, 10, 20)
	fmt.Println(url)
	body, err := MakeAPICall(url)
	if err != nil {
		t.Errorf("Error making API call: %v", err)
	}
	var resp dto.ReviewDetailsResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Errorf("Error unmarshalling response: %v", err)
	}

	t.Logf("Response body: %s", func() string { b, _ := json.Marshal(resp); return string(b) }())
}

// https://www.adidas.jp/api/models/BSZ08/reviews?bazaarVoiceLocale=ja_JP&feature&includeLocales=ja%2A&limit=10offset=0&sort=newest
// https://www.adidas.jp/api/models/BSZ08/reviews?bazaarVoiceLocale=ja_JP&feature&includeLocales=ja%2A&limit=10&sort=newest&offset=0
