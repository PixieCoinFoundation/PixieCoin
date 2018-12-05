package items

import (
	"encoding/json"
)

import (
	"appcfg"
	. "types"
)

var (
	AllItems []*Item
)

func init() {
	startItemsStr := appcfg.GetString("start_items", "[{\"itemID\":101,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":30},{\"itemID\":102,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":10},{\"itemID\":103,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":20},{\"itemID\":104,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":50},{\"itemID\":105,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":25},{\"itemID\":106,\"count\":1,\"maxCount\":99,\"priceType\":2,\"price\":10}]")

	json.Unmarshal([]byte(startItemsStr), &AllItems)
}

func GetItemPrice(itemID int) (priceType int, price int) {
	for _, v := range AllItems {
		if v.ItemID == itemID {
			price = v.Price
			priceType = v.PriceType
			break
		}
	}
	return
}
