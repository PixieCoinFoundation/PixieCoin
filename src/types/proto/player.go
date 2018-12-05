package proto

const (
	PLAYER_MSG_CLASS uint16 = 2000
)

const (
	HTTP_PLAYER_DEL_USER string = "api/delUser"

	HTTP_PLAYER_DISCONNECT        string = "api/disconnect"
	HTTP_PLAYER_LOGIN             string = "api/loginGFServer"
	HTTP_PLAYER_HEART_BEAT        string = "api/heartBeat"
	HTTP_PLAYER_GET_ALL_USER_DATA string = "api/getAllUserData"
	HTTP_PLAYER_SET_PLAYER_HEAD   string = "api/setPlayerHead"
	HTTP_PLAYER_SET_NICKNAME      string = "api/setNickname"

	HTTP_PLAYER_RECOVER_HEART string = "api/recoverHeart"
	HTTP_PLAYER_USE_HEART     string = "api/useHeart"

	HTTP_PLAYER_BUY_CLOTHES       string = "api/buyClothes"
	HTTP_PLAYER_BATCH_BUY_CLOTHES string = "api/batchBuyClothes"
	HTTP_PLAYER_SELL_CLOTHES      string = "api/sellClothes"
	HTTP_PLAYER_YUYUE_GIFT        string = "api/yuyueGift"

	// cos
	HTTP_PLAYER_GET_OPEN_COS    string = "api/getOpenCos"
	HTTP_PLAYER_UPLOAD_COS_ITEM string = "api/uploadCosItem"
	HTTP_PLAYER_GET_COS_ITEMS   string = "api/getCosItems"
	HTTP_PLAYER_ADD_ITEM_SCORE  string = "api/addItemScore"
	HTTP_PLAYER_GET_COS_PARTY   string = "api/getCosParty"

	// custom
	HTTP_PLAYER_UPLOAD_CLOTHES       string = "api/uploadClothes"
	HTTP_PLAYER_GET_CUSTOMS          string = "api/getCustoms"
	HTTP_PLAYER_GET_CUSTOMS_SALE_BAN string = "api/getCustomsSaleBan" //非卖品列表
	HTTP_PLAYER_UP_CUSTOM            string = "api/upCustom"
	HTTP_PLAYER_BATCH_UP_CUSTOM      string = "api/batchUpCustom"
	HTTP_PLAYER_DOWN_CUSTOM          string = "api/downCustom"
	HTTP_PLAYER_COPY_REPORT          string = "api/copyReport"

	// iap
	HTTP_PLAYER_BUY_GOLD      string = "api/buyGold"
	HTTP_PLAYER_GET_ORDER     string = "api/getOrder"
	HTTP_PLAYER_CONFIRM_ORDER string = "api/confirmOrder"
	HTTP_PLAYER_BUY_DIAMOND   string = "api/buyDiamond"

	// mail
	HTTP_PLAYER_GET_MAIL_LIST_BY_PAGE string = "api/getMailListByPage"
	HTTP_PLAYER_GET_MAIL_GIFT         string = "api/getMailGift"
	HTTP_PLAYER_READ_MAIL             string = "api/readMail"
	HTTP_PLAYER_GET_UNREAD_MAIL_COUNT string = "api/getUnreadMailCount"
	HTTP_PLAYER_PROCESS_BATCH_MAIL    string = "api/processBatchMail"

	// record
	HTTP_PLAYER_SAVE_RECORD string = "api/saveRecord"

	// help
	HTTP_PLAYER_ADD_HELP         string = "api/addHelp"
	HTTP_PLAYER_GET_HELP         string = "api/getHelp"
	HTTP_PLAYER_ADD_HELP_COMMENT string = "api/addHelpComment"
	HTTP_PLAYER_GET_HELP_COMMENT string = "api/getHelpComment"

	HTTP_PLAYER_ADD_SCRIPT    string = "api/addScript"
	HTTP_PLAYER_DELETE_SCRIPT string = "api/deleteScript"

	HTTP_PLAYER_ADD_ENDING string = "api/addEnding"

	// pk
	HTTP_PLAYER_PK_MATCHING   string = "api/pkMatching"
	HTTP_PLAYER_ADD_PK_POINTS string = "api/addPKPoints"
	HTTP_PLAYER_ADD_PK_OP     string = "api/addPKOP"
	HTTP_PLAYER_PK_REWARD     string = "api/PKReward"
	HTTP_PLAYER_ADD_PK_COUNT  string = "api/addPKCount"
	HTTP_PLAYER_GET_PK_BONUS  string = "api/getPKBonus"
	HTTP_PLAYER_GET_PK_RANK   string = "api/getPKRank"

	// task
	HTTP_PLAYER_UPDATE_TASK     string = "api/updateTask"
	HTTP_PLAYER_GET_TASK_REWARD string = "api/getTaskReward"

	// item
	HTTP_PLAYER_ITEM_INIT    string = "api/itemInit"
	HTTP_PLAYER_BUY_ITEM     string = "api/buyItem"
	HTTP_PLAYER_USE_ITEM_102 string = "api/useItem102"
	HTTP_PLAYER_USE_ITEM_103 string = "api/useItem103"
	HTTP_PLAYER_USE_ITEM_104 string = "api/useItem104"
	HTTP_PLAYER_USE_ITEM_105 string = "api/useItem105"

	// cos comment
	HTTP_PLAYER_ADD_COS_COMMENT  string = "api/addCosComment"
	HTTP_PLAYER_GET_COS_COMMENT  string = "api/getCosComment"
	HTTP_PLAYER_GET_MY_COS_ITEMS string = "api/getMyCosItems"
	HTTP_PLAYER_GET_LEADER_COS   string = "api/getLeaderCos"
	HTTP_PLAYER_ADD_TOP_ITEM     string = "api/addTopItem"

	// 体力
	HTTP_PLAYER_RECOVER_TILI string = "api/recoverTili"
	HTTP_PLAYER_BUY_TILI     string = "api/buyTili"

	// 套装
	HTTP_PLAYER_SET_SUIT string = "api/setSuit"

	// 设置suit
	HTTP_PLAYER_GET_S_SUIT_STATUS string = "api/getSSuitStatus"
	HTTP_PLAYER_SET_S_SUIT_STATUS string = "api/setSSuitStatus"

	// 设计师部分
	HTTP_PLAYER_CUSTOM_NEW_SOLD_READ     string = "api/customNewSoldRead"
	HTTP_PLAYER_CUSTOM_GET_DESIGNER_INFO string = "api/customGetDesignerInfo"

	HTTP_PLAYER_CUSTOM_GET_DESIGNER_TOP_AMOUNT      string = "api/customGetDesignerTopAmount"
	HTTP_PLAYER_CUSTOM_GET_DESIGNER_TOP_SALE        string = "api/customGetDesignerTopSale"
	HTTP_PLAYER_CUSTOM_GET_DESIGNER_TOP_DESIGN_COIN string = "api/customGetDesignerTopDesignCoin"
	HTTP_PLAYER_CUSTOM_GET_DESIGNER_TOP_GOLD        string = "api/customGetDesignerTopGold"
	HTTP_PLAYER_CUSTOM_GET_DESIGNER_TOP_DIAMOND     string = "api/customGetDesignerTopDiamond"

	HTTP_PLAYER_CUSTOM_PULL_WALLET string = "api/customPullWallet"
	HTTP_PLAYER_USE_TILI           string = "api/useTili"
	HTTP_PLAYER_REPLAY_LEVEL       string = "api/replayLevel"
	HTTP_PLAYER_BUY_ALL_CLOTHES    string = "api/buyAllClothes"

	// 抽奖
	HTTP_PLAYER_GET_BONUS    string = "api/getBonus"
	HTTP_PLAYER_UNLOCK_LEVEL string = "api/unlockLevel"
	HTTP_PLAYER_SWITCH_LINE  string = "api/switchLine"

	// 设计师上传撤销
	HTTP_PLAYER_CUSTOM_CANCEL_UPLOAD string = "api/customCancelUpload"
	// 推送给玩家
	HTTP_PLAYER_PUSH_TO_PLAYER string = "api/pushToPlayer"
	// 获取分享奖励
	HTTP_PLAYER_GET_SHARE_REWARD string = "api/getShareReward"
	// 关注吐槽或者取消吐槽
	HTTP_PLAYER_ADD_HELP_FOLLOW string = "api/addHelpFollow"

	// 微博
	HTTP_PLAYER_SET_TWITTER string = "api/setTwitter"
	HTTP_PLAYER_GET_TWITTER string = "api/getTwitter"
	// 上传文件
	HTTP_PLAYER_UPLOAD_FILE string = "api/uploadFiles"

	// 添加联系人
	HTTP_PLAYER_ADD_CONTACT      string = "api/addContact"
	HTTP_PLAYER_SET_CONTACT_READ string = "api/setContactRead"

	// 标签功能
	HTTP_PLAYER_GET_TAG_LIST  string = "api/getTagList"
	HTTP_PLAYER_UNLOCK_TAG    string = "api/unlockTag"
	HTTP_PLAYER_SAVE_TAG_LIST string = "api/saveTagList"
	// 后台切回
	HTTP_PLAYER_RECONNECT string = "api/reconnect"

	//skip plot
	HTTP_PLAYER_SKIP_PLOT string = "api/skipPlot"

	//magic cash
	HTTP_PLAYER_MAGIC_CASH string = "api/magicCash"

	//friend
	HTTP_PLAYER_FRIEND_REQUEST string = "api/friendRequest"
	HTTP_PLAYER_ADD_FRIEND     string = "api/addFriend"
	HTTP_PLAYER_REMOVE_FRIEND  string = "api/removeFriend"
	HTTP_PLAYER_SEARCH_FRIEND  string = "api/searchFriend"
	HTTP_PLAYER_QUERY_WARDROBE string = "api/queryWardrobe"
	HTTP_PLAYER_LIST_FRIEND    string = "api/listFriend"

	//clothes piece swap
	HTTP_PLAYER_CLOTHES_PIECE_REQUEST string = "api/cloPieceRequest"
	HTTP_PLAYER_SWAP_CLOTHES_PIECE    string = "api/swapCloPiece"

	//notice
	HTTP_PLAYER_GET_NOTICE_LIST string = "api/getNoticeList"
	HTTP_PLAYER_DELETE_NOTICE   string = "api/deleteNotice"
	HTTP_PLAYER_GET_CLO_PIECE   string = "api/getCloPiece"

	//EVENT
	HTTP_PLAYER_GET_EVENTS string = "api/getEvents"

	//magic speed
	HTTP_PLAYER_MAGIC_SPEED string = "api/magicSpeed"

	//combine clothes
	HTTP_PLAYER_COMBINE_CLOTHES = "api/combineClothes"

	//planet unlock
	HTTP_PLAYER_UPLOCK_PLANET = "api/unlockPlanet"

	//event reward
	HTTP_PLAYER_GET_EVENT_REWARD = "api/getEventReward"
	HTTP_PLAYER_CHANGE_STAMP     = "api/changeStamp"

	//client log
	HTTP_PLAYER_ADD_CLIENT_LOG = "api/addClientLog"

	HTTP_PLAYER_GET_RECORD_HISTORY_CLOTHES = "api/getRecordHistoryClothes"

	//party
	HTTP_PLAYER_ADD_PARTY                  = "api/addParty"
	HTTP_PLAYER_GET_PARTY                  = "api/getParty"
	HTTP_PLAYER_ADD_PARTY_ITEM             = "api/addPartyItem"
	HTTP_PLAYER_GET_PARTY_ITEMS            = "api/getPartyItems"
	HTTP_PLAYER_RECOVER_GENERAL            = "api/recoverGeneral"
	HTTP_PLAYER_GET_RANDOM_PARTY_ITEM      = "api/getRandomPartyItem"
	HTTP_PLAYER_GET_RANDOM_PARTY_ITEM_PAIR = "api/getRandomPartyItemPair"
	HTTP_PLAYER_MARK_PARTY_ITEM            = "api/markPartyItem"
	HTTP_PLAYER_ADD_PARTY_ITEM_CMT         = "api/addPartyItemCmt"
	HTTP_PLAYER_GET_PARTY_ITEM_CMT         = "api/getPartyItemCmt"
	HTTP_PLAYER_GET_MARK_PARTY_ITEM_REWARD = "api/getMarkReward"
	HTTP_PLAYER_REFRESH_PARTY_ITEM         = "api/refreshPartyItem"
	HTTP_PLAYER_REPLY_PARTY                = "api/replyParty"
	HTTP_PLAYER_QUERY_PARTY_QINMI_USED     = "api/qPartyQinmi"

	//gift pack
	HTTP_PLAYER_GET_GIFT_PACK = "api/getGiftPack"

	HTTP_PLAYER_GET_EVENT_PROCESS = "api/getEventProgress"

	//board
	HTTP_PLAYER_ADD_BOARD_MSG = "api/addBoardMsg"
	HTTP_PLAYER_GET_BOARD_MSG = "api/getBoardMsg"
	HTTP_PLAYER_DEL_BOARD_MSG = "api/delBoardMsg"

	HTTP_PLAYER_SAVE_DEFAULT_CLOTHES  = "api/saveDefClo"
	HTTP_PLAYER_SET_BOARD_READ        = "api/setBoardRead"
	HTTP_PLAYER_GET_RANDOM_FRIEND     = "api/getRandomFriend"
	HTTP_PLAYER_PARTY_INVITE          = "api/inviteParty"
	HTTP_PLAYER_BUY_BACK_PK_CLOTHES   = "api/buyBackPKClo"
	HTTP_PLAYER_GET_ONE_CUSTOM        = "api/getOneCustom"
	HTTP_PLAYER_GET_BATCH_CUSTOM      = "api/getBatchCustom"
	HTTP_PLAYER_GET_MONTH_CARD_REWARD = "api/getMCReward"
	HTTP_PLAYER_INVITE_PLAYER         = "api/invitePlayer"

	HTTP_PLAYER_GRANT_CLOP = "api/grantClop"
	//AD
	HTTP_PLAYER_QUERY_AD_HANDLE = "api/queryAdvertisingHandle"
	HTTP_PLAYER_CONFIG_REWARD   = "api/adConfigReward"

	HTTP_PLAYER_JOIN_ADVERISINGPLACE    = "api/joinAdverisingPlace"
	HTTP_PLAYER_AUCTION_ADVERISINGPLACE = "api/auctionAdverisingPlace"
	HTTP_PLAYER_CONTRIBUTION_RANK       = "api/contributionRank"
	HTTP_PLAYER_GET_AUCTION_CUSTOMS     = "api/getAuctionCustoms"

	//DesignerZone

	HTTP_PLAYER_QUERY_DESIGNERZONE = "api/designerZone"
	HTTP_PLAYER_QUERY_TOPICDETAIL  = "api/topicDetail"
	HTTP_PLAYER_QUERY_ALL_TOPIC    = "api/allTopic"
	HTTP_PLAYER_QUERY_PUBLICBUY    = "api/publicBuy"

	//suit
	HTTP_PLAYER_ADD_SUIT    = "api/addSuit"
	HTTP_PLAYER_GET_SUIT    = "api/getSuit"
	HTTP_PLAYER_DELETE_SUIT = "api/deleteSuit"

	//HOPE
	HTTP_PLAYER_SEND_HOPE   = "api/sendHope"
	HTTP_PLAYER_GET_HOPE    = "api/getHope"
	HTTP_PLAYER_HELP_HOPE   = "api/helpHope"
	HTTP_PLAYER_FINISH_HOPE = "api/finishHope"

	HTTP_PLAYER_GET_YUYUE_GIFT = "api/getYuyueGift"
	HTTP_PLAYER_IAP_RETURN     = "api/iapReturn"

	HTTP_PALYER_INDEX_BANNERLIST = "api/bannerList"

	//new designer
	HTTP_PLAYER_FOLLOW_DESIGNER    = "api/followDesigner"
	HTTP_PLAYER_GET_DESIGN_MSG     = "api/getDesignMsg"
	HTTP_PLAYER_GET_DESIGN_MSG_NEW = "api/getDesignMsgNew"
	HTTP_PLAYER_LIST_FOLLOW        = "api/listFollow"
	HTTP_PLAYER_MARK_DESIGN_MSG    = "api/markDesignMsg"
	HTTP_PLAYER_CHANGE_DESIGN_DESC = "api/changeDesignDesc"

	HTTP_PLAYER_PRODUCT_CLOTHES   = "api/productClothes"
	HTTP_PLAYER_DECOMPOSE_CLOTHES = "api/decomposeClothes"
	HTTP_PLAYER_CHECK_TASK        = "api/checkTask"

	HTTP_PLAYER_LOGIN_RCLOUD = "api/loginRCloud"

	HTTP_PLAYER_CHANGE_CART = "api/changeCart"

	//guild
	HTTP_PLAYER_CREATE_GUILD                 = "api/createGuild"
	HTTP_PLAYER_GET_GUILD                    = "api/getGuild"
	HTTP_PLAYER_APPLY_GUILD                  = "api/applyGuild"
	HTTP_PLAYER_LIST_APPLY_GUILD             = "api/listGuildApply"
	HTTP_PLAYER_PASS_APPLY_GUILD             = "api/passGuildApply"
	HTTP_PLAYER_REJECT_APPLY_GUILD           = "api/rejectGuildApply"
	HTTP_PLAYER_LIST_GUILD                   = "api/listGuild"
	HTTP_PLAYER_PASS_ALL_APPLY               = "api/passAllApply"
	HTTP_PLAYER_REJECT_ALL_APPLY             = "api/rejectAllApply"
	HTTP_PLAYER_SEARCH_GUILD                 = "api/searchGuild"
	HTTP_PLAYER_KICK_GUILD_MEMBER            = "api/kickGuildMember"
	HTTP_PLAYER_DISBAND_GUILD                = "api/disbandGuild"
	HTTP_PLAYER_TRANSFER_GUILD               = "api/transferGuild"
	HTTP_PLAYER_QUIT_GUILD                   = "api/quitGuild"
	HTTP_PLAYER_LIST_GUILD_MEMBER            = "api/listGuildMember"
	HTTP_PLAYER_OFFER_GUILD_CLOTHES          = "api/offerGuildClothes"
	HTTP_PLAYER_LIST_GUILD_CLOTHES           = "api/listGuildClothes"
	HTTP_PLAYER_CONFIG_GUILD                 = "api/configGuild"
	HTTP_PLAYER_CANCEL_GUILD_CLOTHES         = "api/cancelGuildClothes"
	HTTP_PLAYER_REWARD_GUILD_CLOTHES         = "api/rewardGuildClothes"
	HTTP_PLAYER_MATCH_GUILD_WAR              = "api/matchGuildWar"
	HTTP_PLAYER_GUILD_WAR_INFO               = "api/guildWarInfo"
	HTTP_PLAYER_GUILD_WAR_SET_DEFEND_CLOTHES = "api/setDefendClothes"
	HTTP_PLAYER_QUERY_GUILD_MEMBER           = "api/queryGuildMember"
	HTTP_PLAYER_GUILD_LOTTERY                = "api/guildLottery"
	HTTP_PLAYER_GUILD_LOTTERY_STATUS         = "api/guildLotteryStatus"
	HTTP_PLAYER_GUILD_ATTACK_REWARD          = "api/attackReward"
	HTTP_PLAYER_GUILD_NOTIFY_DEFEND          = "api/notifyDefend"
	HTTP_PLAYER_ADD_GUILD_BOARD_MSG          = "api/addGuildBoardMsg"
	HTTP_PLAYER_GET_GUILD_BOARD_MSG          = "api/getGuildBoardMsg"
	HTTP_PLAYER_DEL_GUILD_BOARD_MSG          = "api/delGuildBoardMsg"
	HTTP_PLAYER_SET_CONFIG_CN                = "api/setConfigCN"
	HTTP_PLAYER_BUY_SCENE                    = "api/buyScene"
	HTTP_PLAYER_GET_THEATER_RANK             = "api/getTheaterRank"
	HTTP_PLAYER_BUY_THEATER                  = "api/buyTheater"
	HTTP_PLAYER_SEARCH_DESIGNER              = "api/searchDesigner"
	HTTP_PLAYER_CONFIG_GUILD_VP              = "api/configGuildVP"
	HTTP_PLAYER_STICK_CUSTOM                 = "api/stickCustom"
	HTTP_PLAYER_CANCEL_HOPE                  = "api/cancelHope"
	HTTP_PLAYER_BUY_DESIGN_COIN_PRICE_TOKEN  = "api/buyDCPriceToken"

	//KOR
	HTTP_PLAYER_GET_SUIT_REWARD   = "api/getSuitReward"
	HTTP_PLAYER_ADD_INVITE        = "api/addInvite"
	HTTP_PLAYER_COUNT_INVITE      = "api/countInvite"
	HTTP_PLAYER_GET_INVITE_REWARD = "api/getInviteReward"
	HTTP_PLAYER_UNREGISTER        = "api/unregister"
	HTTP_PLAYER_SET_CONFIG        = "api/setConfig"

	//rt party
	HTTP_PLAYER_MATCH_RT_PARTY      = "api/matchRTParty"
	HTTP_PLAYER_SYNC_RT_PARTY       = "api/syncRTParty"
	HTTP_PLAYER_OFFER_SUBJECT       = "api/offerRTSubject"
	HTTP_PLAYER_VOTE_SUBJECT        = "api/voteRTSubject"
	HTTP_PLAYER_UPLOAD_DRESS        = "api/uploadRTDress"
	HTTP_PLAYER_VOTE_DRESS          = "api/voteRTDress"
	HTTP_PLAYER_CANCEL_MATCH        = "api/cancelMatchRTParty"
	HTTP_PLAYER_LEAVE_HOST          = "api/leaveRTPartyHost"
	HTTP_PLAYER_GET_RT_PARTY_REWARD = "api/getRTPartyReward"
	HTTP_PLAYER_SEND_RT_CHAT        = "api/sendRTChat"

	//eth
	HTTP_PLAYER_GET_ETH_PXC_BALANCE   = "api/getBalance"
	HTTP_PLAYER_LIST_ETHEREUM_LOG     = "api/listEthereumLog"
	HTTP_PLAYER_LIST_TRANSFER_OUT_LOG = "api/listTransferOutLog"
	HTTP_PLAYER_ESTIMATE_ETH_GAS      = "api/estimateGas"
	HTTP_PLAYER_TRANSFER_OUT_ETH      = "api/transferOutETH"
	HTTP_PLAYER_ETH_SEND_VERIFY       = "api/sendETHVerify"
	HTTP_PLAYER_SET_ETH_PAY_PASSWORD  = "api/setETHPayPwd"
)
