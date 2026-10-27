package constants

type ContextKey int

const (
	ContextKeyChannel ContextKey = iota
	ContextKeyCity
	ContextKeyEnabled
	ContextKeyFeed
	ContextKeyJob
	ContextKeyLevel
	ContextKeyOrder
	ContextKeyServer
	ContextKeyTwitter
)
