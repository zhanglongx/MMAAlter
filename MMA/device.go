// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package mma

type device struct {
	// globaldevicestatus.id
	id string

	// globaldevicestatus.ip
	ip string

	// globaldevicestatus.devaudiorecvport
	recvPort string
}

// link device -> to
func (d *device) link(unit string, to *device) error {
	return nil
}

// dislink device -> to
func (d *device) dislink(unit string, to *device) error {
	return nil
}
