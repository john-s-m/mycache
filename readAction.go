package main

import (
	"bufio"
	"fmt"
)

func readAction( scanner *bufio.Scanner, action *byte, key *int, value *string ) bool {
	if ( scanner.Scan() )	{
		_, err := fmt.Sscanf(scanner.Text(), "%c%d%s", action, key, value )
		if ( err == nil ) {
			return( true )
		}
	}
	return( false )
}
