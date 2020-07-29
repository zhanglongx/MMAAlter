// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package mma

import (
	"database/sql"
	"fmt"
	"net"
	"regexp"

	// MySql driver
	_ "github.com/go-sql-driver/mysql"
)

// MMASqlName is the default MMA DB Name
var MMASqlName = "mmasystem"

// MMAGlobaldevicestatus -> globaldevicestatus
var MMAGlobaldevicestatus = "globaldevicestatus"

// db is the struct for db connection
type db struct {
	// MMA MySql ip (local)
	ip net.IP

	// MMA DB user
	user string

	// MMA DB pwd
	password string

	// MMA MySql sql
	sql *sql.DB
}

// open a MMA database
func (d *db) open() error {
	var err error
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s", d.user, d.password, d.ip.String(), MMASqlName)
	d.sql, err = sql.Open("mysql", url)
	if err != nil {
		return err
	}

	return nil
}

// close a MMA database
func (d *db) close() {
	d.sql.Close()
}

// getAllDevices select from local database, and return devices
func (d *db) getAllDevices() ([]device, error) {
	// Execute the query
	rows, err := d.sql.Query("SELECT id,ip,devmcport,devworksta FROM " + MMAGlobaldevicestatus)
	if err != nil {
		return nil, err
	}

	var devices []device
	for rows.Next() {
		dev := device{}

		if err := rows.Scan(&dev.id, &dev.ip, &dev.recvPort, &dev.devworksta); err != nil {
			return nil, err
		}

		devices = append(devices, dev)

		// fmt.Printf("global devices: %s, %s, %d\n", dev.id, dev.ip, dev.recvPort)
	}

	return devices, nil
}

// updateDB update local databaSe
func (d *db) updateDB(table string, key string, sql string) error {
	if table != "globaldevicestatus" {
		return nil
	}

	query := fmt.Sprintf("SELECT * from %s WHERE id = '%s'", table, key)
	rows, err := d.sql.Query(query)
	if err != nil {
		return err
	}

	if !rows.Next() {
		query := fmt.Sprintf("INSERT INTO %s(id) VALUES('%s')", table, key)
		if _, execErr := d.sql.Exec(query); execErr != nil {
			return execErr
		}
	}

	// fmt.Printf(sql + "\n")

	re := regexp.MustCompile("and serverip=.*$")
	sql = re.ReplaceAllString(sql, "")

	// FIXME: transaction?
	_, execErr := d.sql.Exec(sql + ";")
	if execErr != nil {
		return execErr
	}

	return nil
}
