package redirect

import (
	"net/http"
	"strings"
)

func ToLoginWithError(w http.ResponseWriter, r *http.Request, errorMessage string) {
	errorParams := strings.Split(errorMessage, " ")
	errorMessageNormilized := strings.Join(errorParams, "+")
	http.Redirect(w, r, "/auth/login?error="+errorMessageNormilized, http.StatusTemporaryRedirect)
}

func ToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
}
