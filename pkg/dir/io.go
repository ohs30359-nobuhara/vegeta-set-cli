package dir

import (
	"os"
	"path/filepath"
)

// Scan 指定したパス配下のファイル一覧を取得する
func Scan(path string) ([]string, error) {
	var fileList []string
	if e := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileList = append(fileList, p)
		}
		return nil
	}); e != nil {
		return nil, e
	}

	return fileList, nil
}
