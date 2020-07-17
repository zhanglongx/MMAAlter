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

// MMA is the struct for MMA connection
type MMA struct {
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
func (m *MMA) Open() error {
	var err error
	url := fmt.Sprintf("%s:%s@tcp(%s)/%s", m.User, m.Password, m.IP.String(), MMASqlName)
	m.db, err = sql.Open("mysql", url)
	if err != nil {
		return err
	}

	return nil
}

// Close a MMA database
func (m *MMA) close() {
	m.close()
}

// GetDevices select from local database, and return devices
func (m *MMA) GetDevices() ([]DevicesInfo, error) {
	// Execute the query
	rows, err := m.db.Query("SELECT id,ip FROM " + MMAGlobaldevicestatus)
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
