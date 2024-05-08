// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package middleware

import (
	"context"
	"time"

	logCtx "github.com/dedyf5/resik/ctx/log"
	"github.com/dedyf5/resik/entities/config"
	"github.com/rs/xid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Interceptor struct {
	appModule config.Module
	log       *logCtx.Log
}

func NewInterceptor(appModule config.Module, log *logCtx.Log) *Interceptor {
	return &Interceptor{
		appModule: appModule,
		log:       log,
	}
}

func (i *Interceptor) logCtx(ctx context.Context) (*context.Context, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "Unable to read metadata")
	}

	correlationID := xid.New().String()
	newMeta := meta.Copy()
	newMeta.Append(logCtx.CorrelationIDKeyContext.String(), correlationID)
	newMeta.Append(logCtx.CorrelationIDKeyXHeader.String(), correlationID)
	i.log.CorrelationID = correlationID
	newCtx := metadata.NewIncomingContext(ctx, newMeta)

	return &newCtx, nil
}

func (i *Interceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	newCtx, err := i.logCtx(ctx)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	resHandler, err := handler(*newCtx, req)
	logger := logCtx.NewGRPC(i.appModule, i.log, start, info.FullMethod, req, resHandler, err)
	logger.Write()
	return resHandler, err
}
