package models

type VariantOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Dimensions struct {
	Length int    `json:"length"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Units  string `json:"units"`
}

type WholesalePricing struct {
	Min struct {
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
	} `json:"min"`
	Max struct {
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
	} `json:"max"`
}
type Image struct {
	ID         string   `json:"id"`
	Source     string   `json:"source"`
	Position   int      `json:"position"`
	VariantIDs []string `json:"variantIds"`
}
type GoogleProductCategory struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}

type Data struct {
	ID                      string                  `json:"id"`
	CompanyID               string                  `json:"companyId"`
	Title                   string                  `json:"title"`
	TitleTranslations       *map[string]string      `json:"titleTranslations"`
	Brand                   string                  `json:"brand"`
	Description             string                  `json:"description"`
	DescriptionTranslations *map[string]string      `json:"descriptionTranslations"`
	SellerReference         string                  `json:"sellerReference"`
	Variants                *[]Variant              `json:"variants"`
	Tags                    *[]string               `json:"tags"`
	Images                  *[]Image                `json:"images"`
	OptionNames             *[]string               `json:"optionNames"`
	GoogleProductCategory   *GoogleProductCategory  `json:"googleProductCategory"`
	Attributes              *map[string]interface{} `json:"attributes"`
	SEOTitle                string                  `json:"seoTitle"`
	Created                 string                  `json:"created"`
	Updated                 string                  `json:"updated"`
	WholesalePricing        WholesalePricing        `json:"wholesalePricing"`
	Metafields              *map[string]interface{} `json:"metafields"`
}
type Variant struct {
	ID              string              `json:"id"`
	SKU             string              `json:"sku"`
	Title           string              `json:"title"`
	InventoryAmount int                 `json:"inventoryAmount"`
	RetailPrice     float64             `json:"retailPrice"`
	RetailCurrency  string              `json:"retailCurrency"`
	BasePrice       float64             `json:"basePrice"`
	BaseCurrency    string              `json:"baseCurrency"`
	CompareAtPrice  string              `json:"compareAtPrice"`
	Barcode         string              `json:"barcode"`
	BarcodeType     string              `json:"barcodeType"`
	Options         []map[string]string `json:"options"`
	SkipCount       bool                `json:"skipCount"`
	Weight          int                 `json:"weight"`
	WeightUnits     string              `json:"weightUnits"`
	Dimensions      struct {
		Length int    `json:"length"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Units  string `json:"units"`
	} `json:"dimensions"`
	Attributes *map[string]interface{} `json:"attributes"`
	Metafields *map[string]interface{} `json:"metafields"`
}
type ProductList struct {
	HasMore  bool        `json:"hasMore"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Data     []Data      `json:"data"`
	Error    interface{} `json:"error"`
}

type Attributes struct {
	Price                  Attribute `json:"price"`
	Quantity               Attribute `json:"quantity"`
	RippedJeansWithPockets Attribute `json:"ripped jeans with pockets"`
	Variant                Attribute `json:"variant"`
}

type Attribute struct {
	Value string `json:"value"`
}

type Product struct {
	Data  Data  `json:"data"`
	Error []any `json:"error"`
}
