package types

type BoardMsg struct {
	ID             int64
	Owner          string
	Author         string
	AuthorNickname string
	Head           string
	ReplyTo        string
	Content        string
	Time           int64
}

type BoardMsgView struct {
	BoardMsg
	Times string
}

type GuildBoardMsg struct {
	ID             int64
	Author         string
	AuthorNickname string
	Head           string
	ReplyTo        string
	Content        string
	Time           int64
}
