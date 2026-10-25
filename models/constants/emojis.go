package constants

const (
	EmojiIDFirst    EmojiMiscID = "first"
	EmojiIDGame     EmojiMiscID = "dofusGame"
	EmojiIDGithub   EmojiMiscID = "github"
	EmojiIDGlobal   EmojiMiscID = "global"
	EmojiIDLast     EmojiMiscID = "last"
	EmojiIDNext     EmojiMiscID = "next"
	EmojiIDPrevious EmojiMiscID = "previous"
	EmojiIDScroll   EmojiMiscID = "scroll"
	EmojiIDTwitter  EmojiMiscID = "twitter"

	EmojiTypeMisc EmojiType = "Miscellaneous"

	EmojiTypeCity   EmojiEntity = "City"
	EmojiTypeJob    EmojiEntity = "Job"
	EmojiTypeOrder  EmojiEntity = "Order"
	EmojiTypeServer EmojiEntity = "Server"
)

type EmojiMiscID string
type EmojiType string
type EmojiEntity EmojiType
