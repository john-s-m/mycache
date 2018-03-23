package main

import (
	"fmt"
	"os"
	"bufio"
)


func initCache( initFile string, cacheItems map [int]cacheItem ) error {
	var initLine  string
	var key       int
	var ival      int
	var fval      float64
	var c         byte
	var sval      *string
	var cItem     *cacheItem

	fmt.Println( initFile )
	f, err := os.Open(initFile)
	if err != nil {
		return(err)
	}
	defer f.Close()


	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		cItem = new( cacheItem )
		cItem.readCh = make( chan int )
		cItem.writeCh = make( chan int )
		initLine = scanner.Text()
		fmt.Println( "init:", initLine )
		fmt.Println( "scanner:", scanner.Text() );
		_, err = fmt.Sscanf( scanner.Text(), "%c", &c )
		switch c {
		    case 'i', 'd':
 			_, err = fmt.Sscanf( scanner.Text(), "%c%d%d", &c, &key, &ival )
			if ( err != nil ) {
				continue
			}
			cItem.value = ival
			fmt.Println( "KV: ", key, cItem.value );

		    case 'f':
			_, err = fmt.Sscanf( scanner.Text(), "%c%d%f", &c, &key, &fval )
			if ( err != nil ) {
				continue
			}
			cItem.value = fval

		    case 's':
			sval = new( string )
			_, err = fmt.Sscanf( scanner.Text(), "%c%d%s", &c, &key, sval )
			if ( err != nil ) {
				continue
			}
			cItem.value = *sval
		}
		cacheItems[key] = *cItem
	}
	fmt.Println( "cache: ", cacheItems )
	return( err )
}
