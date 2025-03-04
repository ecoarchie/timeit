package httpv1

import (
	"fmt"
	"net/http"
)

func errorResponse(w http.ResponseWriter, status int, message any) {
	mes := map[string]any{"error": message}
	err := writeJSON(w, status, mes, nil)
	if err != nil {
		w.WriteHeader(500)
	}
}

func serverErrorResponse(w http.ResponseWriter, err error) {
	message := fmt.Errorf("server error: %w", err)
	errorResponse(w, http.StatusInternalServerError, message.Error())
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, http.StatusNotFound, message)
}

func methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, http.StatusMethodNotAllowed, message)
}

func failedValidationResponse(w http.ResponseWriter, errors map[string]string) {
	errorResponse(w, http.StatusUnprocessableEntity, errors)
}
