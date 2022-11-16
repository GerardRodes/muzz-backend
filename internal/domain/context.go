package domain

import "context"

var contextUserIDKey = struct{}{}

func ContextWithUserID(parent context.Context, userID uint32) context.Context {
	return context.WithValue(parent, contextUserIDKey, userID)
}

func UserIDFromContext(ctx context.Context) uint32 {
	return ctx.Value(contextUserIDKey).(uint32)
}
