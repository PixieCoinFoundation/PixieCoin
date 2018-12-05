package orders

import (
	"appcfg"
	"constants"
	"dao"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/garyburd/redigo/redis"
	"rao"
)

import (
	. "logger"
	. "types"
)

func FindO(orderUID string) (o Order, err error) {
	if o, err = getR(orderUID); err != nil {
		return getDB(orderUID)
	}
	return
}

func getDB(orderUID string) (o Order, err error) {
	var rw string
	if err = dao.GetOrderStmt.QueryRow(orderUID).Scan(&o.OrderId, &o.Username, &o.DiamondId, &o.OrderTime, &o.PayTime, &o.Status, &o.Diamond, &o.Code, &o.RmbCost, &o.Channel, &o.Info, &o.OrderUID, &rw); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("can't find order in db:", orderUID)
		} else {
			Err(err)
		}
		return
	}

	return
}

// //CleanKorOrderLog 清楚订单日志
// func CleanKorOrderLog(ds string) (bool, error) {
// 	if len(ds) < 1 {
// 		Info("ds can not empty:", ds)
// 		return false, errors.New("ds can not empty")
// 	}
// 	rows, err := dao.CleanKorOrderlogStmt.Exec(ds)
// 	if err != nil {
// 		Info("can't find data in db:", ds)
// 		return false, err
// 	}
// 	if num, err := rows.RowsAffected(); err != nil {
// 		return false, err
// 	} else if num > 0 {
// 		return true, nil
// 	}
// 	return false, errors.New("delect data fail or database not exist")

// }

// //QueryKorOrderlog 查询韩国某小时的日志
// func QueryKorOrderlog(lim int, et int64) (ds string) {
// 	rows, err := dao.QueryKorOrderLogStmt.Query(lim)
// 	if err != nil {
// 		Err(err)
// 		return
// 	} else {
// 		for rows.Next() {
// 			var tid int64
// 			var fbt string
// 			if err := rows.Scan(&tid, &fbt); err != nil {
// 				Err(err)
// 				return
// 			}
// 			if tid < et {
// 				ds = fbt
// 			}
// 		}
// 	}
// 	return
// }

func getR(orderUID string) (o Order, err error) {
	o.OrderUID = orderUID

	conn := rao.GetConn()
	defer conn.Close()

	key := constants.ORDER_HASH_PREFIX + orderUID
	var values []interface{}
	if values, err = redis.Values(conn.Do("HMGET", key, "username", "diamond_id", "order_time", "pay_time", "status", "diamond", "code", "rmb_cost", "channel", "info")); err != nil {
		Err(err)
		return
	}

	if values, err = redis.Scan(values, &o.Username); err != nil {
		Err(err)
		return
	} else if o.Username == "" {
		err = errors.New("can't find order in redis")
		return
	}
	if values, err = redis.Scan(values, &o.DiamondId); err != nil {
		Err(err)
		return
	}
	if values, err = redis.Scan(values, &o.OrderTime); err != nil {
		Err(err)
		return
	}
	if values, err = redis.Scan(values, &o.PayTime); err != nil {
		Err(err)
		return
	}
	if values, err = redis.Scan(values, &o.Status); err != nil {
		Err(err)
		return
	}
	if values, err = redis.Scan(values, &o.Diamond); err != nil {
		Err(err)
		return
	}
	if values, err = redis.Scan(values, &o.Code); err != nil {
		Err(err)
		return
	}
	if values, err = redis.Scan(values, &o.RmbCost); err != nil {
		Err(err)
		return
	}
	if values, err = redis.Scan(values, &o.Channel); err != nil {
		Err(err)
		return
	}
	if values, err = redis.Scan(values, &o.Info); err != nil {
		Err(err)
		return
	}

	return
}

func addDB(o *Order) (id int64, err error) {
	// go order_batch.AddOrder(o)
	var res sql.Result
	if res, err = dao.AddOrderStmt.Exec(o.Username, o.DiamondId, o.OrderTime, o.PayTime, o.Status, o.Diamond, o.Code, o.RmbCost, o.Channel, o.Info, o.OrderUID, ""); err != nil {
		Err(err)
		return
	} else if id, err = res.LastInsertId(); err != nil {
		Err(err)
		return
	}

	return
}

func addR(o *Order) (err error) {
	c := rao.GetConn()
	defer c.Close()

	key := constants.ORDER_HASH_PREFIX + o.OrderUID

	// c.Send("MULTI")
	c.Send("HSET", key, "username", o.Username)
	c.Send("HSET", key, "diamond_id", o.DiamondId)
	c.Send("HSET", key, "order_time", o.OrderTime)
	c.Send("HSET", key, "pay_time", o.PayTime)
	c.Send("HSET", key, "status", o.Status)
	c.Send("HSET", key, "diamond", o.Diamond)
	c.Send("HSET", key, "code", o.Code)
	c.Send("HSET", key, "rmb_cost", o.RmbCost)
	c.Send("HSET", key, "channel", o.Channel)
	c.Send("HSET", key, "info", o.Info)

	c.Send("EXPIRE", key, "1200")
	if err = c.Flush(); err != nil {
		Err(err)
		return
	}

	// Info("SUCCESS insert order to redis username:", o.Username, "order id:", o.OrderId)
	return nil
}

type orderInfo struct {
	OrderId       int64       `bson:"orderId" json:"orderId"`
	Username      string      `bson:"username" json:"username"`
	DiamondId     int         `bson:"diamondID" json:"diamondID"`
	OrderTime     int64       `bson:"orderTime" json:"orderTime"`
	PayTime       int64       `bson:"payTime" json:"payTime"`
	Status        OrderStatus `bson:"status" json:"status"`
	Diamond       int         `bson:"diamond" json:"diamond"`
	Code          string      `bson:"code" json:"code"`
	RmbCost       float64     `bson:"rmbcost" json:"rmbcost"`
	Channel       string      `bson:"channel" json:"channel"`
	Info          string      `bson:"info" json:"info"`
	DeviceId      string      `bson:"deviceId" json:"info"`
	OrderUID      string      `bson:"orderUID" json:"info"`
	ServerAddress string      `bson:"serverAddress" json:"info"`
}

func InsertO(obj *Order, deviceId string) (uid string, err error) {
	uid = obj.OrderUID

	var iid int64
	if iid, err = addDB(obj); err != nil {
		return
	}

	//gmlog
	i := orderInfo{
		OrderId:       iid,
		Username:      obj.Username,
		DiamondId:     obj.DiamondId,
		OrderTime:     obj.OrderTime,
		PayTime:       obj.PayTime,
		Status:        obj.Status,
		Diamond:       obj.Diamond,
		Code:          obj.Code,
		RmbCost:       obj.RmbCost,
		Channel:       obj.Channel,
		Info:          obj.Info,
		DeviceId:      deviceId,
		OrderUID:      uid,
		ServerAddress: appcfg.GetAddress(),
	}
	data, _ := json.Marshal(i)
	GMLog(constants.C1_PLAYER, constants.C2_PAY, constants.C3_DEFAULT, obj.Username, string(data))

	if err = addR(obj); err != nil {
		return
	}

	return
}

func updateInDB(orderUID string, payTime int64, status OrderStatus, channel, info, reward string) (err error) {
	var res sql.Result
	var effect int64
	if res, err = dao.UpdateOrderStmt.Exec(payTime, status, channel, info, reward, orderUID); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect <= 0 {
		err = errors.New("update order effect 0")
	}

	return
}

func updateInRedis(orderUID string, payTime int64, status OrderStatus, channel string, info string) (err error) {
	key := constants.ORDER_HASH_PREFIX + orderUID

	c := rao.GetConn()
	defer c.Close()

	c.Send("HSET", key, "pay_time", payTime)
	c.Send("HSET", key, "status", status)
	c.Send("HSET", key, "channel", channel)
	c.Send("HSET", key, "info", info)

	if err = c.Flush(); err != nil {
		Err(err)
		return
	}

	return
}

func GetPlayerOrder(username string, status OrderStatus) (res []Order, err error) {
	var rows *sql.Rows
	if rows, err = dao.GetPlayerOrderStmt.Query(username, status); err != nil {
		Err(err)
		return
	} else {
		defer rows.Close()
		res = make([]Order, 0)
		for rows.Next() {
			var o Order
			var rs string
			if err = rows.Scan(&o.OrderId, &o.Username, &o.DiamondId, &o.OrderTime, &o.PayTime, &o.Status, &o.Diamond, &o.Code, &o.RmbCost, &o.Channel, &o.Info, &o.OrderUID, &rs); err != nil {
				Err(err)
				return
			}

			res = append(res, o)
		}
	}
	return
}

//GetPlayerOrderCode 获取订单码
func GetPlayerOrderCode(paycode string) (exist bool, err error) {
	var (
		id   int
		code string
	)
	if err = dao.GetPlayerOrderCodeisExistStmt.QueryRow(paycode).Scan(&id, &code); err != nil {
		if err.Error() == constants.SQL_NO_DATA_ERR_MSG {
			Info("can't find order in db:", paycode)
			return false, nil
		}
		Err(err)
		return true, err
	}
	if len(code) > 1 {
		return true, err
	}
	return false, nil
}

//AddPlayerOrderCode 增加苹果收据表信息
func AddPlayerOrderCode(paycode string) error {
	_, err := dao.AddPlayerOrderCodeUniqueStmt.Exec(paycode)
	return err
}

func UpdateStatus(orderUID string, status OrderStatus, formerStatus OrderStatus) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao.UpdateOrderStatusStmt.Exec(status, orderUID, formerStatus); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		success = true
		//update in redis
		key := constants.ORDER_HASH_PREFIX + orderUID

		c := rao.GetConn()
		defer c.Close()

		// c.Send("HSET", key, "pay_time", payTime)
		c.Send("HSET", key, "status", status)
		// c.Send("HSET", key, "channel", channel)
		// c.Send("HSET", key, "info", info)

		if err = c.Flush(); err != nil {
			Err(err)
			return
		}

		return
	}
	return
}

func UpdateStatusInfo(orderUID string, info string, status OrderStatus, formerStatus OrderStatus) (success bool, err error) {
	var res sql.Result
	var effect int64
	if res, err = dao.UpdateOrderStatusInfoStmt.Exec(status, info, orderUID, formerStatus); err != nil {
		Err(err)
		return
	} else if effect, err = res.RowsAffected(); err != nil {
		Err(err)
		return
	} else if effect > 0 {
		success = true
		//update in redis
		key := constants.ORDER_HASH_PREFIX + orderUID

		c := rao.GetConn()
		defer c.Close()

		// c.Send("HSET", key, "pay_time", payTime)
		c.Send("HSET", key, "status", status)
		// c.Send("HSET", key, "channel", channel)
		c.Send("HSET", key, "info", info)

		if err = c.Flush(); err != nil {
			Err(err)
			return
		}

		return
	}
	return
}
func Update(orderUID string, payTime int64, status OrderStatus, channel, info, username, reward string) (err error) {
	if err = updateInDB(orderUID, payTime, status, channel, info, reward); err != nil {
		return
	}

	// d := OrderUpdateExtra{
	// 	OrderUID: orderUID,
	// 	PayTime:  payTime,
	// 	Status:   int(status),
	// 	Channel:  channel,
	// 	Info:     info,
	// }
	// data, eerr := json.Marshal(d)
	// if eerr != nil {
	// 	Err(err)
	// }
	// //gm log
	// GMLog(constants.C1_PLAYER, constants.C2_PAY, constants.C3_DEFAULT, username, string(data))

	if err = updateInRedis(orderUID, payTime, status, channel, info); err != nil {
		return
	}

	return
}

// type SimpleOrder struct {
// 	Username  string
// 	OrderTime string
// 	PayTime   string
// 	DiamondId int
// 	RmbCost   int
// 	Info      string
// }

// func GetOrders(status OrderStatus, st int64, et int64) (res []SimpleOrder, err error) {
// 	var rows *sql.Rows
// 	if rows, err = dao.GetAllOrderStmt.Query(status, st, et); err != nil {
// 		Err(err)
// 		return
// 	} else {
// 		res = make([]SimpleOrder, 0)
// 		for rows.Next() {
// 			var o SimpleOrder
// 			if err = rows.Scan(&o.Username, &o.OrderTime, &o.PayTime, &o.DiamondId, &o.RmbCost, &o.Info); err != nil {
// 				Err(err)
// 				return
// 			}

// 			res = append(res, o)
// 		}
// 	}
// 	return
// }
