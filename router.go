package lib

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"https://github.com/IminaDanzi/Docker-Project/logger"
	"https://github.com/IminaDanzi/Docker-Project/rqctx"

	"go.uber.org/zap"
)

// NewRouter returns a new instance of router
func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.NotFoundHandler = notFound
	r.MethodNotAllowedHandler = notFound

	r.HandleFunc("/",
		contextify(postToDatabaseHandler)).
		Methods("POST")

	return r
}

func contextify(h func(ctx *rqctx.Context) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := rqctx.NewContext(w, r, logger.DefaultLogger)

		// open DB or log, fail and return
		mainCtx := context.Background()
		err := ctx.OpenDB(mainCtx)
		if err != nil {
			ctx.Logger.Error("failed to open DB", zap.Error(err))
			marshalAndWrite(ctx.Writer, "Internal server error", http.StatusInternalServerError)
			return
		}

		// panic recovery block
		defer recoverPanic(ctx)

		// Execute the request with context and send the response according to whether the handling errored
		err = h(ctx)

		switch {
		case err != nil:
			ctx.Logger.Error("Internal server error", zap.String("uuid", ctx.UUID.String()))
			if err := ctx.Rollback(); err != nil {
				ctx.Logger.Error("rollback failed", zap.Error(err))
			}
			if err := marshalAndWrite(ctx.Writer, err, http.StatusInternalServerError); err != nil {
				ctx.Logger.Error("sending error failed", zap.Error(err))
			}

		default:
			ctx.Logger.Info("200 OK", zap.String("uuid", ctx.UUID.String()))
			if err := ctx.Commit(); err != nil {
				ctx.Logger.Error("commit failed", zap.Error(err))
			}
			if err := marshalAndWrite(ctx.Writer, "OK", http.StatusOK); err != nil {
				ctx.Logger.Error("sending OK failed", zap.Error(err))
			}
		}
	}
}

var recoverPanic = func(ctx *rqctx.Context) {
	err := recover()
	if err != nil {
		e, ok := err.(error)
		if !ok {
			err = fmt.Errorf("%#v\n%s", err, string(debug.Stack()))
		}
		ctx.Logger.Error("Recovered handler panic", zap.Error(e))

		//This will let you see stack traces from the recover in the log
		debug.PrintStack()
		rErr := ctx.Rollback()
		if err != nil {
			ctx.Logger.Error("failed to rollback error", zap.Error(rErr))
		}
	}
}

var notFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
	_, _ = w.Write([]byte(`{"message": "Not found"}`))
})
