package utils

type NavigationMap[K comparable, V any] struct {
	internalMap        map[K]V
	keysInsertionOrder []K
}

func CreateNavigationMap[K comparable, V any]() *NavigationMap[K, V] {
	return &NavigationMap[K, V]{
		internalMap:        make(map[K]V),
		keysInsertionOrder: make([]K, 0),
	}
}

func (n *NavigationMap[K, V]) Put(key K, value V) {
	n.internalMap[key] = value
	n.keysInsertionOrder = append(n.keysInsertionOrder, key)
}

func (n *NavigationMap[K, V]) Get(key K) V {
	return n.internalMap[key]
}

func (n *NavigationMap[K, V]) Contains(keys ...K) bool {
	for _, key := range keys {
		if _, containedInIntermapMap := n.internalMap[key]; !containedInIntermapMap {
			return false
		}
	}

	return true
}

func (n *NavigationMap[K, V]) Keys() []K {
	return n.keysInsertionOrder[:]
}

func (n *NavigationMap[K, V]) Size() int {
	return len(n.keysInsertionOrder)
}
