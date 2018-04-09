package main

import (
	"bufio"
	"fmt"
	"os"
	"mycache/cacheMgr"
	"strings"
	"time"
	"math/rand"
)


func NewFileActor( dataFilePrefix string, goRoutineId int ) (*ActionItem, *os.File) {
	var dataFile string

	dataFile = fmt.Sprintf( "%s%d.dat", dataFilePrefix, goRoutineId )
	f, err := os.Open(dataFile)
	if err != nil {
		return nil, nil
	}

	pA := new( ActionItem )
//	pA.ActionValues = make( ActionInfo )
	pA.Injector = fileInjector
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	pA.Parameters = scanner
	return pA, f

}

func NewRandomActor() *ActionItem {
	pA := new( ActionItem )
	pA.Parameters = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	pA.Injector = randomInjector
//	pA.ActionValues = make( ActionInfo )
	return pA
}

func serializedActor( cm *cacheMgr.CacheMap, pAction *ActionItem, goRoutineId int, doneCh chan int ) {
	for {
		localAction := *pAction
		err := pAction.Injector( &localAction.Parameters, &localAction.ActionValues )
		if ( strings.Compare( err.Error(), "EOF" ) == 0 ) {
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
//		localAction := *pAction
//		err := pAction.Injector( &localAction.Parameters, &localAction.ActionValues )
		err := pAction.Injector( &(pAction.Parameters), &(pAction.ActionValues) )
		if ( err != nil && strings.Compare( err.Error(), "EOF" ) == 0 ) {
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
