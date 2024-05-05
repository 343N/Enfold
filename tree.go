package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/shirou/gopsutil/disk"
)

type CompressionType int

const (
	NONE CompressionType = iota
	XPRESS_4K
	XPRESS_8K
	XPRESS_16K
	LZMA
)

type Directory struct {
	Name           string
	ParentPath     string
	IsDrive        bool
	SubDirectories []*Directory
	Files          []*File
	Compression    CompressionType
	AutoCompress   bool
	Analyzed       bool
	AnalyzedSize   int64
}

func NewDirectory() *Directory {
	return &Directory{
		Name:         "",
		ParentPath:   "",
		IsDrive:      false,
		Compression:  NONE,
		AutoCompress: false,
		Analyzed:     false,
	}
}

func (dir *Directory) AnalyzeDirectory(recurseSteps bool) {
	if !dir.Analyzed {
		return
	}
	if dir.Name == "" {
		return
	}

	filepath.WalkDir(dir.Name, func(path string, d fs.DirEntry, err error) error {
		logErr(err)
		if d.IsDir() {
			dir := NewDirectory()
			dir.Name = d.Name()
			dir.ParentPath = path
			dir.IsDrive = false
		} else {
			dInfo, err := d.Info()
			if err != nil {
				logErr(err)
				return nil
			}
			f := NewFile()
			f.Name = d.Name()
			f.ParentPath = path
			f.Extension = filepath.Ext(d.Name())
			f.Size = dInfo.Size()
		}

		return nil
	})
}

type File struct {
	Name        string
	ParentPath  string
	Extension   string
	Size        int64
	Compression CompressionType
}

func NewFile() *File {
	f := &File{
		Name:        "",
		ParentPath:  "",
		Extension:   "",
		Size:        0,
		Compression: NONE,
	}
	return f
}

var processing = true

func ResumeProcessing() {
	fmt.Println("Disabling processing")
	processing = false
}

func EnableProcessing() {
	fmt.Println("Enabling processing")
	processing = true
}

/**
 * GetDirectory returns a list of directories in the given path
 */
func (a *App) GetDirectory(dirPath string, recursive bool) []*Directory {
	if !processing {
		return make([]*Directory, 0)
	}
	if dirPath == "" || dirPath == "/" {
		dirs := make([]*Directory, 0)
		partitions, _ := disk.Partitions(false)
		for _, partition := range partitions {
			d := NewDirectory()
			d.Name = partition.Mountpoint
			d.IsDrive = true
			d.Files = make([]*File, 0)
			d.SubDirectories = make([]*Directory, 0)
			filepath.WalkDir(d.Name, func(path string, d fs.DirEntry, err error) error {
				logErr(err)
				if d.IsDir() {
					dir := NewDirectory()
					dir.Name = d.Name()
					dir.ParentPath = path
					dir.IsDrive = false
					dirs = append(dirs, dir)
				} else {
					dInfo, err := d.Info()
					if err != nil {
						logErr(err)
						continue
					}
					f := &File{
						Name:       d.Name(),
						ParentPath: path,
						Extension:  filepath.Ext(d.Name()),
						Size:       dInfo.Size(),
					}
					dirs.Files = append(dirs.Files, f)
				}
			})

		}

	}
	dirPath = filepath.Clean(dirPath)
	fmt.Println("Getting directory: " + dirPath)

	root := NewDirectory()
	root.Name = filepath.Base(dirPath)
	root.ParentPath = filepath.Dir(dirPath)

}

func (a *App) GetDirectoryChildren(dirPath string, recursive bool) []*Directory {
	if !processing {
		return make([]*Directory, 0)
	}
	fmt.Println("Getting directory: " + dirPath)
	var dirs = make([]*Directory, 0)
	if dirPath == "" || dirPath == "/" {
		partitions, _ := disk.Partitions(false)
		for _, partition := range partitions {
			d := NewDirectory()
			d.Name = partition.Mountpoint
			d.IsDrive = true
			if recursive {
				d.SubDirectories = a.GetDirectory(d.ParentPath+"/"+d.Name, recursive)
			}
			dirs = append(dirs, d)
		}
	} else {
		s_FS := os.DirFS(dirPath)
		res, err := fs.ReadDir(s_FS, ".")
		logErr(err)
		for _, entry := range res {
			if entry.IsDir() {
				d := NewDirectory()
				d.Name = entry.Name()
				if recursive {
					d.SubDirectories = a.GetDirectory(d.ParentPath+"/"+d.Name, recursive)
				}
				dirs = append(dirs, d)
			} else {
				f := &File{
					Name:       entry.Name(),
					ParentPath: dirPath,
				}
				f = append(dirs.Files, f)
			}
		}
	}
	fmt.Printf("Returning %d directories on %s.\n", len(dirs), dirPath)
	return dirs
}
