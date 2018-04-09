package main

import (
	"testing"
	"reflect"
	"strings"
)

func TestNewRandomActor( t *testing.T) {
	pAction := NewRandomActor()
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
	if ( strings.Compare(s, "*rand.Rand" ) != 0 ) {
		t.Errorf( "NewRandomActor failed, expected Parameter to have a \"*rand.Rand\" got \"%s\" \n", s )
	}
}

func TestrandomInjector(t *testing.T) {
	pAction := NewRandomActor()

	for i :=0; i < 100; i++ {
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
	}
}
	
