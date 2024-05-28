package navitas

import (
	"net/http"
)

func (n *Navitas) SessionLoad(next http.Handler) http.Handler {
	n.InfoLog.Println("SessionLoad called")
	return n.Session.LoadAndSave(next)
}
