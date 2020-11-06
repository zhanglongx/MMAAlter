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
	"time"
)

// CMDPORT is WS cmd port
const CMDPORT = ":7001"

type svr struct {
	db        *db
	enumerate func(devices []Device)
}

func (s *svr) Open() {
	go s.Listen()
}

func (s *svr) Listen() {
	ln, err := net.Listen("tcp", CMDPORT)
	if err != nil {
		panic(err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}

		go func() {
			isEnumerate := false

			for true {
				buffer := make([]byte, 512*1024)
				if _, err := conn.Read(buffer); err != nil {
					break
				}

				for _, b := range strings.Split(string(buffer), "\n") {
					if strings.Contains(b, "_MMAREQ_") {
						// TODO
					} else if strings.Contains(b, "_MMASYNC_") || strings.Contains(b, "_MMAEVENT_") {

						re := regexp.MustCompile(`_TABLENAME_:=(.*)_PRIMARYKEYVALUE_:=(.*)_SQL_:=(.*)`)
						if matched := re.FindSubmatch([]byte(b)); matched != nil {
							if err := s.db.updateDB(string(matched[1]), string(matched[2]), string(matched[3])); err != nil {
								fmt.Printf("%q\n", err)
							}
						}

						isEnumerate = true
					} else if strings.Contains(b, "_MMACMD_") {

						re := regexp.MustCompile(`_MMACMD_#_CMD_:=(.*)`)
						if matched := re.FindSubmatch([]byte(b)); matched != nil {
							var ack string
							if ack, err = cmdRoutine(string(matched[1])); err != nil {
								fmt.Printf("%q\n", err)
							}

							ack = ack + "\r"
							if _, err := conn.Write([]byte(ack)); err != nil {
								return
							}
						}
					}
				}
			}

			if isEnumerate {
				go s.delayedEnumerate()
			}
		}()
	}
}

func (s *svr) Close() {
	// TODO
}

func cmdRoutine(cmd string) (string, error) {
	join := func(s string) string {
		return "_MMAINFO_#_" + s + "_:=d41d8cd98f00b204e9800998ecf8427e"
	}

	switch cmd {
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
		return "_MMAINFO_#_" + "ACTIVED", nil
	}

	return "", errors.New("Not support cmd")
}

func (s *svr) delayedEnumerate() {
	// XXX: delayed 30 sec to read DB, to Ensure DB is updated
	time.Sleep(time.Duration(30) * time.Second)

	// TODO: delayed Enumerate
	if devs, err := s.db.getAllDevices(); err != nil {
		fmt.Printf("%v\n", err)
	} else {
		s.enumerate(devs)
	}
}
