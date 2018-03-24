package main

import (
	"fmt"
)


type cacheItem struct {
	value interface{}
	readCh chan int
	writeCh chan int
}

func newCacheItem() *cacheItem {
	var cItem *cacheItem
	cItem = new(cacheItem)
	cItem.readCh = make(chan int)
	cItem.writeCh = make(chan int)
	return( cItem )
}

func main() {
	var sharedData map [int]cacheItem
	var sdRead chan int
	var sdWrite chan int
	var done [10]chan int
	var ec error
	var i int

	sharedData = make( map [int]cacheItem )
	ec = initCache( "initData.dat", sharedData )
	fmt.Println( sharedData )
	if ( ec != nil ) {
		fmt.Println( "Failed to read initialization data:", ec.Error() )
		return
	}

	for i=0; i<10; i++ {
		go sharedActor( sharedData, "dataFile", i, &done[i], &sdRead, &sdWrite )
		ec = nil
	}

	for i=0; i<10; i++ {
		<-done[i]
	}

//	PrintSharedData( sharedData )
}
