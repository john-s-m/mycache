package main

import (
	"bufio"
	"fmt"
	"os"
	"mycache/cacheMgr"
)

func sharedActor( cm *cacheMgr.CacheMap, dataFilePrefix string, goRoutineId int, doneCh chan int ) {
	var dataFile string

	dataFile = fmt.Sprintf( "%s%d.dat", dataFilePrefix, goRoutineId )
	f, err := os.Open(dataFile)
	if err != nil {
		return
	}
	defer f.Close()

//	fmt.Printf( "Opened %s\n", dataFile )
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var value    interface{}
	var key      int
	var action   byte

	for readAction( scanner, &action, &key, &value ) {
		switch action {
		case 'r':
			fmt.Printf( "read results  key: %d    value:%v\n", key, cm.Reader(key) )
			
		case 'w':
			cm.Writer( key, value )
		}
	}
	doneCh<- 1;
}
