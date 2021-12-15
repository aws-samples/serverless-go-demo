package types

type Product struct {
	Id    string  `dynamodbav:"id" json:"id"`
	Name  string  `dynamodbav:"name" json:"name"`
	Price float64 `dynamodbav:"price" json:"price"`
}

type ProductRange struct {
	Products []Product `json:"products"`
	Next     *string   `json:"next,omitempty"`
}
