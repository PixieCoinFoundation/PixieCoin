package language

type LanguageDetail struct {
	Chinese string
	Korean  string
}

var languageMap map[string]*LanguageDetail = map[string]*LanguageDetail{
	"cos1": &LanguageDetail{
		Chinese: "感谢你参与Cosplay【%s】",
		Korean:  "【%s】파티에 참여해주셔서 감사드려요!",
	},
	"cos2": &LanguageDetail{
		Chinese: "恭喜你进入Cosplay【%s】排名前%d%%",
		Korean:  "【%s】파티 TOP %d%% 진입을 축하드려요!",
	},
	"cos3": &LanguageDetail{
		Chinese: "cosplay结算奖励",
		Korean:  "파티결산보상",
	},
	"cos4": &LanguageDetail{
		Chinese: "前十名",
		Korean:  "10위 내",
	},
	"cos5": &LanguageDetail{
		Chinese: "第一名",
		Korean:  "1위",
	},
	"cos6": &LanguageDetail{
		Chinese: "第二名",
		Korean:  "2위",
	},
	"cos7": &LanguageDetail{
		Chinese: "第三名",
		Korean:  "3위",
	},
	"cos8": &LanguageDetail{
		Chinese: "第三名",
		Korean:  "3위",
	},
	"cos9": &LanguageDetail{
		Chinese: "恭喜您进入cosplay【%s】%s", //例如：恭喜你进入cosplay【哈哈哈】前10名
		Korean:  "【%s】파티 TOP %s 진입을 축하드려요!",
	},

	"design1": &LanguageDetail{
		Chinese: "设计师新品通知",
		Korean:  "",
	},
	"design2": &LanguageDetail{
		Chinese: "设计师【%s】发布了新作品【%s】，快去游戏里看看吧~",
		Korean:  "",
	},
	"design3": &LanguageDetail{
		Chinese: "补货通知",
		Korean:  "매진알림",
	},
	"design4": &LanguageDetail{
		Chinese: "你的作品卖得太好啦，小菊喊你去补货~",
		Korean:  "등록하신 의상이 너무 핫한걸요~? 어서와서 재고를 채워주세요~",
	},
	"design5": &LanguageDetail{
		Chinese: "作品完全下架通知",
		Korean:  "디자인 삭제 안내",
	},
	"design6": &LanguageDetail{
		Chinese: "你的作品【%s】已完全下架，回收金币收益：%d，钻石收益：%d，设计师经验:%d。",
		Korean:  "디자이너님의 [%s] 디자인이 운영정책에 따라 삭제 처리되었습니다. 회수 골드: [%d] ; 쥬얼: [%d] ; 디자이너 경험치: [%d]",
	},

	"friend1": &LanguageDetail{
		Chinese: "好友申请通知",
		Korean:  "",
	},
	"friend2": &LanguageDetail{
		Chinese: "%s（等级：%d，VIP：%d）申请成为你的好友",
		Korean:  "",
	},

	"hope1": &LanguageDetail{
		Chinese: "心愿完成通知",
		Korean:  "",
	},
	"hope2": &LanguageDetail{
		Chinese: "您的心愿已被达成，快去游戏中领取吧~",
		Korean:  "",
	},

	"party1": &LanguageDetail{
		Chinese: "舞会【%s】举办奖励",
		Korean:  "[%s] 파티 개최 보상",
	},
	"party2": &LanguageDetail{
		Chinese: "您举办的【%s】舞会顺利结束，举办奖励请查收~",
		Korean:  "[%s] 파티가 종료되었어요. 보상을 확인해 주세요!",
	},
	"party3": &LanguageDetail{
		Chinese: "您举办的【%s】舞会顺利结束，参与玩家数超过100的奖励请查收~",
		Korean:  "[%s] 파티가 종료되었어요. 참여자 수가 100명을 넘어 보상을 드릴게요!",
	},
	"party4": &LanguageDetail{
		Chinese: "舞会【%s】结算奖励",
		Korean:  "[%s] 파티 결산 보상",
	},
	"party5": &LanguageDetail{
		Chinese: "您参加的【%s】舞会顺利结束，恭喜您获得第%d名（该排名段并列搭配数：%d）。",
		Korean:  "[%s] 파티가 종료되었어요. %d위를 축하드려요（공동 순위자 수：%d）!",
	},
	// "party6": &LanguageDetail{
	// 	Chinese: "和舞伴：",
	// 	Korean:  "파트너와의:",
	// },
	"party7": &LanguageDetail{
		Chinese: "舞会【%s】门票收入",
		Korean:  "[%s] 파티 입장료",
	},
	"party8": &LanguageDetail{
		Chinese: "您举办的【%s】舞会顺利结束，门票收入请查收~",
		Korean:  "[%s] 파티가 종료되었어요. 입장료를 확인해 주세요! ",
	},
	"party9": &LanguageDetail{
		Chinese: "您被邀请参加的【%s】舞会顺利结束，恭喜您和舞伴获得第%d名（该排名段并列搭配数：%d）。",
		Korean:  "초대로 참석한 [%s] 파티가 종료되었어요.  %d위를 축하드려요（공동 순위자 수：%d）!",
	},
	"party10": &LanguageDetail{
		Chinese: "好友：%s 邀请你成为舞会【%s】的舞伴，快去参加吧：）",
		Korean:  "친구 %s님이 [%s] 파티에 초대했어요!",
	},
	"party11": &LanguageDetail{
		Chinese: "%s邀请你成为舞伴",
		Korean:  "%s님이 파트너 초대를 보냈어요!",
	},
	// "party12": &LanguageDetail{
	// 	Chinese: "您参加的【%s】舞会顺利结束，恭喜您%s获得第%d名。",
	// 	Korean:  "참여하신【%s】파티가 종료 되었어요.%s %d위를 축하드려요.",
	// },
	"party13": &LanguageDetail{
		Chinese: "好友%s向你要一朵花",
		Korean:  "친구%s 님이 꽃을 요청했어요",
	},
	"party14": &LanguageDetail{
		Chinese: "%s正在参加舞会【%s】，并向你要了一朵花。",
		Korean:  "%s님이 【%s】파티에 참가 중이에요. 당신에게 꽃을 요청했어요. ",
	},
	"party15": &LanguageDetail{
		Chinese: "好友%s邀请你参加舞会",
		Korean:  "친구 %s님이 파티에 초대했어요. ",
	},
	"party16": &LanguageDetail{
		Chinese: "%s邀请你参加舞会【%s】。",
		Korean:  "%s님이 【%s】파티에 초대했어요.",
	},
	"party17": &LanguageDetail{
		Chinese: "您参加的【%s】舞会顺利结束，恭喜您获得第%d名。",
		Korean:  "참여하신【%s】파티가 종료 되었어요. %d위를 축하드려요.",
	},
	"party18": &LanguageDetail{
		Chinese: "您参加的【%s】舞会顺利结束，恭喜您和舞伴：%s获得第%d名。",
		Korean:  "참여하신【%s】파티가 종료 되었어요.파트너%s와의 %d위를 축하드려요.",
	},
	"party19": &LanguageDetail{
		Chinese: "舞会关闭通知",
		Korean:  "코스튬파티 삭제 안내",
	},
	"party20": &LanguageDetail{
		Chinese: "舞会【%s】已被管理员强制关闭",
		Korean:  "참여하신 【%s】가(이) 부적절한 파티 주제로 확인되어 삭제처리 되었습니다.",
	},
	"guild1": &LanguageDetail{
		Chinese: "公会加入成功",
		Korean:  "길드가입성공",
	},
	"guild2": &LanguageDetail{
		Chinese: "您向【%s】提交的入会申请已被批准！",
		Korean:  "【%s】길드에 가입되었습니다.",
	},
	"guild3": &LanguageDetail{
		Chinese: "移除公会通知",
		Korean:  "길드제명알림",
	},
	"guild4": &LanguageDetail{
		Chinese: "您已被会长移出了妖精公会【%s】，您可以继续加入其它公会。",
		Korean:  "【%s】길드에서 제명되었습니다. 타 길드에 가입신청 가능합니다.",
	},
	"guild5": &LanguageDetail{
		Chinese: "公会挑战赛开始",
		Korean:  "길드대전시작",
	},
	"guild6": &LanguageDetail{
		Chinese: "您所在的妖精公会向【%s】公会发起了公会挑战赛，赶快去做准备吧~",
		Korean:  "당신이 속한 길드에서 【%s】길드에 대전을 신청하였습니다. 서둘러 대전 준비해 주세요~",
	},
	"guild7": &LanguageDetail{
		Chinese: "【%s】公会向您所在的公会发起了公会挑战赛，赶快去做准备吧~",
		Korean:  "【%s】길드에서 당신이 속한 길드에 대전을 신청하였습니다. 서둘러 대전 준비해 주세요~",
	},
	"guild8": &LanguageDetail{
		Chinese: "公会挑战赛排名奖励",
		Korean:  "길드대전랭킹보상",
	},
	"guild9": &LanguageDetail{
		Chinese: "公会抽奖奖励",
		Korean:  "길드뽑기보상",
	},
	"guild10": &LanguageDetail{
		Chinese: "由于你和你的妖精公会本周足够活跃，公会成员【%s】为本公会抽取到一件【%s】奖励！在游戏中继续加油吧~",
		Korean:  "당신이 속한 길드의 활약도가 뛰어나 길드원 【%s】이 길드뽑기에서 【%s】보상을 획득하였습니다. 계속 분발해주세요~",
	},
	"guild11": &LanguageDetail{
		Chinese: "公会战设置防守提醒",
		Korean:  "길드전방어알림설정",
	},
	"guild12": &LanguageDetail{
		Chinese: "公会战正在应战准备阶段，快去设置防守阵容为公会贡献力量吧~",
		Korean:  "현재 길드전 준비단계입니다. 방어진영을 설정하여 길드전 승리에 기여합시다~",
	},
	"guild13": &LanguageDetail{
		Chinese: "公会挑战赛赛季排名嘉奖",
		Korean:  "길드전시즌랭킹보상",
	},
	"guild14": &LanguageDetail{
		Chinese: "您所在的公会在上赛季公会挑战赛中获得了第%d名的好成绩，特送上排名奖励，这个赛季继续加油哦~",
		Korean:  "당신이 속한 길드가 지난 시즌 길드전에서 %d위의 성적을 거둬 이에 랭킹보상을 보내드립니다. 다음 시즌에도 파이팅 해 주세요~",
	},
	"guild15": &LanguageDetail{
		Chinese: "公会挑战赛获胜",
		Korean:  "길드전승리",
	},
	"guild16": &LanguageDetail{
		Chinese: "恭喜您所在公会在本次挑战赛获得胜利！感谢你对公会的无私贡献，请查收获胜嘉奖！",
		Korean:  "축하드립니다! 이번 길드전에서 승리하여 보상을 보내드립니다. 길드에 아낌없는 기여 감사드립니다!",
	},
	"guild17": &LanguageDetail{
		Chinese: "公会挑战赛赛季段位奖励",
		Korean:  "길드전시즌참여보상",
	},
	"guild18": &LanguageDetail{
		Chinese: "上赛季公会挑战赛已结束，请查收公会段位奖励~",
		Korean:  "지난 시즌 길드전이 종료되었습니다. 참여보상을 확인해 주세요~",
	},
	"guild19": &LanguageDetail{
		Chinese: "%s在公会回复了你",
		Korean:  "",
	},
	"guild20": &LanguageDetail{
		Chinese: "%s在公会留言板回复了你：【%s】",
		Korean:  "",
	},
	"guild21": &LanguageDetail{
		Chinese: "公会解散通知",
		Korean:  "",
	},
	"guild22": &LanguageDetail{
		Chinese: "您所在的公会已解散，请继续选择别的公会加入哦~",
		Korean:  "",
	},
	"guild23": &LanguageDetail{
		Chinese: "公会转让通知",
		Korean:  "",
	},
	"guild24": &LanguageDetail{
		Chinese: "您被%s设置为新的公会会长，赶快上线跟你的公会小伙伴们打个招呼吧~",
		Korean:  "",
	},
	"guild25": &LanguageDetail{
		Chinese: "设置副会长通知",
		Korean:  "",
	},
	"guild26": &LanguageDetail{
		Chinese: "你在公会中被设置为副会长，副会长可以发起挑战赛哦，快去公会中查看吧~",
		Korean:  "",
	},

	"cp1": &LanguageDetail{
		Chinese: "很可惜！您发送的碎片交换请求无人帮助，请查收退回的服装碎片:)",
		Korean:  "의상조각 도움을 받지 못했어요. 요청하셨던 의상조각을 확인해 주세요!",
	},
	"cp2": &LanguageDetail{
		Chinese: "碎片交换请求超时",
		Korean:  "의상조각 요청 시간이 초과되었어요.",
	},
	"cp3": &LanguageDetail{
		Chinese: "您发送的碎片交换请求已达成，请查收服装碎片:)",
		Korean:  "의상조각 도움을 받았어요. 의상조각을 확인해 주세요!",
	},
	"cp4": &LanguageDetail{
		Chinese: "碎片交换达成超时",
		Korean:  "의상조각 교환 시간이 초과되었어요.",
	},

	"pk1": &LanguageDetail{
		Chinese: "赛季PK级别提升奖励",
		Korean:  "코디대결 시즌 칭호 보상",
	},
	"pk2": &LanguageDetail{
		Chinese: "恭喜,您的赛季PK等级提升为:【%s】,请查收奖励!",
		Korean:  "축하드려요! 이번 시즌 코디대결에서 [%s] 칭호를 달성했어요. 보상을 확인해 주세요!",
	},
	"pk3": &LanguageDetail{
		Chinese: "首次参与PK奖励",
		Korean:  "첫 코디대결 참여 보상",
	},
	"pk4": &LanguageDetail{
		Chinese: "首次参与PK,如果获胜请再接再厉,即便失败也请不要气馁哟!",
		Korean:  "첫 코디대결 참여를 축하드려요! 앞으로도 승리가 가득하시길 바랄게요.",
	},
	"pk5": &LanguageDetail{
		Chinese: "PK周结算-排名奖",
		Korean:  "코디대결 주간 순위 결산 보상",
	},
	"pk6": &LanguageDetail{
		Chinese: "恭喜您在上周的PK玩法中获得第%d名",
		Korean:  "축하드려요! 지난주 코디대결에서 %d위를 달성했어요.",
	},
	"pk7": &LanguageDetail{
		Chinese: "PK赛季结算-排名奖",
		Korean:  "코디대결 시즌 순위 결산 보상 ",
	},
	"pk8": &LanguageDetail{
		Chinese: "恭喜您在上赛季的PK玩法中获得第%d名",
		Korean:  "축하드려요! 지난 시즌 코디대결에서 %d위를 달성했어요. ",
	},
	"pk9": &LanguageDetail{
		Chinese: "PK结算奖励",
		Korean:  "코디대결 결산 보상 ",
	},
	"pk10": &LanguageDetail{
		Chinese: "上周",
		Korean:  "지난주",
	},
	"pk11": &LanguageDetail{
		Chinese: "上赛季",
		Korean:  "지난 시즌",
	},
	"pk12": &LanguageDetail{
		Chinese: "恭喜您在%sPK活动中获得称号：【%s】，结算奖励请查收！记得多多参与PK玩法哟：）。",
		Korean:  "축하드려요! %s 코디대결 이벤트에서 [%s] 칭호를 획득했어요. 결산 보상을 확인해 주세요!",
	},
	"pk13": &LanguageDetail{
		Chinese: "学徒",
		Korean:  "학생",
	},
	"pk14": &LanguageDetail{
		Chinese: "师兄",
		Korean:  "보조스탭",
	},
	"pk15": &LanguageDetail{
		Chinese: "出师",
		Korean:  "스탭",
	},
	"pk16": &LanguageDetail{
		Chinese: "裁缝",
		Korean:  "주니어",
	},
	"pk17": &LanguageDetail{
		Chinese: "设计助理",
		Korean:  "시니어",
	},
	"pk18": &LanguageDetail{
		Chinese: "设计师",
		Korean:  "디자이너",
	},
	"pk19": &LanguageDetail{
		Chinese: "时尚达人",
		Korean:  "대가",
	},
	"pk20": &LanguageDetail{
		Chinese: "国际大师",
		Korean:  "거장",
	},
	"pk21": &LanguageDetail{
		Chinese: "服装泰斗",
		Korean:  "명장",
	},
	"pk22": &LanguageDetail{
		Chinese: "衣圣",
		Korean:  "전설",
	},

	"o1": &LanguageDetail{
		Chinese: "公测预约礼包",
		Korean:  "사전예약 선물",
	},
	"o2": &LanguageDetail{
		Chinese: "感谢你对妖精的衣橱的喜爱，祝游戏愉快~",
		Korean:  "사랑해주셔서 감사드립니다! 항상 행복하세요! ",
	},
	"o3": &LanguageDetail{
		Chinese: "补偿通知",
		Korean:  "보상알림",
	},
	"o4": &LanguageDetail{
		Chinese: "玩家你好，非常抱歉，您之前购买的设计师作品：【%s】因涉嫌抄袭已从您的衣橱中删除，请查收补偿。",
		Korean:  "안내 말씀드립니다. 이전에 구입하신 디자이너 의상 [%s] 이(가) 운영정책에 따라 회수 처리 되었습니다. 보상 내역을 확인 부탁드립니다. ",
	},
	"o5": &LanguageDetail{
		Chinese: "服装删除通知",
		Korean:  "의상삭제알림",
	},
	"o6": &LanguageDetail{
		Chinese: "玩家你好，非常抱歉，您之前购买的设计师作品：【%s】因涉嫌：抄袭/侵权/其他原因，已从您的衣橱中删除，请查收补偿，给您带来的不便深感抱歉。",
		Korean:  "죄송합니다. 기존에 구입하신 디자이너 작품 【%s】이 운영정책에 따라 삭제처리 되었습니다. 해당 의상은 이미 옷장에서 삭제되었으며 이에 보상을 보내드릴 예정입니다. 우편함을 확인해주세요.",
	},
	"o7": &LanguageDetail{
		Chinese: "%s回复了你",
		Korean:  "%s이(가) 답장하였습니다.",
	},
	"o8": &LanguageDetail{
		Chinese: "%s回复了你：【%s】",
		Korean:  "%s이(가) 답장하였습니다:【%s】",
	},
	"o9": &LanguageDetail{
		Chinese: "礼包码奖励",
		Korean:  "쿠폰보상",
	},
	"o10": &LanguageDetail{
		Chinese: "兑换成功，请领取奖励",
		Korean:  "교환 성공. 보상을 받으세요.",
	},
	"o11": &LanguageDetail{
		Chinese: "妖精管理局",
		Korean:  "유나의 옷장",
	},
	"o12": &LanguageDetail{
		Chinese: "妖精公会管理局",
		Korean:  "유나의 옷장",
	},
	"o13": &LanguageDetail{
		Chinese: "月卡奖励",
		Korean:  "월정액보상",
	},
	"o14": &LanguageDetail{
		Chinese: "请查收月卡套装奖励",
		Korean:  "월정액 세트보상을 확인해 주세요",
	},
	"o15": &LanguageDetail{
		Chinese: "购物车商品库存紧张",
		Korean:  "장바구니에 있는 상품이 곧 품절되요",
	},
	"o16": &LanguageDetail{
		Chinese: "您购物车中的作品【%s】库存紧张，快去购买吧~",
		Korean:  "장바구니에 있는【%s】의상이 곧 품절되요. 어서 구매하세요~",
	},
	"o17": &LanguageDetail{
		Chinese: "购物车商品补货啦",
		Korean:  "장바구니에 있는 의상재고가 채워졌어요~",
	},
	"o18": &LanguageDetail{
		Chinese: "您购物车中的作品【%s】已补货，快去购买吧~",
		Korean:  "장바구니에 있는【%s】의상의 재고가 채워졌어요. 어서 구매하세요~ ",
	},
	"o19": &LanguageDetail{
		Chinese: "%d月%d日%s点到%s点",
		Korean:  "%d월%d일%s시부터%s시까지",
	},

	"ad1": &LanguageDetail{
		Chinese: "设计师打CALL日榜奖励",
		Korean:  "응원하기 일일순위 보상",
	},
	"ad2": &LanguageDetail{
		Chinese: "设计师打CALL月榜奖励",
		Korean:  "응원하기 월순위 보상",
	},
	"ad3": &LanguageDetail{
		Chinese: "设计师竞拍花费返还通知",
		Korean:  "디자이너 경매비용 반환알림",
	},
	"ad4": &LanguageDetail{
		Chinese: "您当日为设计师%s打CALL贡献值排名第%d。为了答谢您的支持，设计师送给您%s,请笑纳！这些都是不会上架的非卖品哦~",
		Korean:  "오늘 디자이너 %s에 대한 응원 공헌순위에서 %d등을 달성했어요. 해당 디자이너가 감사의 선물로 %s을 선물했으니, 기쁘게 받아주세요~! 또한 해당 의상은 모두 판매불가한 비매품이니 유의해주세요~",
	},
	"ad5": &LanguageDetail{
		Chinese: "您当月为设计师%s打CALL贡献值排名第%d。为了答谢您的支持，设计师送给您%s,请笑纳！这些都是不会上架的非卖品哦~",
		Korean:  "이번달 디자이너 %s에 대한 응원 공헌순위에서 %d등을 달성했어요. 해당 디자이너가 감사의 선물로 %s을 선물했으니, 기쁘게 받아주세요~! 또한 해당 의상은 모두 판매불가한 비매품이니 유의해주세요~",
	},
	"ad6": &LanguageDetail{
		Chinese: "亲爱的设计师,您好很遗憾您的【%s】 在%s的广告竞拍中未能竞拍成功,竞拍出价和加价(扣除5%%的手续费)已退还。",
		Korean:  "친애하는 디자이너님. 정말 아쉽게도【%s】의상이 %s기간에 진행된 Hot스팟 경매에 낙찰받지 못했어요. 경매 중 지불한 출고가와 인상가(등록비5%는 제외)는 다시 반환해 드릴게요~",
	},
	"ad7": &LanguageDetail{
		Chinese: "设计师月度嘉奖",
		Korean:  "디자이너 월공헌 보상",
	},
	"ad8": &LanguageDetail{
		Chinese: "亲爱的设计师,新的一月已经开始，请继续加油哦",
		Korean:  "친애하는 디자이너님~ 새로운 달이 시작됬어요~ 계속 화이팅해 주세요!",
	},
	"ad9": &LanguageDetail{
		Chinese: "广告位竞拍成功",
		Korean:  "Hot스팟 낙찰성공",
	},
	"ad10": &LanguageDetail{
		Chinese: "亲爱的设计师,恭喜您的作品【%s】 在%s的广告竞拍中竞拍成功。",
		Korean:  "축하드립니다~디자이너님! 등록하신【%s】의상이 %s기간에 진행된 Hot스팟 낙찰에 성공했어요.",
	},

	"pxc1": &LanguageDetail{
		Chinese: "设计师作品购买成功",
		Korean:  "디자이너의상 구매성공",
	},
	"pxc2": &LanguageDetail{
		Chinese: "亲爱的玩家你好，您向设计师【%s】购买的%d件作品【%s】已成功，请查收！",
		Korean:  "안녕하세요~ 【%s】디자이너에게서 구매한【%s】의상 %d벌 구매에 성공하였습니다",
	},
	"pxc3": &LanguageDetail{
		Chinese: "亲爱的玩家你好，您向设计师【%s】批量购买的作品已成功，请查收！",
		Korean:  "안녕하세요~ 【%s】디자이너의 의상 대량구매에 성공하였습니다. 확인 후 아이템을 수령해주세요",
	},
	"pxc4": &LanguageDetail{
		Chinese: "设计师作品购买失败",
		Korean:  "디자이너의상 구매 실패",
	},
	"pxc5": &LanguageDetail{
		Chinese: "亲爱的玩家你好，您向设计师【%s】购买的作品处理失败。",
		Korean:  "안녕하세요~ 【%s】디자이너의 의상 구입처리에 실패하였습니다",
	},
	"pxc6": &LanguageDetail{
		Chinese: "妖精的衣橱-验证码",
		Korean:  "유나의옷장-인증번호",
	},
	"pxc7": &LanguageDetail{
		Chinese: "亲爱的玩家你好，您正在设置资金密码，对应的验证码为：【%s】。官方不会以任何理由向您所要验证码，请勿将此验证码告知他人！",
		Korean:  "안녕하세요~ 현재 설정 중이신 인출 비밀번호의 인증번호는 [%s] 입니다. 관리자가 별도로 인증번호를 요청하는 상황은 없으며, 또한 타인에게 인증번호를 알려주시면 안돼요~",
	},
	"pixie_suit": &LanguageDetail{
		Chinese: "当前搭配和【%s】套装相同～",
		Korean:  "",
	},
}
