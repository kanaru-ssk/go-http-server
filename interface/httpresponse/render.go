package httpresponse

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func RenderJson(ctx context.Context, w http.ResponseWriter, statusCode int, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("Content-Length", strconv.Itoa(len(b)))
	w.WriteHeader(statusCode)
	if _, err := w.Write(b); err != nil {
		slog.ErrorContext(ctx, "httpresponse.RenderJson: http.ResponseWriter.Write:", "err", err)
	}
}
