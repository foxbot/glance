// Package main for makestate is a command to make a state file
// (we're just using this for a single migration when moving to the state system)
package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"io/ioutil"
)

var shards int
var state int

func init() {
	flag.IntVar(&shards, "shards", 0, "total shards")
	flag.IntVar(&state, "state", 0, "default state")
}

func main() {
	flag.Parse()

	m := make(map[int]int, shards)
	for i := 0; i < shards; i++ {
		m[i] = state
	}

	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)

	err := e.Encode(m)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("makestate.dat", b.Bytes(), 666)
	if err != nil {
		panic(err)
	}
	println()
}
