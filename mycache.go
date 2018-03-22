package main

import (
	"fmt"
)


type cacheItem struct {
	value interface{}
	readCh chan
	writeCh chan
}



func main() {
	var sharedData map [int]cacheItem
	var sdReadCh chan
	var sdWrite chan
	var done [10]chan;
	var ec error;

	ec = initCache( "initData.dat", sharedData );
	if ( ec != nil ) {
		Println( "Failed to read initialization data: %s", ec.Error() );
		return;
	}

	for i=0; i<10; i++ {
		go sharedActor( "dataFile", i, done[i], sdRead, sdWrite );
	}

	for i=0; i<10; i++ {
		<-done[i]
	}

	PrintSharedData( sharedData );
}
