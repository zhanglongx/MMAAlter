// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package mma

import (
	"errors"
	"net"
	"regexp"
)

type cmd struct {
}

func (c *cmd) Open() {
	go listen()
}

func (c *cmd) Close() {
	// TODO:
}

func listen() {
	ln, err := net.Listen("tcp", ":7001")
	if err != nil {
		panic(err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			for true {
				buffer := make([]byte, 512*1024)
				if _, err := conn.Read(buffer); err != nil {
					return
				}

				re := regexp.MustCompile(`_MMACMD_#_CMD_:=(.*)`)
				if matched := re.FindSubmatch(buffer); matched != nil {
					if ack, err := cmdRoutine(matched[1]); err == nil {
						if _, err := conn.Write(ack); err != nil {
							return
						}
					}
				}
			}
		}()
	}
}

func cmdRoutine(cmd []byte) ([]byte, error) {
	cmdStr := string(cmd)
	switch cmdStr {
	case "GETDEVICESTACHECKSUM":
		return []byte("_MMAINFO_#_" + cmdStr + "_:=d41d8cd98f00b204e9800998ecf8427e"), nil
	case "GETUSERSTACHECKSUM":
		return []byte("_MMAINFO_#_" + cmdStr + "_:=d41d8cd98f00b204e9800998ecf8427e"), nil
	case "GETUNITCHECKSUM":
		return []byte("_MMAINFO_#_" + cmdStr + "_:=d41d8cd98f00b204e9800998ecf8427e"), nil
	}

	return nil, errors.New("Not support cmd")
}
