// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package mma

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

// DEVENC -> globaldevicestatus.devworksta
const (
	DEVLOST  = -1
	DEVEMPTY = 0
	DEVDEC   = 1
	DEVENC   = 10
)

type link struct {
	// centerIP
	center net.IP

	lock sync.RWMutex
}

var (
	errDevStaError = errors.New("Dev Sta Error")
)

func (l *link) open() error {
	return nil
}

func (l *link) close() {
}

func (l *link) encStart(d *Device, to *Device) error {
	// TODO: device type
	if d.Devworksta != DEVEMPTY {
		return errDevStaError
	}

	return l.setLinkSta(d.Devunitid, d.ID, to.IP, to.RecvPort, 10)
}

func (l *link) encStop(d *Device) error {
	if d.Devworksta != DEVENC {
		return errDevStaError
	}

	return l.setLinkSta(d.Devunitid, d.ID, "0", 0, 10)
}

func (l *link) decStart(d *Device, from *Device) error {
	// TODO: device type
	if d.Devworksta != DEVEMPTY {
		return errDevStaError
	}

	return l.setLinkSta(d.Devunitid, d.ID, from.IP, from.RecvPort, 1)

}

func (l *link) decStop(d *Device) error {
	if d.Devworksta != DEVDEC {
		return errDevStaError
	}

	return l.setLinkSta(d.Devunitid, d.ID, "0", 0, 1)
}

func (l *link) setLinkSta(unit string, id string, ip string, recvPort int, t int) error {

	l.lock.Lock()

	defer l.lock.Unlock()

	conn, err := net.Dial("tcp", l.center.String()+":7001")
	if err != nil {
		return err
	}

	defer conn.Close()

	var cmd []byte
	result := make([]byte, 1024)

	cmd = []byte("_MMACMD_#_CMD_:=CHECKACTIVE\r")
	if _, err := conn.Write(cmd); err != nil {
		return err
	}

	if _, err := conn.Read(result); err != nil {
		return err
	}

	cmd = []byte("DEVCOMMAND:SETLINKSTA\r")
	if _, err := conn.Write(cmd); err != nil {
		return err
	}

	if _, err := conn.Read(result); err != nil {
		return err
	}

	cmd = []byte(fmt.Sprintf("_DEVID_:=%s_DEVWORKIP_:=%s_DEVWORKPORT_:=%d_DEVWORKTYPE_:=%d_SETBYUNITID_:=%s\r",
		id, ip, recvPort, t, unit))

	if _, err := conn.Write(cmd); err != nil {
		return err
	}

	if _, err := conn.Read(result); err != nil {
		return err
	}

	return nil
}
