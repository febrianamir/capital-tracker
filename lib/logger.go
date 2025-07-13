package lib

import (
	"context"
	"time"

	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLLogger struct {
	logger.Interface
	Ctx      context.Context
	Repo     string
	Env      string
	SQL      string
	Rows     int64
	Duration time.Duration
}

func (l *SQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	l.SQL = sql
	l.Rows = rows
	l.Duration = time.Since(begin)

	if l.Duration >= 200*time.Millisecond || l.Env == "development" { // >= 200ms is considered slow query
		log.Println(l.SQL)
	}
	l.Interface.Trace(ctx, begin, fc, err)
}

func WithSQLLogger(ctx context.Context, db *gorm.DB, opts ...func(*SQLLogger)) (*gorm.DB, *SQLLogger) {
	sqlLogger := &SQLLogger{
		Interface: logger.Default.LogMode(logger.Silent),
		Ctx:       ctx,
	}

	for _, opt := range opts {
		opt(sqlLogger)
	}

	return db.Session(&gorm.Session{
		Logger: sqlLogger,
	}), sqlLogger
}

func WithRepo(name string) func(*SQLLogger) {
	return func(l *SQLLogger) {
		l.Repo = name
	}
}

func WithEnv(env string) func(*SQLLogger) {
	return func(l *SQLLogger) {
		l.Env = env
	}
}
