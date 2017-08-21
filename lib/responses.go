package lib

type APIStatsPricesItem struct {
	City         string `json:"city"`
	SellPriceMin int    `json:"sell_price_min"`
	SellPriceMax int    `json:"sell_price_max"`
	BuyPriceMin  int    `json:"buy_price_min"`
	BuyPriceMax  int    `json:"buy_price_max"`
}
