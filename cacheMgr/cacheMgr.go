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

func NewCacheMap() *CacheMap {
	cm := new( CacheMap )
	cm.SharedMap = make( map [int]CacheItem )
	cm.readKeyCh = make( chan int )
	cm.readCh = make( chan interface{} )
	cm.writeCh = make( chan CacheKeyValue )
	cm.insertCh = make( chan CacheKeyValue )
	go cm.serializer()
	return( cm )
}

	
func (cm CacheMap) Reader( key int ) interface{} {
	cm.readKeyCh<- key
	return <-cm.readCh
}

func (cm CacheMap) Writer( key int, val interface{} ) error {
	if ( cm.Reader(key) == nil ) {
		return( cm.Inserter( key, val ) )
	} else {
		return( cm.Updater( key, val ) )
	}
}

func (cm CacheMap) Updater( key int, val interface{} ) error {
	keyVal := CacheKeyValue{ key, val }
	cm.writeCh <- keyVal
	return nil
}

func (cm CacheMap) Inserter( key int, val interface{} ) error {
	keyVal := CacheKeyValue{ key, val }
	cm.insertCh <- keyVal
	return nil
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

