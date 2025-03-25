package lib

import (
	"bytes"
	"database/sql"
	"html/template"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestDbConnect (t *testing.T) {
	const MARIADB_DATABASE = "mc15"
	const MARIADB_HOST = "localhost"
	const MARIADB_PASSWORD = "mc15"
	const MARIADB_PORT = 3306
	const MARIADB_USER = "mc15"
	var dbData = map[string]any {
		"db": MARIADB_DATABASE,
		"host": MARIADB_HOST,
		"password": MARIADB_PASSWORD,
		"port": MARIADB_PORT,
		"username": MARIADB_USER,
	}
	sConnStringTmpl := "{{.username}}:{{.password}}@tcp({{.host}}:{{.port}})/{{.db}}"
	tmpl, err := template.New("connString").Parse(sConnStringTmpl)
	assert.NoError(t, err)
	buf := new(bytes.Buffer)
 	err = tmpl.Execute(buf, dbData)
	assert.NoError(t, err)
	connString := buf.String()
	db, err := sql.Open("mysql", connString)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	var nUsers int
	row := db.QueryRow("select count(*) from user")
	err = row.Scan(&nUsers)
	assert.NoError(t, err)
	assert.Equal(t, 0, nUsers)
}
