package cacheMgr

type CacheItem struct {
	Value interface{}
	readCh chan interface{}
	writeCh chan interface{}
}

type CacheKeyValue struct {
	Key int
	Value interface{}
}

type CacheMap struct {
	SharedMap map [int]CacheItem
	readCh chan interface{}
	readKeyCh chan int
	writeCh chan CacheKeyValue
	insertCh chan CacheKeyValue
}

type SafeItemAccess interface {
	Reader() interface{}
	Writer( interface{} )
	serializer()
}

type SafeMapAccess interface {
	Reader( key int ) interface{}
	Writer( key int, value interface{} )
	Inserter( key int, ci *CacheItem )
	serializer()
}

func NewCacheItem( val interface{}, goRtn bool ) *CacheItem {
	var cItem *CacheItem
	cItem = new(CacheItem)
	cItem.Value = val
	if goRtn {
		cItem.readCh = make(chan interface{} )
		cItem.writeCh = make(chan interface{} )
		go cItem.serializer()
	}
	return( cItem )
}

func (ci CacheItem) Reader() interface{} {
	return <- ci.readCh
}

func (ci CacheItem) Writer( val interface{} ) {
	ci.writeCh <- val
}

func (ci CacheItem) serializer () {
	for {
		select {
			case ci.Value = <-ci.writeCh:
			case ci.readCh <- ci.Value:
		}
	}
}

func NewCacheMap() *CacheMap {
	cm := new( CacheMap )
	cm.SharedMap = make( map [int]CacheItem )
	cm.readCh = make( chan interface{} )
	cm.readKeyCh = make( chan int )
	cm.writeCh = make( chan CacheKeyValue )
	cm.insertCh = make( chan CacheKeyValue )
	go cm.serializer()
	return( cm )
}

	
func (cm CacheMap) Reader( key int ) interface{} {
	cm.readKeyCh<- key
	return <-cm.readCh
}

func (cm CacheMap) Writer( key int, val interface{} ) {
	keyVal := CacheKeyValue{ key, val }
	cm.writeCh <- keyVal
}

func (cm CacheMap) Inserter( key int, val interface{} ) {
	keyVal := CacheKeyValue{ key, val }
	cm.insertCh <- keyVal
}

func (cm CacheMap) serializer () {
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

