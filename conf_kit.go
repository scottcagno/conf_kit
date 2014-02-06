// -----------
// conf_kit.go ::: config file kit & helpers
// -----------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package conf_kit

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

// config structure
type Config struct {
	Path string
	Data  map[string]interface{}
	sync.Mutex
}

// create new config structure instance
// return pointer to instance
func NewConfig(path string) *Config {
	return &Config{
		Path: path,
        Data: getData(path),
	}
}

// watches config file, updates Data map every time file changes
func (self *Config) Watch(n int64, debug bool) {
	// check to see if file has changed since n seconds ago
	if changed(self.Path, n) {
		// file changed; lock mutex to make Data map threadsafe
		self.Lock()
		// update Data map with new Data from file
		self.Data = getData(self.Path)
		// unlock mutex and proceed
		self.Unlock()
		// check debug status, if true print info
		if debug {
			log.Printf("%q changed, updated information: %+v\n", self.Path, self.Data)
		}
	}
	// wait for n seconds then recurse; aka: call this func again from within itself
	time.AfterFunc(time.Duration(n)*time.Second, func() { self.Watch(n, debug) })
}

// returns updated file Data as map
func getData(path string) map[string]interface{} {
	// read file Data (json []byte) into b
	b, err := ioutil.ReadFile(path)
	// err check
	if err != nil {
		// if err; print timestamp and err, then panic
		log.Panic(err)
	}
	// no err; create empty map as placeholder
	var m map[string]interface{}
	// unmarshal (convert) b (json []byte) into map
	err = json.Unmarshal(b, &m)
	// err check
	if err != nil {
		// if err; print timestamp and err, then panic
		log.Panic(err)
	}
	// no err; return map 
	return m
}

// returns true if file has been modified within n seconds, else return false
func changed(path string, n int64) bool {
	// gather file status
	fstat, err := os.Stat(path)
	// err check
	if err != nil {
		// if err; print timestamp and err, then panic
		log.Panic(err)
	}
	// no err; eval and return (true) file been modified within the last n seconds
	return fstat.ModTime().Unix() >= time.Now().Unix()-n
}
