package types

import ()

type KuaikanTransData struct {
	AppID        string `json:"app_id"`
	WaresID      int    `json:"wares_id"`
	OutOrderID   string `json:"out_order_id"`
	OpenUID      string `json:"open_uid"`
	OutNotifyURL string `json:"out_notify_url"`
}

type OrderStatus int

type MHRPayInfo struct {
	AppID     string            `json:"appId"`
	BizName   string            `json:"bizName"`
	BizArgs   MHRPayInfoBizArgs `json:"bizArgs"`
	BuyerID   string            `json:"buyerId"`
	Channel   string            `json:"channel"`
	ProxyType int               `json:"proxyType"`
	Cart      MHRPayInfoCart    `json:"cart"`
}

type MHRPayInfoBizArgs struct {
	GameID          string `json:"gameId"`
	MerchantOrderID string `json:"merchantOrderId"`
}

type MHRPayInfoCart struct {
	Amount MHRPayInfoAmount  `json:"amount"`
	Items  []MHRPayInfoItem  `json:"items"`
	Extra  map[string]string `json:"extra"`
}

type MHRPayInfoAmount struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}

type MHRPayInfoItem struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       MHRPayInfoAmount  `json:"price"`
	Quantity    int               `json:"quantity"`
	Extra       map[string]string `json:"extra"`
}

const (
	ORDER_WAIT    OrderStatus = 0 // 等待处理
	ORDER_SUCCEED OrderStatus = 1 // 购买成功，并且已经给到玩家
	ORDER_FAILED  OrderStatus = 2 // 购买失败
	ORDER_CANCEL  OrderStatus = 3 // 取消订单
	ORDER_REFUND  OrderStatus = 4 // 退款成功
	ORDER_PAID    OrderStatus = 5 // webhook返回结果，已付款，但是还没给货到玩家
	ORDER_REVOKE  OrderStatus = 6 //玩家申请取消
)

type Order struct {
	OrderId   int64       `bson:"orderId" json:"orderId"`
	Username  string      `bson:"username" json:"username"`
	DiamondId int         `bson:"diamondID" json:"diamondID"`
	OrderTime int64       `bson:"orderTime" json:"orderTime"`
	PayTime   int64       `bson:"payTime" json:"payTime"`
	Status    OrderStatus `bson:"status" json:"status"`
	Diamond   int         `bson:"diamond" json:"diamond"`
	Code      string      `bson:"code" json:"code"`
	RmbCost   float64     `bson:"rmbcost" json:"rmbcost"`
	Channel   string      `bson:"channel" json:"channel"`
	Info      string      `bson:"info" json:"info"`
	OrderUID  string      `bson:"orderUID" json:"info"`
}

type RawOrder struct {
	DiamondID  string  `bson:"diamondID" json:"diamondID"`
	OrderTime  int64   `bson:"orderTime" json:"orderTime"`
	PayTime    int64   `bson:"payTime" json:"payTime"`
	OrderTimeS string  `bson:"orderTimeS" json:"orderTimeS"`
	PayTimeS   string  `bson:"payTimeS" json:"payTimeS"`
	Status     string  `bson:"status" json:"status"`
	Diamond    int     `bson:"diamond" json:"diamond"`
	Code       string  `bson:"code" json:"code"`
	RmbCost    float64 `bson:"rmbcost" json:"rmbcost"`
	Channel    string  `bson:"channel" json:"channel"`
	Info       string  `bson:"info" json:"info"`
	OrderUID   string  `bson:"orderUID" json:"info"`
	Reward     string  `bson:"reward" json:"reward"`
	RawStatus  int     `bson:"rawStatus" json:"reward"`
	// 订单商品名称
	OrderProduceName string `bson:"rawStatus" json:"orderproducename"`
}

type OrderReward struct {
	DiamondID  int
	Tili       int
	Gold       int64
	Diamond    int
	Clothes    string
	VIPSuit    string
	DesignCoin int

	//extra info
	PlatformOrderID string
}
