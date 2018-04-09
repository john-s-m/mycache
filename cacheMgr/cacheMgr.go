package cacheMgr

import (
	"errors"
	"fmt"
)

type CacheItem struct {
	Value interface{}
	readCh chan interface{}
	writeCh chan interface{}
}

type CacheKeyValue struct {
	Key int
	Value interface{}
}

type SafeItemAccess interface {
	Reader() interface{}
	Writer( interface{} ) error
	serializer()
}

type CacheMap struct {
	SharedMap map [int]CacheItem
	readKeyCh chan int
	readCh chan interface{}
	writeCh chan CacheKeyValue
	insertCh chan CacheKeyValue
}

type SafeMapAccess interface {
	Reader( key int ) interface{}
	Writer( key int, value interface{} ) error
	Updater( key int, value interface{} ) error
	Inserter( key int, value interface{} ) error
	serializer()
}

type CacheMapMultiplex struct {
	SharedMap map [int]CacheItem
	lockChArray []chan int
	unlockChArray []chan int
	writeLockCh chan int
	writeUnlockCh chan int
}

type SafeMapMultiplexAccess interface {
	ReadLock( readerId int )
	ReadUnlock( readerId int )
	WriteLock( readerId int )
	WriteUnlock( readerId int )
	StartAllReaders()
}

func NewCacheItem( val interface{}, needSerialization bool ) *CacheItem {
	var cItem *CacheItem
	cItem = new(CacheItem)
	cItem.Value = val
	if needSerialization {
		cItem.readCh = make(chan interface{} )
		cItem.writeCh = make(chan interface{} )
		go cItem.serializer()
	}
	return( cItem )
}

func (ci CacheItem) Reader() interface{} {
	return <- ci.readCh
}

func (ci CacheItem) Writer( val interface{} ) error {
	ci.writeCh <- val
	return nil
}

func (ci CacheItem) serializer () {
	for {
		select {
			case ci.Value = <-ci.writeCh:
			case ci.readCh <- ci.Value:
		}
	}
}

func NewCacheMapStruct() *CacheMap {
	cm := new( CacheMap )
	cm.SharedMap = make( map [int]CacheItem )
	cm.readKeyCh = make( chan int )
	cm.readCh = make( chan interface{} )
	cm.writeCh = make( chan CacheKeyValue )
	cm.insertCh = make( chan CacheKeyValue )
	return( cm )
}

func NewCacheMap() *CacheMap {
	cm := NewCacheMapStruct()
	go cm.serializer()
	return( cm )
}

	
func (cm *CacheMap) Reader( key int ) interface{} {
	cm.readKeyCh<- key
	return <-cm.readCh
}

func (cm *CacheMap) Writer( key int, val interface{} ) error {
	if ( cm.Reader(key) == nil ) {
		return( cm.Inserter( key, val ) )
	} else {
		return( cm.Updater( key, val ) )
	}
}

func (cm *CacheMap) Updater( key int, val interface{} ) error {
	keyVal := CacheKeyValue{ key, val }
	cm.writeCh <- keyVal
	return nil
}

func (cm *CacheMap) Inserter( key int, val interface{} ) error {
	keyVal := CacheKeyValue{ key, val }
	cm.insertCh <- keyVal
	return nil
}

func (cm *CacheMap) serializer () {
	var keyVal CacheKeyValue
	var key int
	for {
		select {
		case keyVal = <-cm.writeCh:
			ci := cm.SharedMap[keyVal.Key]
			ci.Value = keyVal.Value
			cm.SharedMap[keyVal.Key] = ci
			
		case key = <-cm.readKeyCh:
			cm.readCh<- cm.SharedMap[key].Value

		case keyVal = <-cm.insertCh:
			ci := NewCacheItem( keyVal.Value, false )
			cm.SharedMap[keyVal.Key] = *ci
		}
	}
}

func NewCacheMapMultiplex( ) *CacheMapMultiplex {
	var cmm = new( CacheMapMultiplex )
	cmm.SharedMap = make( map [int]CacheItem )
	cmm.writeLockCh = make( chan int )
	cmm.writeUnlockCh = make( chan int )
	return( cmm )
}

func (cmm *CacheMapMultiplex) AddReader() ( int, error ) {
	cmm.lockChArray = append( cmm.lockChArray, make( chan int ) )
	cmm.unlockChArray = append( cmm.unlockChArray, make( chan int ) )
	i := len(cmm.lockChArray) - 1
	return i, nil
}
	
func (cmm *CacheMapMultiplex) unlockForward( inputCh chan int, outputCh chan int ) {
	for {
		select {
		case <-inputCh:
		case outputCh <- 1:
		}
	}
}

func (cmm *CacheMapMultiplex) ReadLock( readerId int ) error {
	if ( readerId >= len( cmm.lockChArray ) ) {
		return( errors.New( "cacheMgr: Reaader Id is out of range" ) )
	}
	<- cmm.lockChArray[ readerId ]
	return nil
}

func (cmm *CacheMapMultiplex) ReadUnlock( readerId int ) error {
	if ( readerId >= len( cmm.unlockChArray ) ) {
		return( errors.New( "cacheMgr: Reaader Id is out of range" ) )
	}
	if ( len( cmm.unlockChArray[ readerId ] ) > 0 ) {
		return( errors.New( "cacheMgr: Reader Id already unlocked" ) )
	}
	
	cmm.unlockChArray[ readerId ] <- 1
	return nil
}

func (cmm *CacheMapMultiplex) WriteLock( readerId int ) error {
	<- cmm.writeLockCh
	
	var secondPassCh []chan int

	for  i, ch := range cmm.lockChArray { // lock what you can on this pass without blocking

		if ( i == readerId ) {
			continue
		}
		select {
		case <- ch:
		default:
			secondPassCh = append( secondPassCh, ch )
		}
	}

	for _, ch := range secondPassCh { // block and lock everything else on this pass
		<- ch
	}

	return nil
}

func (cmm *CacheMapMultiplex) WriteUnlock( readerId int ) error {
	for i, _ := range cmm.unlockChArray { 
		if ( i == readerId ) {
			continue
		}
		cmm.unlockChArray[i] <- 1
	}

	cmm.writeUnlockCh <- 1
	return nil
}

func (cmm *CacheMapMultiplex) Reader( key int, readerId int ) ( interface{}, error ) {
	err := cmm.ReadLock( readerId )
	if ( err != nil ) {
		return nil, err
	}
	defer cmm.ReadUnlock( readerId )
	return cmm.SharedMap[ key ].Value, nil
	
}

func (cmm *CacheMapMultiplex) Writer( key int, val interface{}, readerId int ) error {
	err := cmm.WriteLock( readerId )
	if ( err != nil ) {
		return err
	}

	defer cmm.WriteUnlock( readerId )
	
	if ( cmm.SharedMap[key].Value == nil ) {
		ci := NewCacheItem( val, false )
		cmm.SharedMap[key] = *ci
	} else {
		ci := cmm.SharedMap[key]
		ci.Value = val
		cmm.SharedMap[key] = ci
	}
	return nil
}

func (cmm *CacheMapMultiplex) Updater( key int, val interface{}, readerId int ) error {
	return cmm.Writer( key, val, readerId )
}

func (cmm *CacheMapMultiplex) Inserter( key int, val interface{}, readerId int ) error {
	return cmm.Writer( key, val, readerId )
}

func (cmm *CacheMapMultiplex) StartAllRoutines() error {
	for i, _ := range cmm.lockChArray {
		go cmm.unlockForward( cmm.unlockChArray[i], cmm.lockChArray[i] )
	}

	go cmm.unlockForward( cmm.writeUnlockCh, cmm.writeLockCh )
	
	return ( cmm.WriteUnlock( -1 ) )
}

func doNothing() {
	fmt.Println( "doNothing" );
}
