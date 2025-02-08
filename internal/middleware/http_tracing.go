package middleware

import (
	"net/http"

	"github.com/bravo1goingdark/jaegaurd/internal/tracing"
)

func HTTPTracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tr := tracing.Tracer()
		ctx, span := tr.Start(r.Context(), r.URL.Path)
		defer span.End()

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
