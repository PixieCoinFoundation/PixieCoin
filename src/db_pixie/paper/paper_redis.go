package paper

import (
	"constants"
	"github.com/garyburd/redigo/redis"
	. "logger"
	. "pixie_contract/api_specification"
	"rao"
)

func AddPaperRedis(paperID int64, au, an, ah string, as int, clothesType, partType, cname, desc, paperExtra, paperFile string, status PaperStatus) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	p := BasePaper{
		PaperID:        paperID,
		AuthorUsername: au,
		AuthorNickname: an,
		AuthorHead:     ah,
		AuthorSex:      as,

		ClothesType: clothesType,
		PartType:    partType,
		Cname:       cname,
		Desc:        desc,
		Extra:       paperExtra,
		File:        paperFile,
		Status:      int(status),
	}

	key := rao.GetPaperKey(paperID)
	if _, err = conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(&p)...); err != nil {
		Err(err)
		return
	}

	return
}

func UpdatePaperVerifyRedis(paperID int64, star int, tag1, tag2, style string) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetPaperKey(paperID)
	if _, err = conn.Do("HMSET", key, "Star", star, "Tag1", tag1, "Tag2", tag2, "Style", style, "Status", PAPER_STATUS_WAIT_PUBLISH); err != nil {
		Err(err)
		return
	}

	return
}

func UpdatePaperStatusRedis(paperID int64, status PaperStatus) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetPaperKey(paperID)
	if _, err = conn.Do("HSET", key, "Status", status); err != nil {
		Err(err)
		return
	}

	return
}

func UpdatePaperOwnerRedis(paperID int64, ou, on, oh string, os int) (err error) {
	conn := rao.GetConn()
	defer conn.Close()

	key := rao.GetPaperKey(paperID)
	if _, err = conn.Do("HMSET", key, "OwnerUsername", ou, "OwnerNickname", on, "OwnerHead", oh, "OwnerSex", os, "Status", PAPER_STATUS_OCCUPY); err != nil {
		Err(err)
		return
	}

	return
}

func BatchGetPaperFromRedis(paperIDList []int64) (err error, list []*BasePaper) {
	list = make([]*BasePaper, 0)

	if len(paperIDList) <= 0 {
		return
	}

	conn := rao.GetConn()
	defer conn.Close()

	for _, id := range paperIDList {
		key := rao.GetPaperKey(id)
		conn.Send("HGETALL", key)
	}

	if err = conn.Flush(); err != nil {
		Err(err)
		return
	}

	var values []interface{}
	for i := 0; i < len(paperIDList); i++ {
		var p BasePaper

		if values, err = redis.Values(conn.Receive()); err != nil {
			if err.Error() == constants.RDS_NO_DATA_ERR_MSG {
				Info("no info for getParty")
			} else {
				Err(err)
			}
			return
		} else if err = redis.ScanStruct(values, &p); err != nil {
			Err(err)
			return
		}

		if p.PaperID > 0 {
			list = append(list, &p)
		}
	}

	return
}
