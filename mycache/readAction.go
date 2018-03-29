package main

import (
	"bufio"
	"strings"
	"strconv"
	"fmt"
)

func readAction( scanner *bufio.Scanner, action *byte, key *int, value *interface{} ) bool {
	if ( ! scanner.Scan() )	{
		return( false )
	}

	var typeBytes []byte
	inputText := strings.Trim( scanner.Text(), " " )
	slices := strings.SplitN( inputText, " ", 4 )
	for i, s := range slices {
		switch i {
		case 0:
			actBytes := []byte( s )
			*action = actBytes[0]

		case 1:
			*key, _ = strconv.Atoi(s)

		case 2:
			typeBytes = []byte( s )

		case 3:
			switch typeBytes[0] {
			case 'i', 'd':
				*value, _ = strconv.Atoi( s )

			case 'f':
				*value, _ = strconv.ParseFloat( s, 64 )

			case 's':
				*value = s
			}
		}
	}
			
	fmt.Printf( "action: %c key: %d value:%v\n", *action, *key, *value )
	return( true )
}
