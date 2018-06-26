package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"git.betfavorit.cf/vadim.tsurkov/kuberweb/libhttp"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/models"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/redisService"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

func PageRedis(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := r.Context().Value("sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "kuberweb-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)

	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	pgs := NewPageSysContent()
	redisClient := r.Context().Value("RedisCache").(*redisService.RedisCache)

	r.ParseForm()
	if r.Method == "POST" {
		db := r.Context().Value("db").(*sqlx.DB)
		h := models.NewHistoryUserActions(db, currentUser)

		rls := redisClient.FlushAll().RedisInfos
		for _, v := range rls {
			if v.Err != "" {
				pgs.AddDangerText(fmt.Sprintf("%v has error %v", v.Names, v.Err))
				err := h.SaveActionRedisFlush(false, v.Err, v.Names)
				if err != nil {
					pgs.AddDangerText("db hist err:" + err.Error())
				}
			} else {
				pgs.AddSuccesText(fmt.Sprintf("cache flush result %v in %v", v.Res, v.Names))
				err := h.SaveActionRedisFlush(true, v.Res, v.Names)
				if err != nil {
					pgs.AddDangerText("db hist err:" + err.Error())
				}
			}
		}
	}

	rinfos := redisClient.GetRedisInfo()

	data := struct {
		CurrentUser    *models.UserRow
		PageSysContent *PageSysContent
		RedisInfos     *redisService.RedisInfos
	}{
		currentUser,
		pgs,
		rinfos,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/redis/redis.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}
