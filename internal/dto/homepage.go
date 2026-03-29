package dto

type HomeHeroData struct {
	StylingBy     string `json:"stylingBy"`
	PhotographyBy string `json:"photographyBy"`
	Description   string `json:"description"`
	CuratorImage  string `json:"curatorImage"`
}

type CuratorSpotlight struct {
	CuratorName  string `json:"curatorName"`
	CuratorImage string `json:"curatorImage"`
	CuratorID    string `json:"curatorId"`
}

type TrendingToday struct {
	VideoLink   string `json:"videoLink"`
	Description string `json:"description"`
}

type CuratorData struct {
	Description  string `json:"description"`
	CuratorName  string `json:"curatorName"`
	CuratorImage string `json:"curatorImage"`
	CuratorID    string `json:"curatorId"`
}

type BestProduct struct {
	BrandName    string `json:"brandName"`
	ProductName  string `json:"productName"`
	Price        string `json:"price"`
	ProductID    string `json:"productId"`
	ProductImage string `json:"productImage"`
}

type BestWhiteTee struct {
	TeeImage    string `json:"teeImage"`
	Price       string `json:"price"`
	CuratorName string `json:"curatorName"`
	ProductName string `json:"productName"`
	ProductID   string `json:"productId"`
}

type HomePage struct {
	HomeHeroData     []HomeHeroData     `json:"homeHeroData"`
	CuratorSpotlight []CuratorSpotlight `json:"curatorSpotlight"`
	TrendingToday    []TrendingToday    `json:"trendingToday"`
	CuratorsData     []CuratorData      `json:"curatorsData"`
	BestOfDotShop    []BestProduct      `json:"bestOfDotShop"`
	BestWhiteTees    []BestWhiteTee     `json:"bestWhiteTees"`
}
