package routineEngine

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

type Engine interface {
	Get() unsafe.Pointer
	Set(p unsafe.Pointer)
}

type engineMap struct {
	Engine
	_id_map map[int64]unsafe.Pointer
	lock    *sync.RWMutex
}

func getGoroutineId() int64 {
	buffer := make([]byte, 32)
	runtime.Stack(buffer, false)
	if parts := split(string(buffer), " "); parts[0] == "goroutine" && len(parts) > 1 {
		if goid, err := strconv.ParseInt(parts[1], 10, 0); err == nil {
			return goid
		}
	}
	return -1
}

func (g *engineMap) init() *engineMap {
	g._id_map = map[int64]unsafe.Pointer{}
	g.lock = new(sync.RWMutex)
	return g
}

func (g *engineMap) Get() unsafe.Pointer {
	gid := getGoroutineId()
	if gid == -1 {
		return nil
	}
	g.lock.RLock()
	result := unsafe.Pointer(nil)
	if r, found := g._id_map[gid]; found {
		result = r
	}
	g.lock.RUnlock()
	return result
}

func (g *engineMap) Set(p unsafe.Pointer) {
	gid := getGoroutineId()
	if gid == -1 {
		return
	}
	g.lock.Lock()
	if p != nil {
		g._id_map[gid] = p
	} else if _, exist := g._id_map[gid]; exist {
		delete(g._id_map, gid)
	}
	g.lock.Unlock()
}

var _gidmap *engineMap

func Get() Engine {
	return _gidmap
}

func init() {
	_gidmap = (&engineMap{}).init()
}

func split(s, sep string) []string {
	sepLen := len(sep)
	if sepLen == 0 {
		return []string{s}
	}
	count := 0
	index := 0
	for i := strings.Index(s[index:], sep); index < len(s); i = strings.Index(s[index:], sep) {
		if i != 0 {
			count++
		}
		if i < 0 {
			break
		}
		index += i + sepLen
	}
	if count == 0 {
		return []string{}
	}
	r := make([]string, count)
	index = 0
	for i := 0; i < count; i++ {
		sepIndex := 0
		for sepIndex = strings.Index(s[index:], sep); sepIndex == 0; sepIndex = strings.Index(s[index:], sep) {
			index += sepLen
		}
		if sepIndex > 0 {
			r[i] = string(s[index : index+sepIndex])
		} else {
			r[i] = string(s[index:])
		}
		index += sepIndex
	}
	return r
}