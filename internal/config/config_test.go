package config

import (
	"net/url"
	"testing"
)

func TestDatabaseDSN(t *testing.T) {
	tests := []struct {
		name string
		db   Database
	}{
		{
			name: "plain",
			db: Database{
				Host:     "localhost",
				Port:     "5432",
				User:     "shop",
				Password: "secret",
				Name:     "shop",
				SSLMode:  "disable",
			},
		},
		{
			name: "special characters in password",
			db: Database{
				Host:     "localhost",
				Port:     "5432",
				User:     "shop",
				Password: "p@ss:w/rd?#",
				Name:     "shop",
				SSLMode:  "disable",
			},
		},
		{
			name: "ipv6host",
			db: Database{
				Host:     "::1",
				Port:     "5432",
				User:     "shop",
				Password: "secret",
				Name:     "shop",
				SSLMode:  "require",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.db.DSN())
			if err != nil {
				t.Fatalf("Database.DSN() returned error: %v", err)
			}
			if u.Scheme != "postgres" {
				t.Errorf("Database.DSN() returned scheme %v, want postgres", u.Scheme)
			}
			if got := u.User.Username(); got != tt.db.User {
				t.Errorf("Database.DSN() returned username %v, want %v", got, tt.db.User)
			}
			if got, ok := u.User.Password(); !ok || got != tt.db.Password {
				t.Errorf("Database.DSN() returned password %v, want %v", got, tt.db.Password)
			}
			if got := u.Hostname(); got != tt.db.Host {
				t.Errorf("Database.DSN() returned host %v, want %v", got, tt.db.Host)
			}
			if got := u.Port(); got != tt.db.Port {
				t.Errorf("Database.DSN() returned port %v, want %v", got, tt.db.Port)
			}
			if want := "/" + tt.db.Name; u.Path != want {
				t.Errorf("Database.DSN() returned path %v, want %v", u.Path, want)
			}
			if got := u.Query().Get("sslmode"); got != tt.db.SSLMode {
				t.Errorf("Database.DSN() returned sslmode %v, want %v", got, tt.db.SSLMode)
			}
		})
	}
}

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_USER", "shop")
	t.Setenv("DATABASE_PASSWORD", "secret")
	t.Setenv("DATABASE_NAME", "shop")
	t.Setenv("JWT_SECRET", "this-is-a-very-long-secret-key-for-testing")
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr bool
	}{
		{name: "default ok", env: nil, wantErr: false},
		{name: "max conns zero", env: map[string]string{"DATABASE_MAX_CONNS": "0"}, wantErr: true},
		{name: "min conns negative", env: map[string]string{"DATABASE_MIN_CONNS": "-1"}, wantErr: true},
		{name: "min conns zero", env: map[string]string{"DATABASE_MIN_CONNS": "0"}, wantErr: false},
		{name: "min greater than max", env: map[string]string{"DATABASE_MIN_CONNS": "10", "DATABASE_MAX_CONNS": "5"}, wantErr: true},
		{name: "min equals max", env: map[string]string{"DATABASE_MIN_CONNS": "5", "DATABASE_MAX_CONNS": "5"}, wantErr: false},
		{name: "max conns not a number", env: map[string]string{"DATABASE_MAX_CONNS": "abc"}, wantErr: true},
		{name: "short jwt secret", env: map[string]string{"JWT_SECRET": "short"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setRequiredEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			cfg, err := Load()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			if cfg == nil {
				t.Fatalf("Load() returned nil config without error")
			}
		})
	}
}
