// Go MySQL Driver - A MySQL-Driver for Go's database/sql package
//
// Copyright 2012 The Go-MySQL-Driver Authors. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at http://mozilla.org/MPL/2.0/.

package mysql

import (
	"os"
	"time"
)

type mysqlTx struct {
	mc      *mysqlConn
	txStack []byte
	cfg     *Config
	closed  bool
}

func (tx *mysqlTx) Commit() (err error) {
	if tx.mc == nil || tx.mc.closed.IsSet() {
		return ErrInvalidConn
	}
	err = tx.mc.exec("COMMIT")
	tx.mc = nil
	tx.closed = true
	return
}

func (tx *mysqlTx) Rollback() (err error) {
	if tx.mc == nil || tx.mc.closed.IsSet() {
		return ErrInvalidConn
	}
	err = tx.mc.exec("ROLLBACK")
	tx.mc = nil
	tx.closed = true
	return
}

func (tx *mysqlTx) startLeakCheckTimer() {

	if !tx.cfg.LeakDetectionEnabled {
		return
	}
	go func() {
		time.Sleep(*tx.cfg.LeakTimeout)
		tx.leakCheck()
	}()
}

func (tx *mysqlTx) leakCheck() {

	if !tx.closed {
		os.Stderr.WriteString(ErrConnectionLeak.Error())
		os.Stderr.Write(tx.txStack)
		if tx.cfg.PanicOnLeak {
			panic(ErrConnectionLeak)
		}
	}

}
