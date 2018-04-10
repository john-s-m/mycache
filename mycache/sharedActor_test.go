package main

import (
	"testing"
	"mycache/cacheMgr"
)

func TestMultiplexActor(t *testing.T) {
	var actionArray = []ActionInfo {
		{'w', 1, 1},
		{'w', 2, 2},
		{'w', 3, 3},
		{'w', 4, 4},
		{'w', 5, 5},
		{'w', 6, 6},
		{'w', 7, 7},
		{'w', 8, 8},
		{'w', 9, 9},
		{'w', 10, 10},
		{'r', 1, -1},
		{'r', 2, -1},
		{'r', 3, -1},
		{'r', 4, -1},
		{'r', 5, -1},
		{'r', 6, -1},
		{'r', 7, -1},
		{'r', 8, -1},
		{'r', 9, -1},
		{'r', 10, -1} }

	pAction := NewArrayActor( actionArray )
	if ( pAction == nil ) {
		t.Errorf( "TestmultiplexActor failed on allocating NewArrayActor, got nil pointer back\n" );
		return
	}

	cmm := cacheMgr.NewCacheMapMultiplex()
	if ( cmm == nil ) {
		t.Errorf( "TestmultiplexActor failed on allocating NewCacheMapMultiplex, got nil pointer back\n" );
		return
	}

	cmm.AddReader()
	cmm.StartAllRoutines()

	doneCh := make(chan int)
	go multiplexActor( cmm, pAction, 0, doneCh )
	<- doneCh
//	cmm.TerminateForwardGR()
	
	for i:=1; i<=10; i++ {
		if  ( cmm.SharedMap[i].Value.(int) != i ) {
			t.Errorf( "TestmultiplexActor - data mismatch for key %d: %v\n", i, cmm.SharedMap[i].Value )
		}
	}
}
	
