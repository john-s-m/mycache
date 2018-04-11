package main

import (
	"bufio"
	"fmt"
	"os"
	"mycache/cacheMgr"
	"errors"
)

func initCache(initFile string, cacheItems map[int]cacheMgr.CacheItem) error {
	var key int
	var ival int
	var fval float64
	var c byte
	var sval *string
	var cItem *cacheMgr.CacheItem

	f, err := os.Open(initFile)
	if err != nil {
		return (err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		_, err = fmt.Sscanf(scanner.Text(), "%c", &c)
		switch c {
		case 'i', 'd':
			_, err = fmt.Sscanf(scanner.Text(), "%c%d%d", &c, &key, &ival)
			if err != nil {
				continue
			}
			cItem = cacheMgr.NewCacheItem( ival, false )

		case 'f':
			_, err = fmt.Sscanf(scanner.Text(), "%c%d%f", &c, &key, &fval)
			if err != nil {
				continue
			}
			cItem = cacheMgr.NewCacheItem( fval, false )

		case 's':
			sval = new(string)
			_, err = fmt.Sscanf(scanner.Text(), "%c%d%s", &c, &key, sval)
			if err != nil {
				continue
			}
			cItem = cacheMgr.NewCacheItem( *sval, false )

		default:
			s := fmt.Sprintf( "Invalid initialization file format: %s", initFile )
			return( errors.New( s ) )
		}
		cacheItems[key] = *cItem
	}
	fmt.Println("initial cache: ", cacheItems)
	return (err)
}
