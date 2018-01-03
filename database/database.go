package database

import (
	"bytes"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DB struct {
	Conn     *gorm.DB
	User     string
	Name     string
	SSLMode  string
	Password string
	ConnStr  string
}

// Option is functional option for DB
type Option func(d *DB)

// WithName sets name
func WithName(name string) Option {
	return func(d *DB) {
		d.Name = name
	}
}

// WithUser sets user
func WithUser(user string) Option {
	return func(d *DB) {
		d.User = user
	}
}

// WithConnStr sets ConnStr
func WithConnStr(connStr string) Option {
	return func(d *DB) {
		d.ConnStr = connStr
	}
}

// NewDB creates database connection
func NewDB(options ...Option) (*DB, error) {
	d := &DB{
		Name:    "kommunicake_development",
		SSLMode: "disable",
	}
	d.ConnStr = d.formatConnStr()

	for _, option := range options {
		option(d)
	}

	conn, err := gorm.Open("postgres", d.ConnStr)
	// defer conn.Close()

	if err != nil {
		return nil, err
	}
	d.Conn = conn

	return d, nil
}

func (d *DB) formatConnStr() string {
	var buffer bytes.Buffer
	fields := map[string]string{
		"user=":     d.User,
		"dbname=":   d.Name,
		"sslmode=":  d.SSLMode,
		"password=": d.Password,
	}
	for k, v := range fields {
		if len(v) > 0 {
			buffer.WriteString(k)
			buffer.WriteString(v)
			buffer.WriteString(" ")
		}
	}
	return buffer.String()
}
