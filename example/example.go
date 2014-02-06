// ----------
// example.go ::: example usage for conf_kit
// ----------
// Copyright (c) 2013-Present, Scott Cagno. All rights reserved.
// This source code is governed by a BSD-style license.

package main

import (
	"conf_kit"
	"fmt"
)

/*
 * 	INSTRUCTION
 * 	-----------
 *  1) compile and execute this file
 *	2) open config/mycfg.json
 * 	3) simply save or change data and save
 * 	4) update notification (debug is true) should print within n (3) seconds
 */

func main() {

	// setup new config location
	cfg := conf_kit.NewConfig("config/mycfg.json")
	// run watch method; check the file every 3 seconds, debug mode set to true
	cfg.Watch(3, false)
    fmt.Printf("%+v\n", cfg.Data)
    fmt.Printf("%v\n", cfg.Data["name"])


	// this is used to halt program; otherwise program will exit right after the watch function runs
	func() { fmt.Println("save/change config file, any key will exit..."); var n int; fmt.Scanln(&n) }()
}
