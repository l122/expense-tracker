package authVerifier

import (
	"context"
)

type key int

const (
	verifierKey key = 0
)

func NewContext(ctx context.Context, verifier string) context.Context {
	return context.WithValue(ctx, verifierKey, verifier)
}

func FromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(verifierKey).(string)
	return token, ok
}
