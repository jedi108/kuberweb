package handlers

import (
	"html/template"
	"net/http"

	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/domain/deployments"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/kubService"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/libhttp"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/models"

	"strings"

	"strconv"

	"fmt"

	"log"

	"github.com/gorilla/sessions"
)

type ErrorsForm struct {
	Errors []string
}

type KeyValue struct {
	hasValue bool
	r        *http.Request
	Key      string
	Value    string
	ValInt   uint64
}

func (kv *KeyValue) validate(pgs *PageSysContent) (*KeyValue, bool) {
	if kv.hasValue == false {
		return kv, true
	}
	i, err := strconv.ParseUint(kv.Value, 10, 8)
	if err != nil {
		pgs.AddDangerText(err.Error())
		return kv, false
	}
	kv.ValInt = i
	return kv, true
}

func GetKeyValue(r *http.Request) *KeyValue {
	keyValue := &KeyValue{}
	r.ParseForm()
	if len(r.Form) != 0 {
		for k, v := range r.Form {
			value := strings.Join(v, "")
			if value != "" {
				keyValue.hasValue = true
				keyValue.Key = k
				keyValue.Value = value
				break
			}
		}
	}
	return keyValue
}

func PageDeployments(w http.ResponseWriter, r *http.Request) {
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

	apiDeps, err := serviceKubernetes.GetDeployments()
	if err != nil {
		pgs.AddDangerText(err.Error())
		log.Printf("Error GetDeployments %v ", err.Error())
	}

	if kv, ok := GetKeyValue(r).validate(pgs); ok == true && kv.hasValue {
		dresp, err := serviceKubernetes.ScaleBy(kv.Key, kv.ValInt)
		if err != nil {
			pgs.AddDangerText(err.Error())
		} else {
			pgs.AddSuccesText(fmt.Sprintf("%v set scale to %v", kv.Key, kv.Value))
			pgs.AddSuccesText(fmt.Sprintf("actual replicas: %v, desired replicas: %v", dresp.ActualReplicas, dresp.DesiredReplicas))
		}
	}

	data := struct {
		CurrentUser    *models.UserRow
		Deps           *deployments.Deployments
		PageSysContent *PageSysContent
	}{
		currentUser,
		apiDeps,
		pgs,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/kubernetes/deployments.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}
