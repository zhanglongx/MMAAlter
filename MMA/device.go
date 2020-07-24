// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package mma

import (
	"fmt"
	"net"
	"time"
)

type device struct {
	// globaldevicestatus.id
	id string

	// globaldevicestatus.ip
	ip string

	// globaldevicestatus.devaudiorecvport
	recvPort int
}

// link device -> to
func (d *device) link(center net.IP, unit string, to *device) error {
	conn, err := net.Dial("tcp", center.String()+":7001")
	if err != nil {
		return err
	}

	defer conn.Close()

	var cmd []byte
	result := make([]byte, 1024)

	cmd = []byte("DEVCOMMAND:SETLINKSTA")
	if _, err := conn.Write(cmd); err != nil {
		return err
	}

	if _, err := conn.Read(result); err != nil {
		return err
	}

	// tempz
	cmd = []byte(fmt.Sprintf("_DEVID_:=%s_DEVIP_:=%s_DEVWORKPORT_:=%d_DEVWORKTYPE_:=%d_SETBYUNITID_=%s",
		to.id, to.ip, to.recvPort, 10, unit))
	if _, err := conn.Write(cmd); err != nil {
		return err
	}

	if _, err := conn.Read(result); err != nil {
		return err
	}

	// tempz
	time.Sleep(time.Duration(10) * time.Second)

	cmd = []byte("DEVCOMMAND:SETLINKSTA")
	if _, err := conn.Write(cmd); err != nil {
		return err
	}

	if _, err := conn.Read(result); err != nil {
		return err
	}

	// tempz
	cmd = []byte(fmt.Sprintf("_DEVID_:=%s_DEVIP_:=%s_DEVWORKPORT_:=%d_DEVWORKTYPE_:=%d_SETBYUNITID_=%s",
		d.id, d.ip, to.recvPort, 1, unit))
	if _, err := conn.Write(cmd); err != nil {
		return err
	}

	if _, err := conn.Read(result); err != nil {
		return err
	}

	return nil
}

// dislink device -> to
func (d *device) dislink(center net.IP, unit string, to *device) error {
	conn, err := net.Dial("tcp", center.String()+":7001")
	if err != nil {
		return err
	}

	defer conn.Close()

	return nil
}