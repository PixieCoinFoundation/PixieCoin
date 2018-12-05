package rt_party

import (
	"time"
	. "types"
)

// func ListChat(hostID int64) ([]RTPartyChat, error) {
// 	return listChatRedis(hostID,nil)
// }

func SendChat(username, content string, hostID int64, typee int, t int64) ([]RTPartyChat, error) {
	return sendChatRedis(username, content, hostID, typee, t)
}

func OfferSubject(username, subject string, hostID int64, now time.Time) (err error) {
	return offerSubjectRedis(username, subject, hostID, now)
}

func VoteSubject(self, u1, s1, u2, s2 string, hostID int64, now time.Time) (err error) {
	return voteSubjectRedis(self, u1, s1, u2, s2, hostID, now)
}

func UploadDress(username string, dyeCnt int, dyeMap map[string][7][4]float64, clothes []string, modelNo int, backgroundID, img, dressWord string, dressWordType int, hostID int64, now time.Time) (err error) {
	return setJoinerDressRedis(hostID, username, dyeCnt, dyeMap, clothes, modelNo, backgroundID, img, dressWord, dressWordType, now)
}

func VoteDress(self, u1, u2 string, hostID int64, now time.Time) (err error) {
	return voteDressRedis(self, u1, u2, hostID, now)
}

func CreateHost(username, nickname, head string, level, vip int, now time.Time, defaultSubjects []RTPartySubject) (hostID int64, err error) {
	//create in db
	if hostID, err = createHostDB(username, now.Unix()); err != nil {
		return
	}

	joiner := RTPartyHostJoiner{
		Username: username,
		Nickname: nickname,
		Head:     head,
		Level:    level,
		VIP:      vip,
	}
	joiners := []*RTPartyHostJoiner{&joiner}

	//set in redis
	h := RTPartyHost{
		ID:              hostID,
		StartTime:       now.Unix(),
		DefaultSubjects: defaultSubjects,
		Joiners:         joiners,
	}

	err = createHostRedis(&h, now)

	return
}

func JoinHost(username, nickname, head string, level, vip int, hostID int64, now time.Time) (err error) {
	return joinHostRedis(username, nickname, head, level, vip, hostID, now)
}

func LeaveHostBeforeBegin(username string, hostID int64, now time.Time) (err error) {
	return leaveHostBeforeBeginRedis(username, hostID, now)
}

func CancelAndLeaveHost(username string, hostID int64, now time.Time) (err error) {
	return cancelAndLeaveHostRedis(username, hostID, now)
}

//查看未满人房间
func CheckEmptyHost(now time.Time) (createHost bool, hostID int64, err error) {
	return checkEmptyHostRedis(now)
}

func Sync(hid int64, needExtra bool, syncUser string, now time.Time) (h RTPartyHost, err error) {
	return getHostRedis(hid, nil, needExtra, syncUser, now)
}

func SyncWithChat(hid int64, needExtra bool, syncUser string, now time.Time) (h RTPartyHost, chats []RTPartyChat, err error) {
	return getHostAndChatRedis(hid, needExtra, syncUser, now)
}

func Flush(h *RTPartyHost) error {
	return setHostRedis(h, nil)
}
