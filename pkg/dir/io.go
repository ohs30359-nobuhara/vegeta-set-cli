package dir

import (
	"ohs30359/vegeta-cli/pkg/file"
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

// Copy ディレクトリを複製する
func Copy(src, dest string) error {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, path[len(src):])
		if info.IsDir() {
			err := os.MkdirAll(destPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			_, err := file.Copy(path, destPath)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
