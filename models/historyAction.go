package models

import "github.com/jmoiron/sqlx"

type HistoryActionRow struct {
	ID     int64  `db:"id"`
	Action string `db:"action"`
}

type HistoryAction struct {
	Base
}

func NewHistoryAction(db *sqlx.DB) *HistoryActions {
	return &HistoryActions{
		//Base{
		//	db:    db,
		//	table: "history_action",
		//	hasID: true,
		//},
	}
}
