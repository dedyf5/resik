// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package middleware provides gRPC unary server interceptors for logging, language
// resolution, JWT authentication, error formatting, and rate limiting.
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
	"github.com/dedyf5/resik/ctx/lang/term"
	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/entities/config"
	"github.com/dedyf5/resik/internal/identity"
	resPkg "github.com/dedyf5/resik/pkg/response"
	"github.com/dedyf5/resik/utils/ratelimit"
	statusUtil "github.com/dedyf5/resik/utils/status"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Role identifies an authorization requirement for a gRPC method.
type Role string

const (
	// RoleValidToken indicates that the gRPC method requires a valid JWT bearer token.
	RoleValidToken Role = "ValidToken"

	// healthService is the fully qualified gRPC service path for health checks.
	healthService string = "/health.HealthService"
	// HealthzGetPath is the fully gRPC method path for healthz readiness/liveness endpoint.
	HealthzGetPath string = healthService + "/HealthzGet"
)

// Interceptor provides gRPC unary unary interceptors for request tracking, localization,
// authentication, structured logging, and rate limiting.
type Interceptor struct {
	module      config.Module
	auth        config.Auth
	resolver    identity.IdentityResolver
	limiter     ratelimit.Limiter
	log         *logCtx.Log
	methodRoles map[string][]Role
}

// NewInterceptor constructs a new gRPC Interceptor configured with module settings,
// authentication config, identity resolver, rate limiter, and base logger.
func NewInterceptor(module config.Module, auth config.Auth, resolver identity.IdentityResolver, limiter ratelimit.Limiter, log *logCtx.Log) *Interceptor {
	return &Interceptor{
		module:      module,
		auth:        auth,
		resolver:    resolver,
		limiter:     limiter,
		log:         log,
		methodRoles: methodRoles(),
	}
}

// methodRoles returns a map of fully qualified gRPC method names to their required roles.
func methodRoles() map[string][]Role {
	const organizationService = "/organization.OrganizationService/"
	const transactionService = "/transaction.TransactionService/"
	const userService = "/user.UserService/"
	return map[string][]Role{
		organizationService + "OrganizationDetailGet": {RoleValidToken},
		organizationService + "OrganizationListGet":   {RoleValidToken},
		transactionService + "TransactionPost":        {RoleValidToken},
		transactionService + "TransactionPut":         {RoleValidToken},
		transactionService + "TransactionDelete":      {RoleValidToken},
		transactionService + "TransactionDetailGet":   {RoleValidToken},
		transactionService + "TransactionListGet":     {RoleValidToken},
		transactionService + "OrganizationOmzetGet":   {RoleValidToken},
		transactionService + "BranchOmzetGet":         {RoleValidToken},
		userService + "TokenRefreshGet":               {RoleValidToken},
	}
}

// logCtx ensures a correlation ID exists on c and populates i.log with the correlation ID and request path.
func (i *Interceptor) logCtx(c context.Context, path string) (*context.Context, error) {
	correlationID, newCtx := logCtx.EnsureCorrelationID(c)

	i.log.CorrelationID = correlationID
	i.log.Path = path

	return &newCtx, nil
}

// writeLogger records a gRPC audit/request log entry and converts err into a detailed gRPC error.
func (i *Interceptor) writeLogger(c context.Context, start time.Time, fullMethod string, req any, res any, err error) (any, error) {
	logger := logCtx.NewGRPC(c, i.log, start, fullMethod, req, res, err)
	logger.Write()
	return res, statusUtil.GRPCErrorWithDetails(err, i.log)
}

// langCtx reads incoming gRPC metadata to resolve the requested language, attaches a localized
// langCtx.Lang instance to the context, and returns the updated context and language config.
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

// validateJWT enfores RoleValidToken when configured for fullMethod by parsing the Bearer token
// from incoming gRPC metadata and storing the resulting claims in the returned context.
func (i *Interceptor) validateJWT(c context.Context, lang *langCtx.Lang, signatureKey string, fullMethod string) (*context.Context, error) {
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
	claim, err := jwtCxt.AuthClaimsFromString(token, signatureKey, c, i.resolver, lang)
	if err != nil {
		return nil, err
	}
	newCtx := context.WithValue(c, jwtCxt.AuthClaimsKey, claim)

	return &newCtx, nil
}

// errorMetadata returns a nil context and an internal server error status when gRPC metadata cannot be read.
func (i *Interceptor) errorMetadata() (*context.Context, error) {
	return nil, resPkg.NewStatusMessage(http.StatusInternalServerError, "unable to read metadata", nil)
}

// Unary is a gRPC unary server interceptor that initializes tracking and caller context,
// resolves language preferences, validates JWT credentials, invokes handler, and logs the outcome.
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

	return i.writeLogger(*tokenC, start, info.FullMethod, req, res, err)
}

// RateLimit is a gRPC unary server interceptor that applies rate limiting per request key,
// sets rate-limit response headers, and rejects excess requests with HTTP 429 status details.
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
			term.TooManyRequests.Localize(lang.Localizer),
			nil,
		)
	}

	if err := grpc.SetHeader(c, md); err != nil {
		return nil, resPkg.NewStatusError(http.StatusInternalServerError, err)
	}

	return handler(c, req)
}
