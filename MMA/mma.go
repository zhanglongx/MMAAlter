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
	// Erl
	go epmd()
	go erlDP()

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
	devices, err := m.db.getDevices()
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
	devices, err := m.db.getDevices()
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

func epmd() {
	ln, err := net.Listen("tcp", ":4369")
	if err != nil {
		panic(err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		{
			buffer := make([]byte, 1024)
			if _, err := conn.Read(buffer); err != nil {
				panic(err)
			}

			epmdAck := []byte{0x77, 0x00, 0x11, 0x12, 0x4d, 0x00, 0x00, 0x05, 0x00,
				0x05, 0x00, 0x0c, 0x31, 0x31, 0x5f, 0x31, 0x31, 0x5f, 0x31, 0x31, 0x5f, 0x31, 0x30, 0x39, 0x00}

			// EPMD_PORT2_REQ
			if buffer[2] == 0x7A {
				if _, err := conn.Write(epmdAck); err != nil {
					panic(err)
				}
			}
		}
	}
}

func erlDP() {
	ln, err := net.Listen("tcp", ":4370")
	if err != nil {
		panic(err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		{
			buffer := make([]byte, 1024)
			if _, err := conn.Read(buffer); err != nil {
				panic(err)
			}

			erlDPAck := []byte{0x00, 0x03, 0x73, 0x6f, 0x6b}

			// Version: R6
			if buffer[3] == 0x00 && buffer[4] == 0x05 {
				if _, err := conn.Write(erlDPAck); err != nil {
					panic(err)
				}
			}
		}
	}
}
