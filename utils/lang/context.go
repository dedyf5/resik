// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package lang

import "context"

type langKey string

const (
	ContextKey langKey = "lang"
)

func (c langKey) String() string {
	return string(c)
}

func ContextWithLang(ctx context.Context, lang *Lang) context.Context {
	return context.WithValue(ctx, ContextKey, lang)
}

func ContextLang(ctx context.Context) *Lang {
	langCtx := ctx.Value(ContextKey)
	if langRes, ok := langCtx.(*Lang); ok {
		return langRes
	}
	return nil
}
