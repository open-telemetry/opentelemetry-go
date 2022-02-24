package internal

import (
	"fmt"
	"path"
	"strings"
)

// CleanPath returns cleaned URL path. Replace with default path if path is nil
func CleanPath(URLPath string, defaultPath string) string {
	tmp := strings.TrimSpace(URLPath)
	if tmp == "" {
		return defaultPath
	} else {
		tmp = path.Clean(tmp)
		if !path.IsAbs(tmp) {
			tmp = fmt.Sprintf("/%s", tmp)
		}
	}
	return tmp
}
