package main

import (
	"fmt"
	"mycache/cacheMgr"
)


func main() {
	var cm *cacheMgr.CacheMap
	var done map[int]chan int
	var ec error
	var i int
	var count int = 10

	done = make( map[int] chan int )
	cm = cacheMgr.NewCacheMap()
	ec = initCache( "initData.dat", cm.SharedMap )
	fmt.Println( cm.SharedMap )
	if ( ec != nil ) {
		fmt.Println( "Failed to read initialization data:", ec.Error() )
		return
	}

	for i=0; i<count; i++ {
		done[i] = make( chan int )
		go sharedActor( cm, "dataFile", i, done[i] )
	}

	for i=0; i<count; i++ {
		<-done[i]
	}

	fmt.Println( cm.SharedMap )
}
