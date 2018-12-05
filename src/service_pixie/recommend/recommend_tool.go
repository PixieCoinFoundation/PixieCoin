package recommend

import (
	"fmt"
	"math/rand"
	. "pixie_contract/api_specification"
	"time"
)

func splitCustomPoolByDate(list []*ItemsPaper, tNow time.Time) (pool1, pool2, pool3, pool4, pool5, pool6, pool7 []*ItemsPaper) {
	pool1 = make([]*ItemsPaper, 0)
	pool2 = make([]*ItemsPaper, 0)
	pool3 = make([]*ItemsPaper, 0)
	pool4 = make([]*ItemsPaper, 0)
	pool5 = make([]*ItemsPaper, 0)
	pool6 = make([]*ItemsPaper, 0)
	pool7 = make([]*ItemsPaper, 0)

	for _, v := range list {
		tDuration := tNow.Unix() - v.AddTime
		tDay := float64(tDuration / 86400)
		if tDay >= 0 && tDay < 7 { //一周之内
			pool1 = append(pool1, v)
		} else if tDay >= 7 && tDay < 14 { //
			pool2 = append(pool2, v)
		} else if tDay >= 14 && tDay < 21 {
			pool3 = append(pool3, v)
		} else if tDay >= 21 && tDay < 28 {
			pool4 = append(pool4, v)
		} else if tDay >= 28 && tDay < 56 {
			pool5 = append(pool5, v)
		} else if tDay >= 48 && tDay < 84 {
			pool6 = append(pool6, v)
		} else {
			pool7 = append(pool7, v)
		}
	}
	return
}

func getRandomClothes(p1, p2, p3, p4, p5, p6, p7 []*ItemsPaper, size int) (res []*ItemsPaper) {
	gotMap := make(map[string]int)
	res = make([]*ItemsPaper, 0)

	var total, num1, num2, num3, num4, num5, num6, num7 int
	size1 := len(p1)
	size2 := len(p2)
	size3 := len(p3)
	size4 := len(p4)
	size5 := len(p5)
	size6 := len(p6)
	size7 := len(p7)
	if size1 > 0 {
		total += 25
		num1 = total
	}

	if size2 > 0 {
		total += 20
		num2 = total
	}

	if size3 > 0 {
		total += 16
		num3 = total
	}

	if size4 > 0 {
		total += 13
		num4 = total
	}

	if size5 > 0 {
		total += 12
		num5 = total
	}

	if size6 > 0 {
		total += 9
		num6 = total
	}

	if size7 > 0 {
		total += 5
		num7 = total
	}

	if total <= 0 {
		return
	}

	if size1+size2+size3+size4+size5+size6+size7 <= size {
		res = append(res, p1...)
		res = append(res, p2...)
		res = append(res, p3...)
		res = append(res, p4...)
		res = append(res, p5...)
		res = append(res, p6...)
		res = append(res, p7...)
		return
	}

	for i := 0; i < size*2; i++ {
		var r *ItemsPaper
		var index int
		var token string
		randValue := rand.Intn(total)

		if size1 > 0 && randValue < num1 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size1)
				r = p1[index]
			}

			token = fmt.Sprintf("%d%d", 1, index)
		} else if size2 > 0 && randValue < num2 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size2)
				r = p2[index]
			}

			token = fmt.Sprintf("%d%d", 2, index)
		} else if size3 > 0 && randValue < num3 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size3)
				r = p3[index]
			}

			token = fmt.Sprintf("%d%d", 3, index)
		} else if size4 > 0 && randValue < num4 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size4)
				r = p4[index]
			}

			token = fmt.Sprintf("%d%d", 4, index)
		} else if size5 > 0 && randValue < num5 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size5)
				r = p5[index]
			}

			token = fmt.Sprintf("%d%d", 5, index)
		} else if size6 > 0 && randValue < num6 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size6)
				r = p6[index]
			}

			token = fmt.Sprintf("%d%d", 6, index)
		} else if size7 > 0 && randValue < num7 {
			for j := 0; j < 5; j++ {
				index = rand.Intn(size7)
				r = p7[index]
			}

			token = fmt.Sprintf("%d%d", 7, index)
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
func getRandomDesigner(pool []*ItemsPaper) ItemsPaper {
	size := len(pool)
	if size > 0 {
		index := rand.Intn(size)
		return *pool[index]
	} else {
		return ItemsPaper{}
	}
}
