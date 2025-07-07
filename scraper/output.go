package scraper

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"simple-go-crawler/model"
	"sync"
)

type Output struct {
	Products   []*model.Product
	mutex      sync.Mutex
	ProductMap map[string]struct{} // Optional: if you want to access products by ID
}

var GlobalOutput *Output

func init() {
	GlobalOutput = &Output{
		Products:   make([]*model.Product, 0),
		ProductMap: make(map[string]struct{}),
	}
}

func (o *Output) AddProduct(product *model.Product) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.Products = append(o.Products, product)
	o.ProductMap[product.Id] = struct{}{}
}

func (o *Output) GetProducts() []*model.Product {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	products := make([]*model.Product, len(o.Products))
	copy(products, o.Products)
	return products
}

// check if a product with the given ID already exists
func (o *Output) ProductExists(productId string) bool {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	_, exists := o.ProductMap[productId]
	return exists
}

func (o *Output) PrintToExcelSheet() {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	f := excelize.NewFile()
	sheet := "Products"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)

	// Headers
	headers := []string{
		"Name", "Title", "Description", "Itemized Description",
		"Size Chart", "Category", "Price", "Available Sizes",
		"Image URLs", "Details URL", "Avg Rating", "Total Reviews", "Recommend %",
		"Coordinate Products",
	}

	// Add review headers for 5 reviews
	for i := 1; i <= 5; i++ {
		headers = append(headers,
			fmt.Sprintf("Review%d_Username", i),
			fmt.Sprintf("Review%d_Rating", i),
			fmt.Sprintf("Review%d_Title", i),
			fmt.Sprintf("Review%d_Desc", i),
			fmt.Sprintf("Review%d_Date", i),
		)
	}

	// Write headers to sheet
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	// Write product data
	for rowIdx, p := range o.Products {
		row := rowIdx + 2 // Start from row 2

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), p.Name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), p.TitleOfDesc)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), p.Description)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), join(p.ItemizedDescription))
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), join(p.SizeChart))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), p.Category)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), p.Price)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), join(p.AvailableSizes))
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), join(p.ImagesURLs))
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), p.DetailsURL)

		if p.AllReviews != nil {
			f.SetCellValue(sheet, fmt.Sprintf("K%d", row), p.AllReviews.AverageRating)
			f.SetCellValue(sheet, fmt.Sprintf("L%d", row), p.AllReviews.NumberOfReviews)
			f.SetCellValue(sheet, fmt.Sprintf("M%d", row), p.AllReviews.RecommendRating)
		}

		if p.Coordinates != nil {
			var names []string
			for _, set := range p.Coordinates.Coordinates {
				for _, cp := range set.CoordinateProducts {
					names = append(names, cp.Name+"("+cp.Price+")")
				}
			}
			f.SetCellValue(sheet, fmt.Sprintf("N%d", row), join(names))
		}

		// Add up to 5 review entries
		if p.AllReviews != nil && len(p.AllReviews.GeneralReviews) > 0 {
			maxReviews := 5
			for i := 0; i < len(p.AllReviews.GeneralReviews) && i < maxReviews; i++ {
				r := p.AllReviews.GeneralReviews[i]
				startCol := 15 + i*5 // Starting column index for this review block (O = 15)

				f.SetCellValue(sheet, fmt.Sprintf("%s%d", excelColumn(startCol), row), r.Username)
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", excelColumn(startCol+1), row), r.Rating)
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", excelColumn(startCol+2), row), r.ReviewTitle)
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", excelColumn(startCol+3), row), r.ReviewDesc)
				f.SetCellValue(sheet, fmt.Sprintf("%s%d", excelColumn(startCol+4), row), r.ReviewDate)
			}
		}
	}

	// Save the file
	if err := f.SaveAs("products.xlsx"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Excel file created: products.xlsx")
}

// Helper to join string slices with newline
func join(slice []string) string {
	result := ""
	for i, v := range slice {
		if i > 0 {
			result += "\n"
		}
		result += v
	}
	return result
}

// Convert column number to Excel letter (e.g., 1 -> A, 27 -> AA)
func excelColumn(colIndex int) string {
	colName, _ := excelize.ColumnNumberToName(colIndex)
	return colName
}
