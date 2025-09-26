package app

import (
	"backend/pkg/logger"
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func (app *App) addMiddleware() {
	app.Router.Use(middleware.RequestID)
	app.Router.Use(middleware.RealIP)
	app.Router.Use(Logger(app.Logger))
	app.Router.Use(middleware.Recoverer)
}

type httpRequest struct {
	protocol     string
	path         string
	method       string
	duration     time.Duration
	status       int
	bytesWritten int
	reqeustID    string
	scheme       string
	user         string
	body         string
}

func Logger(l logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			var requestBody string
			if r.Header.Get("Content-Type") == "application/json" {
				defer r.Body.Close()
				body, err := io.ReadAll(r.Body)
				if err != nil {
					l.Error("Error reading request body", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				r.Body = io.NopCloser(bytes.NewBuffer(body))

				requestBody = string(body)
			}

			t1 := time.Now()
			defer func() {
				request := httpRequest{
					protocol:     r.Proto,
					path:         r.URL.Path,
					method:       r.Method,
					duration:     time.Since(t1),
					status:       ww.Status(),
					bytesWritten: ww.BytesWritten(),
					reqeustID:    middleware.GetReqID(r.Context()),
					scheme:       r.URL.Scheme,
					user:         r.URL.User.String(),
					body:         requestBody,
				}

				l.Info("Served", request)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
