// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package middleware

import (
	"context"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	jwtCxt "github.com/dedyf5/resik/ctx/jwt"
	langCtx "github.com/dedyf5/resik/ctx/lang"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/entities/config"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/ratelimit"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Role string

const (
	RoleValidToken Role = "ValidToken"

	healthService  string = "/health.HealthService"
	HealthzGetPath string = healthService + "/HealthzGet"
)

type Interceptor struct {
	module      config.Module
	auth        config.Auth
	limiter     ratelimit.Limiter
	log         *logCtx.Log
	methodRoles map[string][]Role
}

func NewInterceptor(module config.Module, auth config.Auth, limiter ratelimit.Limiter, log *logCtx.Log) *Interceptor {
	return &Interceptor{
		module:      module,
		auth:        auth,
		limiter:     limiter,
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

func (i *Interceptor) logCtx(c context.Context, path string) (*context.Context, error) {
	correlationID, newCtx := logCtx.EnsureCorrelationID(c)

	i.log.CorrelationID = correlationID
	i.log.Path = path

	return &newCtx, nil
}

func (i *Interceptor) writeLogger(c context.Context, start time.Time, fullMethod string, req any, res any, err error) (any, error) {
	logger := logCtx.NewGRPC(c, i.log, start, fullMethod, req, res, err)
	logger.Write()
	return res, err
}

func (i *Interceptor) langCtx(c context.Context, langDefault language.Tag) (*context.Context, *langCtx.Lang, error) {
	meta, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, nil, resPkg.NewStatusMessage(http.StatusInternalServerError, "unable to read metadata", nil)
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
	if !slices.Contains(i.methodRoles[fullMethod], RoleValidToken) {
		return &c, nil
	}
	if len(meta["authorization"]) != 1 {
		return nil, resPkg.NewStatusCode(http.StatusUnauthorized)
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
	return nil, resPkg.NewStatusMessage(http.StatusInternalServerError, "unable to read metadata", nil)
}

func (i *Interceptor) Unary(c context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	callerHolder := &logCtx.CallerHolder{}

	c = context.WithValue(c, logCtx.KeyCallerHolderContext, callerHolder)

	logC, err := i.logCtx(c, info.FullMethod)
	if err != nil {
		return handler, err
	}

	langC, langRes, err := i.langCtx(*logC, i.module.LangDefault)
	if err != nil {
		return i.writeLogger(c, start, info.FullMethod, req, nil, err)
	}

	tokenC, err := i.validateJWT(*langC, langRes, i.auth.SignatureKey, info.FullMethod)
	if err != nil {
		return i.writeLogger(c, start, info.FullMethod, req, nil, err)
	}

	res, err := handler(*tokenC, req)

	return i.writeLogger(c, start, info.FullMethod, req, res, err)
}

func (i *Interceptor) RateLimit(c context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if info.FullMethod == HealthzGetPath {
		return handler(c, req)
	}

	key, err := i.limiter.GetKeyGRPC(c)
	if err != nil {
		return nil, err
	}

	res, err := i.limiter.Take(c, key)
	if err != nil {
		return nil, err
	}

	md := metadata.Pairs(
		"x-ratelimit-limit", strconv.FormatInt(res.Limit, 10),
		"x-ratelimit-remaining", strconv.FormatInt(res.Remaining, 10),
		"x-ratelimit-reset", strconv.FormatInt(res.Reset, 10),
	)

	lang, err := langCtx.FromContext(c)
	if err != nil {
		return nil, err
	}

	if res.Reached {
		md.Append("retry-after", strconv.FormatInt(res.RetryAfter, 10))
		if err := grpc.SetHeader(c, md); err != nil {
			return nil, resPkg.NewStatusError(http.StatusInternalServerError, err)
		}

		return nil, resPkg.NewStatusMessage(
			http.StatusTooManyRequests,
			lang.GetByMessageID("too_many_requests"),
			nil,
		)
	}

	if err := grpc.SetHeader(c, md); err != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, err)
	}

	return handler(c, req)
}
