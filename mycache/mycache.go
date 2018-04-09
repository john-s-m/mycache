package main

import (
	"fmt"
	"mycache/cacheMgr"
	"os"
	"strings"
)

func main() {
	var useSerializer bool = true
	var useRandomizer bool = false

	for a := range os.Args {
		switch {
		case strings.Compare( os.Args[a], "-m") == 0 :
			useSerializer = false
		case strings.Compare( os.Args[a], "-r") == 0 :
			useRandomizer = true
		}
	}
	

	var done map[int]chan int
	var count int = 10
	var ec error
	var i int
	var openFileList []*os.File

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

	var pActionList []*ActionItem

	fmt.Printf( "Args: Serial:%t Rand:%t\n", useSerializer, useRandomizer )
	
	for i = 0; i < count; i++ {
		var pAction *ActionItem
		
		if useRandomizer {
			pAction = NewRandomActor()
			fmt.Printf( "RandomActor: %v\n", pAction )
			if ( pAction == nil ) {
				fmt.Println( "Failed to initialize random number actor" )
				return
			}
		} else {
			var pFile *os.File
			pAction, pFile = NewFileActor( "dataFile", i )
			if ( pAction == nil || pFile == nil ) {
				fmt.Printf( "failed to open datafile%d.dat\n", i )
				return
			}
			openFileList = append( openFileList, pFile )
		}
		pActionList = append( pActionList, pAction )
	}
	
	if useSerializer {
		fmt.Println( "Using Serializer" )

		for i = 0; i < count; i++ {
			go serializedActor(cm, pActionList[ i ], i, done[i])
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
			go multiplexActor( cmm, pActionList[ i ], i, done[i] )
		}
	}
	
	for i = 0; i < count; i++ {
		<-done[i]
	}

	for f := range openFileList {
		openFileList[f].Close()
	}
	
	fmt.Println(*mapPointer)
}
