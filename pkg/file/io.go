package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Status struct {
	Name     string
	FullPath string
}

// Copy ファイルを複製する
func Copy(src, dest string) (Status, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return Status{}, err
	}
	defer srcFile.Close()

	destinationPath := dest
	fileName := ""

	// .が存在しない (拡張子がない) 場合はコピー元のファイル名を複製する
	if !hasFileExtension(dest) {
		fileName = filepath.Base(src)
		destinationPath = filepath.Join(dest, fileName)
	}

	destFile, err := os.Create(destinationPath)
	if err != nil {
		return Status{}, err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return Status{}, err
	}

	return Status{
		Name:     fileName,
		FullPath: destinationPath,
	}, nil
}

func Write(path string, buffer []byte) error {
	// "target"で指定されるファイルを作成
	targetFile, e := os.Create(path)
	targetFile.Write(buffer)
	if e != nil {
		return e
	}

	if e := targetFile.Close(); e != nil {
		return e
	}
	return nil
}

// hasFileExtension 拡張子が含まれているかを判定する
func hasFileExtension(filename string) bool {
	ext := filepath.Ext(filename)
	return ext != "" && ext != "." && !strings.Contains(ext, "/") && !strings.Contains(ext, "\\")
}
