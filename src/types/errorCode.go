package types

type ErrorCode int

const (
	NO_ERROR ErrorCode = 0

	ERR_UNKNOW     ErrorCode = -1
	ERR_NO_NETWORK ErrorCode = -2
	ERR_MAINTAIN   ErrorCode = -3
	ERR_PERSIST    ErrorCode = -4

	ERR_CHANNEL_LIMIT          ErrorCode = -5
	ERR_DELETING_USER          ErrorCode = -6
	ERR_DELED_USER             ErrorCode = -7
	ERR_LOGIN_GAME_PARAM_EMPTY ErrorCode = -8
	ERR_GF_TOKEN_WRONG         ErrorCode = -9
	ERR_CHANNEL_NOT_SUPPORT    ErrorCode = -10

	ERR_REG_USERNAME_EXISTS   ErrorCode = 1000
	ERR_REG_QUERY_USER_FAILED ErrorCode = 1002
	ERR_REG_FORMAT_WRONG      ErrorCode = 1003

	ERR_LOGIN_FIELD_WRONG             ErrorCode = 1001
	ERR_LOGIN_FORMAT_WRONG            ErrorCode = 1004
	ERR_LOGIN_GET_ACCESS_TOKEN_FAILED ErrorCode = 1005
	ERR_LOGIN_NO_PLAYER_FOUND         ErrorCode = 1006

	// 网络相关，和其它http服务器通信时
	// ERR_NET_BAD_REQUEST   ErrorCode = 2001
	// ERR_NET_SEND_FAILED   ErrorCode = 2002
	// ERR_NET_READ_FAILED   ErrorCode = 2003
	// ERR_CONNECT_REDIS_ERR ErrorCode = 2004 // redis错误

	// 游戏服务器相关LUA,
	// ERR_LUA_WRONG_PARAMS       ErrorCode = 3000
	// ERR_LUA_CALL_FUNC          ErrorCode = 3001
	// ERR_LUA_ACCESS_TOKEN_WRONG ErrorCode = 3002

	// token过期
	ERR_SHARE_FILE_LENGTH              ErrorCode = 3000
	ERR_SHARE_ADDRESS_NOT_CONFIG       ErrorCode = 3001
	ERR_ACCESS_TOKEN_WRONG             ErrorCode = 3002
	ERR_SESSION_TIME_OUT               ErrorCode = 3003
	ERR_FORCE_LOG_OUT                  ErrorCode = 3004
	ERR_SHARE_OUTPUT_FILE_LENGTH_WRONG ErrorCode = 3009

	ERR_NO_HANDLE_FOUND     ErrorCode = 4000
	ERR_HANDLE_PARAMS_WRONG ErrorCode = 4001
	ERR_WRONG_PARAMS        ErrorCode = 4002
	ERR_HANDLE_TIME_OUT     ErrorCode = 4003

	// 服装相关
	ERR_GOLD_NOT_ENOUGH               ErrorCode = 41000 // 金币不足
	ERR_CLOTHES_GOT                   ErrorCode = 41001 // 衣服已经有了
	ERR_PAY_GOLD_FAILED               ErrorCode = 41002 // 扣钱失败
	ERR_ADD_CLOTHES_FAILED_GOLD       ErrorCode = 41003 // 金币添加衣服失败
	ERR_ADD_CLOTHES_FAILED_NO_GOLD    ErrorCode = 41004 // 金币添加衣服失败
	ERR_ADD_CLOTHES_FAILED_DIAMOND    ErrorCode = 41005 // 钻石添加衣服失败
	ERR_ADD_CLOTHES_FAILED_NO_DIAMOND ErrorCode = 41006 // 钻石添加衣服失败
	ERR_PAY_DIAMOND_FAILED            ErrorCode = 41007 // 扣除钻石失败
	ERR_DIAMOND_NOT_ENOUGH            ErrorCode = 41008 // 钻石不够
	ERR_NOT_FOR_SALE                  ErrorCode = 41009 // 非卖品
	ERR_NO_SUCH_CLOTHES               ErrorCode = 41010 // 没有这个衣服
	ERR_BUY_ALL_CLOTHES               ErrorCode = 41011 // 购买所有衣服失败

	ERR_GET_BONUS_GOLD_NOT_ENOUGH    ErrorCode = 41012 // 摇奖金币不够
	ERR_GET_BONUS_DIAMOND_NOT_ENOUGH ErrorCode = 41013 // 摇奖钻石不够
	ERR_GET_BONUS                    ErrorCode = 41014 // 摇奖失败
	ERR_GET_BONUS_CHANGE_MAGIC       ErrorCode = 41015
	ERR_SELL_CLOTHES                 ErrorCode = 41016
	ERR_UNKNOWN_PRICE_TYPE           ErrorCode = 41017
	ERR_CURRENCY_NOT_ENOUGH          ErrorCode = 41018 //货币不足
	ERR_BUY_CLOTHES                  ErrorCode = 41019
	ERR_CLOTH_PRICE_CHANGE           ErrorCode = 41020
	ERR_FORMULA_NOT_EXIST            ErrorCode = 41021
	ERR_FORMULA_NOT_ENOUGH           ErrorCode = 41022
	ERR_DECOMPOSR_CLOTHES            ErrorCode = 41023
	ERR_CUSTOM_SIZE_LIMIT            ErrorCode = 41024
	ERR_SUIT_ALREADY_REWARDED        ErrorCode = 41025
	ERR_NOT_HAVE_SUIT                ErrorCode = 41026
	ERR_BATCH_UP_EMPTY               ErrorCode = 41027
	ERR_CUSTOM_CANCEL_STATUS_WRONG   ErrorCode = 41028
	ERR_NOT_IN_SHOP                  ErrorCode = 41029
	ERR_LOTTERY_NOTTIME              ErrorCode = 41030
	ERR_LOTTERY_EMPTY                ErrorCode = 41031
	ERR_PLAYER_NOT_IN_LIST           ErrorCode = 41032 //设计师不在白名单
	ERR_CUSTOM_STATUS_WRONG          ErrorCode = 41033
	ERR_CLOTHES_LIST_EMPTY           ErrorCode = 41034

	// 记录相关
	ERR_SAVE_RECORD                     ErrorCode = 42000 // 保存记录错误
	ERR_UNLOCK_LEVEL_DIAMOND_NOT_ENOUGH ErrorCode = 42001 // 解锁关卡钻石不足
	ERR_UNLOCK_LEVEL                    ErrorCode = 42002 // 解锁关卡失败
	ERR_SWITCH_LINE_DIAMOND_NOT_ENOUGH  ErrorCode = 42003 // 切换线路钻石不足
	ERR_SWITCH_LINE                     ErrorCode = 42004 // 切换线路失败
	ERR_UNLOCK_PLANET                   ErrorCode = 42005
	ERR_LEVEL_NOT_OPEN                  ErrorCode = 42006
	ERR_THEATER_DAY_LIMIT               ErrorCode = 42007
	ERR_LEVEL_NOT_EXIST                 ErrorCode = 42008
	ERR_SPECIAL_MISSION_DAY_LIMIT       ErrorCode = 42009
	ERR_THEATER_REPLAY10                ErrorCode = 42010
	ERR_SPECIAL_MISSION_REPLAY10        ErrorCode = 42011
	ERR_GET_THEATER_RANK                ErrorCode = 42012
	ERR_BUY_THEATER_LIMIT               ErrorCode = 42013
	ERR_THERTER_DAY_LIMIT_NOT_FULL      ErrorCode = 42014
	ERR_THEATER_NOT_OPEN                ErrorCode = 42015
	ERR_UPDATE_CUSTOM_STICK_TIME        ErrorCode = 42030

	// 玩家信息
	ERR_SET_NICK_NAME           ErrorCode = 43000
	ERR_GET_NICK_NAME           ErrorCode = 43001
	ERR_HEART_BEAT              ErrorCode = 43002
	ERR_SET_HEAD                ErrorCode = 43003
	ERR_RECOVER_HEART           ErrorCode = 43004
	ERR_USE_HEART               ErrorCode = 43005
	ERR_RECOVER_TILI            ErrorCode = 43006
	ERR_SET_NICKNAME_LENGTH     ErrorCode = 43007
	ERR_SUIT_LIMIT              ErrorCode = 43008
	ERR_LEVEL_NOT_ENOUGH        ErrorCode = 43009
	ERR_NOT_YY                  ErrorCode = 43010
	ERR_ALREADY_YY_GIFT         ErrorCode = 43011
	ERR_IAP_ALREADY_RETURN      ErrorCode = 43012
	ERR_NO_IAP_RETURN           ErrorCode = 43013
	ERR_MODEL_ILLEGAL           ErrorCode = 43014
	ERR_HAS_BACKGROUND          ErrorCode = 43015
	ERR_BACKGROUND_NOT_FOR_SALE ErrorCode = 43016
	ERR_BUY_BACKGROUND          ErrorCode = 43017
	ERR_NOT_HAVE_BACKGROUND     ErrorCode = 43018

	//magic
	ERR_MAGIC_CASH_MAGIC_NOT_ENOUGH = 44000

	// 邮件相关
	ERR_MAIL_WRONG_PARAMS    ErrorCode = 50000 // 错误的参数
	ERR_MAIL_FIND_FAILED     ErrorCode = 50001 // 查询邮件错误
	ERR_MAIL_FIND_ONE_FAILED ErrorCode = 50002 // 查询单条邮件错误
	ERR_MAIL_ACCEPT_ALREADY  ErrorCode = 50003 // 邮件礼物已经接收过了
	ERR_MAIL_ACCEPT_ERROR    ErrorCode = 50004 // 接受邮件礼物时，出错，需要重试
	ERR_DELETE_MAIL          ErrorCode = 50005 // 删除邮件错误
	ERR_READ_MAIL            ErrorCode = 50006 // 设置邮件已读错误
	ERR_GET_UNREAD_COUNT     ErrorCode = 50007 // 设置邮件已读错误

	// 上传衣服
	ERR_UPLOAD_FAILED      ErrorCode = 60000 // 上传失败
	ERR_FIND_CUSTOM_FAILED ErrorCode = 60001 // 查询失败
	ERR_SET_CUSTOM_NEWSOLD ErrorCode = 60002 // 设置设计师服装无新卖出
	ERR_GET_DESIGNER_INFO  ErrorCode = 60003 // 获取设计师信息失败

	ERR_DESIGNER_PULL_WALLET                    ErrorCode = 60004 // 取出设计师钱包里的钱失败
	ERR_CANCEL_CUSTOM                           ErrorCode = 60005 // 玩家取消设计失败
	ERR_UPLOAD_CUSTOM_DAY_LIMIT                 ErrorCode = 60006
	ERR_UPLOAD_DAY_LIMIT                        ErrorCode = 60007
	ERR_UP_CUSTOM                               ErrorCode = 60008
	ERR_DOWN_CUSTOM                             ErrorCode = 60009
	ERR_GET_ONE_CUSTOM                          ErrorCode = 60010
	ERR_WRONG_CLOTHES_ORIGIN                    ErrorCode = 60011
	ERR_ADD_COPY_REPORT                         ErrorCode = 60012
	ERR_DAY_COPY_REPORT_LIMIT                   ErrorCode = 60013
	ERR_DAY_SHENSU_LIMIT                        ErrorCode = 60014
	ERR_NOT_CUSTOM                              ErrorCode = 60015
	ERR_PLAYER_FIND_CUSTOM_FAILED               ErrorCode = 60016
	ERR_CUSTOM_TAG_COLLISION                    ErrorCode = 60017
	ERR_CUSTOM_INVENTORY_NOT_EMPTY              ErrorCode = 60018
	ERR_INVENTORY_EMPTY                         ErrorCode = 60019
	ERR_LOCK_CUSTOM                             ErrorCode = 60020
	ERR_CUSTOM_PRICE_CHANGE                     ErrorCode = 60021
	ERR_MONEY_TYPE_WRONG                        ErrorCode = 60022
	ERR_SEARCH_DESIGNER                         ErrorCode = 60023
	ERR_SEARCH_DESIGNER_USERNAME                ErrorCode = 60024
	ERR_DESIGN_COIN_PRICE_TOKEN_NOT_ENOUGH      ErrorCode = 60025
	ERR_BUY_DESIGN_COIN_PRICE_TOKEN_MONTH_LIMIT ErrorCode = 60026
	ERR_CART_SIZE_FULL                          ErrorCode = 60027
	ERR_ADD_CART                                ErrorCode = 60028
	ERR_DEL_CART                                ErrorCode = 60029
	ERR_CLOTHES_FORMAT_WRONG                    ErrorCode = 60030
	ERR_CART_DO_NOT_CONTAIN_CLOTHES             ErrorCode = 60031
	ERR_CART_DO_NOT_ENOUGH_CLOTHES              ErrorCode = 60032
	ERR_INCRE_CART                              ErrorCode = 60033
	ERR_DECRE_CART                              ErrorCode = 60034

	// 内购
	ERR_SAVE_ORDER                       ErrorCode = 70000 // 保存订单失败
	ERR_CREATE_ORDER                     ErrorCode = 70001 // 创建订单失败
	ERR_ORDER_NO_ITEM                    ErrorCode = 70002 // 没有id对应的内购项目
	ERR_ORDER_ID_WRONG                   ErrorCode = 70003 // 错误的order id
	ERR_CONFIRM_ORDER_ADD_DIAMOND_FAILED ErrorCode = 70004 // 添加钻石错误
	ERR_NO_MONEY_ITEM                    ErrorCode = 70005 // 没有id对应的金币购买条目
	ERR_IAP_DIAMOND_NOT_ENOUGH           ErrorCode = 70006 // 钻石不够
	ERR_IAP_BUY_MONEY_ERR                ErrorCode = 70007 // 购买金币失败
	ERR_INVALID_RECEIPT                  ErrorCode = 70008 // 错误的收据

	ERR_IAP_BUY_TILI_ERR    ErrorCode = 70009 // 购买体力错误
	ERR_NO_TILI_ITEM        ErrorCode = 70010 // 没有id对应的体力
	ERR_IAP_TILI_NOT_ENOUGH ErrorCode = 70011 // 没钱买体力

	ERR_ORDER_FAILED  ErrorCode = 70012 // 付款失败
	ERR_ORDER_REFUND  ErrorCode = 70013 // 已经退款
	ERR_ORDER_WAITING ErrorCode = 70014 // 等待中

	ERR_USE_TILI_NOT_ENOUGH ErrorCode = 70015 // 体力不够了
	ERR_BUY_TILI_DAY_LIMIT  ErrorCode = 70016

	ERR_ORDER_CLOSED                   ErrorCode = 70017
	ERR_MONTH_CARD_REWARD              ErrorCode = 70018
	ERR_TINY_PACK_NOT_EXIST            ErrorCode = 70019
	ERR_TINY_PACK_BUY_LIMIT            ErrorCode = 70020
	ERR_VIP_NOT_ENOUGH                 ErrorCode = 70021
	ERR_GEN_VIVO_ORDER                 ErrorCode = 70022
	ERR_CHECK_HUAWEI_ORDER             ErrorCode = 70023
	ERR_CHECK_KUAIKAN_ORDER            ErrorCode = 70024
	ERR_CHECK_VIVO_ORDER               ErrorCode = 70025
	ERR_CHECK_UC_ORDER                 ErrorCode = 70026
	ERR_CHECK_XIAOMI_ORDER             ErrorCode = 70027
	ERR_CHECK_SINA_ORDER               ErrorCode = 70028
	ERR_CHECK_4399_ORDER               ErrorCode = 70029
	ERR_CHECK_BUKA_ORDER               ErrorCode = 70030
	ERR_CHECK_OPPO_ORDER               ErrorCode = 70031
	ERR_CHECK_MEITU_ORDER              ErrorCode = 70032
	ERR_CHECK_PAPA_ORDER               ErrorCode = 70033
	ERR_CHECK_THIRD_ORDER              ErrorCode = 70034
	ERR_CHECK_MHR_ORDER                ErrorCode = 70035
	ERR_CHECK_360_ORDER                ErrorCode = 70036
	ERR_GEN_MGR_SIGN                   ErrorCode = 70037
	ERR_CHECK_BAIDU_ORDER              ErrorCode = 70038
	ERR_DUPLICATE_APPLE_RECEIPT_APPLE  ErrorCode = 70039
	ERR_DUPLICATE_APPLE_RECEIPT_GOOGLE ErrorCode = 70040

	PLATFORM_APPLE  = "apple"
	PLATFORM_GOOGLE = "google"
	// Cosplay
	ERR_GET_COSPLAY           ErrorCode = 80000 // 查询Cosplay信息失败
	ERR_UPLOAD_COSPLAY        ErrorCode = 80001 // 上传CospItem失败
	ERR_GET_COS_ITEM_LIST     ErrorCode = 80002 // 查询CosItemList失败
	ERR_COS_ADD_ITEM_SCORE    ErrorCode = 80003 // 添加分数失败
	ERR_ADD_COS_COMMENT       ErrorCode = 80004 // 添加评论失败
	ERR_GET_COS_COMMENT       ErrorCode = 80005 // 获取评论失败
	ERR_COS_CLOSED            ErrorCode = 80006 // 当前cos已经结束
	ERR_GET_MY_COS_ITEM_LIST  ErrorCode = 80007 // 查询我的cos失败
	ERR_ADD_TOP_COS_ITEM      ErrorCode = 80008
	ERR_GET_COS_ITEM_SCORE    ErrorCode = 80009
	ERR_COSPLAY_CANNOT_INVITE ErrorCode = 80010

	// 吐槽
	ERR_TUCAO_GET_INIT_DATA      ErrorCode = 90001 // 获取吐槽失败
	ERR_ADD_HELP                 ErrorCode = 100001
	ERR_GET_HELP                 ErrorCode = 100002
	ERR_ADD_HELP_COMMENT         ErrorCode = 100003
	ERR_GET_HELP_COMMENT         ErrorCode = 100004
	ERR_ADD_HELP_FOLLOW_TOO_MANY ErrorCode = 100005
	ERR_ADD_HELP_FOLLOW          ErrorCode = 100006
	ERR_DEL_HELP_FOLLOW          ErrorCode = 100007
	ERR_ADD_HELP_DAY_LIMIT       ErrorCode = 100008

	// script
	ERR_ADD_SCRIPT ErrorCode = 110001
	ERR_DEL_SCRIPT ErrorCode = 110002

	// ending
	ERR_ADD_ENDING ErrorCode = 110003

	// pk
	ERR_ADD_PK_POINTS   ErrorCode = 110004
	ERR_ADD_PKOP_FAILED ErrorCode = 110005
	ERR_PK_ERWARD       ErrorCode = 110006
	ERR_ADD_PK_COUNT    ErrorCode = 110007

	// task
	ERR_UPDATE_TASK     ErrorCode = 110008
	ERR_GET_TASK_REWARD ErrorCode = 110009

	// item
	ERR_UPDATE_ITEM               ErrorCode = 110010
	ERR_BUY_ITEM_MONEY_NOT_ENOUGH ErrorCode = 110011 // 购买道具钱不够
	ERR_BUY_ITEM                  ErrorCode = 110012 // 购买道具失败
	ERR_USE_ITEM_102              ErrorCode = 110013 // 使用道具102失败
	ERR_USE_ITEM_103              ErrorCode = 110014 // 使用道具103失败
	ERR_USE_ITEM_104              ErrorCode = 110015 // 使用道具104失败
	ERR_USE_ITEM_105              ErrorCode = 110016 // 使用道具105失败

	ERR_GET_PK_BONUS                  ErrorCode = 110017
	ERR_GET_S_SUIT_STATUS             ErrorCode = 110018 // 获取关卡的s过关搭配的购买状态失败
	ERR_SET_S_SUIT_STATUS             ErrorCode = 110019
	ERR_SET_S_SUIT_DIAMOND_NOT_ENOUGH ErrorCode = 110020 // 钻石不够

	ERR_SHARE_TOO_MANY ErrorCode = 110021
	ERR_SHARE          ErrorCode = 110022

	ERR_SET_TWEETS ErrorCode = 110023 // 保存tweet失败
	ERR_GET_TWEETS ErrorCode = 110024 // 获取tweet失败

	ERR_UPLOAD_FILE ErrorCode = 110025 // 上传文件失败

	ERR_UNLOCK_TAG_NOT_ENOUGH_DIAMOND ErrorCode = 110026 // 解锁标签没有钱
	ERR_UNLOCK_TAG                    ErrorCode = 110027 // 解锁标签失败

	// DESIGNER
	ERR_UPDATE_DESIGNER ErrorCode = 110028

	//err parse request
	ERR_PARSE_REQUEST_JSON ErrorCode = 110029

	ERR_GET_UID           = 110030
	ERR_BAN_USER          = 110031
	ERR_USER_LOCKED       = 110032
	ERR_GET_IPPORT        = 110033
	ERR_GET_STMT          = 110034
	ERR_UNBAN_USER        = 110035
	ERR_PLAYER_NOT_BANNED = 110036

	ERR_ADD_BOARD_MSG ErrorCode = 110037
	ERR_GET_BOARD_MSG ErrorCode = 110038
	ERR_DEL_BOARD_MSG ErrorCode = 110039

	ERR_NO_PK_LOSE_CLOTHES ErrorCode = 110040

	ERR_GET_PK_RANK            ErrorCode = 110041
	ERR_GET_PK_RANK_TYPE_WRONG ErrorCode = 110042

	ERR_PK_CLOSING               ErrorCode = 110043
	ERR_UNKNOWN_FOLLOW_TYPE      ErrorCode = 110044
	ERR_FOLLOW                   ErrorCode = 110045
	ERR_UNFOLLOW                 ErrorCode = 110046
	ERR_GET_DESIGN_MSG           ErrorCode = 110047
	ERR_GET_FENSI                ErrorCode = 110048
	ERR_GET_GUANZHU              ErrorCode = 110049
	ERR_UNKNOWN_GET_FOLLOW_TYPE  ErrorCode = 110050
	ERR_DIANZAN                  ErrorCode = 110051
	ERR_UNDIANZAN                ErrorCode = 110052
	ERR_PINGLUN                  ErrorCode = 110053
	ERR_UNKNOWN_MARK_DESIGN_TYPE ErrorCode = 110054
	ERR_CHANGE_DESIGNER_DESC     ErrorCode = 110055
	ERR_TASK_NOT_EXISTS          ErrorCode = 110056

	// BILIBILI
	ERR_BILIBILI_CHECK_RECEIPT      ErrorCode = 120001 // 查询B站订单出错
	ERR_BILIBILI_ORDER_NOT_FINISHED ErrorCode = 120002 // B站订单未完成

	//friend
	ERR_SEND_FRIEND_REQUEST         = 130000
	ERR_ALREADY_FRIEND              = 130001
	ERR_FRIEND_SIZE_FULL            = 130002
	ERR_ADD_FRIEND                  = 130003
	ERR_REMOVE_FRIEND               = 130004
	ERR_SEARCH_FRIEND               = 130005
	ERR_QUERY_CLOTHES               = 130006
	ERR_SEND_SWAP_CLO_PIECE_REQUEST = 130007
	ERR_SWAP_CLO_PIECE              = 130008
	ERR_PIECE_NOT_ENOUGH            = 130009
	ERR_SWAP_CLO_PIECE_LIMIT        = 130010
	ERR_SWAP_CLO_PIECE_COUNT_LIMIT  = 130011
	ERR_ADD_SELF                    = 130012
	ERR_NOT_FRIEND                  = 130013
	ERR_ALREADY_SEND                = 130014
	ERR_GET_RANDOM_FRIEND           = 130015
	ERR_ALREADY_INVITE              = 130016
	ERR_LIST_FRIEND                 = 130017
	ERR_DEST_FRIEND_SIZE_FULL       = 130018
	ERR_QUERY_KOR_USERNAME          = 130019
	ERR_GET_FRIEND_QINMI            = 130020

	//notice
	ERR_GET_NOTICE_LIST = 140000
	ERR_DELETE_NOTICE   = 140001
	ERR_GET_NOTICE      = 140002
	ERR_GET_CLO_PIECE   = 140003
	ERR_REJECT_NOTICE   = 140004

	//event
	ERR_GET_EVENTS             = 150000
	ERR_EVENT_NOT_OPEN         = 150001
	ERR_EVENT_NO_REWARD        = 150002
	ERR_EVENT_PROGRESS_ILLEGAL = 150003
	ERR_EVENT_GET_REWARD       = 150004
	ERR_EVENT_TIME_ILLEGAL     = 150005
	ERR_EVENT_REWARD_ALREADY   = 150006
	ERR_CHANGE_STAMP           = 150007

	//party
	ERR_GET_PARTYS                 = 160000
	ERR_ADD_PARTY                  = 160001
	ERR_GET_PARTY_ITEMS            = 160002
	ERR_PARTY_CLOSED               = 160003
	ERR_GET_PARTY                  = 160004
	ERR_ADD_PARTY_PLAYER           = 160005
	ERR_GET_PARTY_ITEM             = 160006
	ERR_PARTY_ALREADY_JOINED       = 160007
	ERR_PARTY_NO_JOIN_CNT          = 160008
	ERR_PARTY_NO_MARK_CNT          = 160009
	ERR_GET_RANDOM_PARTY_ITEM      = 160010
	ERR_MARK_PARTY_ITEM            = 160011
	ERR_ADD_PARTY_ITEM_CMT         = 160012
	ERR_GET_PARTY_ITEM_CMT         = 160013
	ERR_PARTY_SIZE_LIMIT           = 160014
	ERR_ALREADY_HOST_CASUAL_PARTY  = 160015
	ERR_MARK_CNT_NOT_ENOUGH        = 160016
	ERR_PARTY_FULL                 = 160017
	ERR_SEND_FLOWER_FAIL           = 160018
	ERR_PARTY_SINGLE               = 160019
	ERR_PARTY_ITEM_HAS_PARTNER     = 160020
	ERR_PARTY_ITEM_FORBID_PARTNER  = 160021
	ERR_REFRESH_PARTY_ITEM         = 160022
	ERR_INVITE_PARTNER             = 160023
	ERR_PARTNER_NOT_ALLOW          = 160024
	ERR_UNKNOWN_PARTNER_TYPE       = 160025
	ERR_UNKNOWN_PARTY_MONEY_TYPE   = 160026
	ERR_ADD_PARTY_ITEM             = 160027
	ERR_QUERY_PARTY_JOIN           = 160028
	ERR_INVITE_EMPTY               = 160029
	ERR_MARK_SELF                  = 160030
	ERR_LOCK_PARTY_ITEM            = 160031
	ERR_QUERY_PARTY_QINMI_USED     = 160032
	ERR_GET_RANDOM_PARTY_ITEM_PAIR = 160032
	ERR_MARK_PARTY_PARAM_WRONG     = 160034

	//gift pack
	ERR_GIFT_PACK_NOT_EXIST         = 170000
	ERR_GIFT_CODE_ILLEGAL           = 170001
	ERR_GET_GIFT                    = 170002
	ERR_GIFT_PACK_NOT_OPEN          = 170003
	ERR_GIFT_PACK_CHANNEL_NOT_MATCH = 170004

	//hope
	ERR_SEND_HOPE_TOO_MUCH      = 180000
	ERR_SEND_HOPE_REQUIRE_EMPTY = 180001
	ERR_SEND_HOPE               = 180002
	ERR_GET_HOPES               = 180003
	ERR_GET_HOPE                = 180004
	ERR_HOPE_STATUS_WRONG       = 180005
	ERR_CANCEL_HOPE             = 180006

	//guild
	ERR_CREATE_GUILD                           = 190000
	ERR_GET_GUILD                              = 190001
	ERR_SEND_APPLY                             = 190002
	ERR_LIST_APPLY                             = 190003
	ERR_PASS_APPLY                             = 190004
	ERR_REJECT_APPLY                           = 190005
	ERR_GET_RANDOM_GUILD                       = 190006
	ERR_PLAYER_ALREADY_IN_GUILD                = 190007
	ERR_NOT_IN_GUILD                           = 190008
	ERR_SEARCH_GUILD                           = 190009
	ERR_KICK_GUILD_MEMBER                      = 190010
	ERR_KICK_SELF                              = 190011
	ERR_DISBAND_GUILD                          = 190012
	ERR_TRANSFER_GUILD                         = 190013
	ERR_QUIT_GUILD                             = 190014
	ERR_LIST_MEMBER                            = 190015
	ERR_NOT_HAVE_CLOTHES                       = 190016
	ERR_ADD_CLOTHES                            = 190017
	ERR_LIST_CLOTHES                           = 190018
	ERR_NOT_OWN_GUILD                          = 190019
	ERR_UPDATE_MEMBER_WAR                      = 190020
	ERR_UPDATE_GUILD_DESC                      = 190021
	ERR_GUILD_OWNER                            = 190022
	ERR_GUILD_LOCK                             = 190023
	ERR_DEL_CLOTHES                            = 190024
	ERR_SHARE_REWARDED                         = 190025
	ERR_SHARE_NOT_ENOUGH                       = 190026
	ERR_GET_PLAYER_GUILD                       = 190027
	ERR_GUILD_WAR_TIME_WRONG                   = 190028
	ERR_WAR_MEMBER_NOT_ENOUGH                  = 190029
	ERR_MATCH_GUILD_WAR_LATER                  = 190030
	ERR_GET_WAR_LOG                            = 190031
	ERR_OPPONENT_WRONG                         = 190032
	ERR_SET_DEFEND_CLOTHES                     = 190033
	ERR_SET_DEFEND_HOUR_WRONG                  = 190034
	ERR_GET_DEFEND_SIZE                        = 190035
	ERR_ADD_GUILD_ACTIVITY                     = 190036
	ERR_ADD_PK_DATA                            = 190037
	ERR_NO_ATTACK_CNT                          = 190038
	ERR_GUILD_LOTTERY_INDEX                    = 190039
	ERR_QUIT_GUILD_CD                          = 190040
	ERR_GUILD_LOTTERY                          = 190041
	ERR_WAR_DATE_WRONG                         = 190042
	ERR_WIN_MEDAL_SIZE_WRONG                   = 190043
	ERR_SEASON_LAST_DAY                        = 190044
	ERR_GUILD_DISBAND                          = 190045
	ERR_MATCH_GUILD_WAR_PANIC                  = 190046
	ERR_QUERY_SHARE_IN_DEFEND                  = 190047
	ERR_GUILD_WAR_ON_GOING                     = 190048
	ERR_GUILD_CHANGED                          = 190049
	ERR_APPLY_TOO_QUICK                        = 190050
	ERR_KICK_WAR_MEMBER                        = 190051
	ERR_GUILD_SETTLE_SEASON_RUNNING            = 190052
	ERR_GUILD_SIZE_LIMIT                       = 190053
	ERR_ADD_GUILD_BOARD_MSG                    = 190054
	ERR_GET_GUILD_BOARD_MSG                    = 190055
	ERR_DEL_GUILD_BOARD_MSG                    = 190056
	ERR_CONFIG_GUILD_VICE_PRESIDENT            = 190057
	ERR_CONFIG_GUILD_VICE_PRESIDENT_TYPE_WRONG = 190058
	ERR_CONFIG_GUILD_VICE_PRESIDENT_EFFECT     = 190059
	ERR_CONFIG_GUILD_VICE_PRESIDENT_SIZE_FULL  = 190060
	ERR_CONFIG_SELF_GUILD_VICE_PRESIDENT       = 190061
	ERR_LIST_VICE_PRESIDENT                    = 190062

	ERR_ADVERTISING_MSG                = 200000
	ERR_ADVERTISING_JOIN_MSG           = 200001
	ERR_ADVERTISING_REDIS_MSG          = 200002
	ERR_ADVERTISING_PARAMAS_NOT_MATCH  = 200003
	ERR_ADVERTISING_CLOTHES_STATUS     = 200004
	ERR_ADVERTISING_CLOTHES_DB_ERR     = 200005
	ERR_ADVERTISING_REWARD_EXIST       = 200006
	ERR_AD_PLAYER_NOT_JOIN_NOT_AUCTION = 200007
	ERR_AD_REDIS_ERR                   = 200008
	ERR_AD_GETCUSTOMS_ERR              = 200008
	ERR_AD_CLOTHESID_EXIST             = 200009
	ERR_AD_ADD_CONTRIBUTION_ERR        = 200010
	ERR_ADVERTISING_AUCATION_END       = 200011
	ERR_AD_UPDATECONFIG_ERR            = 200012
	ERR_AD_AUTOCONFIG_ERR              = 200013
	ERR_AD_EXISTKEY_ERR                = 200014
	ERR_AD_ERR_PARAMS                  = 200015
	ERR_MONEY_TYPE                     = 200016
	ERR_CANNOT_DOWN_AUCATION_CLOTHES   = 200017
	ERR_CANNOT_DOWN_ONSALE_CLOTHES     = 200018
	ERR_ADD_DESIGN_GOLD                = 200019
	ERR_ADD_DESIGN_DIAM                = 200020
	ERR_GET_CONTRIBUTION_INFO          = 200021
	ERR_CLOTHESS_NOT_PASSED            = 200022
	ERR_NOT_PERRMIT_TO_CONFIG          = 200023

	//rt party
	ERR_MATCH_RT_PARTY         = 210000
	ERR_JOIN_RT_PARTY          = 210001
	ERR_SYNC_PARTY             = 210002
	ERR_CHECK_RT_PARTY_STAGE   = 210003
	ERR_ALREADY_OFFER_SUBJECT  = 210004
	ERR_OFFER_SUBJECT          = 210005
	ERR_ALREADY_VOTE_SUBJECT   = 210006
	ERR_VOTE_SUBJECT           = 210007
	ERR_UPLOAD_DRESS           = 210008
	ERR_ALREADY_VOTE_DRESS     = 210009
	ERR_VOTE_DRESS             = 210010
	ERR_RT_PARAM_WRONG         = 210011
	ERR_VOTE_SELF_DRESS        = 210012
	ERR_RT_PARTY_ALREADY_START = 210013
	ERR_NOT_ATTEND             = 210014
	ERR_CANCEL_MATCH           = 210015
	ERR_MATCH_RT_PARTY_TIMEOUT = 210016
	ERR_JOIN_RT_PARTY_FULL     = 210017
	ERR_MATCH_RT_PARTY_HOSTID  = 210018
	ERR_NOT_RT_PARTY_TIME      = 210019
	RT_RECONNECT               = 210020
	ERR_VOTE_SAME_SUBJECT      = 210021
	ERR_RT_PARTY_NOT_START     = 210022
	ERR_LEAVE_RT_PARTY_HOST    = 210023
	ERR_VOTE_NOT_FINISH_DRESS  = 210024
	ERR_ALREADY_REWARDED_PARTY = 210025
	ERR_PARTY_NOT_END          = 210026
	ERR_STAGE_CHECK_FAIL       = 210027
	ERR_GET_PARTY_STAGE        = 210028
	ERR_SEND_CHAT              = 210029
	ERR_ALREADY_LEAVE_PARTY    = 210030

	//suit套装
	ERR_ADD_SUIT_ERR    = 220010
	ERR_GET_SUIT_ERR    = 220011
	ERR_DELETE_SUTI_ERR = 220012

	NOT_EXIST_TOPIC = 230010

	ERR_TRANSFER_PXC                  = 240000
	ERR_GET_DESIGNER_ETH_ACCOUNT      = 240001
	ERR_GET_ETH_BALANCE               = 240002
	ERR_GET_CONTRACT_BALANCE          = 240003
	ERR_ETH_ACCOUNT_EMPTY             = 240004
	ERR_LIST_PXC_LOG                  = 240005
	ERR_ESTIMATE_GAS                  = 240006
	ERR_DEST_ACCOUNT_EMPTY            = 240007
	ERR_ESTIMATE_TYPE_WRONG           = 240008
	ERR_PARSE_HEX                     = 240009
	ERR_GET_GAS_PRICE                 = 240010
	ERR_EMAIL_EMPTY                   = 240011
	ERR_EMAIL_NOT_MATCH               = 240012
	ERR_SEND_VERIFY_TOO_FAST          = 240013
	ERR_SEND_ETH_VERIFY_CODE          = 240014
	ERR_VERIFY_CODE_NOT_MATCH         = 240015
	ERR_VERIFY_CODE_EXPIRE            = 240016
	ERR_ETH_PAY_PWD_WRONG             = 240017
	ERR_TRANSFER_OUT_PARAM_WRONG_TYPE = 240018
	ERR_TRANSFER_TOO_BIG              = 240019

	//kor
	ERR_ADD_INVITE              = 900000
	ERR_COUNT_INVITE            = 900001
	ERR_PARAM_WRONG             = 900002
	ERR_INVITE_ALREADY_REWARDED = 900003
	ERR_ADD_UNREGISTER          = 900004
)

//gl err code
const (
	ERR_LOGIN_NO_SUCH_USER ErrorCode = 1006

	ERR_EDIT_PWD                ErrorCode = 1007
	ERR_EDIT_PWD_NEW_PWD_LENGTH ErrorCode = 1008
	ERR_EDIT_PWD_DATA_ERR       ErrorCode = 1009

	ERR_LOGIN_TOO_FAST ErrorCode = 1010

	ERR_LOGIN_GF ErrorCode = 1011

	ERR_LOGOUT_USER              ErrorCode = 1012
	ERR_LOGOUT_LAST_SERVER_ALIVE ErrorCode = 1013
	ERR_QUERY_LAST_LOGIN_INFO    ErrorCode = 1014
	ERR_LOGOUT_ALL               ErrorCode = 1015
	ERR_LOGIN_PLAYER_LOCK        ErrorCode = 1016
	ERR_DELETE_PLAYER            ErrorCode = 1017

	GAME_UPDATE_CORE_VERSION      ErrorCode = 2002
	GAME_UPDATE_SCRIPT_VERSION    ErrorCode = 2003
	ERR_PARSE_SERVER_GAME_VERSION ErrorCode = 2004
	GAME_GET_VERSION_ERROR        ErrorCode = 2005

	ERR_BILI_SESSION_VERIFY_FORMAT ErrorCode = 5000
	ERR_BILI_SESSION_VERIFY        ErrorCode = 5001
	ERR_BILI_SESSION_VERIFY_FAILED ErrorCode = 5002 // 通过accesstoken验证uid失败
	ERR_GAME_SERVER_DOWN           ErrorCode = 5003 // 无可用的游戏服务器

	ERR_VIVO_CHECK_FAIL     ErrorCode = 6000
	ERR_HUAWEI_CHECK_FAIL   ErrorCode = 6001
	ERR_KUAIKAN_CHECK_FAIL  ErrorCode = 6002
	ERR_BUKA_CHECK_FAIL     ErrorCode = 6003
	ERR_MEITU_CHECK_FAIL    ErrorCode = 6004
	ERR_PAPA_CHECK_FAIL     ErrorCode = 6005
	ERR_4399_CHECK_FAIL     ErrorCode = 6006
	ERR_UC_CHECK_FAIL       ErrorCode = 6007
	ERR_XIAOMI_CHECK_FAIL   ErrorCode = 6008
	ERR_SINA_CHECK_FAIL     ErrorCode = 6009
	ERR_360_CHECK_FAIL      ErrorCode = 6010
	ERR_OPPO_CHECK_FAIL     ErrorCode = 6011
	ERR_THIRD_CHANNEL_EMPTY ErrorCode = 6012
	ERR_BAIDU_CHECK_FAIL    ErrorCode = 6013

	ERR_KAKAO_CHECK_FAIL ErrorCode = 7000
)

//gm err code
const (
	// 注册
	ERR_REGISTER_INPUT_WRONG   ErrorCode = 1000
	ERR_REGISTER_ADMIN_EXIST   ErrorCode = 1001
	ERR_REGISTER_QUERY_FAILED  ErrorCode = 1002
	ERR_REGISTER_INSERT_FAILED ErrorCode = 1003
	// 登陆
	ERR_LOGIN_INPUT_WRONG  ErrorCode = 1010
	ERR_LOGIN_QUERY_FAILED ErrorCode = 1011
	ERR_LOGIN_NOT_EXIST    ErrorCode = 1012
	ERR_LOGIN_ERROR        ErrorCode = 1013
	// 群发邮件
	ERR_TITLE_LENGTH    ErrorCode = 1020
	ERR_CONTENT_LENGTH  ErrorCode = 1021
	ERR_DIAMOND_GOLD    ErrorCode = 1022
	ERR_COOKIE_USERNAME ErrorCode = 1023
	ERR_QUERY_USERNAME  ErrorCode = 1024
	ERR_SEND_FAIL       ErrorCode = 1025
	ERR_ALL_GET_MAIL_ID ErrorCode = 1026
	// 单发邮件
	ERR_ONE_TITLE_LENGTH    ErrorCode = 1027
	ERR_ONE_CONTENT_LENGTH  ErrorCode = 1028
	ERR_ONE_DIAMOND_GOLD    ErrorCode = 1029
	ERR_ONE_COOKIE_USERNAME ErrorCode = 1030
	ERR_ONE_QUERY_USERNAME  ErrorCode = 1031
	ERR_ONE_SEND_FAIL       ErrorCode = 1032
	ERR_ONE_TO_LENGTH       ErrorCode = 1033
	ERR_ONE_QUERY_ADMINNAME ErrorCode = 1034
	ERR_ONE_GET_MAIL_ID     ErrorCode = 1035

	// 玩家管理
	ERR_USERNAME_LENGTH ErrorCode = 2000
	ERR_FIND_USER_ERR   ErrorCode = 2001

	ERR_FIND_ALL_USER_PAGE   ErrorCode = 2002
	ERR_FIND_ALL_USER_FAILED ErrorCode = 2003

	// 推送管理
	ERR_PUSH_TO_LENGTH          ErrorCode = 3001
	ERR_PUSH_CANNOT_FIND_PLAYER ErrorCode = 3002
	ERR_PUSH_CONTENT_LENGTH     ErrorCode = 3003

	// 查询玩家上传的衣服信息
	ERR_CUSTOM_FIND ErrorCode = 4000
	ERR_SET_PASS    ErrorCode = 4001
	ERR_SET_FAIL    ErrorCode = 4002

	// cos管理
	ERR_COS_TITLE_LEN        ErrorCode = 5000
	ERR_COS_KEYWORD_LEN      ErrorCode = 5001
	ERR_COS_WRONG_TYPE       ErrorCode = 5002
	ERR_COS_WRONG_OPEN_TIME  ErrorCode = 5003
	ERR_COS_WRONG_CLOSE_TIME ErrorCode = 5004
	ERR_ALREADY_RUNNING_ONE  ErrorCode = 5005 // 已经有正在进行的Cosplay了
	ERR_COS_ADD_FAILED       ErrorCode = 5006 // 添加失败
	ERR_COS_FIND_FAILED      ErrorCode = 5007 // 查询失败

	ERR_COS_REWARD_1_LEN ErrorCode = 5008
	ERR_COS_REWARD_2_LEN ErrorCode = 5009
	ERR_COS_REWARD_3_LEN ErrorCode = 5010

	ERR_COS_NO_FILE_UPLOADED ErrorCode = 5011
	ERR_COS_FILE_UPLOADED    ErrorCode = 5012
	ERR_PUB_COS_PARAM        ErrorCode = 5013 // cos参数出错了

	ERR_NEED_LOGIN      ErrorCode = 6000
	GM_ERR_WRONG_PARAMS ErrorCode = 7000

	ERR_PLATFORM_NOT_SUPPORT = 8000
)

// //gv err code
// const (
// 	// 审核管理
// 	ERR_VERIFY_INSERT_FAILED = 4001
// 	ERR_VERIFY_UPDATE_FAILED = 4002
// 	ERR_VERIFY_QUERY_FAILED  = 4003
// )
