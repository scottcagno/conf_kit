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
	fpath string
	data  map[string]interface{}
	mu    sync.Mutex
}

// create new config structure instance
// return pointer to instance
func NewConfig(fpath string) *Config {
	return &Config{fpath: fpath}
}

// watches config file, updates data map every time file changes
func (self *Config) Watch(n int64, debug bool) {
	// check to see if file has changed since n seconds ago
	if changed(self.fpath, n) {
		// file changed; lock mutex to make data map threadsafe
		self.mu.Lock()
		// update data map with new data from file
		self.data = getdata(self.fpath)
		// unlock mutex and proceed
		self.mu.Unlock()
		// check debug status, if true print info
		if debug {
			log.Printf("%q changed. updated map data\n  ~ %+v\n", self.fpath, self.data)
		}
	}
	// wait for n seconds then recurse; aka: call this func again from within itself
	time.AfterFunc(time.Duration(n)*time.Second, func() { self.Watch(n, debug) })
}

// returns updated file data as map
func getdata(fpath string) map[string]interface{} {
	// read file data (json []byte) into b
	b, err := ioutil.ReadFile(fpath)
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
func changed(fpath string, n int64) bool {
	// gather file status
	fstat, err := os.Stat(fpath)
	// err check
	if err != nil {
		// if err; print timestamp and err, then panic
		log.Panic(err)
	}
	// no err; eval and return (true) file been modified within the last n seconds
	return fstat.ModTime().Unix() >= time.Now().Unix()-n
}
