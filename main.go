// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net"
	"time"

	mma "github.com/zhanglongx/MMAAlter/MMA"
)

func main() {
	m := mma.MMA{
		DbIP:      net.IPv4(192, 168, 64, 135),
		Enumerate: enumerate,
	}

	if err := m.Open(); err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(20) * time.Second)

	for true {
		// IP1 := net.IPv4(11, 11, 11, 104)
		// IP2 := net.IPv4(11, 11, 11, 109)
		// if err := m.LinkDevices(IP1, IP2); err != nil {
		// 	fmt.Printf("%q\n", err)

		// 	time.Sleep(time.Duration(10) * time.Second)
		// 	continue
		// }

		time.Sleep(time.Duration(100) * time.Second)

		// FIXME:
		// if err := m.DisLinkDevices(IP1, IP2); err != nil {
		// 	fmt.Printf("%q\n", err)

		// 	time.Sleep(time.Duration(10) * time.Second)
		// 	continue
		// }
	}

	time.Sleep(time.Duration(3600) * time.Second)

	m.Close()
}

func enumerate(devices []mma.Device) {
	for _, d := range devices {
		fmt.Printf("%v\n", d)
	}
}
