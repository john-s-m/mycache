package main

import (
	"bufio"
	"fmt"
	"os"
)

func sharedActor( sharedData map[int]cacheItem, dataFilePrefix string, goRoutineId int, doneCh *chan int , sdRead *chan int, sdWrite *chan int ) {
	var dataFile string

	dataFile = fmt.Sprintf( "%s%d.dat", dataFilePrefix, goRoutineId )
	f, err := os.Open(dataFile)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Printf( "Opened %s\n", dataFile )
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var value    string
	var key      int
	var action   byte

	for readAction( scanner, &action, &key, &value ) {
		switch action {
		case 'r':
/*			*sdRead <- 1
			var isWriteLocked int
			isWriteLock <- *sdWrite
			*sdWrite <- 1
*/
			fmt.Printf( "action: %c  key: %d map: %v\n", action, key, sharedData[key] )
//			<- *sdRead
			
		case 'w':
			fmt.Printf( "action: %c  key: %d value: %v\n", action, key, value )
		}
	}
	*doneCh<- 1;
}
