// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"

	mma "github.com/zhanglongx/MMAAlter/MMA"
)

func main() {
	m := mma.MMA{
		DbIP: net.IPv4(11, 11, 11, 105),
	}

	if err := m.Open(); err != nil {
		panic(err)
	}

	var devices []mma.DevicesInfo
	var err error

	if devices, err = m.GetDevices(); err != nil {
		panic(err)
	}

	if len(devices) == 0 {
		fmt.Printf("No devices. Cascaded?\n")
		return
	}

	IP1 := net.IPv4(11, 11, 11, 104)
	IP2 := net.IPv4(11, 11, 11, 109)
	if err := m.LinkDevices(IP1, IP2); err != nil {
		panic(err)
	}

	m.Close()
}
