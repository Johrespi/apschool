package middleware

import (
	"apschool/internal/ctxkeys"
	"apschool/internal/response"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/httprate"
)

func RateLimitByIP(requestLimit int, windowLength time.Duration) func(http.Handler) http.Handler {
	return httprate.Limit(requestLimit, windowLength, httprate.WithKeyFuncs(httprate.KeyByIP),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			response.TooManyRequests(w)
		}),
	)
}

func RateLimitByUserID(requestLimit int, windowLength time.Duration) func(http.Handler) http.Handler {
	return httprate.Limit(requestLimit, windowLength, httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
		userID, ok := ctxkeys.GetUserID(r.Context())
		if !ok {
			return httprate.KeyByIP(r)
		}
		return strconv.Itoa(userID), nil
	}),
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			response.TooManyRequests(w)
		}),
	)
}
