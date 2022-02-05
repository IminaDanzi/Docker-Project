package rqctx

import (
	"context"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"https://github.com/IminaDanzi/Docker-Project/dat"
	"go.uber.org/zap"
)

// Context contains data specific to a single request
type Context struct {
	*http.Request
	Logger      *zap.Logger
	Start       time.Time
	Writer      http.ResponseWriter
	UUID        uuid.UUID
	dat.Context // data context is anonymous for mapper accessing convenience
}

// NewContext returns a new context for a single request
func NewContext(w http.ResponseWriter, r *http.Request, logger *zap.Logger) *Context {
	requestUUID, err := uuid.NewV1()
	if err != nil {
		logger.Error("could not generate UUID for request", zap.Error(err))
	}

	w.Header().Set("X-Request-ID", requestUUID.String())

	reqCtx := Context{
		Start:   time.Now(),
		Logger:  logger.With(zap.String("requestID", requestUUID.String())),
		Request: r,
		Writer:  w,
		UUID:    requestUUID,
	}
	return &reqCtx
}

//OpenDB sets the database context
func (ctx *Context) OpenDB(mainContext context.Context) error {
	dataContext, err := dat.NewDataContext(mainContext)
	if err != nil {
		return err
	}
	ctx.Context = *dataContext
	return nil
}

// Rollback reverts any changes done within the transaction
func (ctx *Context) Rollback() error {
	return ctx.Context.Rollback()
}

// Commit applies any changes done within the transaction
func (ctx *Context) Commit() error {
	return ctx.Context.Commit()
}
