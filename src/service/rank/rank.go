/*
* 关卡评分
 */
package rank

import (
	"appcfg"
	"common"
	"constants"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

import (
	// jl "jsonloader"
	. "logger"
	"service/clothes"
	// "service/params"
	. "types"
	"xlsx"
)

var (
	AllLevelsP         []*GFLevelP
	LevelMapP          map[string]*GFLevelP
	TheaterLevelCntMap map[int]int

	STAGE_CZ_END_INDEX int
	STAGE_GZ_END_INDEX int
	STAGE_DX_END_INDEX int
	STAGE_SH_END_INDEX int
	STAGE_5P_END_INDEX int

	STAGE_CZ_END_LEVEL_INDEX int
)

func init() {
	if appcfg.GetServerType() != "" && appcfg.GetServerType() != constants.SERVER_TYPE_SNAPSHOT_PLAYER && appcfg.GetServerType() != constants.SERVER_TYPE_SIMPLE_TEST {
		return
	}

	rand.Seed(time.Now().Unix())
	LevelMapP = make(map[string]*GFLevelP, 0)
	AllLevelsP = make([]*GFLevelP, 0)
	TheaterLevelCntMap = make(map[int]int)

	// loadLevels()
	// loadLevelsByXLSX()
	// loadLevelByXlsxs()

	// for _, m := range AllLevelsP {
	// 	if m.Planet == 1 && m.Level > STAGE_CZ_END_LEVEL_INDEX {
	// 		STAGE_CZ_END_LEVEL_INDEX = m.Level
	// 	}

	// 	if m.Planet == 1 && m.No > STAGE_CZ_END_INDEX {
	// 		STAGE_CZ_END_INDEX = m.No
	// 	} else if m.Planet == 2 && m.No > STAGE_GZ_END_INDEX {
	// 		STAGE_GZ_END_INDEX = m.No
	// 	} else if m.Planet == 3 && m.No > STAGE_DX_END_INDEX {
	// 		STAGE_DX_END_INDEX = m.No
	// 	} else if m.Planet == 4 && m.No > STAGE_SH_END_INDEX {
	// 		STAGE_SH_END_INDEX = m.No
	// 	} else if m.Planet == 5 && m.No > STAGE_5P_END_INDEX {
	// 		STAGE_5P_END_INDEX = m.No
	// 	}
	// }

	// Info("CZ END:", STAGE_CZ_END_INDEX, "GZ END:", STAGE_GZ_END_INDEX, "DX END:", STAGE_DX_END_INDEX, "SH END:", STAGE_SH_END_INDEX, "5P END:", STAGE_5P_END_INDEX, "CZ END LEVEL:", STAGE_CZ_END_LEVEL_INDEX)

	// Info("all mission cnt:", len(AllLevelsP))
	// Info("theater level cnt map:", TheaterLevelCntMap)
	// for _, l := range AllLevelsP {
	// 	LevelMapP[l.LevelID] = l
	// }

	// Info("check mission level:", GetLevelMissionLevel("1016a"), GetLevelMissionLevel("1016b"), GetLevelMissionLevel("1016b"), GetLevelMissionLevel("1016d"))
}

// func loadLevels() {
// 	if err := jl.LoadFile("scripts/data/data_mission.json", &AllLevelsP); err != nil {
// 		panic(err)
// 	}
// }

// func CheckLevelOpen(lid string) bool {
// 	if l := getLevelByLevelID(lid); l != nil {
// 		if l.TheaterID != 0 {
// 			//check theater
// 			return params.TheaterOpen(l.TheaterID)
// 		} else {
// 			return true
// 		}
// 	}
// 	return false
// }

func loadLevelByXlsxs() {
	if appcfg.GetLanguage() == constants.KOR_LANGUAGE {
		loadLevelsByXLSX("scripts/data/data_mission_kor.xlsx")
		loadLevelsByXLSX("scripts/data/data_mission_t_1_kor.xlsx")
		loadLevelsByXLSX("scripts/data/data_mission_t_2_kor.xlsx")
		loadLevelsByXLSX("scripts/data/data_mission_t_3_kor.xlsx")
		loadLevelsByXLSX("scripts/data/data_mission_t_4_kor.xlsx")
	} else {
		loadLevelsByXLSX("scripts/data/data_mission.xlsx")
		loadLevelsByXLSX("scripts/data/data_mission_t_1.xlsx")
		loadLevelsByXLSX("scripts/data/data_mission_t_2.xlsx")
		loadLevelsByXLSX("scripts/data/data_mission_t_3.xlsx")
		loadLevelsByXLSX("scripts/data/data_mission_t_4.xlsx")
		loadLevelsByXLSX("scripts/data/data_mission_t_5.xlsx")
		loadLevelsByXLSX("scripts/data/data_mission_t_6.xlsx")
		loadLevelsByXLSX("scripts/data/data_mission_t_7.xlsx")
	}
}

func loadLevelsByXLSX(name string) {
	Info("load file:", name)
	var b3Cnt, tncc int
	xlFile, err := xlsx.OpenFile(name)
	if err != nil {
		Err(err)
		return
	}
	var lidIdx, noIdx, typeIdx, cd1Idx, cd2Idx, gsIdx, gaIdx, gbIdx, gfIdx, bnIdx, thIdx, plnIdx, mlIdx int
	for _, sheet := range xlFile.Sheets {
		if sheet.Name == "Sheet1" {
			for i, row := range sheet.Rows {
				if i == 0 {
					for j, cell := range row.Cells {
						v, _ := cell.String()
						if v == "id" {
							lidIdx = j
						} else if v == "no" {
							noIdx = j
						} else if v == "type" {
							typeIdx = j
						} else if v == "clothes_drop1" {
							cd1Idx = j
						} else if v == "clothes_drop2" {
							cd2Idx = j
						} else if v == "gold_drop_s" {
							gsIdx = j
						} else if v == "gold_drop_a" {
							gaIdx = j
						} else if v == "gold_drop_b" {
							gbIdx = j
						} else if v == "gold_drop_f" {
							gfIdx = j
						} else if v == "branch" {
							bnIdx = j
						} else if v == "theater_id" {
							thIdx = j
						} else if v == "planet" {
							plnIdx = j
						} else if v == "level" {
							mlIdx = j
						}
					}
					Info("mission row index:", lidIdx, noIdx, typeIdx, cd1Idx, cd2Idx, gsIdx, gaIdx, gbIdx, gfIdx, bnIdx, thIdx, plnIdx)
					continue
				}
				l := GFLevelP{}
				for j, cell := range row.Cells {
					v, _ := cell.String()

					if j == lidIdx {
						//levelid
						l.LevelID = v
					} else if j == noIdx {
						l.No = getInt(v, "mission no")
					} else if j == typeIdx {
						l.Type = getInt(v, "mission type")
					} else if j == cd1Idx {
						l.ClothesDrop1 = v
					} else if j == cd2Idx {
						l.ClothesDrop2 = v
					} else if j == gsIdx {
						l.GoldDropS = getInt(v, "mission gold drop s")
					} else if j == gaIdx {
						l.GoldDropA = getInt(v, "mission gold drop a")
					} else if j == gbIdx {
						l.GoldDropB = getInt(v, "mission gold drop b")
					} else if j == gfIdx {
						l.GoldDropF = getInt(v, "mission gold drop f")
					} else if j == bnIdx {
						l.Branch = getInt(v, "mission branch")
					} else if j == thIdx {
						l.TheaterID = getInt(v, "mission theater id")
					} else if j == plnIdx {
						l.Planet = getInt(v, "mission planet id")
					} else if j == mlIdx {
						l.Level = getInt(v, "mission level")
					}
				}
				if i == 1 {
					Info("mission sample:", l)
				}
				if l.Branch == 3 {
					b3Cnt++
				}
				if l.TheaterID > 0 {
					TheaterLevelCntMap[l.TheaterID] += 1
					if l.ClothesDrop1 != "" && clothes.CloPieceChance(l.ClothesDrop1) <= 0 {
						Info("theatre clothes:" + l.ClothesDrop1 + ",not in data_chance!")
						tncc++
					}
					if l.ClothesDrop2 != "" && clothes.CloPieceChance(l.ClothesDrop2) <= 0 {
						Info("theatre clothes:" + l.ClothesDrop1 + ",not in data_chance!")
						tncc++
					}
				}

				AllLevelsP = append(AllLevelsP, &l)
			}
		}
	}

	if b3Cnt <= 0 {
		panic("no ending(branch=3) in :" + name)
	}

	if tncc > 0 {
		panic("theatre clothes not in data_chance!")
	}
}

func LoadSpecialMissionFile(file *xlsx.File) (res []GFLevelS) {
	res = make([]GFLevelS, 0)

	var lidIdx, noIdx, typeIdx int
	var sIdx, aIdx, bIdx, fIdx, keyIdx, char1Idx, char2Idx, slineIdx, alineIdx, blineIdx, adjustIdx int
	var warmIdx, formalIdx, tightIdx, brightIdx, darkIdx, cuteIdx, manIdx, toughIdx, nobleIdx, strangeIdx, sexyIdx, sportIdx int
	for _, sheet := range file.Sheets {
		if sheet.Name == "Sheet1" {
			for i, row := range sheet.Rows {
				if i == 0 {
					for j, cell := range row.Cells {
						v, _ := cell.String()
						if v == "id" {
							lidIdx = j
						} else if v == "no" {
							noIdx = j
						} else if v == "s" {
							sIdx = j
						} else if v == "a" {
							aIdx = j
						} else if v == "b" {
							bIdx = j
						} else if v == "f" {
							fIdx = j
						} else if v == "keyword" {
							keyIdx = j
						} else if v == "type" {
							typeIdx = j
						} else if v == "char1" {
							char1Idx = j
						} else if v == "char2" {
							char2Idx = j
						} else if v == "s_line" {
							slineIdx = j
						} else if v == "a_line" {
							alineIdx = j
						} else if v == "b_line" {
							blineIdx = j
						} else if v == "adjust" {
							adjustIdx = j
						} else if v == "warm" {
							warmIdx = j
						} else if v == "formal" {
							formalIdx = j
						} else if v == "tight" {
							tightIdx = j
						} else if v == "bright" {
							brightIdx = j
						} else if v == "dark" {
							darkIdx = j
						} else if v == "cute" {
							cuteIdx = j
						} else if v == "man" {
							manIdx = j
						} else if v == "tough" {
							toughIdx = j
						} else if v == "noble" {
							nobleIdx = j
						} else if v == "strange" {
							strangeIdx = j
						} else if v == "sexy" {
							sexyIdx = j
						} else if v == "sport" {
							sportIdx = j
						}
					}
					// Info("mission row index:", lidIdx, noIdx, typeIdx, cd1Idx, cd2Idx, gsIdx, gaIdx, gbIdx, gfIdx, bnIdx, thIdx, plnIdx)
					continue
				}
				l := GFLevelS{}
				for j, cell := range row.Cells {
					v, _ := cell.String()

					if j == lidIdx {
						//levelid
						l.LevelID = v
					} else if j == noIdx {
						l.No = getInt(v, "mission no")
					} else if j == typeIdx {
						l.Type = getInt(v, "mission type")
					} else if j == sIdx {
						l.S = v
					} else if j == aIdx {
						l.A = v
					} else if j == bIdx {
						l.B = v
					} else if j == fIdx {
						l.F = v
					} else if j == keyIdx {
						l.Keyword = v
					} else if j == char1Idx {
						l.Char1 = getInt(v, "mission char1")
					} else if j == char2Idx {
						l.Char2 = getInt(v, "mission char2")
					} else if j == slineIdx {
						l.SLine = getInt(v, "mission sline")
					} else if j == alineIdx {
						l.ALine = getInt(v, "mission aline")
					} else if j == blineIdx {
						l.BLine = getInt(v, "mission bline")
					} else if j == adjustIdx {
						l.Adjust = getInt(v, "mission adjust")
					} else if j == warmIdx {
						l.Warm = getFloat64(v, "mission warm")
					} else if j == formalIdx {
						l.Formal = getFloat64(v, "mission formal")
					} else if j == tightIdx {
						l.Tight = getFloat64(v, "mission tight")
					} else if j == brightIdx {
						l.Bright = getFloat64(v, "mission bright")
					} else if j == darkIdx {
						l.Dark = getFloat64(v, "mission dark")
					} else if j == cuteIdx {
						l.Cute = getFloat64(v, "mission cute")
					} else if j == manIdx {
						l.Man = getFloat64(v, "mission man")
					} else if j == toughIdx {
						l.Tough = getFloat64(v, "mission tough")
					} else if j == nobleIdx {
						l.Noble = getFloat64(v, "mission noble")
					} else if j == strangeIdx {
						l.Strange = getFloat64(v, "mission strange")
					} else if j == sexyIdx {
						l.Sexy = getFloat64(v, "mission sexy")
					} else if j == sportIdx {
						l.Sport = getFloat64(v, "mission sport")
					}
				}
				res = append(res, l)
			}
		}
	}

	return
}

// func GetLevelStage(levelID string) string {
// 	if l := getLevelByLevelID(levelID); l != nil && l.TheaterID <= 0 {
// 		if l.No <= 17 {
// 			return constants.STAGE_CZ
// 		} else if l.No <= 23 {
// 			return constants.STAGE_GZ
// 		} else if l.No <= 37 {
// 			return constants.STAGE_DX
// 		} else if l.No <= 50 {
// 			return constants.STAGE_SH
// 		} else if l.No <= 62 {
// 			return constants.STAGE_5P
// 		}
// 	}
// 	return ""
// }

func getInt(v string, info string) int {
	if r, err := strconv.Atoi(v); err != nil {
		Err(err, v, info)
		return 0
	} else {
		return r
	}
}

func getFloat64(v string, info string) float64 {
	if r, err := strconv.ParseFloat(v, 64); err != nil {
		Err(err, v, info)
		return 0
	} else {
		return r
	}
}

func drop(cloID string, cpMulti float64, isTheater bool) (bool, int) {
	clo := clothes.GetClothesById(cloID)
	if clo == nil {
		return false, 0
	}

	randValue := rand.Intn(100)

	if chance := clothes.CloPieceChance(cloID); chance > 0 {
		if randValue <= int(float64(chance-1)*cpMulti) {
			return true, clo.Star
		} else {
			return false, clo.Star
		}
	} else {
		if clo.Star == 1 && randValue <= int(49.0*cpMulti) {
			return true, clo.Star
		} else if clo.Star == 2 && randValue <= int(39.0*cpMulti) {
			return true, clo.Star
		} else if clo.Star == 3 && randValue <= int(29.0*cpMulti) {
			return true, clo.Star
		} else if clo.Star == 4 && randValue <= int(19.0*cpMulti) {
			return true, clo.Star
		} else if clo.Star == 5 && randValue <= int(14.0*cpMulti) {
			return true, clo.Star
		} else {
			return false, clo.Star
		}
	}

}

func GetRewardPieces(lid string, curPieces map[string]int, cpMulti float64, isTheater bool) map[string]int {
	res := make(map[string]int, 0)
	level := LevelMapP[lid]

	if level != nil {
		drop1 := level.ClothesDrop1
		drop2 := level.ClothesDrop2

		for i := 0; i < 2; i++ {
			cid := drop1
			if i == 1 {
				cid = drop2
			}

			d, _ := drop(cid, cpMulti, isTheater)

			if cid != "" && d {
				p1Name := fmt.Sprintf("%s:%d", cid, 1)
				p2Name := fmt.Sprintf("%s:%d", cid, 2)
				p3Name := fmt.Sprintf("%s:%d", cid, 3)
				p4Name := fmt.Sprintf("%s:%d", cid, 4)

				res[p1Name] = 0
				res[p2Name] = 0
				res[p3Name] = 0
				res[p4Name] = 0

				p1cnt := curPieces[p1Name]
				p2cnt := curPieces[p2Name]
				p3cnt := curPieces[p3Name]
				p4cnt := curPieces[p4Name]

				pieceNum := 1
				if isTheater {
					pieceNum = 4
				}
				// if s == 3 {
				// 	pieceNum = 10
				// } else if s == 4 {
				// 	pieceNum = 16
				// } else if s == 5 {
				// 	pieceNum = 25
				// }
				for j := 0; j < pieceNum; j++ {
					p := common.GetOneCloPieceReward(p1cnt, p2cnt, p3cnt, p4cnt)
					switch p {
					case 1:
						p1cnt++
						res[p1Name]++
					case 2:
						p2cnt++
						res[p2Name]++
					case 3:
						p3cnt++
						res[p3Name]++
					case 4:
						p4cnt++
						res[p4Name]++
					default:
					}
				}
			}
		}

	}

	return res
}

func GetLevelRewardByRank(rank string) (gold int) {
	// level := getLevelByLevelID(levelID)

	// if level != nil {
	if rank == "f" {
		gold = constants.RECORD_F_GOLD + int((rand.Float32()*0.2-0.1)*float32(constants.RECORD_F_GOLD))
	} else if rank == "b" {
		gold = constants.RECORD_B_GOLD + int((rand.Float32()*0.2-0.1)*float32(constants.RECORD_B_GOLD))
	} else if rank == "a" {
		gold = constants.RECORD_A_GOLD + int((rand.Float32()*0.2-0.1)*float32(constants.RECORD_A_GOLD))
	} else if rank == "s" {
		gold = constants.RECORD_S_GOLD + int((rand.Float32()*0.2-0.1)*float32(constants.RECORD_S_GOLD))
	} else {
		gold = 0
	}
	// } else {
	// fmt.Println("没有找到关卡:", levelID)
	// gold = 0
	// }

	return
}

func GetBranch(levelID string) int {
	if l := getLevelByLevelID(levelID); l != nil {
		return l.Branch
	}
	return 0
}

func GetEndingClothesList(endingLevel string) []string {
	lv := getLevelByLevelID(endingLevel)

	// 获取关卡掉落的套装
	return clothes.GetSuitClothesList(lv.ClothesDrop1)
}

// func GetLevelNoByLevelID(lid string) string {
// 	if l := getLevelByLevelID(lid); l != nil {
// 		return strconv.Itoa(l.No)
// 	}
// 	return ""
// }

func GetLevelNo(lid string) int {
	if l := getLevelByLevelID(lid); l != nil {
		return l.No
	}
	return 0
}

func GetLevelMissionLevel(lid string) int {
	if l := getLevelByLevelID(lid); l != nil {
		return l.Level
	}
	return 0
}

func GetTheaterID(levelID string) int {
	if l := getLevelByLevelID(levelID); l != nil {
		return l.TheaterID
	}
	return 0
}

func getLevelByLevelID(levelID string) (level *GFLevelP) {
	if v, ok := LevelMapP[levelID]; ok {
		level = v
		return
	}

	return
}
