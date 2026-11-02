package utils

import (
	"os"
	"path/filepath"
)

// IsBadgerDir 方法并不可靠，请把数据库文件按照约定存储
func IsBadgerDir(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	// 第一层级，是多个数据库（多 network）
	for _, e := range entries {
		entries, err := os.ReadDir(filepath.Join(path, e.Name()))
		if err != nil {
			return false, err
		}

		// 第二层级才是数据库文件
		for _, e := range entries {
			// 该文件 badger db 专属
			if filepath.Ext(e.Name()) == ".vlog" {
				return true, nil
			}
		}
	}
	return false, nil
}
