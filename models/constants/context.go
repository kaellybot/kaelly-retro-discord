package constants

type ContextKey int

const (
	ContextKeyChannel ContextKey = iota
	ContextKeyCity
	ContextKeyJob
	ContextKeyLevel
	ContextKeyOrder
	ContextKeyServer
)
