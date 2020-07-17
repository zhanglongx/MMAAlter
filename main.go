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
		IP:       net.IPv4(10, 10, 10, 106),
		User:     "root",
		Password: "wisdom",
	}

	if err := m.Open(); err != nil {
		panic(err)
	}

	var devices []mma.Devices
	var err error
	if devices, err = m.GetDevices(); err != nil {
		panic(err)
	}

	if len(devices) > 0 {
		fmt.Printf("%s\n", devices[0].ID)
	}
}
