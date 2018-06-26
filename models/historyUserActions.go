package models

import (
	"fmt"

	"git.betfavorit.cf/backend/logger"

	"github.com/jmoiron/sqlx"
)

type HistoryActions struct {
	Base Base
	User *UserRow
}

const (
	ActionSetScale   uint = 1
	ActionFlushRedis uint = 2
)

func NewHistoryUserActions(db *sqlx.DB, user *UserRow) *HistoryActions {
	return &HistoryActions{
		User: user,
		Base: Base{
			db:    db,
			table: "history_user_actions",
			hasID: true,
		},
	}
}

func (hs *HistoryActions) SaveActionScale(isDone bool, result string, reqparam string) error {
	return hs.insertAction(nil, isDone, ActionSetScale, result, reqparam)
}

func (hs *HistoryActions) SaveActionRedisFlush(isDone bool, result string, reqparam string) error {
	return hs.insertAction(nil, isDone, ActionFlushRedis, result, reqparam)
}

func (hs *HistoryActions) insertAction(tx *sqlx.Tx, isDone bool, actionId uint, result string, reqparam string) error {
	data := make(map[string]interface{})
	data["action_id"] = fmt.Sprintf("%v", actionId)
	data["result"] = result
	data["user_id"] = hs.User.ID
	data["isdone"] = fmt.Sprintf("%v", isDone)
	data["reqparam"] = fmt.Sprintf("%v", reqparam)

	_, err := hs.Base.InsertIntoTable(tx, data)

	if err != nil {
		logger.Errorf("History user action has error:%v", err.Error())
		logger.Debugf("UserId %v cannot save history action by request: %v, error: %v, resp:%v", hs.User.ID, reqparam, err.Error(), result)
	}
	return err
}
