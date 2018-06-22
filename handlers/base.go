// Package handlers provides request handlers.
package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"git.betfavorit.cf/vadim.tsurkov/kuberweb/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type PageSysContent struct {
	AlertDanger struct {
		ClassHidden string
		TextAlert   []string
	}
	AlertSuccess struct {
		ClassHidden string
		TextAlert   []string
	}
}

func NewPageSysContent() *PageSysContent {
	p := &PageSysContent{}
	p.AlertDanger.ClassHidden = "hidden"
	p.AlertSuccess.ClassHidden = "hidden"
	return p
}

func (psc *PageSysContent) AddDangerText(text string) {
	psc.AlertDanger.TextAlert = append(psc.AlertDanger.TextAlert, text)
	psc.AlertDanger.ClassHidden = ""
}

func (psc *PageSysContent) AddSuccesText(text string) {
	psc.AlertSuccess.TextAlert = append(psc.AlertSuccess.TextAlert, text)
	psc.AlertSuccess.ClassHidden = ""
}

func getCurrentUser(w http.ResponseWriter, r *http.Request) *models.UserRow {
	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "kuberweb-session")
	return session.Values["user"].(*models.UserRow)
}

func getIdFromPath(w http.ResponseWriter, r *http.Request) (int64, error) {
	userIdString := mux.Vars(r)["id"]
	if userIdString == "" {
		return -1, errors.New("user id cannot be empty.")
	}

	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return -1, err
	}

	return userId, nil
}
