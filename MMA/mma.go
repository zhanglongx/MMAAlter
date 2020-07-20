// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package mma

import (
	"fmt"
	"net"
)

// MMA is main struct
type MMA struct {
	// DbIP is local MMA IP
	DbIP net.IP

	// mma.db
	db *db
}

// DevicesInfo is the device info
type DevicesInfo struct {
	// DeviceIP
	DeviceIP net.IP
}

// Open a MMA
func (m *MMA) Open() error {
	db := &db{
		ip:       m.DbIP,
		user:     "root",
		password: "wisdom",
	}

	if err := db.open(); err != nil {
		return err
	}

	m.db = db

	fmt.Printf("Create mma successfully\n")

	return nil
}

// GetDevices return DevicesInfo
func (m *MMA) GetDevices() ([]DevicesInfo, error) {
	var devices []device
	var err error
	if devices, err = m.db.getDevices(); err != nil {
		return nil, err
	}

	var Infos []DevicesInfo
	for _, d := range devices {
		Infos = append(Infos, DevicesInfo{
			DeviceIP: net.ParseIP(d.ip),
		})

		fmt.Printf("Get DevIP: %s\n", d.ip)
	}

	return Infos, nil
}

// Close a MMA
func (m *MMA) Close() {
	m.db.close()

	fmt.Printf("Close mma successfully\n")
}
