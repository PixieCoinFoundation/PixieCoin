package servers

import (
	// "appcfg"
	// "common"
	// "constants"
	// "encoding/json"
	// "fmt"
	// "gv_common"
	// "gv_db/designer"
	// "gv_handle/adminHandle"
	// "gv_handle/copyHandle"
	// "gv_handle/customHandle"
	// "gv_handle/designerZoneHandle"
	// "gv_handle/resetHandle"
	// "gv_handle/testHandle"
	// "gv_handle/verifyHandle"
	// . "logger"
	// "net/http"
	// "os"
	// "path/filepath"
	"runtime"
	// "service/files"
	// "strings"
	// "time"
	// "tools"
	// . "types"
)

func StartGVServer() {
	// 启动服务器

	runtime.GOMAXPROCS(runtime.NumCPU()) //running in multicore
}

// // 测试相关
// http.HandleFunc("/GFVServer/echo", testHandle.Echo)

// // 页面--------------------------------------------
// // 管理员页面
// // http.HandleFunc("/GFVServer/register", adminHandle.RegisterPage)
// http.HandleFunc("/GFVServer/login", adminHandle.LoginPage)

// // 评审相关
// http.HandleFunc("/GFVServer/verify", verifyHandle.EnterVerify)
// http.HandleFunc("/GFVServer/verifyOne", verifyHandle.VerifyOne)
// http.HandleFunc("/GFVServer/verifyOneMS", verifyHandle.VerifyOneMS)
// // http.HandleFunc("/GFVServer/doVerify", verifyHandle.DoVerify)
// // http.HandleFunc("/GFVServer/rejectVerify", verifyHandle.RejectVerify)
// http.HandleFunc("/GFVServer/verified", verifyHandle.Verified)
// http.HandleFunc("/GFVServer/rejected", verifyHandle.Rejected)
// http.HandleFunc("/GFVServer/down", verifyHandle.Down)

// http.HandleFunc("/GFVServer/qOne", verifyHandle.QOneHandle)
// http.HandleFunc("/GFVServer/doQOne", verifyHandle.DoQOneHandle)
// http.HandleFunc("/GFVServer/verifyList", verifyHandle.VerifyList)
// http.HandleFunc("/GFVServer/verifyListMS", verifyHandle.VerifyListMS)
// http.HandleFunc("/GFVServer/forceDownClo", copyHandle.ForceDownClothesHandle)
// http.HandleFunc("/GFVServer/getCustom", customHandle.StartProcessOneCustom)
// // http.HandleFunc("/GFVServer/getQueueCnt", customHandle.GetQueueCnt)
// http.HandleFunc("/GFVServer/getCustomList", customHandle.GetCustomList)
// http.HandleFunc("/GFVServer/submitVerify", customHandle.SubmitVerify)

// http.HandleFunc("/GFVServer/resetIcon", resetHandle.ResetIconHandle)
// http.HandleFunc("/GFVServer/doResetIcon", resetHandle.DoResetIconHandle)

// http.HandleFunc("/GFVServer/resetClo", resetHandle.ResetCloHandle)
// http.HandleFunc("/GFVServer/doResetClo", resetHandle.DoResetCloHandle)

// //新增查询设计师作品接口
// http.HandleFunc("/GFVServer/qDesignerWorks", resetHandle.QDesignerWorks)
// http.HandleFunc("/GFVServer/qDesignerClot", customHandle.QueryDesignerClot)

// http.HandleFunc("/GFVServer/banDesigner", copyHandle.BanDesignerPageHandle)
// http.HandleFunc("/GFVServer/doBanDesigner", copyHandle.DoBanDesignerHandle)
// http.HandleFunc("/GFVServer/doQueryBanDesigner", copyHandle.QueryBanDeisgnerHandle)
// http.HandleFunc("/GFVServer/doQueryDesignerBanHistory", copyHandle.QueryDeisgnerBanHistoryHandle)
// http.HandleFunc("/GFVServer/doUnbanDesigner", copyHandle.UnbanDesignerHandle)

// http.HandleFunc("/GFVServer/copyList", copyHandle.CopyListPage)
// http.HandleFunc("/GFVServer/shensuList", copyHandle.ShensuListPage)
// http.HandleFunc("/GFVServer/getCopyList", copyHandle.GetCopyList)
// http.HandleFunc("/GFVServer/processCopy", copyHandle.ProcessCopy)
// http.HandleFunc("/GFVServer/submitCopy", copyHandle.SubmitCopy)
// http.HandleFunc("/GFVServer/processShensu", copyHandle.ProcessShensu)
// http.HandleFunc("/GFVServer/submitShensu", copyHandle.SubmitShensu)
// http.HandleFunc("/GFVServer/hisCopy", copyHandle.HisCopyListPage)
// http.HandleFunc("/GFVServer/doQueryHisCopy", copyHandle.ListHisCopyHandle)
// http.HandleFunc("/GFVServer/doCheckCopy", copyHandle.CheckCopyHandle)

// http.HandleFunc("/GFVServer/sendCustom", customHandle.SendCustomHandle)

// //设计师社区
// http.HandleFunc("/GFVServer/bannerConf", designerZoneHandle.BannerConfig)
// http.HandleFunc("/GFVServer/recommendBanner", designerZoneHandle.RecommendBanner)

// http.HandleFunc("/GFVServer/forRecommendPool", designerZoneHandle.ForRecommendPool)
// http.HandleFunc("/GFVServer/forRecommendDesignerPool", designerZoneHandle.ForRecommendDesignerPool)
// http.HandleFunc("/GFVServer/newForRecommendDesignerPool", designerZoneHandle.NewForRecommendDesignerPool)
// http.HandleFunc("/GFVServer/newForRecommendPool", designerZoneHandle.NewForRecommendPool)

// http.HandleFunc("/GFVServer/recommendPool", designerZoneHandle.RecommendPool)
// http.HandleFunc("/GFVServer/recommendDesignerPool", designerZoneHandle.RecommendDesignerPool)
// http.HandleFunc("/GFVServer/newRecommendDesignerPool", designerZoneHandle.NewRecommendDesignerPool)
// http.HandleFunc("/GFVServer/newRecommendPool", designerZoneHandle.NewRecommendPool)

// //获取设计师社区主页banner信息
// http.HandleFunc("/GFVServer/pageInfo", designerZoneHandle.GetPageInfo)
// //增加设计师主页banner 信息
// http.HandleFunc("/GFVServer/addIndexBanner", designerZoneHandle.AddIndexBanner)
// //修改设计师主页banner 信息
// http.HandleFunc("/GFVServer/updateIndexBanner", designerZoneHandle.UpdateIndexBanner)
// //删除设计师主页banner 信息
// http.HandleFunc("/GFVServer/deleteIndexBanner", designerZoneHandle.DeleteIndexBanner)

// http.HandleFunc("/GFVServer/stickIndexBanner", designerZoneHandle.StickIndexBanner)

// http.HandleFunc("/GFVServer/recommendList", designerZoneHandle.TopicRecommend)
// http.HandleFunc("/GFVServer/topicInfo", designerZoneHandle.GetTopicInfo)
// http.HandleFunc("/GFVServer/addTopic", designerZoneHandle.InsertTopicAjax)
// http.HandleFunc("/GFVServer/editTopic", designerZoneHandle.UpdateTopicAjax)

// http.HandleFunc("/GFVServer/stickTopic", designerZoneHandle.UpdateStickTopicAjax)
// http.HandleFunc("/GFVServer/delTopic", designerZoneHandle.DeleteTopicAjax)

// //获取设计师池的页面信息
// http.HandleFunc("/GFVServer/clothesList", designerZoneHandle.GetClothesPoolPageAjax)
// //
// http.HandleFunc("/GFVServer/designerList", designerZoneHandle.GetDesignerPoolPageAjax)

// http.HandleFunc("/GFVServer/clothesListCondition", designerZoneHandle.GetClothesPoolPageCondition)
// http.HandleFunc("/GFVServer/designerListCondition", designerZoneHandle.GetDesignerPoolPageCondition)

// http.HandleFunc("/GFVServer/addDesignerToPool", designerZoneHandle.AddDesignerToPoolajax)
// // http.HandleFunc("/GFVServer/addDesignerToForPoll", designerZoneHandle.AddDesignerToForPoolajax)

// http.HandleFunc("/GFVServer/delDesignerFromPool", designerZoneHandle.DeleteDesignerFromPoolajax)
// http.HandleFunc("/GFVServer/delDesignerFromForPool", designerZoneHandle.DeleteDesignerFromForPoolajax)

// http.HandleFunc("/GFVServer/addClothesToPool", designerZoneHandle.AddClothesToPoolajax)
// http.HandleFunc("/GFVServer/addClothesToForPool", designerZoneHandle.AddClothesToForPoolajax)
// http.HandleFunc("/GFVServer/stickClothes", designerZoneHandle.StickClothes)
// http.HandleFunc("/GFVServer/stickDesigner", designerZoneHandle.StickDesigner)

// http.HandleFunc("/GFVServer/deleteClothesFromForPool", designerZoneHandle.DelClothesFromForPoolajax)
// http.HandleFunc("/GFVServer/deleteClothesFromPool", designerZoneHandle.DelClothesFromPoolajax)

// http.HandleFunc("/GFVServer/allDesignerList", designerZoneHandle.GetAllDseignerList)
// http.HandleFunc("/GFVServer/allDesignerListAjax", designerZoneHandle.GetDesingerListAjax)
// http.HandleFunc("/GFVServer/addDesignerToForPool", designerZoneHandle.AddDesignerToForPool)

// http.HandleFunc("/GFVServer/forTopicPool", designerZoneHandle.ForTopicPage)
// http.HandleFunc("/GFVServer/topicPool", designerZoneHandle.TopicPage)

// http.HandleFunc("/GFVServer/addForTopicPool", designerZoneHandle.AddClothesToTopicPoolAjax)

// http.HandleFunc("/GFVServer/deleteTopicPool", designerZoneHandle.DeleteClothesToTopicPoolAjax)
// http.HandleFunc("/GFVServer/topicClothesPage", designerZoneHandle.GetTopicAjaxPage)
// http.HandleFunc("/GFVServer/addClothesToTopic", designerZoneHandle.UpdateTopicToPoolAjax)
// // 逻辑--------------------------------------------
// // 管理员相关
// // http.HandleFunc("/GFVServer/doRegister", adminHandle.Register)
// http.HandleFunc("/GFVServer/doLogin", adminHandle.Login)

// //	图片相关
// http.Handle("/GFVServer/pic/", http.StripPrefix("/GFVServer/pic/", http.FileServer(http.Dir("./downloads/"))))

// // 静态资源
// http.Handle("/GFVServer/web_res/", http.StripPrefix("/GFVServer/web_res/", http.FileServer(http.Dir("./web_res/"))))

// addr := appcfg.GetString("bind_ip", "127.0.0.1") + appcfg.GetString("port", ":29960")
// go refreshTopDesignerLoop()

// fmt.Println("Server start listen at:" + addr)
// err := http.ListenAndServe(addr, nil)
// if err != nil {
// 	Err("Server start up error:", err)
// 	return
// }
// }

// func refreshTopDesignerLoop() {
// 	if appcfg.GetLanguage() != "" {
// 		Info("not cn.do not refresh hot designer")
// 		return
// 	}

// 	fa := appcfg.GetString("file_address_prefix", "")
// 	hdls := appcfg.GetString("hot_designer_list", "")
// 	hdl := make([]string, 0)
// 	hdp := appcfg.GetString("hot_designer_platform", "")

// 	if hdp == "" {
// 		Info("no hot designer platform.do not go to refresh loop")
// 	}

// 	if hdls != "" {
// 		if err := json.Unmarshal([]byte(hdls), &hdl); err != nil {
// 			panic(err)
// 		}
// 	} else {
// 		Info("no hot_designer_list.do not go to refresh loop")
// 		return
// 	}

// 	refreshTopDesigner(fa, hdl, hdp)

// 	for {
// 		time.Sleep(3 * time.Hour)

// 		hour := time.Now().Hour()

// 		if hour >= 0 && hour <= 8 {
// 			refreshTopDesigner(fa, hdl, hdp)
// 		} else {
// 			Info("not time for hot designer refresh.")
// 		}
// 	}
// }

// func refreshTopDesigner(fa string, hdl []string, hdp string) {
// 	// fa := appcfg.GetString("file_address_prefix", "")
// 	gv_common.ProcessLock()
// 	defer gv_common.ProcessUnlock()

// 	if appcfg.GetLanguage() != "" {
// 		Info("not cn.do not refresh hot designer")
// 		return
// 	}
// 	if fa == "" {
// 		Info("no file address config.do not refresh hot designer")
// 		return
// 	}

// 	// hdl := params.GetHotDesignerList()
// 	hdl = tools.SliceOutOfOrder(hdl)

// 	if len(hdl) >= constants.CUSTOM_SHARE_DESIGNER_SIZE {
// 		resp := GetHotCustomResp{
// 			Designers: make([]HotDesigner, 0),
// 			Customs:   make([]HotClothes, 0),
// 		}

// 		for _, du := range hdl {
// 			var hd HotDesigner

// 			if len(resp.Designers) < constants.CUSTOM_SHARE_DESIGNER_SIZE {
// 				if un, nn, fc := designer.GetDesignerFromDB(hdp); un != "" {
// 					hd.Name = nn
// 					hd.FanCount = fc
// 				} else {
// 					continue
// 				}
// 			}

// 			hcl := make([]HotClothes, 0)
// 			if cl, head, err := designer.GetHotCustomsFromDB(hdp); err != nil || len(cl) <= 0 {
// 				Info("err or result 0 for hot designer:", du)
// 				continue
// 			} else {
// 				if hd.Head == "" && len(resp.Designers) < constants.CUSTOM_SHARE_DESIGNER_SIZE {
// 					hd.Head = fa + head
// 				}

// 				if len(resp.Designers) >= 2 {
// 					cl = cl[0:1]
// 				}

// 				for _, c := range cl {
// 					if c.CID <= 0 {
// 						continue
// 					}

// 					var hc HotClothes

// 					obName := fmt.Sprintf("%s/hot_%s_%d", files.GetPlatformOSSDirName(), du, c.CID)
// 					if c.CloSet.Main != "" {
// 						if err := decreDownload(c.CloSet.Main, obName, hdp); err != nil {
// 							continue
// 						}
// 					} else {
// 						Info("empty main_g:", du, c.CID)
// 						continue
// 					}

// 					hc.Icon = fa + obName

// 					hc.Cnt = c.BuyCount

// 					hcl = append(hcl, hc)
// 				}
// 			}

// 			if len(resp.Designers) < constants.CUSTOM_SHARE_DESIGNER_SIZE {
// 				hd.HotList = hcl
// 				resp.Designers = append(resp.Designers, hd)
// 			}

// 			if len(resp.Customs) < constants.CUSTOM_SHARE_SIZE && len(hcl) > 0 {
// 				resp.Customs = append(resp.Customs, hcl[0])
// 			}
// 		}

// 		if appcfg.GetBool("local_test", false) {
// 			common.Execute("cp", "-rf", appcfg.GetProjectPathPrefix()+"GFTP/clothes/.", appcfg.GetProjectPathPrefix()+"GFTP/output/")
// 		} else {
// 			if res, err := common.Execute("/bin/bash", appcfg.GetProjectPathPrefix()+"process_file.sh"); err == nil {
// 				Info(res)
// 			} else {
// 				Err(res)
// 			}
// 		}

// 		//upload
// 		if err := filepath.Walk("GFTP/output/", func(path string, info os.FileInfo, err error) error {
// 			if err != nil {
// 				return err
// 			}
// 			if info.IsDir() {
// 				return nil
// 			}

// 			Info(info.Name())
// 			names := strings.Split(info.Name(), "/")
// 			length := len(names)

// 			if length > 0 {
// 				if _, err = files.UploadFileByName("GFTP/output/"+info.Name(), files.GetPlatformOSSDirName()+"/"+names[length-1], 7); err != nil {
// 					return err
// 				}
// 			}

// 			return err
// 		}); err != nil {
// 			return
// 		}

// 		//refresh hot resp in redis at last
// 		if len(resp.Designers) == 2 && len(resp.Customs) > 0 {
// 			designer.SetHotResp(&resp)
// 		}
// 	} else {
// 		Info("not enough hot designer list")
// 	}
// }

// func decreDownload(objName, newObjName, platform string) (err error) {
// 	sfn := newObjName
// 	if strings.Contains(sfn, "/") {
// 		sfns := strings.Split(sfn, "/")
// 		if len(sfns) == 2 {
// 			sfn = sfns[1]
// 		}
// 	}
// 	saveFile := "./GFTP/clothes/" + sfn

// 	if err = files.DecryptDownloadFile(objName, saveFile); err != nil {
// 		return
// 	}
// 	return
// }
