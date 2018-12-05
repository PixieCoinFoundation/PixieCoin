package zk_manager

import (
	"fmt"
	rl "github.com/redis-lock"
	// "github.com/samuel/go-zookeeper/zk"
)

// var zkConn *zk.Conn

// var ZK_ACLS = zk.WorldACL(zk.PermAll)

const (
	LOCK_PREFIX = "/locks/"

	GUILD_LOCK_KEY_PREFIX        = "guild/gl_"
	PLAYER_LOGIN_LOCK_KEY_PREFIX = "player/login_"
	CUSTOM_KEY_PREFIX            = "custom/c_"
	PARTY_ITEM_KEY_PREFIX        = "partyItem/l_"

	RT_PARTY_KEY             = "rtparty/list"
	RT_PARTY_HOST_KEY_PREFIX = "rtparty/host_"
)

// func init() {
// 	if appcfg.GetServerType() == "" || appcfg.GetServerType() == constants.SERVER_TYPE_GL || appcfg.GetServerType() == constants.SERVER_TYPE_GM {
// 		addressStr := appcfg.GetString("zookeeper_address", "")

// 		if appcfg.GetLanguage() == constants.KOR_LANGUAGE && appcfg.GetServerType() == constants.SERVER_TYPE_GM {
// 			addressStr = appcfg.GetString("kor_zookeeper_address", "")
// 		}

// 		if addressStr == "" {
// 			panic("zookeeper address empty")
// 		}
// 		addresses := strings.Split(addressStr, ";")

// 		var err error
// 		if zkConn, _, err = zk.Connect(addresses, time.Second); err != nil {
// 			panic(err)
// 		}

// 		if ok, path := LockGuild(1); ok {
// 			Unlock(path)
// 		} else {
// 			Err("lock test failed!")
// 		}
// 	}
// }

func LockRTPartyHostList(timeout int) (bool, string) {
	return lockRedis(RT_PARTY_KEY), RT_PARTY_KEY
}

func LockRTPartyHost(hostID int64, timeout int) (bool, string) {
	key := getRTPartyHostLockKey(hostID)
	return lockRedis(key), key
}

func LockGuild(gid int64) (bool, string) {
	key := fmt.Sprintf("%s%d", GUILD_LOCK_KEY_PREFIX, gid)
	return lockRedis(key), key
}

func LockPlayerLogin(username string) (bool, string) {
	key := getPlayerLoginLockKey(username)
	return lockRedis(key), key
}

func LockCustom(cid int64) (bool, string) {
	key := fmt.Sprintf("%s%d", CUSTOM_KEY_PREFIX, cid)
	return lockRedis(key), key
}

func LockPartyItem(partyID int64, username string) (bool, string) {
	key := getPartyItemLockKey(partyID, username)
	return lockRedis(key), key
}

// func DelRTPartyHostLock(hostID int64) {
// 	path := LOCK_PREFIX + getRTPartyHostLockKey(hostID)
// 	if err := zk.UnlockPath(zkConn, path); err != nil {
// 		if err == zk.ErrNoNode {
// 			Info("no rt party host lock key:", hostID)
// 		} else {
// 			Err("del path failed:", path, err)
// 		}
// 	}
// }

// func DelPartyItemLock(partyID int64, username string) {
// 	path := LOCK_PREFIX + getPartyItemLockKey(partyID, username)
// 	if err := zk.UnlockPath(zkConn, path); err != nil {
// 		if err == zk.ErrNoNode {
// 			Info("no party item lock key:", partyID, username)
// 		} else {
// 			Err("del path failed:", path, err)
// 		}
// 	}
// }

// func DelPlayerLoginLock(username string) {
// 	path := LOCK_PREFIX + getPlayerLoginLockKey(username)
// 	if err := zk.UnlockPath(zkConn, path); err != nil {
// 		if err == zk.ErrNoNode {
// 			Info("no player login lock key:", username)
// 		} else {
// 			Err("del path failed:", path, err)
// 		}
// 	}
// }

func getRTPartyHostLockKey(hostID int64) string {
	return fmt.Sprintf("%s%d", RT_PARTY_HOST_KEY_PREFIX, hostID)
}

func getPartyItemLockKey(partyID int64, username string) string {
	return fmt.Sprintf("%s%d_%s", PARTY_ITEM_KEY_PREFIX, partyID, username)
}

func getPlayerLoginLockKey(username string) string {
	return fmt.Sprintf("%s%s", PLAYER_LOGIN_LOCK_KEY_PREFIX, username)
}

//zookeeper lock
// func lock(key string, timeout int) (bool, string) {
// 	// if conn := zkPool.Get(); conn != nil {
// 	path := LOCK_PREFIX + key
// 	if pathc, err := zk.LockPathWithTimeout(zkConn, path, ZK_ACLS, timeout); err != nil {
// 		Err("lock key failed:", key, err)

// 		if e, nodeState, err := zkConn.Exists(path); err != nil {
// 			ErrMail("check exists key failed:", key, err)
// 		} else if e && nodeState != nil {
// 			ErrMail("path:", path, "children:", nodeState.NumChildren)
// 		}

// 		return false, pathc
// 	} else {
// 		return true, pathc
// 	}
// }

func Unlock(key string) {
	// if err := zk.UnlockPath(zkConn, path); err != nil {
	// 	Err("unlock path failed:", path, err)
	// } else {
	// 	// Info("unlock path succeed:", path)
	// }

	rl.Unlock(key)
}

func lockRedis(key string) bool {
	ok, _ := rl.Lock(key)

	return ok
}
