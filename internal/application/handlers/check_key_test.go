package handlers

import (
	"bytes"
	"database/sql"
	"net/http"
	"testing"

	"github.com/kyrare/ya-metrics/internal/infrastructure/connection"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHandler_checkRequestKey(t *testing.T) {
	logger := zap.New(nil)
	sugar := *logger.Sugar()
	storage, err := metrics.NewMemStorage("", sugar)
	assert.NoError(t, err, "Не удалось создать storage")

	db, err := connection.New("", sugar)
	assert.NoError(t, err, "Не удалось создать соединение с БД")

	type fields struct {
		storage  metrics.Storage
		checkKey bool
		DB       *sql.DB
		logger   zap.SugaredLogger
	}
	type args struct {
		r *http.Request
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "CheckKey false",
			fields: fields{
				storage:  storage,
				checkKey: false,
				DB:       db,
				logger:   sugar,
			},
			args: args{
				r: new(http.Request),
			},
			want: true,
		},
		{
			name: "CheckKey true and empty header",
			fields: fields{
				storage:  storage,
				checkKey: true,
				DB:       db,
				logger:   sugar,
			},
			args: args{
				r: new(http.Request),
			},
			want: false,
		},
		{
			name: "CheckKey true and correct header",
			fields: fields{
				storage:  storage,
				checkKey: true,
				DB:       db,
				logger:   sugar,
			},
			args: args{
				r: (func() *http.Request {
					body := bytes.NewBuffer([]byte("test"))
					r, _ := http.NewRequest("GET", "/test", body)
					r.Header.Add("HashSHA256", "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08")

					return r
				})(),
			},
			want: true,
		},
		{
			name: "CheckKey true and incorrect header",
			fields: fields{
				storage:  storage,
				checkKey: true,
				DB:       db,
				logger:   sugar,
			},
			args: args{
				r: (func() *http.Request {
					body := bytes.NewBuffer([]byte("test"))
					r, _ := http.NewRequest("GET", "/test", body)
					r.Header.Add("HashSHA256", "incorrect")

					return r
				})(),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				storage:  tt.fields.storage,
				checkKey: tt.fields.checkKey,
				DB:       tt.fields.DB,
				logger:   tt.fields.logger,
			}

			assert.Equal(t, tt.want, h.checkRequestKey(tt.args.r))
		})
	}
}
