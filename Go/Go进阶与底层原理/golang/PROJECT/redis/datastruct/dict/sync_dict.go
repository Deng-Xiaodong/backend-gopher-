package dict

import "sync"

type SysnDict struct {
	m sync.Map
}

func MakeSyncDict() *SysnDict {
	return &SysnDict{}
}
func (dict *SysnDict) Get(key string) (val any, exists bool) {
	val, exists = dict.m.Load(key)
	return
}

func (dict *SysnDict) Len() int {
	lenth := 0
	dict.m.Range(func(key, value any) bool {
		lenth++
		return true
	})
	return lenth
}

func (dict *SysnDict) Put(key string, val any) (result int) {
	_, existed := dict.m.Load(key)
	dict.m.Store(key, val)
	if existed {
		return 0
	}
	return 1
}

func (dict *SysnDict) PutIfAbsent(key string, val any) (result int) {
	_, existed := dict.m.Load(key)
	if existed {
		return 0
	}
	dict.m.Store(key, val)
	return 1
}

func (dict *SysnDict) PutIfExists(key string, val any) (result int) {
	_, existed := dict.m.Load(key)
	if existed {
		dict.m.Store(key, val)
		return 1
	}

	return 0
}

func (dict *SysnDict) Remove(key string) bool {
	_, existed := dict.m.Load(key)
	if existed {
		dict.m.Delete(key)
		return true
	}
	return false

}

//类型转换  强转和断言
func (dict *SysnDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value any) bool {
		return consumer(key.(string), value)
	})
}

func (dict *SysnDict) Keys() []string {
	keys := make([]string, dict.Len())
	i := 0
	dict.m.Range(func(key, value any) bool {
		keys[i] = key.(string)
		i++
		return true
	})
	return keys
}

func (dict *SysnDict) RandomKeys(limit int) []string {
	result := make([]string, limit)
	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value any) bool {
			result[i] = key.(string)
			return false
		})
	}
	return result
}

func (dict *SysnDict) RandomDistinctKeys(limit int) []string {
	result := make([]string, limit)
	i := 0
	dict.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		return true
	})
	return result
}

//ERROR Assignment to the method receiver propagates only to callees but not to callers
func (dict *SysnDict) Clear() {
	*dict = *MakeSyncDict()
}
