package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"io/ioutil"
)

var totalPatrons int
var patronShards int

func init() {
	flag.IntVar(&totalPatrons, "patrons", 0, "total patron bots")
	flag.IntVar(&patronShards, "shards", 0, "total patron shards")
}

func main() {
	flag.Parse()

	var m map[int]int
	raw, err := ioutil.ReadFile("state.dat")
	if err != nil {
		panic(err)
	}
	b := bytes.NewBuffer(raw)
	d := gob.NewDecoder(b)

	err = d.Decode(&m)
	if err != nil {
		panic(err)
	}

	z := make(map[int]map[int]int, totalPatrons+1)
	z[0] = m
	for i := 1; i < totalPatrons+1; i++ {
		z[i] = make(map[int]int, patronShards)
	}

	b = new(bytes.Buffer)
	e := gob.NewEncoder(b)

	err = e.Encode(z)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("migratestate.dat", b.Bytes(), 666)
	if err != nil {
		panic(err)
	}

}
