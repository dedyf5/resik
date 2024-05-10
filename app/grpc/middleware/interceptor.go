// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/dedyf5/resik/ctx"
	jwtCxt "github.com/dedyf5/resik/ctx/jwt"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/entities/config"
	"github.com/dedyf5/resik/pkg/array"
	"github.com/rs/xid"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Role string

const (
	RoleValidToken Role = "ValidToken"
)

type Interceptor struct {
	app         config.App
	auth        config.Auth
	log         *logCtx.Log
	methodRoles map[string][]Role
}

func NewInterceptor(app config.App, auth config.Auth, log *logCtx.Log) *Interceptor {
	return &Interceptor{
		app:         app,
		auth:        auth,
		log:         log,
		methodRoles: methodRoles(),
	}
}

func methodRoles() map[string][]Role {
	const transactionService = "/transaction.TransactionService/"
	return map[string][]Role{
		transactionService + "MerchantOmzetGet": {RoleValidToken},
	}
}

func (i *Interceptor) logCtx(c context.Context) (*context.Context, error) {
	meta, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, status.Error(codes.Internal, "Unable to read metadata")
	}

	correlationID := xid.New().String()
	newMeta := meta.Copy()
	newMeta.Append(logCtx.CorrelationIDKeyContext.String(), correlationID)
	newMeta.Append(logCtx.CorrelationIDKeyXHeader.String(), correlationID)
	i.log.CorrelationID = correlationID
	newCtx := metadata.NewIncomingContext(c, newMeta)

	return &newCtx, nil
}

func (i *Interceptor) langCtx(c context.Context, langDefault language.Tag) (*context.Context, error) {
	meta, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, status.Error(codes.Internal, "Unable to read metadata")
	}
	var langReq *language.Tag = nil
	langKey := langCtx.ContextKey.String()
	if len(meta[langKey]) == 1 {
		if langString := meta[langKey][0]; langString != "" {
			langRes, err := langCtx.GetLanguageAvailable(langString)
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, err.MessageOrDefault())
			}
			langReq = langRes
		}
	}
	newCtx := context.WithValue(c, langCtx.ContextKey, langCtx.NewLang(langDefault, langReq, ""))

	return &newCtx, nil
}

func (i *Interceptor) validateJWT(c context.Context, signatureKey, fullMethod string) (*context.Context, error) {
	meta, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, status.Error(codes.Internal, "Unable to read metadata")
	}
	if i.methodRoles[fullMethod] == nil {
		return &c, nil
	}
	if array.InArray(RoleValidToken, i.methodRoles[fullMethod]) < 0 {
		return &c, nil
	}
	if len(meta["authorization"]) != 1 {
		return nil, status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}
	bearer := meta["authorization"][0]
	token := strings.ReplaceAll(bearer, "Bearer ", "")
	claim, err := jwtCxt.AuthClaimsFromString(token, signatureKey)
	if err != nil {
		return nil, err.GRPC().Err()
	}
	newCtx := context.WithValue(c, jwtCxt.AuthClaimsKey, claim)

	return &newCtx, nil
}

func (i *Interceptor) Unary(c context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	langC, err := i.langCtx(c, i.app.LangDefault)
	if err != nil {
		return nil, err
	}
	logC, err := i.logCtx(*langC)
	if err != nil {
		return nil, err
	}
	tokenC, err := i.validateJWT(*logC, i.auth.SignatureKey, info.FullMethod)
	if err != nil {
		return nil, err
	}
	newCtx := context.WithValue(*tokenC, ctx.KeyFullMethod, info.FullMethod)
	start := time.Now()
	resHandler, err := handler(newCtx, req)
	logger := logCtx.NewGRPC(i.app.Module, i.log, start, info.FullMethod, req, resHandler, err)
	logger.Write()
	return resHandler, err
}
