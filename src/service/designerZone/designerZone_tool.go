package designerZone

import (
	"appcfg"
	"constants"
	// "db/designerzone"
	// . "logger"
	"math/rand"
	// "strconv"
	// "sync"
	"fmt"
	. "logger"
	"time"
	. "types"
)

func init() {
	if appcfg.GetServerType() == constants.SERVER_TYPE_SIMPLE_TEST {
		// testRandom()
		// panic("end")
	}
}

func testRandom() {
	p1 := []*CustomPool{&CustomPool{Custom: Custom{CID: 1}}}
	p2 := []*CustomPool{&CustomPool{Custom: Custom{CID: 2}}}
	p3 := []*CustomPool{&CustomPool{Custom: Custom{CID: 3}}}
	p4 := []*CustomPool{&CustomPool{Custom: Custom{CID: 4}}}
	p5 := []*CustomPool{&CustomPool{Custom: Custom{CID: 5}}}
	p6 := []*CustomPool{&CustomPool{Custom: Custom{CID: 6}}}

	totalTestCnt := 100
	totalSize := 0
	max := 0
	min := 0
	sizeMap := make(map[int]int)

	for i := 0; i < totalTestCnt; i++ {
		size := len(getRandomClothes(p1, p2, p3, p4, p5, p6, 4))

		sizeMap[size] += 1

		totalSize += size

		if max < size {
			max = size
		}

		if size <= min || min == 0 {
			min = size
		}
	}

	Info("avg size", totalSize/totalTestCnt, max, min)
	for k, v := range sizeMap {
		Info(k, v)
	}
}

func splitCustomPoolByDate(list map[string]*CustomPool, tNow time.Time) (pool1, pool2, pool3, pool4, pool5, pool6 []*CustomPool) {
	pool1 = make([]*CustomPool, 0)
	pool2 = make([]*CustomPool, 0)
	pool3 = make([]*CustomPool, 0)
	pool4 = make([]*CustomPool, 0)
	pool5 = make([]*CustomPool, 0)
	pool6 = make([]*CustomPool, 0)

	for _, v := range list {
		tDuration := tNow.Unix() - v.EntryTime
		tWeek := float64(tDuration / 86400)
		if tWeek >= 0 && tWeek < 7 { //一周之内
			pool1 = append(pool1, v)
		} else if tWeek >= 7 && tWeek < 14 { //
			pool2 = append(pool2, v)
		} else if tWeek >= 14 && tWeek < 21 {
			pool3 = append(pool3, v)
		} else if tWeek >= 21 && tWeek < 28 {
			pool4 = append(pool4, v)
		} else if tWeek >= 28 && tWeek < 56 {
			pool5 = append(pool5, v)
		} else if tWeek >= 48 && tWeek < 84 {
			pool6 = append(pool6, v)
		}
	}

	return
}

func getRandomClothes(p1, p2, p3, p4, p5, p6 []*CustomPool, size int) (res []*CustomPool) {
	gotMap := make(map[string]int)
	res = make([]*CustomPool, 0)

	var total, num1, num2, num3, num4, num5, num6 int
	size1 := len(p1)
	size2 := len(p2)
	size3 := len(p3)
	size4 := len(p4)
	size5 := len(p5)
	size6 := len(p6)
	if size1 > 0 {
		total += 20
		num1 = total
	}

	if size2 > 0 {
		total += 16
		num2 = total
	}

	if size3 > 0 {
		total += 16
		num3 = total
	}

	if size4 > 0 {
		total += 16
		num4 = total
	}

	if size5 > 0 {
		total += 16
		num5 = total
	}

	if size6 > 0 {
		total += 16
		num6 = total
	}

	if total <= 0 {
		return
	}

	if size1+size2+size3+size4+size5+size6 <= size {
		res = append(res, p1...)
		res = append(res, p2...)
		res = append(res, p3...)
		res = append(res, p4...)
		res = append(res, p5...)
		res = append(res, p6...)
		return
	}

	for i := 0; i < size*2; i++ {
		var r *CustomPool
		var index int
		var token string
		randValue := rand.Intn(total)

		if size1 > 0 && randValue < num1 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size1)
				r = p1[index]

				if r.Inventory <= 0 {
					continue
				} else {
					break
				}
			}

			token = fmt.Sprintf("%d%d", 1, index)
		} else if size2 > 0 && randValue < num2 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size2)
				r = p2[index]

				if r.Inventory <= 0 {
					continue
				} else {
					break
				}
			}

			token = fmt.Sprintf("%d%d", 2, index)
		} else if size3 > 0 && randValue < num3 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size3)
				r = p3[index]

				if r.Inventory <= 0 {
					continue
				} else {
					break
				}
			}

			token = fmt.Sprintf("%d%d", 3, index)
		} else if size4 > 0 && randValue < num4 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size4)
				r = p4[index]

				if r.Inventory <= 0 {
					continue
				} else {
					break
				}
			}

			token = fmt.Sprintf("%d%d", 4, index)
		} else if size5 > 0 && randValue < num5 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size5)
				r = p5[index]

				if r.Inventory <= 0 {
					continue
				} else {
					break
				}
			}

			token = fmt.Sprintf("%d%d", 5, index)
		} else if size6 > 0 && randValue < num6 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size6)
				r = p6[index]

				if r.Inventory <= 0 {
					continue
				} else {
					break
				}
			}

			token = fmt.Sprintf("%d%d", 6, index)
		}

		if gotMap[token] <= 0 {
			res = append(res, r)
			gotMap[token] = 1
		}
		if len(res) >= size {
			return
		}
	}

	return
}
func getRandomDesigner(pool []*GFDesignerPool) GFDesignerPool {
	size := len(pool)
	if size > 0 {
		index := rand.Intn(size)
		return *pool[index]
	} else {
		return GFDesignerPool{}
	}
}
