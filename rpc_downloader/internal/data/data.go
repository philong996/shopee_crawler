package data

type JobMessage struct {
	Shop_id int `json:shop_id`
	Interval int `json:interval`
}

type Item struct {
	Itemid int `json:itemid`
}

type DetailProduct struct {
	Item Item `json:item`
}

type ShopPage struct {
	Items []Item `json:items`
}

