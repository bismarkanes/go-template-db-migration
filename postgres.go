/*
Copyright © 2023 Bismark <bismark.john.anes@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package migration

import (
	"fmt"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type postgresMigration struct {
	dbHost     string `validate:"required,hostname"`
	dbUsername string `validate:"required"`
	dbPassword string `validate:"required"`
	dbName     string `validate:"required"`
	dbPort     string `validate:"required,hostname_port"`
	useSSL     bool
	sqlPath    string
}

// NewPostgresMigration use this to initialize migration if you are using postgresql database
func NewPostgresMigration(dbHost, dbUsername, dbPassword, dbName, dbPort string, useSSL bool) Migration {
	return &postgresMigration{
		dbHost:     dbHost,
		dbUsername: dbUsername,
		dbPassword: dbPassword,
		dbName:     dbName,
		dbPort:     dbPort,
		useSSL:     useSSL,
	}
}

// SQL files is defined in folder sql of root project
func (pm postgresMigration) generateMigrationDsn() string {
	sslMode := "disable"
	if pm.useSSL {
		sslMode = "enable"
	}
	return "postgres://" + url.QueryEscape(pm.dbUsername) + ":" + url.QueryEscape(pm.dbPassword) + "@" + fmt.Sprintf("%s:%s", pm.dbHost, pm.dbPort) + "/" + pm.dbName + "?sslmode=" + sslMode
}

func (pm postgresMigration) DoMigration(isUp bool) error {
	m, err := migrate.New(
		"file://sql",
		pm.generateMigrationDsn(),
	)

	if err != nil {
		return err
	}

	// execute the migration, up or down?
	if isUp {
		err = m.Up()
	} else {
		err = m.Down()
	}

	if err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}
