package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sort"
	
	"github.com/jung-kurt/gofpdf"
)

func allergic(err error, fmtstr string, args ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, fmtstr, args...)
		os.Exit(1)
	}
}

func enumerateDir(dir string) []string {
	dh, err := os.Open(dir)
	allergic(err, "open %s: %v\n", dir, err)
	fis, err := dh.Readdir(-1)
	allergic(err, "readdir %s: %v\n", dir, err)
	r := make([]string, 0, len(fis))
	for _, fi := range fis {
		if ext := strings.ToLower(filepath.Ext(fi.Name())); ext == ".jpg" || ext == ".jpeg" {
			r = append(r, filepath.Join(dir, fi.Name()))
		}
	}
	return r
}

func main() {
	files := os.Args[1:]
	var dir string
	
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "nessun input\n")
		os.Exit(1)
	}
	
	if len(files) == 1 {
		fi, err := os.Stat(files[0])
		allergic(err, "stat %s: %v\n", files[0], err)
		if fi.IsDir() {
			dir = files[0]
			files = enumerateDir(dir)
		} else {
			dir = filepath.Dir(files[0])
		}
	} else {
		for _, file := range files {
			fi, err := os.Stat(file)
			allergic(err, "stat %s: %v\n", file, err)
			if fi.IsDir() {
				fmt.Fprintf(os.Stderr, "%s è una directory\n", file)
				os.Exit(1)
			}
		}
		dir = filepath.Dir(files[0])
	}
	
	dst := filepath.Join(dir, filepath.Base(dir) + ".pdf")
	sort.Strings(files)
	
	fmt.Printf("%q <- %v\n", dst, files)
	
	pdf := gofpdf.New("P", "mm", "A4", "")
	
	for _, path := range files {
		info := pdf.RegisterImage(path, "")
		w, h := info.Extent()
		
		pdf.AddPageFormat("P", gofpdf.SizeType{ w, h })
		pdf.Image(path, 0, 0, w, h, false, "", 0, "")
	}
	
	err := pdf.OutputFileAndClose(dst)
	allergic(err, "write: %v", err)
}
