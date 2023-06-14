/*
 * mysqldumper
 *
 * Copyright (c) XeLabs
 * Copyright (c) 2023-2030 NeoDB Author
 * GPL License
 *
 */

package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sealdb/mysqldumper/config"
	"github.com/sealdb/mysqlstack/driver"
	"github.com/sealdb/mysqlstack/sqlparser/depends/sqltypes"
	"github.com/sealdb/mysqlstack/xlog"
)

func TestLoader(t *testing.T) {
	log := xlog.NewStdLog(xlog.Level(xlog.DEBUG))
	fakedbs := driver.NewTestHandler(log)
	server, err := driver.MockMysqlServer(log, fakedbs)
	assert.Nil(t, err)
	defer server.Close()
	address := server.Addr()

	// fakedbs.
	{
		fakedbs.AddQueryPattern("create database if not exists `test.?`", &sqltypes.Result{})
		fakedbs.AddQuery("create table `t1-05-11` (`a` int(11) default null,`b` varchar(100) default null) engine=innodb", &sqltypes.Result{})
		fakedbs.AddQueryPattern("use .*", &sqltypes.Result{})
		fakedbs.AddQueryPattern("insert into .*", &sqltypes.Result{})
		fakedbs.AddQueryPattern("drop table .*", &sqltypes.Result{})
		fakedbs.AddQueryPattern("set foreign_key_checks=.*", &sqltypes.Result{})
	}

	args := &config.Config{
		Outdir:          "/tmp/dumpertest",
		User:            "mock",
		Password:        "mock",
		Threads:         16,
		Address:         address,
		IntervalMs:      500,
		OverwriteTables: true,
	}
	// Loader.
	{
		Loader(log, args)
	}
}
