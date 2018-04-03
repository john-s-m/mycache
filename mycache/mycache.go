package main

import (
	"fmt"
	"mycache/cacheMgr"
	"os"
	"strings"
)

func main() {
	var useSerializer bool = true
	if (len(os.Args) > 1) && ( strings.Compare( os.Args[1], "-m") == 0 ) {
		useSerializer = false
	}

	var done map[int]chan int
	var count int = 10
	var ec error
	var i int

	done = make(map[int]chan int)
	for i = 0; i < count; i++ {
		done[i] = make(chan int)
	}

	var cm *cacheMgr.CacheMap
	var cmm *cacheMgr.CacheMapMultiplex
	var mapPointer *map [int]cacheMgr.CacheItem

	if useSerializer {
		cm = cacheMgr.NewCacheMap()
		mapPointer = &cm.SharedMap
	} else {
		cmm = cacheMgr.NewCacheMapMultiplex()
		mapPointer = &cmm.SharedMap
	}
	
	ec = initCache("initData.dat", *mapPointer)
	if ec != nil {
		fmt.Println("Failed to read initialization data:", ec.Error())
		return
	}
	fmt.Println(*mapPointer)

	if useSerializer {
		fmt.Println( "Using Serializer" )

		for i = 0; i < count; i++ {
			go sharedActor(cm, "dataFile", i, done[i])
		}
	} else {
		fmt.Println( "Using Multiplexer" )

		for i = 0; i < count; i++ {
			cmm.AddReader()
		}
		fmt.Println( "starting" );
		cmm.StartAllRoutines()
		fmt.Println( "all routines started" );

		for i = 0; i < count; i++ {
			go multiplexActor( cmm, "dataFile", i, done[i] )
		}
	}
	
	for i = 0; i < count; i++ {
		<-done[i]
	}

	fmt.Println(*mapPointer)
}
