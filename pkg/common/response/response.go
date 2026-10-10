package response

import "net/http"

// ErrorVm represents standard RFC7807/Spring style error view model
type ErrorVm struct {
	StatusCode string `json:"statusCode"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
}

func NotFound(w http.ResponseWriter, r *http.Request, detail string) {
	WriteJSON(w, http.StatusNotFound, ErrorVm{
		StatusCode: "404",
		Title:      "Not found",
		Detail:     detail,
	})
}

func BadRequest(w http.ResponseWriter, r *http.Request, detail string) {
	WriteJSON(w, http.StatusBadRequest, ErrorVm{
		StatusCode: "400",
		Title:      "Bad request",
		Detail:     detail,
	})
}

func InternalError(w http.ResponseWriter, r *http.Request, detail string) {
	WriteJSON(w, http.StatusInternalServerError, ErrorVm{
		StatusCode: "500",
		Title:      "Internal server error",
		Detail:     detail,
	})
}
