package main

import (
	"fmt"
	"os"
)


func initCache( initFile string, cacheItems map [int] ) error {
	var initLine  [1000]byte
	var lineLen   int
	var key       int
	var ival      int
	var fval      float64
	var sval      [1000]byte

	f, err := os.Open(filePath)
	defer f.Close()

	if err != nil {
		return(err)
	}

	for _ ; err == nil; lineLen, err = fmt.Fscanln( f, initLine ) {
		if ( err != nil ) break;
		if ( lineLen < 4 ) continue;
		
		_, err = Sscanf( initLine, "%c", &c );
		switch c {
		    'i', 'd':
 			_, err = fmt.Sscanf( initLine, "%c%d%d", &c, &key, &ival );
			if ( err != nil ) continue;
			cacheItems[key] = ival;

		    'f':
			_, err = fmt.Sscanf( initLine, "%c%d%f", &c, &key, &fval );
			if ( err != nil ) continue;
			cacheItems[key] = fval;

		    's':
			_, err = fmt.Sscanf( initLine, "%c%d%s", &c, &key, sval );
			if ( err != nil ) continue;
			s = new( [len(sval) + 1]byte )
			copy( s, sval )
			cacheItems[key] = *s
		}
	}
	return( err )
}
