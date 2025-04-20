package xdb

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"sync"
)

var (
	ErrUrlEmpty = fmt.Errorf("url is empty, you can set url or 'file://' on local host")
	ErrNotFound = fmt.Errorf("not found")
)

type OptionFunc func(*Option)

type Option struct {
	filePath       []string
	Urls           []string
	SecurityAccess bool
}

type XDB struct {
	locker  sync.RWMutex
	handle  *os.File
	Options *Option
}

func (o *Option) WithUrl(u ...string) {
	o.Urls = u
}

func (o *Option) WithSecurityAccess(security bool) {
	o.SecurityAccess = security
}

func NewXDB(opts ...OptionFunc) (*XDB, error) {
	options := &Option{}
	for _, opt := range opts {
		opt(options)
	}
	for _, u := range options.Urls {
		if u == "" {
			return nil, ErrUrlEmpty
		}
		parse, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		if parse.Scheme == "file" {
			options.filePath = parse.Path
		}
	}
	handle, err := os.OpenFile(options.filePath, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	return &XDB{Options: options, handle: handle}, nil
}

func (xdb *XDB) Close() {
	if xdb.handle != nil {
		err := xdb.handle.Close()
		if err != nil {
			return
		}
	}
}

func (xdb *XDB) SearchByStr(str string) (string, error) {
	ip := net.ParseIP(str)
	if ip == nil {

	}
	ip.To4()
	return xdb.Search(str)
}

func (xdb *XDB) Search(ip uint32) (string, error) {
	if xdb.Options.SecurityAccess {
		xdb.locker.Lock()
		defer xdb.locker.Unlock()
	}
}
