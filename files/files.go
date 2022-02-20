package files

import (
	//"bytes"
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"fmt"
)

// fileExists returns whether the given file or directory exists or not.
// Taken from https://stackoverflow.com/a/10510783
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func UnpackFiles(files embed.FS, directory string) error {
	err := os.MkdirAll(directory, 0755)
	if err != nil {
		return err
	}

	wlk := walker{
		base: directory,
		filesystem: files,
	}

	err = fs.WalkDir(files, ".", wlk.walkUnpack)
	if err != nil {
		return err
	}

	return nil
}

func UnpackFileBytes(data []byte, dest string) error {
	if fileExists(dest) {
		return nil
	}

	fmt.Println("Unpacking", dest)

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	//in := bytes.NewReader(data)
	//_, err = io.Copy(out, in)
	_, err = out.Write(data)
	if err != nil {
		return err
	}
	return nil
}

type walker struct {
	base string
	filesystem fs.FS
}

func (w *walker) walkUnpack(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	fullname := filepath.Join(w.base, path)

	if fileExists(fullname) {
		return nil
	}

	if d.IsDir() {
		info, err := d.Info()
		if err != nil {
			return err
		}

		err = os.Mkdir(fullname, info.Mode())
		if err != nil {
			return err
		}
		return nil
	}

	fmt.Println("Unpacking", path)
	out, err := os.Create(fullname)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := w.filesystem.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}
