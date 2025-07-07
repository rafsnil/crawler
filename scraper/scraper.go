package scraper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"math/rand"
	"simple-go-crawler/common"
	"simple-go-crawler/constant"
	"simple-go-crawler/dto"
	"simple-go-crawler/model"
	"strings"
	"time"
)

func ExtractProductData(prodId string, counter *ProductCounter) {
	// Get Product Details
	prodDetails, err := getProductDetails(prodId)
	if err != nil {
		log.Printf("Error getting product details: %v", err)
		return
	}

	// Get Review
	reviewDetails, err := getReviewDetails(prodDetails.ModelNumber)
	if err != nil {
		log.Printf("Error getting review details: %v", err)
	}

	// Get Size Chart
	var sizeChartDetails *dto.SizeChartDetailsResponse
	if prodDetails.AttributeList.SizeChartId != constant.NO_SIZE_CHART {
		sizeChartDetails, err = getSizeChartDetails(prodDetails.Id, prodDetails.AttributeList.SizeChartId)
		if err != nil {
			log.Printf("Error getting size chart details: %v", err)
		}
	}

	// Get All Coordinates
	coordinateDetails, err := getCoordinateDetails(prodDetails.Id, prodDetails.ModelNumber)
	if err != nil {
		log.Printf("Error getting coordinate details: %v", err)
	}
	// TODO:
	// For each StyleId:
	// 1. Hit the LookBack web page (as this data is server-side rendered)
	// 2. Extract the data from the script in HTML and convert it to JSON
	// 3. Extract the data from the JSON and convert it the required output
	var lookBookResponse []*dto.CoordinateProductsSetDto
	if coordinateDetails.Total > 0 {
		// gets all sets of coordinates
		lookBookResponse, err = extractLookBookData(&coordinateDetails)
		if err != nil {
			log.Printf("Error getting look book data: %v", err)
		}
	}

	productData := aggregateAllData(prodDetails, reviewDetails, sizeChartDetails, coordinateDetails, lookBookResponse)
	GlobalOutput.Products = append(GlobalOutput.Products, productData)
	counter.Increment()
}

func aggregateAllData(prodDetails *dto.ProductDetailsResponse, reviewDetails *dto.ReviewDetailsResponse, sizeChartDetails *dto.SizeChartDetailsResponse, coordinateDetails dto.CoordinatesResponse, lookBookResponse []*dto.CoordinateProductsSetDto) *model.Product {

	detailsUrl := common.GetProductDetailsWebURL(prodDetails.Id)

	// Available Sizes
	availableSizes := make([]string, 0)
	for _, v := range prodDetails.VariationList {
		availableSizes = append(availableSizes, v.Size)
	}

	// Image Urls
	imageUrls := make([]string, 0)
	for _, vl := range prodDetails.ViewList {
		imageUrls = append(imageUrls, vl.ImageUrl)
	}

	// Size Chart
	var sizeChartInfo []string
	if sizeChartDetails != nil {
		sizeChartInfo = extractSizeChartStrings(sizeChartDetails)
	}

	// Reviews
	var reviews *model.Review
	if reviewDetails != nil {
		reviews = makeReviewDataModel(reviewDetails)
	}

	// Coordinates
	var coordinates *model.AllCoordinates
	if coordinateDetails.Total > 0 {
		coordinates = makeCoordinatesDataModel(coordinateDetails, lookBookResponse)
	}

	productData := model.Product{
		Name:                prodDetails.ProductDescription.Name,
		TitleOfDesc:         prodDetails.ProductDescription.TitleOfDesc,
		Description:         prodDetails.ProductDescription.Description,
		ItemizedDescription: prodDetails.ProductDescription.ItemizedDescription,
		SizeChart:           sizeChartInfo,
		Category:            prodDetails.AttributeList.Category,
		Price:               fmt.Sprintf("%d", prodDetails.PriceInformation.CurrentPrice),
		AvailableSizes:      availableSizes,
		ImagesURLs:          imageUrls,
		DetailsURL:          detailsUrl,
		AllReviews:          reviews,
		Coordinates:         coordinates,
	}
	return &productData
}

func makeCoordinatesDataModel(coordinateDetails dto.CoordinatesResponse, lookBookResponse []*dto.CoordinateProductsSetDto) *model.AllCoordinates {
	allCoords := &model.AllCoordinates{
		TotalCoordinates: coordinateDetails.Total,
		Coordinates:      make([]*model.CoordinateSet, 0, len(lookBookResponse)),
	}
	// Build a map from product ID to detailed info from lookBookResponse for quick lookup
	productDetailsMap := make(map[string]*dto.Products)
	for _, coordSet := range lookBookResponse {
		for _, prod := range coordSet.Products {
			productDetailsMap[prod.Id] = prod
		}
	}

	for _, coordSetDto := range lookBookResponse {
		coordSet := &model.CoordinateSet{
			CoordinateProducts: make([]*model.CoordinateProduct, 0, len(coordSetDto.Products)),
		}

		for _, prodDto := range coordSetDto.Products {
			coordProd := &model.CoordinateProduct{
				ProductNumber: prodDto.Id,
				Name:          prodDto.Name,
				Price:         prodDto.Price,
				ImageURL:      prodDto.ImageUrl,
				DetailsURL:    common.GetProductDetailsWebURL(prodDto.Id),
			}
			coordSet.CoordinateProducts = append(coordSet.CoordinateProducts, coordProd)
		}

		allCoords.Coordinates = append(allCoords.Coordinates, coordSet)
	}

	return allCoords
}

func makeReviewDataModel(details *dto.ReviewDetailsResponse) *model.Review {
	if details == nil {
		return nil
	}

	review := &model.Review{
		NumberOfReviews: fmt.Sprintf("%d", details.TotalResults),
		GeneralReviews:  make([]*model.GeneralReview, 0, len(details.ReviewDetailsList)),
	}

	var ratingSum float32

	for _, r := range details.ReviewDetailsList {
		ratingSum += r.Rating

		review.GeneralReviews = append(review.GeneralReviews, &model.GeneralReview{
			ReviewId:    r.Id,
			Username:    r.Username,
			Rating:      fmt.Sprintf("%.1f", r.Rating),
			ReviewTitle: r.Title,
			ReviewDesc:  r.Desc,
			ReviewDate:  r.Date,
		})
	}

	// Compute average rating if there are reviews
	if len(details.ReviewDetailsList) > 0 {
		avg := ratingSum / float32(len(details.ReviewDetailsList))
		review.AverageRating = fmt.Sprintf("%.2f", avg)
	} else {
		review.AverageRating = "N/A"
	}

	review.RecommendRating = "N/A"

	return review
}

func extractSizeChartStrings(resp *dto.SizeChartDetailsResponse) []string {
	for _, section := range resp.Data {
		if section.Name != "garment-measurement" || section.Components == nil || section.Components.Table == nil {
			continue
		}

		table := section.Components.Table
		if len(table) < 2 || len(table[0]) < 2 {
			continue
		}

		// table[0] = ["", "28", "30", ...]
		sizes := table[0][1:] // skip first empty cell

		sizeData := make([][]string, len(sizes)) // one entry per size

		// Loop through each property row
		for _, row := range table[1:] {
			if len(row) < 2 {
				continue
			}
			property := row[0] // e.g., "ウエスト"
			for i, value := range row[1:] {
				sizeData[i] = append(sizeData[i], fmt.Sprintf("%s=%s", property, value))
			}
		}

		// Build final strings
		var result []string
		for i, size := range sizes {
			entry := fmt.Sprintf("%s (%s)", size, strings.Join(sizeData[i], ", "))
			result = append(result, entry)
		}

		return result
	}
	return nil
}

// Extracts the LookBook data from the HTML page (The Actual Required Coordinates Data)
func extractLookBookData(details *dto.CoordinatesResponse) ([]*dto.CoordinateProductsSetDto, error) {
	list := make([]*dto.CoordinateProductsSetDto, 0)
	for _, coordinate := range details.ProductList {
		url := common.GetCoordinateLookBookWebURL(coordinate.Id)
		// TODO: change
		_, htmlBody, err := common.GetWebPage(url)
		if err != nil {
			log.Printf("Error getting web page: %v", err)
			return nil, err
		}

		coordinateData, err := getLookBookScriptData(htmlBody)
		if err != nil {
			log.Printf("Error getting look book data: %v", err)
			return nil, err
		}
		list = append(list, coordinateData)
	}
	return list, nil
}

// getLookBookScriptData parses the HTML node and extracts the JSON inside the <script> tag with data-mf-id="lookbook-microfrontend"
// returns 1 set of coordinate data for 1 product
func getLookBookScriptData(body *html.Node) (*dto.CoordinateProductsSetDto, error) {
	var result dto.CoordinateProductsSetDto

	var findScript func(*html.Node) *html.Node
	findScript = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, attr := range n.Attr {
				if attr.Key == "data-mf-id" && attr.Val == constant.LOOKBOOK_SSR_DATA_IDENTIFIER {
					return n
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if res := findScript(c); res != nil {
				return res
			}
		}
		return nil
	}

	scriptNode := findScript(body)
	if scriptNode == nil {
		return nil, errors.New("script tag with data-mf-id not found")
	}

	if scriptNode.FirstChild == nil {
		return nil, errors.New("script tag found but no content")
	}

	jsonData := scriptNode.FirstChild.Data

	decoder := json.NewDecoder(bytes.NewBufferString(jsonData))
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
func getCoordinateDetails(prodId string, modelNumber string) (dto.CoordinatesResponse, error) {
	url := common.GetCoordinatesAPIURL(prodId, modelNumber)
	d := time.Duration(rand.Intn(3)+3) * time.Second
	fmt.Println("sleeping for", d, "before next API call For Coordinates")
	respJsonBytes, err := common.MakeAPICall(url)
	if err != nil {
		log.Printf("Error making API call for CoordinatesUrl %s", url)
		return dto.CoordinatesResponse{}, err
	}
	var resp dto.CoordinatesResponse
	if err := json.Unmarshal(respJsonBytes, &resp); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return dto.CoordinatesResponse{}, err
	}
	return resp, nil
}

func getSizeChartDetails(productId, sizeChartId string) (*dto.SizeChartDetailsResponse, error) {
	sizeChartUrl := common.GetSizeChartAPIURL(productId, sizeChartId)
	d := time.Duration(rand.Intn(3)+3) * time.Second
	fmt.Println("sleeping for", d, "before next API call For Size Charts")
	sizeChartRespJsonBytes, err := common.MakeAPICall(sizeChartUrl)
	if err != nil {
		log.Printf("Error making API call for SizeChartUrl %s", sizeChartUrl)
		return nil, err
	}
	fmt.Println(string(sizeChartRespJsonBytes))
	var sizeChartResp dto.SizeChartDetailsResponse
	if err := json.Unmarshal(sizeChartRespJsonBytes, &sizeChartResp); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return nil, err
	}
	return &sizeChartResp, nil
}

func getProductDetails(prodId string) (*dto.ProductDetailsResponse, error) {
	//productId := getProductIdFromURL(url)
	productDetailsURL := common.GetProductDetailsAPIURL(prodId)
	d := time.Duration(rand.Intn(3)+3) * time.Second
	fmt.Println("sleeping for", d, "before next API call For Product Details")
	prodDetailsRespJsonBytes, err := common.MakeAPICall(productDetailsURL)
	if err != nil {
		log.Printf("Error making API call for ProductUrl %s", productDetailsURL)
		return nil, err
	}
	var prodDetailsResp dto.ProductDetailsResponse
	if err := json.Unmarshal(prodDetailsRespJsonBytes, &prodDetailsResp); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return nil, err
	}

	return &prodDetailsResp, nil
}

func getReviewDetails(modelNumber string) (*dto.ReviewDetailsResponse, error) {
	reviewDetailsRespList, err := extractAllReviews(modelNumber, 10, 0)
	if err != nil {
		return nil, err
	}
	finalReviewData := aggregateAllReviews(reviewDetailsRespList)
	return finalReviewData, nil
}

func aggregateAllReviews(list []*dto.ReviewDetailsResponse) *dto.ReviewDetailsResponse {
	count := 0
	var output dto.ReviewDetailsResponse
	output.TotalResults = 0
	for _, reviewData := range list {
		output.TotalResults = reviewData.TotalResults
		count = count + len(reviewData.ReviewDetailsList)
		output.ReviewDetailsList = append(output.ReviewDetailsList, reviewData.ReviewDetailsList...)
	}
	fmt.Println("Total reviews: ", count)
	return &output
}

// This API has a cap of sending only 10 data at once.
// Hence, the recursive calls
func extractAllReviews(modelNumber string, limit, offset int) ([]*dto.ReviewDetailsResponse, error) {
	reviewUrl := common.GetReviewAPIURL(modelNumber, limit, offset)
	d := time.Duration(rand.Intn(3)+3) * time.Second
	fmt.Println("sleeping for", d, "before next API call For All Reviews")
	reviewDetailsRespBytes, err := common.MakeAPICall(reviewUrl)
	if err != nil {
		log.Printf("Error making API call for ReviewUrl %s: %v", reviewUrl, err)
		return nil, fmt.Errorf("API call failed for offset %d: %w", offset, err)
	}

	var reviewDetailsResp dto.ReviewDetailsResponse
	if err := json.Unmarshal(reviewDetailsRespBytes, &reviewDetailsResp); err != nil {
		log.Printf("Error unmarshalling response: %v", err)
		return nil, fmt.Errorf("unmarshalling failed for offset %d: %w", offset, err)
	}

	if len(reviewDetailsResp.ReviewDetailsList) == 0 {
		return []*dto.ReviewDetailsResponse{}, nil
	}

	time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)

	remainingResponses, err := extractAllReviews(modelNumber, limit, offset+limit)
	if err != nil {
		return nil, err
	}

	allCollectedResponses := append([]*dto.ReviewDetailsResponse{&reviewDetailsResp}, remainingResponses...)

	return allCollectedResponses, nil
}

func getProductIdFromURL(url string) string {
	splits := strings.Split(url, "/")
	ansSlice := strings.Split(splits[len(splits)-1], ".")
	return ansSlice[0]
}
