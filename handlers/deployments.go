package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/domain/deployments"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/kubService"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/libhttp"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/models"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

const maxValueScale = 5

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
	if kv.ValInt > maxValueScale {
		pgs.AddDangerText(fmt.Sprintf("Max value for Scale is %v", maxValueScale))
		return kv, false
	}
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
		logger.Errorf("Error GetDeployments %v ", err.Error())
		kubService.ReloadAuth()
	}

	if kv, ok := GetKeyValue(r).validate(pgs); ok == true && kv.hasValue {
		db := r.Context().Value("db").(*sqlx.DB)
		h := models.NewHistoryUserActions(db, currentUser)

		dresp, err := serviceKubernetes.ScaleBy(kv.Key, kv.ValInt)
		dreq := fmt.Sprintf("%v=%v", kv.Key, kv.ValInt)
		if err != nil {
			pgs.AddDangerText(err.Error())
			err := h.SaveActionScale(false, err.Error(), dreq)
			if err != nil {
				pgs.AddDangerText("db hist err:" + err.Error())
			}
		} else {
			resultText := fmt.Sprintf("actual: %v, desired: %v", dresp.ActualReplicas, dresp.DesiredReplicas)
			pgs.AddSuccesText(fmt.Sprintf("%v set scale to %v", kv.Key, kv.Value))
			pgs.AddSuccesText(resultText)
			err := h.SaveActionScale(true, resultText, dreq)
			if err != nil {
				pgs.AddDangerText("db hist err:" + err.Error())
			}
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
