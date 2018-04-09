package main

import (
	"math/rand"
	"fmt"
	"bytes"
	"strconv"
	"bufio"
	"io"
	"errors"
)


type ActionInfo struct {
	Action byte
	Key int
	Value interface{}
}

type InjectorFunction func(*interface{},*ActionInfo )(error)

type ActionItem struct {
	Parameters interface{}
	ActionValues ActionInfo
	Injector InjectorFunction
}

func randomInjector( parms *interface{}, pA *ActionInfo )(error) {
	p := *parms
	seededPtr, ok := p.(*rand.Rand)

	if false {
		fmt.Println("Never gets here, just want to import fmt")
	}
	
	if ! ok {
		return errors.New( "randomInjector: Could not convert parms *interface{} to a rand.Rand" )
	}

	if seededPtr == nil {
		return errors.New( "randomInjector: nil Parameters/*rand.Rand value" )
	}

	rValue := seededPtr.Int()
	pA.Action = 'r'
	if ( ( rValue % 4 ) == 0 ) {
		pA.Action = 'w'
		rValue = seededPtr.Int() % 10
		switch {
		case rValue <= 4:
			pA.Value = seededPtr.Int() % 10000

		case rValue < 8:
			pA.Value = seededPtr.Float64()

		default:
			var buffer bytes.Buffer
			buffer.WriteString("String ")
			buffer.WriteString(strconv.Itoa( seededPtr.Int() % 100 ) )
			pA.Value = buffer.String()
		}
	}

	rValue = seededPtr.Int()
	pA.Key = rValue % 100

	rValue = seededPtr.Int()
	if ( rValue % 10000 ) == 0 {
		fmt.Printf( "exit random: %d\n", rValue )
		return ( io.EOF )
	}

	return nil
}


func fileInjector( parms *interface{}, pA *ActionInfo )(error) {
	p := *parms
	scanner, ok := p.(*bufio.Scanner)
	if ! ok {
		return errors.New( "fileInjector: Could not convert parms *interface{} to a bufio.Scanner" )
	}

	if scanner == nil {
		return errors.New( "fileInjector: bufio.Scanner is nil" )
	}
	
	if ! readAction( scanner, &pA.Action, &pA.Key, &pA.Value ) {
		return io.EOF
	}

	return nil
}