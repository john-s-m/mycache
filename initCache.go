package main

import (
	"bufio"
	"fmt"
	"os"
)

func initCache(initFile string, cacheItems map[int]cacheItem) error {
	var key int
	var ival int
	var fval float64
	var c byte
	var sval *string
	var cItem *cacheItem

	fmt.Println(initFile)
	f, err := os.Open(initFile)
	if err != nil {
		return (err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		cItem = newCacheItem()

		_, err = fmt.Sscanf(scanner.Text(), "%c", &c)
		switch c {
		case 'i', 'd':
			_, err = fmt.Sscanf(scanner.Text(), "%c%d%d", &c, &key, &ival)
			if err != nil {
				continue
			}
			cItem.value = ival

		case 'f':
			_, err = fmt.Sscanf(scanner.Text(), "%c%d%f", &c, &key, &fval)
			if err != nil {
				continue
			}
			cItem.value = fval

		case 's':
			sval = new(string)
			_, err = fmt.Sscanf(scanner.Text(), "%c%d%s", &c, &key, sval)
			if err != nil {
				continue
			}
			cItem.value = *sval
		}
		cacheItems[key] = *cItem
	}
	fmt.Println("initial cache: ", cacheItems)
	return (err)
}
