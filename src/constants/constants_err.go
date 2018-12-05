package constants

import (
	"errors"
)

var PartySizeErr = errors.New("party size limit")
var ParsePlayerErr = errors.New("error parse player")
var GMLogServerErr = errors.New("gmlog server num not match")
var GuildSizeErr = errors.New("guild size limit")
var GuildWarSizeErr = errors.New("guild war size limit")
var GuildClothesSizeErr = errors.New("guild clothes size limit")
var LogUploadDoneErr = errors.New("log already uploaded")
var MatchCntErr = errors.New("try match later")
var AlreadyInGuildErr = errors.New("player already in guild")
var PlayerNotBuyCustomErr = errors.New("player do not buy custom")
var OutofprintCustomEffectErr = errors.New("out of print custom effect 0")
var DelCustomEffectErr = errors.New("del custom effect 0")
var CustomCancelStatusWrong = errors.New("custom status wrong")
var DesignerWalletNotEnoughErr = errors.New("designer wallet not enough")
var FinishCustomCopyStatusErr = errors.New("finsh custom copy status")
var GuildVicePresidentSizeErr = errors.New("guild vice president size full")
var GuildVicePresidentConfigErr = errors.New("guild vice president config type wrong")
var GuildVicePresidentConfigEffectErr = errors.New("guild vice president config effect 0")
var NoPartyOpenErr = errors.New("no party open")
var ClosePartyEffectErr = errors.New("close party effect")
var PlayerDeledErr = errors.New("player already be deleted")

var GetContractBalanceErr = errors.New("get contract balance error.maybe contract address wrong.")

var GuildNotExistErr = errors.New("guild not exist")

var PaperVerifyNoReward = errors.New("paper verify no reward")
var PARAMS_CAN_NOT_EMPTY = errors.New("params can not empty")
var ErrorNilData = errors.New("data can not empty")

var PlayerNotFoundErr = errors.New("player not found")
var UpdatePlayerEffectZeroErr = errors.New("update player effect 0")
var UpdatePaperVerifiedInfoEffectZero = errors.New("update paper verifiedInfo effect 0")
var UpdateVerifyCountEffectZero = errors.New("update verify count effect 0")
var UpdatePaperSupportEffectZero = errors.New("update verify support effect 0")

//rt party error
var UnknownMatchRTPartyErr = errors.New("unknown match rt party")
var LockRTPartyListErr = errors.New("lock party list error")
var LockRTPartyHostErr = errors.New("lock party host error")
var PartyHostFormatErr = errors.New("party host format error")
var AlreadyOfferSubjectErr = errors.New("already offer subject")
var AlreadyVoteSubjectErr = errors.New("already vote subject")
var AlreadyVoteDressErr = errors.New("already vote dress")
var PartyAlreadyStartErr = errors.New("party already start")
var PartyNotAttendErr = errors.New("not in party")
var PartyAlreadyLeftErr = errors.New("already leave party")
var SubjectEmptyErr = errors.New("subject empty")
var RTPartyHostFullErr = errors.New("host full")
var RTPartyHostDeadErr = errors.New("host dead")
var RTPartyHostNotBeginErr = errors.New("host not begin")
var VoteToUndoneDressErr = errors.New("vote to undone dress")

//JoinAdverisingPlace 参加广告位类型不匹配
var ParamsNotMatch = errors.New("JoinAdverisingPlace clothes type with type not match")

var ProcessFileError = errors.New("process file error")

var UploadFileErr = errors.New("upload file error")

//pixie
var PaperStatusWrong = errors.New("paper status wrong")
var PaperAuthorOrOwnerWrong = errors.New("paper owner or author wrong")
var PaperSqlEffectZero = errors.New("paper sql effect 0")
var PaperTradeEffectZero = errors.New("paper log insert id  < 1 ")
var FileSizeZero = errors.New("file size 0")
var PaperTradeLogSaleCountEffectZero = errors.New("update paper trade log effect 0")
var ClothesPricingEffectZero = errors.New("update clothes pricing effect 0")
var ClothesPricingPriceNotChange = errors.New("update clothes pricing price not change")
var ClothesPriceChangeEffectZero = errors.New("update clothes price change effect 0")
var UpdatePaperNotifyEffectZero = errors.New("update paper notify effect 0")
var PixieNoPaperForVerify = errors.New("no paper for verify")
var PaperVerifySizeZero = errors.New("paper verify size 0")

var CheckVersionDownloadUrlWrong = errors.New("check version download url wrong")
var UnknownThirdChannel = errors.New("unknown third channel")

var BuildUserNotMatch = errors.New("build username not match")
var NotOccupyLand = errors.New("not occupy land")

var CancelPaperProductEffectZero = errors.New("cancel paper product effect 0")
var UpdatePaperAuctionFailUnreadZero = errors.New("update paper auction_fail_unread effect 0")
var UpdatePaperOccupyAuctionFailUnreadZero = errors.New("update paper occupy auction_fail_unread effect 0")
var UnknownSortTypeErr = errors.New("unknown sort type")
var AuctionPriceIllegal = errors.New("auction price illegal")
var TradeNotExist = errors.New("trade not exist")
var LandStatusWrong = errors.New("land status wrong")
var LandBuildTypeWrong = errors.New("land build type wrong")
var LandSqlEffectZero = errors.New("land sql effect 0")
var LandBusinessTypeWrong = errors.New("land business type wrong")
var DoNotHaveClothesWrong = errors.New("do not have clothes")
var SuitNotClothesWrong = errors.New("not id in suit")
var TypeNotMatchParams = errors.New("info not match params")
var SuitCountNotMatchSearchCount = errors.New("suit clothes count not match search count")
var InsertClothesSaleInfoEffectZero = errors.New("insert clothes sale info effect 0")
var SellClothesEffectZero = errors.New("sell clothes effect 0")

var CirculationNotSet = errors.New("circulation max not set")
var CirculationChange = errors.New("circulation max change")
var CirculationLimit = errors.New("circulation max limit")
var SequenceNotMatchTrade = errors.New("trade sequence not match")

var PlayerNotExist = errors.New("player not exist")
