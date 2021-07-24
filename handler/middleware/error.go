package middleware

import "net/http"

func serverError(w http.ResponseWriter) {
	http.Error(w, "internal server error", 500)
}

func unAuthorizeError(w http.ResponseWriter) {
	http.Error(w, "unauthorized", 401)
}

func forbiddenError(w http.ResponseWriter) {
	http.Error(w, "forbidden", 403)
}

func denyError(w http.ResponseWriter) {
	http.Error(w, "already denied", 403)
}

func banError(w http.ResponseWriter) {
	http.Error(w, "already ban", 403)
}

func deleteError(w http.ResponseWriter) {
	http.Error(w, "already deleted", 404)
}

func invalidRoleError(w http.ResponseWriter) {
	http.Error(w, "invalid user role", 403)
}

func maintenanceError(w http.ResponseWriter) {
	http.Error(w, "maintenance", 503)
}
