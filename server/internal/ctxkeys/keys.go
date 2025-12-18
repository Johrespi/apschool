package ctxkeys

import "context"

type contextKey string

const UserIDKey contextKey = "userID"

func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}
