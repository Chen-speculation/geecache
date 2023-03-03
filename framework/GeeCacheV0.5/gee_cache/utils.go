package gee_cache

// GetGroup returns the named gee_cache previously created with NewGroup,
//
//	or nil if there's no such gee_cache.
func GetGroup(name string) IGroup {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}
