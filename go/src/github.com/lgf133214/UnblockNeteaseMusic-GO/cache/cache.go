package cache

import (
	"github.com/lgf133214/UnblockNeteaseMusic-GO/common"
	"sync"
)

var cache sync.Map

func Put(key interface{}, value interface{}) {
	cache.Store(key, value)
}
func GetSong(key interface{}) (common.Song, bool) {
	var song common.Song
	if value, ok := cache.Load(key); ok {
		song = value.(common.Song)
		return song, ok
	}
	return song, false
}
func Delete(key interface{}) {
	cache.Delete(key)
}
