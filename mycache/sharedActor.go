package main

import (
	"fmt"
	"mycache/cacheMgr"
)



func serializedActor( cm *cacheMgr.CacheMap, pAction *ActionItem, goRoutineId int, doneCh chan int ) {
	for {
		err := pAction.Injector( &(pAction.Parameters), &(pAction.ActionValues) )
		if ( err != nil ) {
			fmt.Println( err.Error() )
			break
		}

		switch pAction.ActionValues.Action {
		case 'r':
			fmt.Printf( "read results  key: %d    value:%v\n", pAction.ActionValues.Key, cm.Reader(pAction.ActionValues.Key) )
			
		case 'w':
			cm.Writer( pAction.ActionValues.Key, pAction.ActionValues.Value )
		}
	}
	doneCh<- 1;
}


func multiplexActor( cmm *cacheMgr.CacheMapMultiplex, pAction *ActionItem, readerId int, doneCh chan int ) {
	for {
		err := pAction.Injector( &(pAction.Parameters), &(pAction.ActionValues) )
		if ( err != nil ) {
			fmt.Println( err.Error() )
			break
		}

		fmt.Printf( "a: %c   k:%d  v: %d\n", pAction.ActionValues.Action, pAction.ActionValues.Key, pAction.ActionValues.Value )

		switch pAction.ActionValues.Action {
		case 'r':
			val, err := cmm.Reader(pAction.ActionValues.Key, readerId)
			if ( err == nil ) {
				fmt.Printf( "read results  key: %d    value:%v\n", pAction.ActionValues.Key, val )
			}
			
		case 'w':
			cmm.Writer( pAction.ActionValues.Key, pAction.ActionValues.Value, readerId )
		}
	}
	doneCh<- 1;
}
