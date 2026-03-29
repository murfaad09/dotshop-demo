package dto

type FavCurator struct {
	CuratorName       string `json:"curatorName"`
	Description       string `json:"description"`
	NumberOfFollowers string `json:"numberOfFollowers"`
	CuratorImage      string `json:"curatorImage"`
	ProductImage1     string `json:"productImage1"`
	ProductImage2     string `json:"productImage2"`
	CuratorID         string `json:"curatorId"`
}

type FavCurators struct {
	Curators []FavCurator `json:"favCurators"`
}
