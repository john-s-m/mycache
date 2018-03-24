package main

import (
	"bufio"
	"fmt"
	"os"
)

func sharedActor( sharedData map[int]cacheItem, dataFilePrefix string, goRoutineId int, doneCh *chan int , sdRead *chan int, sdWrite *chan int ) {
	var dataFile string

	fmt.Sprintf( dataFile, "%s%d.dat", dataFilePrefix, goRoutineId )
	fmt.Println(dataFile)
	f, err := os.Open(dataFile)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var value    string
	var key      int
	var action   byte

	for readAction( scanner, &action, &key, &value ) {
		switch action {
		case 'r':
			fmt.Printf( "action: %c  key: %d map: %v\n", action, key, sharedData[key] )

		case 'w':
			fmt.Printf( "action: %c  key: %d value: %v\n", action, key, value )
		}
	}
	*doneCh<- 1;
}
