package middleware

import (
	"context"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/utils"
)

type ctxCorrelationID string

const CorrelationIDKey ctxCorrelationID = "correlation_id"
const CorrelationIDHeader = "X-Correlation-ID"

func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(CorrelationIDHeader)
		if id == "" {
			id = utils.NewID()
		}
		w.Header().Set(CorrelationIDHeader, id)
		ctx := context.WithValue(r.Context(), CorrelationIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetCorrelationID(ctx context.Context) string {
	id, ok := ctx.Value(CorrelationIDKey).(string)
	if !ok {
		return ""
	}
	return id
}
