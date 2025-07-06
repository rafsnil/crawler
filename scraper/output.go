package scraper

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"simple-go-crawler/model"
	"sync"
)

type Output struct {
	Products []*model.Product
	mutex    sync.Mutex
}

var GlobalOutput *Output

func init() {
	GlobalOutput = &Output{
		Products: make([]*model.Product, 0),
	}
}

func (o *Output) AddProduct(product *model.Product) {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	o.Products = append(o.Products, product)
}

func (o *Output) GetProducts() []*model.Product {
	o.mutex.Lock()
	defer o.mutex.Unlock()
	products := make([]*model.Product, len(o.Products))
	copy(products, o.Products)
	return products
}

func (o *Output) PrintToExcelSheet() {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	f := excelize.NewFile()
	sheet := "Products"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)

	// Header
	headers := []string{
		"Name", "Title", "Description", "Itemized Description",
		"Size Chart", "Category", "Price", "Available Sizes",
		"Image URLs", "Details URL", "Avg Rating", "Total Reviews", "Recommend %",
		"Coordinate Products",
	}
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	// Write product data
	for rowIdx, p := range o.Products {
		row := rowIdx + 2 // start from row 2
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
			// Concatenate all coordinate product names
			var names []string
			for _, set := range p.Coordinates.Coordinates {
				for _, cp := range set.CoordinateProducts {
					names = append(names, cp.Name+"("+cp.Price+")")
				}
			}
			f.SetCellValue(sheet, fmt.Sprintf("N%d", row), join(names))
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
