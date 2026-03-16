package responses

import (
	"net/http"
	"simple-rest-api/internal/handlers/types"
	"simple-rest-api/pkg/http/render"
)

func Error(w http.ResponseWriter, r *http.Request, status int, msg string) {
	render.Encode(w, r, status, types.ErrorResponse{Error: msg})
}

func ValidationError(w http.ResponseWriter, r *http.Request, problems map[string]string) {
	render.Encode(w, r, http.StatusBadRequest, map[string]any{
		"error":    "validation failed",
		"problems": problems,
	})
}
