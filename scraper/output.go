package scraper

import (
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
