package paper

import (
	"appcfg"
	"github.com/garyburd/redigo/redis"
	. "logger"
	. "pixie_contract/api_specification"
	"rao"
)

func init() {
	if appcfg.GetServerType() == "" && appcfg.GetBool("init_paper_redis", false) && appcfg.GetBool("main_server", false) {
		initPaperRedis()
	}
}

func initPaperRedis() {
	conn := rao.GetConn()
	defer conn.Close()

	start := 0
	batchSize := 500

	for {
		Info("init paper redis query:", start, batchSize)
		if err, res := ListDesignPaper(start, batchSize); err != nil {
			panic(err)
		} else {
			if len(res) <= 0 {
				break
			}

			for _, p := range res {
				key := rao.GetPaperKey(p.PaperID)
				conn.Send("DEL", key)

				if p.Status != int(PAPER_STATUS_DELETED) && p.Status != int(PAPER_STATUS_FAIL) {
					conn.Send("HMSET", redis.Args{}.Add(key).AddFlat(&p)...)
				}
			}
		}

		if err := conn.Flush(); err != nil {
			panic(err)
		}

		start += batchSize
	}
	Info("init done")
}
