/**
 * Copyright 2019 gd Author. All rights reserved.
 * Author: Chuck1024
 */

package mysqldb

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

const defaultCharSet = "utf8mb4"

func (c *MysqlClient) initObjForMysqlDb(dbConfPath string) error {
	dbConfRealPath := dbConfPath
	if dbConfRealPath == "" {
		return errors.New("dbConf not set in g_cfg")
	}

	if !strings.HasSuffix(dbConfRealPath, ".ini") {
		return errors.New("dbConf not an ini file")
	}

	dbConf, err := ini.Load(dbConfRealPath)
	if err != nil {
		return err
	}

	if err = c.initDbs(dbConf); err != nil {
		return err
	}
	return nil
}

func (c *MysqlClient) initDbs(f *ini.File) error {
	m := f.Section("Mysql")
	s := f.Section("MysqlSlave")

	masterIp := m.Key("master_ip").String()
	masterPort := m.Key("master_port").String()
	userWrite := m.Key("user_write").String()
	passWrite := m.Key("pass_write").String()

	userRead := m.Key("user_read").String()
	passRead := m.Key("pass_read").String()

	db := m.Key("db").String()
	masterProxy, _ := m.Key("master_is_proxy").Bool()
	slaveIp := s.Key("slave_ip").String()
	slavePort := s.Key("slave_port").String()
	slaveProxy, _ := s.Key("slave_is_proxy").Bool()

	timeout := m.Key("timeout").String()
	if timeout == "" {
		timeout = "5s"
	} else if !strings.HasSuffix(timeout, "s") {
		timeout += "s"
	}

	connTimeout := m.Key("connTimeout").String()
	if connTimeout == "" {
		connTimeout = "1s"
	} else if !strings.HasSuffix(timeout, "s") {
		connTimeout += "s"
	}

	maxOpen, err := m.Key("max_open").Int()
	if err != nil {
		maxOpen = 100
	}
	maxIdle, err := m.Key("max_idle").Int()
	if err != nil {
		maxIdle = 1
	}

	enableSqlSafeUpdates, err := m.Key("enable_sql_safe_updates").Bool()
	if err != nil {
		enableSqlSafeUpdates = false
	}

	masterIps := strings.Split(masterIp, ",")
	connMasters := make([]string, 0)
	for _, masterIpVal := range masterIps {
		if masterIpVal == "" {
			continue
		}

		connMaster := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s", userWrite, passWrite, masterIp, masterPort, db, connTimeout, timeout, timeout)
		if enableSqlSafeUpdates {
			connMaster = connMaster + "&sql_safe_updates=1"
		}

		connMasters = append(connMasters, connMaster)
	}

	slaveIps := strings.Split(slaveIp, ",")
	connSlaves := make([]string, 0)
	for _, slaveIpVal := range slaveIps {
		if slaveIpVal == "" {
			continue
		}
		connSlaves = append(connSlaves, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s", userRead, passRead, slaveIp, slavePort, db, connTimeout, timeout, timeout))
	}

	glSuffix := m.Key("glSuffix").String()
	to, _ := time.ParseDuration(timeout)
	return c.initMainDbsMaxOpen(connMasters, connSlaves, maxOpen, maxIdle, glSuffix, to, masterProxy, slaveProxy)
}

type CommonDbConf struct {
	DbName      string
	ConnTime    string // connect timeout
	ReadTime    string // read timeout
	WriteTime   string // write timeout
	MaxOpen     int    // connect pool
	MaxIdle     int    // max idle connect
	MaxLifeTime int64  // connect life time
	glSuffix    string
	Master      *DbConnectConf
	Slave       *DbConnectConf
}

type DbConnectConf struct {
	Addrs                []string
	User                 string
	Pass                 string
	CharSet              string // default utf8mb4
	ClientFoundRows      bool   // 对于update操作,若更改的字段值跟原来值相同,当clientFoundRows为false时,sql执行结果会返回0;当clientFoundRows为true,sql执行结果返回1
	IsProxy              bool
	EnableSqlSafeUpdates bool // (safe update mode)，该模式不允许没有带WHERE条件的更新语句
}

func (c *MysqlClient) initDbsWithCommonConf(dbConf *CommonDbConf) error {
	if dbConf == nil {
		return errors.New("dbConf is nil")
	}
	if dbConf.Master == nil || len(dbConf.Master.Addrs) == 0 {
		return errors.New("masterAddr is nil")
	}
	if dbConf.DbName == "" {
		return errors.New("dbName is nil")
	}

	connTimeout := dbConf.ConnTime
	if connTimeout == "" {
		connTimeout = "200ms"
	}
	readTimeout := dbConf.ReadTime
	if readTimeout == "" {
		readTimeout = "500ms"
	}
	writeTimeout := dbConf.WriteTime
	if writeTimeout == "" {
		writeTimeout = "500ms"
	}
	maxOpen := dbConf.MaxOpen
	if maxOpen <= 0 {
		maxOpen = 100
	}
	maxIdle := dbConf.MaxIdle
	if maxIdle <= 0 {
		maxIdle = 1
	}

	connMasters, err := c.getReadWriteConnectString(dbConf.Master, connTimeout, readTimeout, writeTimeout, dbConf.DbName)
	if err != nil {
		return err
	}

	if len(connMasters) == 0 {
		return errors.New("no valid master ip found")
	}

	connSlave, err := c.getReadWriteConnectString(dbConf.Slave, connTimeout, readTimeout, writeTimeout, dbConf.DbName)
	if err != nil {
		return err
	}

	slaveIsProxy := false
	if dbConf.Slave != nil {
		slaveIsProxy = dbConf.Slave.IsProxy
	}

	to, err := time.ParseDuration(readTimeout)
	if err != nil {
		return fmt.Errorf("init mysqldb invalid duration %v", readTimeout)
	}

	return c.initMainDbsMaxOpen(connMasters, connSlave, maxOpen, maxIdle, dbConf.glSuffix, to, dbConf.Master.IsProxy, slaveIsProxy)
}

func (c *MysqlClient) getConnectString(conf *DbConnectConf, connTimeout, optTimeout int64, dbname string) ([]string, error) {
	if conf == nil || len(conf.Addrs) == 0 {
		return nil, nil
	}

	if conf.CharSet == "" {
		conf.CharSet = defaultCharSet
	}

	conStrs := make([]string, 0, len(conf.Addrs))
	for _, host := range conf.Addrs {
		if host != "" {
			var constr string
			if conf.ClientFoundRows {
				constr = fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&charset=%s&clientFoundRows=true",
					conf.User, conf.Pass, host, dbname, connTimeout, optTimeout, optTimeout, conf.CharSet)
			} else {
				constr = fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%ds&readTimeout=%ds&writeTimeout=%ds&charset=%s",
					conf.User, conf.Pass, host, dbname, connTimeout, optTimeout, optTimeout, conf.CharSet)
			}

			if conf.EnableSqlSafeUpdates {
				constr = constr + "&sql_safe_updates=1"
			}

			conStrs = append(conStrs, constr)
		}
	}
	return conStrs, nil
}

func (c *MysqlClient) getReadWriteConnectString(conf *DbConnectConf, connTimeout, readTimeout, writeTimeout string, dbname string) ([]string, error) {
	if conf == nil || len(conf.Addrs) == 0 {
		return nil, nil
	}

	if conf.CharSet == "" {
		conf.CharSet = defaultCharSet
	}

	constrs := make([]string, 0, len(conf.Addrs))
	for _, host := range conf.Addrs {
		if host != "" {
			var constr string
			if conf.ClientFoundRows {
				constr = fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&charset=%s&clientFoundRows=true",
					conf.User, conf.Pass, host, dbname, connTimeout, readTimeout, writeTimeout, conf.CharSet)
			} else {
				constr = fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&charset=%s",
					conf.User, conf.Pass, host, dbname, connTimeout, readTimeout, writeTimeout, conf.CharSet)
			}

			if conf.EnableSqlSafeUpdates {
				constr = constr + "&sql_safe_updates=1"
			}

			constrs = append(constrs, constr)
		}
	}

	return constrs, nil
}
