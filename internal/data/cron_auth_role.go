package data

import (
	"context"
	"cron/internal/basic/db"
)

type CronAuthRoleData struct {
	ctx       context.Context
	db        *db.MyDB
	tableName string
}

func NewCronAuthRoleData(ctx context.Context) *CronAuthRoleData {
	return &CronAuthRoleData{
		ctx:       ctx,
		db:        db.New(ctx),
		tableName: "cron_auth_role",
	}
}
