package zk_manager

import (
	"fmt"
)

const (
	PIXIE_PAPER_KEY_PREFIX         = "pixie/paper_"
	PIXIE_PAPER_PRODUCT_KEY_PREFIX = "pixie_product/paper_"
	PIXIE_LAND_KEY_PREFIX          = "pixie_land/land_"
	// PIXIE_LAND_BUILDING_KEY_PREFIX = "pixie_land/land_building_"
)

func getPixiePaperLockKey(paperID int64) string {
	return fmt.Sprintf("%s%d", PIXIE_PAPER_KEY_PREFIX, paperID)
}

func LockPixiePaper(paperID int64) (bool, string) {
	key := getPixiePaperLockKey(paperID)
	return lockRedis(key), key
}

// func DelPixiePaperLock(paperID int64) {
// 	path := LOCK_PREFIX + getPixiePaperLockKey(paperID)
// 	if err := zk.UnlockPath(zkConn, path); err != nil {
// 		if err == zk.ErrNoNode {
// 			Info("no pixie paper lock key when del", paperID)
// 		} else {
// 			Err("del paper lock path failed:", path, err)
// 		}
// 	}
// }

func getPixiePaperProductLockKey(paperID int64) string {
	return fmt.Sprintf("%s%d", PIXIE_PAPER_PRODUCT_KEY_PREFIX, paperID)
}

func LockPixiePaperProduct(paperID int64) (bool, string) {
	key := getPixiePaperProductLockKey(paperID)
	return lockRedis(key), key
}

// func DelPixiePaperProductLock(paperID int64) {
// 	path := LOCK_PREFIX + getPixiePaperProductLockKey(paperID)
// 	if err := zk.UnlockPath(zkConn, path); err != nil {
// 		if err == zk.ErrNoNode {
// 			Info("no pixie paper product lock key when del", paperID)
// 		} else {
// 			Err("del paper product lock path failed:", path, err)
// 		}
// 	}
// }

func getPixieLandLockKey(landID int64) string {
	return fmt.Sprintf("%s%d", PIXIE_LAND_KEY_PREFIX, landID)
}

func LockPixieLand(landID int64) (bool, string) {
	key := getPixieLandLockKey(landID)
	return lockRedis(key), key
}

// func DelPixieLandLock(landID int64) {
// 	path := LOCK_PREFIX + getPixieLandLockKey(landID)
// 	if err := zk.UnlockPath(zkConn, path); err != nil {
// 		if err == zk.ErrNoNode {
// 			Info("no pixie land lock key when del", landID)
// 		} else {
// 			Err("del pixie land lock path failed:", path, err)
// 		}
// 	}
// }
