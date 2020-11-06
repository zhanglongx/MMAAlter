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

// Device is the main struct for table GlobalDeviceStatus row
type Device struct {
	// globaldevicestatus.id
	ID string

	// globaldevicestatus.ip
	IP string

	// globaldevicestatus.devaudiorecvport
	RecvPort int

	// globaldevicestatus.devworksta
	Devworksta int

	// globaldevicestatus.devunitid
	Devunitid string

	// globaldevicestatus.devvideosendip
	Devvideosendip string
}

// MMA is main struct
type MMA struct {
	// DbIP is local MMA IP
	DbIP net.IP

	// Callback
	Enumerate func(devices []Device)

	// mma.db
	db *db

	// mma.SetLinkSta
	l link
}

var (
	errIPNotFound = errors.New("Device not available")
)

// Open a MMA
func (m *MMA) Open() error {

	db := &db{
		ip:       m.DbIP,
		user:     "root",
		password: "123",
	}

	if err := db.open(); err != nil {
		return err
	}

	m.db = db

	wsSvr := svr{
		db:        db,
		enumerate: m.Enumerate,
	}

	wsSvr.Open()

	m.l = link{
		center: net.IPv4(192, 168, 64, 135),
	}

	fmt.Printf("Create mma successfully\n")

	return nil
}

// LinkDevices Link IP1 -> IP2
func (m *MMA) LinkDevices(IP1 net.IP, IP2 net.IP) error {
	devices, err := m.db.getAllDevices()
	if err != nil {
		return err
	}

	var dev1, dev2 Device
	for _, dev := range devices {
		if IP1.Equal(net.ParseIP(dev.IP)) && dev.Devworksta == DEVEMPTY {
			dev1 = dev
		}

		if IP2.Equal(net.ParseIP(dev.IP)) && dev.Devworksta == DEVEMPTY {
			dev2 = dev
		}
	}

	if dev1.ID == "" || dev2.ID == "" {
		return errIPNotFound
	}

	if err := m.l.encStart(&dev1, &dev2); err != nil {
		return err
	}

	if err := m.l.decStart(&dev2, &dev1); err != nil {
		return err
	}

	// TODO: check devWorkSta

	return nil
}

// DisLinkDevices DisLink IP1 -> IP2
func (m *MMA) DisLinkDevices(IP1 net.IP, IP2 net.IP) error {
	devices, err := m.db.getAllDevices()
	if err != nil {
		return err
	}

	var dev1, dev2 Device
	for _, dev := range devices {
		if IP1.Equal(net.ParseIP(dev.IP)) && dev.Devworksta == DEVENC {
			dev1 = dev
		}

		if IP2.Equal(net.ParseIP(dev.IP)) && dev.Devworksta == DEVDEC {
			dev2 = dev
		}
	}

	if dev1.ID == "" || dev2.ID == "" {
		return errIPNotFound
	}

	if err := m.l.encStop(&dev1); err != nil {
		return err
	}

	if err := m.l.decStop(&dev2); err != nil {
		return err
	}

	// TODO: check devWorkSta

	return nil
}

// Close a MMA
func (m *MMA) Close() {
	m.l.close()

	m.db.close()

	fmt.Printf("Close mma successfully\n")
}
