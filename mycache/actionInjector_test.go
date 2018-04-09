package main

import (
	"testing"
	"reflect"
	"strings"
)

func TestNewRandomActor( t *testing.T) {
	pAction := NewRandomActor(10)
	if ( pAction == nil ) {
		t.Errorf( "NewRandomActor failed, got nil pointer back\n" );
		return
	}

/*
	if ( pAction.Injector == nil ) {
		t.Errorf( "NewRandomActor failed, expected Injector to be %q got %q \n", randomInjector, pAction.Injector )
	}
*/
	var s string = reflect.TypeOf( pAction.Parameters ).String()
	if ( strings.Compare(s, "*main.RandomInjectorParameters" ) != 0 ) {
		t.Errorf( "NewRandomActor failed, expected Parameter to have a \"*main.RandomInjectorParameters\" got \"%s\" \n", s )
	}

	pParms, _ := pAction.Parameters.(*RandomInjectorParameters)
	if ( pParms.remainingEvents != 10 ) {
		t.Errorf( "NewRandomActor failed, expected Parameter.remainingEvents to be 10 got %d \n", pParms.remainingEvents )
	}

	if ( pParms.randomizer == nil ) {
		t.Errorf( "NewRandomActor failed to initialize the randomizer\n")
	}
}

func TestrandomInjector(t *testing.T) {
	pAction := NewRandomActor(100)

	var i int = 0
	for  {
		if ( i > 100 ) {
			t.Errorf( "TestrandomInjector failed, Injector did not return EOF\n" )
		}

		err := pAction.Injector( &(pAction.Parameters), &(pAction.ActionValues) )

		if ( err != nil ) {
			if ( strings.Compare(err.Error(), "EOF" ) != 0 ) {
				t.Errorf( "TestrandomInjector failed, expected values error: %s \n", err.Error() )
			} else {
				break
			}
		}

		if ( pAction.ActionValues.Key < 0 ) || ( pAction.ActionValues.Key > 99  ) {
			t.Errorf( "TestrandomInjector - Key out of range: %d\n", pAction.ActionValues.Key )
		}
		if ( pAction.ActionValues.Action != 'r' ) && ( pAction.ActionValues.Action != 'w' ) {
			t.Errorf( "TestrandomInjector - Invalid Action: %d\n", pAction.ActionValues.Action )
		}
		if ( pAction.ActionValues.Value == nil ) {
			t.Errorf( "TestrandomInjector - nil Value\n" )
		}
		i++
	}
}
