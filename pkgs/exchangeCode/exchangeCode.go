package exchangeCode

import (
	"context"
)

type key int

const (
	codeKey key = 0
)

func NewContext(ctx context.Context, code string) context.Context {
	return context.WithValue(ctx, codeKey, code)
}

func FromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(codeKey).(string)
	return token, ok
}
