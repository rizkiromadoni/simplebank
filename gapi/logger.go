package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	c context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (res interface{}, err error) {
	startTime := time.Now()
	res, err = handler(c, req)
	elapsedTime := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.
		Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("elapsed_time", elapsedTime).
		Msg("GRPC")

	return
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Write(body []byte) (int, error) {
	r.Body = body
	return r.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, r)
		elapsedTime := time.Since(startTime)

		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.
			Str("protocol", "http").
			Str("method", r.Method).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("elapsed_time", elapsedTime).
			Msg("HTTP")
	})
}
