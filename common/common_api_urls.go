package common

import "fmt"

// GetReviewAPIURL generates a URL to fetch the review data for a specific product model from the Adidas API.
// e.g. https://www.adidas.jp/api/product/JH7615/reviews?bazaarVoiceLocale=ja_JP&feature&includeLocales=ja%%2A&offset=0&sort=newest
func GetReviewAPIURL(modelNumber string, limit, offset int) string {
	return fmt.Sprintf("https://www.adidas.jp/api/models/%s/reviews?bazaarVoiceLocale=ja_JP&feature&includeLocales=ja%%2A&limit=%d&offset=%d&sort=newest", modelNumber, limit, offset)
}

// GetProductDetailsAPIURL generates a URL to fetch detailed product information for a specific product model from the Adidas API.
// e.g. https://www.adidas.jp/api/product/IQ1401
func GetProductDetailsAPIURL(productId string) string {
	return fmt.Sprintf("https://www.adidas.jp/api/products/%s", productId)
}

// GetSizeChartAPIURL generates a URL to fetch the size chart data for a specific product model from the Adidas API.
// e.g. // https://www.adidas.jp/size-chart/api/size-chart/size-chart-size-m_bottoms?locale=ja_JP&productId=JC6719
func GetSizeChartAPIURL(sizeChartId, productId string) string {
	return fmt.Sprintf("https://www.adidas.jp/size-chart/api/size-chart/%s/?locale=ja_JP&productId=%s", productId, sizeChartId)
}

// GetCoordinatesAPIURL generates a URL to fetch the coordinates for a specific product model from the Adidas API.
// e.g. https://www.adidas.jp/recs/api/styles-list?articleNumber=JC6719&articleModelNumber=KPV99&locale=ja_JP&pageSize=100&pageNumber=1
func GetCoordinatesAPIURL(productId, modelNumber string) string {
	return fmt.Sprintf("https://www.adidas.jp/recs/api/styles-list?articleNumber=%s&articleModelNumber=%s&locale=ja_JP&pageSize=100&pageNumber=1", productId, modelNumber)
}

// GetCoordinateLookBookWebURL generates a URL to fetch the coordinates for a specific product model from the Adidas website.
// e.g. https://www.adidas.jp//lookbook/25308915
func GetCoordinateLookBookWebURL(styleId string) string {
	return fmt.Sprintf("https://www.adidas.jp//lookbook/%s", styleId)
}

// GetProductDetailsWebURL generates a URL to fetch the web page of a specific product model from the Adidas website.
// e.g. https://shop.adidas.jp/products/IQ1401
func GetProductDetailsWebURL(productId string) string {
	return fmt.Sprintf("https://shop.adidas.jp/products/%s", productId)
}
