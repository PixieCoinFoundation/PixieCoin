package pushlist

import (
	"constants"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	. "logger"
	"rao"
	. "types"
)

func AddBatchSinglePushJob(jobs []*SinglePushJob) {
	if len(jobs) <= 0 {
		Info("batch add single push job size 0")
		return
	}

	conn := rao.GetConn()
	defer conn.Close()

	for _, job := range jobs {
		if job.To == "" && job.Content == "" && job.Title == "" {
			ErrMail("push content empty:", job.To, job.Title, job.Content)
			continue
		}

		jb, _ := json.Marshal(*job)

		conn.Send("RPUSH", rao.GetSinglePushJobListKey(), "'"+string(jb)+"'")
	}

	if err := conn.Flush(); err != nil {
		Err(err)
	}
}

func AddSinglePushJob(job *SinglePushJob) {
	if job.To == "" && job.Content == "" && job.Title == "" {
		ErrMail("push content empty:", job.To, job.Title, job.Content)
		return
	}

	conn := rao.GetConn()
	defer conn.Close()

	jb, _ := json.Marshal(*job)

	if _, err := conn.Do("RPUSH", rao.GetSinglePushJobListKey(), "'"+string(jb)+"'"); err != nil {
		Err(err)
	}
}

func PopSinglePushJob() (job *SinglePushJob, err error) {
	conn := rao.GetConn()
	defer conn.Close()

	var js string
	if js, err = redis.String(conn.Do("LPOP", rao.GetSinglePushJobListKey())); err != nil {
		if err.Error() != constants.RDS_NO_DATA_ERR_MSG {
			Err(err)
		} else {
			Info("no single push job list")
		}
		return
	} else if len(js) > 2 {
		js = js[1 : len(js)-1]

		job = &SinglePushJob{}
		if err = json.Unmarshal([]byte(js), job); err != nil {
			Err(err)
			return
		}
	}

	return
}
