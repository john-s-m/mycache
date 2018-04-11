package main

import (
	"math/rand"
	"fmt"
	"bytes"
	"strconv"
	"bufio"
	"io"
	"errors"
	"time"
	"os"
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

func NewFileActor( dataFilePrefix string, goRoutineId int ) (*ActionItem, *os.File) {
	var dataFile string

	dataFile = fmt.Sprintf( "%s%d.dat", dataFilePrefix, goRoutineId )
	f, err := os.Open(dataFile)
	if err != nil {
		return nil, nil
	}

	pA := new( ActionItem )
	pA.Injector = fileInjector
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	pA.Parameters = scanner
	return pA, f

}

type RandomInjectorParameters struct {
	randomizer        *rand.Rand
	remainingEvents   int
}

func NewRandomActor( eventCount int ) *ActionItem {
	pA := new( ActionItem )
	pParms := new( RandomInjectorParameters )
	pParms.remainingEvents = eventCount
	pParms.randomizer = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	pA.Parameters = pParms
	pA.Injector = randomInjector
	return pA
}

func randomInjector( parms *interface{}, pA *ActionInfo )(error) {
	p := *parms
	pParms, ok := p.(*RandomInjectorParameters)

	if ! ok {
		return errors.New( "randomInjector: Could not convert parms *interface{} to a *RandomInjectorParamters" )
	}

	if ( pParms.remainingEvents == 0 ) {
		return ( io.EOF )
	}
	pParms.remainingEvents--

	seededPtr := pParms.randomizer

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
			pA.Value = seededPtr.Float64() * 10000

		default:
			var buffer bytes.Buffer
			buffer.WriteString("String ")
			buffer.WriteString(strconv.Itoa( seededPtr.Int() % 100 ) )
			pA.Value = buffer.String()
		}
	}

	rValue = seededPtr.Int()
	pA.Key = rValue % 100

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

type ArrayInjectorParameters struct {
	actionList  []ActionInfo
	pos         int
}

func NewArrayActor( preparedActionArrray []ActionInfo ) *ActionItem {
	pA := new( ActionItem )
	pParms := new( ArrayInjectorParameters )
	pParms.actionList = preparedActionArrray
	pParms.pos = 0
	pA.Parameters = pParms
	pA.Injector = arrayInjector
	return pA
}

func arrayInjector( parms *interface{}, pA *ActionInfo )(error) {
	p := *parms
	pParms, ok := p.(*ArrayInjectorParameters)

	if ! ok {
		return errors.New( "arrayInjector: Could not convert parms *interface{} to a *ArrayInjectorParamters" )
	}

	if ( pParms.pos == len(pParms.actionList) ) {
		return ( io.EOF )
	}

	*pA = pParms.actionList[pParms.pos]
	pParms.pos++
	return nil
}

func doNothing002() {
	fmt.Println( "doNothing" ); // just so it doesn't complain about importing fmt when I remove all prints
}

	
