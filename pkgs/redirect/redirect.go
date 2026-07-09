package redirect

import (
	"net/http"
	"strings"
)

func SetRedirectToLoginWithError(w http.ResponseWriter, r *http.Request, errorMessage string) {
	errorParams := strings.Split(errorMessage, " ")
	errorMessageNormilized := strings.Join(errorParams, "+")
	http.Redirect(w, r, "/login?error="+errorMessageNormilized, http.StatusTemporaryRedirect)
}
