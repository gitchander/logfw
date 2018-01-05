package logfw

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type backupFormat struct {
	dir    string
	prefix string
	ext    string

	format string
}

func (bf *backupFormat) SetFileName(fileName string) {
	dir, name := filepath.Split(fileName)
	ext := filepath.Ext(name)
	prefix := name[:len(name)-len(ext)]
	*bf = backupFormat{
		dir:    dir,
		prefix: prefix,
		ext:    ext,
	}
}

const defaultFormat = "%s_%d%s"

func (bf *backupFormat) SetNumberLen(n int) {
	if (1 < n) && (n < 10) {
		bf.format = "%s_%0" + strconv.Itoa(n) + "d%s"
	} else {
		bf.format = defaultFormat
	}
}

func (bf *backupFormat) BackupName(number int) string {
	format := bf.format
	if format == "" {
		format = defaultFormat
	}
	name := fmt.Sprintf(format, bf.prefix, number, bf.ext)
	return filepath.Join(bf.dir, name)
}

func (bf *backupFormat) ParseNumber(fileName string) (number int) {
	if strings.HasPrefix(fileName, bf.prefix) && strings.HasSuffix(fileName, bf.ext) {
		s := fileName[len(bf.prefix) : len(fileName)-len(bf.ext)]
		if (len(s) >= 2) && (s[0] == '_') {
			n, err := strconv.Atoi(s[1:])
			if err == nil {
				return n
			}
		}
	}
	return -1
}

func renameAndRotate(fileName string, maxBackups int) (rotateDone chan struct{}, err error) {

	var bf backupFormat
	bf.SetFileName(fileName)

	err = os.Rename(fileName, bf.BackupName(0))
	if err != nil {
		return nil, err
	}

	rotateDone = make(chan struct{})
	go func() {
		err = rotate(bf, maxBackups)
		close(rotateDone)
		if err != nil {
			panic(err)
		}
	}()

	return rotateDone, nil
}

func rotate(bf backupFormat, maxBackups int) error {

	files, err := backupFiles(&bf)
	if err != nil {
		return err
	}

	sort.Sort(sort.Reverse(byNumber(files)))

	needSetLen := true
	for _, file := range files {
		oldpath := filepath.Join(bf.dir, file.name)
		newnum := file.number + 1
		if newnum > maxBackups {
			err := os.Remove(oldpath)
			if err != nil {
				return err
			}
			continue
		}
		if needSetLen {
			bf.SetNumberLen(countDigits(newnum))
			needSetLen = false
		}
		newpath := bf.BackupName(newnum)
		err = os.Rename(oldpath, newpath)
		if err != nil {
			return err
		}
	}

	return nil
}

func countDigits(x int) int {
	var (
		c = 1
		d = 10
	)
	for d <= x {
		c++
		d *= 10
	}
	return c
}

type backupFileInfo struct {
	name   string
	number int
}

type byNumber []backupFileInfo

func (p byNumber) Len() int {
	return len(p)
}

func (p byNumber) Less(i, j int) bool {
	return p[i].number < p[j].number
}

func (p byNumber) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func backupFiles(bf *backupFormat) ([]backupFileInfo, error) {
	dir := bf.dir
	if dir == "" {
		dir = "." // current directory
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	fis := make([]backupFileInfo, 0, len(files))
	for _, file := range files {
		name := file.Name()
		if n := bf.ParseNumber(name); n != -1 {
			fi := backupFileInfo{
				name:   name,
				number: n,
			}
			fis = append(fis, fi)
		}
	}
	return fis, nil
}
