// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package mma

import (
	"errors"
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

var (
	errIPNotFound = errors.New("IP not found")
)

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

	cmd := cmd{
		db: db,
	}

	go cmd.Listen()

	fmt.Printf("Create mma successfully\n")

	return nil
}

// GetDevices return DevicesInfo
func (m *MMA) GetDevices() ([]DevicesInfo, error) {
	devices, err := m.db.getAllDevices()
	if err != nil {
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

// LinkDevices Link IP1 -> IP2
func (m *MMA) LinkDevices(IP1 net.IP, IP2 net.IP) error {
	devices, err := m.db.getAllDevices()
	if err != nil {
		return err
	}

	var dev1, dev2 device
	for _, dev := range devices {
		if IP1.Equal(net.ParseIP(dev.ip)) {
			dev1 = dev
		}

		if IP2.Equal(net.ParseIP(dev.ip)) {
			dev2 = dev
		}
	}

	if dev1.id == "" || dev2.id == "" {
		return errIPNotFound
	}

	// tempz
	center := net.IPv4(11, 11, 11, 105)
	unit := "ab65950e_cb21_40e6_9517_2bbd32281ebc"
	if err := dev1.link(center, unit, &dev2); err != nil {
		return err
	}

	return nil
}

// Close a MMA
func (m *MMA) Close() {
	m.db.close()

	fmt.Printf("Close mma successfully\n")
}
