package main

import (
	"mycache/cacheMgr"	
	"testing"
	"strings"
)

func TestinitCache( t *testing.T) {
	initFile := "testfile.dat"

	cm := cacheMgr.NewCacheMapStruct()

	err := initCache( initFile, cm.SharedMap )
	if ( err != nil ) && ( strings.Compare( err.Error(), "EOF" ) != 0 ) {
		t.Errorf( "TestinitCache - got Error initializing Cache: %s\n", err.Error() )
		return
	}

	for r := range cm.SharedMap {
		if ( r < 0 ) || ( r > 100 ) {
			t.Errorf( "TestinitCache - Key out of range: %d\n", r )
		}
		if ( cm.SharedMap[r].Value == nil ) {
			t.Errorf( "TestinitCache - nil Value for key: %d\n", r )
		}
	}
}
