package navitas

import (
	"net/http"
	"strconv"

	"github.com/justinas/nosurf"
)

func (n *Navitas) SessionLoad(next http.Handler) http.Handler {
	n.InfoLog.Println("SessionLoad called")
	return n.Session.LoadAndSave(next)
}

func (n *Navitas) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(n.config.cookie.secure)

	csrfHandler.ExemptGlob("/api/*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   n.config.cookie.domain,
	})

	return csrfHandler
}
