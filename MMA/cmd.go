// Copyright 2020 Longxiao Zhang <zhanglongx@gmail.com>.
// All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package mma

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

// CMDPORT is WS cmd port
const CMDPORT = ":7001"

type cmd struct {
	db *db
}

func (c *cmd) Open() {
	// TODO
}

func (c *cmd) Listen() {
	ln, err := net.Listen("tcp", CMDPORT)
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

				for _, b := range strings.Split(string(buffer), "\n") {
					if strings.Contains(b, "_MMAREQ_") {
						// TODO
					} else if strings.Contains(b, "_MMASYNC_") {

						re := regexp.MustCompile(`_SQL_:=(.*)`)
						if matched := re.FindSubmatch([]byte(b)); matched != nil {
							if err := c.db.updateDB(string(matched[1])); err != nil {
								fmt.Printf("%q\n", err)
							}
						}

					} else if strings.Contains(b, "_MMACMD_") {

						re := regexp.MustCompile(`_MMACMD_#_CMD_:=(.*)`)
						if matched := re.FindSubmatch([]byte(b)); matched != nil {
							if ack, err := cmdRoutine(matched[1]); err == nil {
								if _, err := conn.Write(ack); err != nil {
									panic(err)
								}
							}
						}

					}
				}
			}
		}()
	}
}

func (c *cmd) Close() {
	// TODO
}

func cmdRoutine(cmd []byte) ([]byte, error) {
	join := func(s string) []byte {
		return []byte("_MMAINFO_#_" + s + "_:=d41d8cd98f00b204e9800998ecf8427e\r")
	}

	cmdStr := string(cmd)
	switch cmdStr {
	case "GETDEVICESTACHECKSUM":
		return join("DEVSTATUSCHECKSUM"), nil
	case "GETUSERSTACHECKSUM":
		return join("USERSTATUSCHECKSUM"), nil
	case "GETUNITCHECKSUM":
		return join("UNITCHECKSUM"), nil
	case "GETDEVGRANTCHECKSUM":
		return join("DEVGRANTCHECKSUM"), nil
	case "GETDEVLINKCHECKSUM":
		return join("DEVLINKSTACHECKSUM"), nil
	case "GETVIRTUALDEVCHECKSUM":
		return join("VIRTUALDEVCHECKSUM"), nil
	case "CHECKACTIVE":
		return []byte("_MMAINFO_#_" + "ACTIVED"), nil
	}

	return nil, errors.New("Not support cmd")
}
