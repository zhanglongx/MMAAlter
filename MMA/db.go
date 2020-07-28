// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package mma

import (
	"database/sql"
	"fmt"
	"net"
	"sync"
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

	// stmt
	stmt *sql.Stmt

	lock sync.RWMutex
}

// open a MMA database
func (d *db) open() error {
	var err error
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s", d.user, d.password, d.ip.String(), MMASqlName)
	d.sql, err = sql.Open("mysql", url)
	if err != nil {
		return err
	}

	d.stmt, err = d.sql.Prepare("?")
	if err != nil {
		return err
	}

	return nil
}

// close a MMA database
func (d *db) close() {
	d.stmt.Close()

	d.sql.Close()
}

// getAllDevices select from local database, and return devices
func (d *db) getAllDevices() ([]device, error) {
	d.lock.Lock()

	defer d.lock.Unlock()

	// Execute the query
	rows, err := d.sql.Query("SELECT id,ip,devmcport FROM " + MMAGlobaldevicestatus)
	if err != nil {
		return nil, err
	}

	var devices []device
	for rows.Next() {
		dev := device{}

		if err := rows.Scan(&dev.id, &dev.ip, &dev.recvPort); err != nil {
			return nil, err
		}

		devices = append(devices, dev)

		// fmt.Printf("global devices: %s, %s, %d\n", dev.id, dev.ip, dev.recvPort)
	}

	return devices, nil
}

// updateDB update local databaSe
func (d *db) updateDB(sql string) error {
	d.lock.Lock()

	defer d.lock.Unlock()

	// Execute
	if _, err := d.sql.Exec(sql); err != nil {
		return err
	}

	return nil
}
