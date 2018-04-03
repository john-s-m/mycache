package main

import (
	"bufio"
	"fmt"
	"os"
)

func initCache_test(initFile string, cacheItems *map[int]cacheItem) bool {
	fmt.Println(initFile)
	f, err := os.Open(initFile)
	if err != nil {
		return (err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var key int
		var ival int
		var fval float64
		var c byte
		var sval *string

		_, err = fmt.Sscanf(scanner.Text(), "%c", &c)
		switch c {
		case 'i', 'd':
			_, err = fmt.Sscanf(scanner.Text(), "%c%d%d", &c, &key, &ival)
			if err != nil {
				continue
			}
			if ( *cachItems[key].value != ival || *cachItems[key].readCh == nil || *cachItems[key].writeCh == nil ) {
				return( false )
			}

		case 'f':
			_, err = fmt.Sscanf(scanner.Text(), "%c%d%f", &c, &key, &fval)
			if err != nil {
				continue
			}
			if ( *cachItems[key] != fval || *cachItems[key].readCh == nil || *cachItems[key].writeCh == nil ) {
				return( false )
			}

		case 's':
			sval = new(string)
			_, err = fmt.Sscanf(scanner.Text(), "%c%d%s", &c, &key, sval)
			if err != nil {
				continue
			}
			if ( *cachItems[key] != sval  || *cachItems[key].readCh == nil || *cachItems[key].writeCh == nil ) {
				return( false )
			}
		}
	}
	return( true )
}
