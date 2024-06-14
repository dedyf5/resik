// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dedyf5/resik/ctx"
	jwtCxt "github.com/dedyf5/resik/ctx/jwt"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/entities/config"
	"github.com/dedyf5/resik/pkg/array"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/rs/xid"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
	const merchantService = "/merchant.MerchantService/"
	const transactionService = "/transaction.TransactionService/"
	const userService = "/user.UserService/"
	return map[string][]Role{
		merchantService + "MerchantPost":        {RoleValidToken},
		merchantService + "MerchantPut":         {RoleValidToken},
		merchantService + "MerchantDelete":      {RoleValidToken},
		merchantService + "MerchantListGet":     {RoleValidToken},
		transactionService + "MerchantOmzetGet": {RoleValidToken},
		transactionService + "OutletOmzetGet":   {RoleValidToken},
		userService + "TokenRefreshGet":         {RoleValidToken},
	}
}

func (i *Interceptor) logCtx(c context.Context) (*context.Context, error) {
	meta, ok := metadata.FromIncomingContext(c)
	if !ok {
		return i.errorMetadata()
	}

	correlationID := xid.New().String()
	newMeta := meta.Copy()
	newMeta.Append(logCtx.CorrelationIDKeyContext.String(), correlationID)
	newMeta.Append(logCtx.CorrelationIDKeyXHeader.String(), correlationID)
	i.log.CorrelationID = correlationID
	newCtx := metadata.NewIncomingContext(c, newMeta)

	return &newCtx, nil
}

func (i *Interceptor) writeLogger(start time.Time, fullMethod string, req any, res any, err error) (any, error) {
	logger := logCtx.NewGRPC(i.app.Module, i.log, start, fullMethod, req, res, err)
	logger.Write()
	return res, err
}

func (i *Interceptor) langCtx(c context.Context, langDefault language.Tag) (*context.Context, *langCtx.Lang, error) {
	meta, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, nil, &resPkg.Status{
			Code:    http.StatusInternalServerError,
			Message: "Unable to read metadata",
		}
	}
	var langReq *language.Tag = nil
	langKey := langCtx.ContextKey.String()
	if len(meta[langKey]) == 1 {
		if langString := meta[langKey][0]; langString != "" {
			langRes, err := langCtx.GetLanguageAvailable(langString)
			if err != nil {
				return nil, nil, err
			}
			langReq = langRes
		}
	}
	langRes := langCtx.NewLang(langDefault, langReq, "")
	newCtx := context.WithValue(c, langCtx.ContextKey, langRes)

	return &newCtx, langRes, nil
}

func (i *Interceptor) validateJWT(c context.Context, lang *langCtx.Lang, signatureKey, fullMethod string) (*context.Context, error) {
	meta, ok := metadata.FromIncomingContext(c)
	if !ok {
		return i.errorMetadata()
	}
	if i.methodRoles[fullMethod] == nil {
		return &c, nil
	}
	if array.InArray(RoleValidToken, i.methodRoles[fullMethod]) < 0 {
		return &c, nil
	}
	if len(meta["authorization"]) != 1 {
		return nil, &resPkg.Status{
			Code: http.StatusUnauthorized,
		}
	}
	bearer := meta["authorization"][0]
	token := strings.ReplaceAll(bearer, "Bearer ", "")
	claim, err := jwtCxt.AuthClaimsFromString(token, signatureKey, lang)
	if err != nil {
		return nil, err
	}
	newCtx := context.WithValue(c, jwtCxt.AuthClaimsKey, claim)

	return &newCtx, nil
}

func (i *Interceptor) errorMetadata() (*context.Context, error) {
	return nil, &resPkg.Status{
		Code:    http.StatusInternalServerError,
		Message: "Unable to read metadata",
	}
}

func (i *Interceptor) Unary(c context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	methodCtx := context.WithValue(c, ctx.KeyFullMethod, info.FullMethod)

	logC, err := i.logCtx(methodCtx)
	if err != nil {
		return handler, err
	}

	langC, langRes, err := i.langCtx(*logC, i.app.LangDefault)
	if err != nil {
		return i.writeLogger(start, info.FullMethod, req, nil, err)
	}

	tokenC, err := i.validateJWT(*langC, langRes, i.auth.SignatureKey, info.FullMethod)
	if err != nil {
		return i.writeLogger(start, info.FullMethod, req, nil, err)
	}

	res, err := handler(*tokenC, req)

	return i.writeLogger(start, info.FullMethod, req, res, err)
}
