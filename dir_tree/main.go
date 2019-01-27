package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	//fmt.Println("*" + strings.Repeat("\t", 1), "test")
	//panic("end")

	start := time.Now()
	out := os.Stdout

	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}

	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles, -1)
	if err != nil {
		panic(err.Error())
	}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Elapsed: ", elapsed)
}

func dirTree(out *os.File, path string, printFiles bool, level int) error {
	directories, err := ioutil.ReadDir(path)
	filesCount := len(directories)
	//currentDirNum := 0

	if err != nil {
		fmt.Println(err)
	}

	/*fmt.Println(len(directories));
	panic("die");*/

	//fmt.Println(level)

	if filesCount > 0 {
		level++
	} else {
		level--
	}

	for _, dir := range directories {
		//currentDirNum := index + 1
		//fmt.Println("├───"+strings.Repeat("	", level), dir.Name(), "level", level, "files", filesCount, "idx ", currentDirNum)

		if !printFiles && !dir.IsDir() {
			continue
		}

		printLines(dir.Name(), level, filesCount)

		if dir.IsDir() {
			dirPath := path + string(os.PathSeparator) + dir.Name()
			err = dirTree(out, dirPath, printFiles, level)
		}
	}

	return err
}

func printLines(dirName string, level int, filesCount int) {
	switch level {
	case 0:
		fmt.Println("├───"+strings.Repeat("", level), dirName)
	case 1:
		fmt.Println(strings.Repeat("\t", level)+"├───", dirName)
	default:
		fmt.Println(strings.Repeat("│\t", level)+"├───", dirName)
	}

	/*if level > 1 {
		fmt.Println(strings.Repeat("│\t", level) + "├───", dirName)
	} else if level == 1 {
		fmt.Println(strings.Repeat("\t", level) + "├───", dirName)
	} else {
		fmt.Println("├───" + strings.Repeat("", level), dirName)
	}*/
}
