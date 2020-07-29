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

	// mma.SetLinkSta
	l link
}

// DevicesInfo is the device info
type DevicesInfo struct {
	// DeviceIP
	DeviceIP net.IP
}

var (
	errIPNotFound = errors.New("Device not available")
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

	wsSvr := svr{
		db: db,
	}

	go wsSvr.Listen()

	m.l = link{
		center: net.IPv4(11, 11, 11, 105),
		unit:   "ab65950e_cb21_40e6_9517_2bbd32281ebc",
	}

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
		if IP1.Equal(net.ParseIP(dev.ip)) && dev.devworksta == DEVEMPTY {
			dev1 = dev
		}

		if IP2.Equal(net.ParseIP(dev.ip)) && dev.devworksta == DEVEMPTY {
			dev2 = dev
		}
	}

	if dev1.id == "" || dev2.id == "" {
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

	var dev1, dev2 device
	for _, dev := range devices {
		if IP1.Equal(net.ParseIP(dev.ip)) && dev.devworksta == DEVENC {
			dev1 = dev
		}

		if IP2.Equal(net.ParseIP(dev.ip)) && dev.devworksta == DEVDEC {
			dev2 = dev
		}
	}

	if dev1.id == "" || dev2.id == "" {
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
