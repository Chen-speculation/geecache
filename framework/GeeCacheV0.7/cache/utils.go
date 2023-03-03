package cache

func CloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func DecodeBasePath(basePath string) string {
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
