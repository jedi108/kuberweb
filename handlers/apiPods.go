package handlers

import (
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/libhttp"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/models"

	"html/template"
	"net/http"

	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/domain/pods"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/kubService"
	"github.com/gorilla/sessions"
)

func PagePods(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "kuberweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)

	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	pgs := NewPageSysContent()
	serviceKubernetes := r.Context().Value("serviceKubernetes").(*kubService.ServiceKubernetes)

	apiPod, err := serviceKubernetes.GetPods()
	if err != nil {
		pgs.AddDangerText(err.Error())
	}

	data := struct {
		CurrentUser    *models.UserRow
		ApiPods        *pods.ApiPod
		PageSysContent *PageSysContent
	}{
		currentUser,
		apiPod,
		pgs,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/kubernetes/pods.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}
