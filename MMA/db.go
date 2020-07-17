// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package mma

import (
	"database/sql"
	"fmt"
	"net"

	// MySql driver
	_ "github.com/go-sql-driver/mysql"
)

// MMASqlName is the default MMA DB Name
var MMASqlName = "mmasystem"

// MMAGlobaldevicestatus -> globaldevicestatus
var MMAGlobaldevicestatus = "globaldevicestatus"

// DB is the struct for DB connection
type DB struct {
	// MMA MySql IP (local)
	IP net.IP

	// MMA DB User
	User string

	// MMA DB pwd
	Password string

	// MMA MySql db
	db *sql.DB
}

// DevicesInfo is the device info
type DevicesInfo struct {
	// mmasystem.globaldevicestatus.id
	ID string

	// mmasystem.globaldevicestatus.ip
	IP net.IP
}

// Open a MMA database
func (d *DB) Open() error {
	var err error
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s", d.User, d.Password, d.IP.String(), MMASqlName)
	d.db, err = sql.Open("mysql", url)
	if err != nil {
		return err
	}

	return nil
}

// Close a MMA database
func (d *DB) close() {
	d.close()
}

// GetDevices select from local database, and return devices
func (d *DB) GetDevices() ([]DevicesInfo, error) {
	// Execute the query
	rows, err := d.db.Query("SELECT id,ip FROM " + MMAGlobaldevicestatus)
	if err != nil {
		return nil, err
	}

	var devices []DevicesInfo
	for rows.Next() {
		dev := DevicesInfo{}
		var IP string

		if err := rows.Scan(&dev.ID, &IP); err != nil {
			return nil, err
		}

		dev.IP = net.ParseIP(IP)
		devices = append(devices, dev)
	}

	return devices, nil
}
