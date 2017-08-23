package lib

import "time"

type APIStatsPricesItem struct {
	ItemID           string    `json:"item_id"`
	City             string    `json:"city"`
	SellPriceMin     int       `json:"sell_price_min"`
	SellPriceMinDate time.Time `json:"sell_price_min_date"`
	SellPriceMax     int       `json:"sell_price_max"`
	SellPriceMaxDate time.Time `json:"sell_price_max_date"`
	BuyPriceMin      int       `json:"buy_price_min"`
	BuyPriceMinDate  time.Time `json:"buy_price_min_date"`
	BuyPriceMax      int       `json:"buy_price_max"`
	BuyPriceMaxDate  time.Time `json:"buy_price_max_date"`
}
