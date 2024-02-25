package game2048

var ShopItemList = []*ShopItem{
	{
		ShopItemID:  1,
		Name:        "消除道具",
		Description: "随机消除一个2或4的格子",
		BagItemID:   ItemIDRemovePiece,
		Count:       1,
		Price:       100,
	},
	{
		ShopItemID:  2,
		Name:        "金币翻倍道具",
		Description: "接下来的1次操作获得金币翻倍",
		BagItemID:   ItemIDDoubleMoney,
		Count:       1,
		Price:       50,
	},
	{
		ShopItemID:  3,
		Name:        "v我50",
		Description: "饿饿，饭饭~",
		BagItemID:   ItemIDV50,
		Count:       1,
		Price:       50,
	},
	{
		ShopItemID:  4,
		Name:        "flag8",
		Description: "获取 flag8 内容",
		BagItemID:   ItemIDFlag8,
		Count:       1,
		Price:       10000,
	},
	{
		ShopItemID:  5,
		Name:        "flagB",
		Description: "获取 flagB 内容",
		BagItemID:   ItemIDFlagB,
		Count:       1,
		Price:       999063388,
	},
}

type ShopItem struct {
	ShopItemID  int    `json:"shop_item_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BagItemID   int    `json:"bag_item_id"`
	Count       int64  `json:"count"`
	Price       int64  `json:"price"`
}
