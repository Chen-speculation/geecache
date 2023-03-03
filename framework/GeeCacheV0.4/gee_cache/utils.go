package gee_cache

// GetGroup returns the named gee_cache previously created with NewGroup, or nil if there's no such gee_cache.
func GetGroup(name string) IGroup {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func decodeBasePath(basePath string) string {
	if len(basePath) == 0 {
		return "/"
	}
	if basePath[0] != '/' {
		basePath = "/" + basePath
	}
	if basePath[len(basePath)-1] != '/' {
		basePath = basePath + "/"
	}
	return basePath
}
