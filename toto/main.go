package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/toto/cmd"
	"github.com/spf13/toto/config"
)

type infoTable struct {
    Name string 
    TimeModif time.Time 
    CheckSum string 
}

func calculateChecksum(filePath string) (string, error) {
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
    if err != nil {
		fmt.Printf("Cannot open %s to calculate checksum\n", filePath)
        return "", err
    }
    defer file.Close()
    hash := md5.New()
    if _, err := io.Copy(hash, file); err != nil {
		fmt.Printf("Cannot create checksum for %s\n", filePath)
        return "", err
    }
    return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

var storedFileInfo []infoTable

func readPrevData(dest string) error{

    dataFile, err1 := os.OpenFile(dest + "\\data.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
    if err1 != nil {
		fmt.Printf("Cannot open/create data.txt\n")
        return err1
	}

	reader := bufio.NewReader(dataFile)
	i := 0
	for {
		line, err := reader.ReadString('\n')
		if len(line) == 0 {
			break
		}
		if err != nil {
			fmt.Printf("Error in reading line %d of data.txt\n", i)
			return err
		}
		fields := strings.Split(line, ",")
		parsedTime, err := time.Parse("2006-01-02 15:04:05", fields[1])
		if err != nil {
			fmt.Printf("Error parsing time in line %d of data.txt: %v\n", i, err)
			return err
		}
		newInfo := infoTable{
			Name: fields[0],
			TimeModif: parsedTime,
			CheckSum: fields[2],
		}
		storedFileInfo = append(storedFileInfo, newInfo)
	}
	dataFile.Close()
	return nil
}

func changesMade(info os.FileInfo, dataFile *os.File, fileChecksum string) (bool, error, bool) {
    index := -1
    found := false
    for i, fileInfo := range storedFileInfo {
        if fileInfo.Name == info.Name() {
            index = i
            found = true
            break
        }
    }
	if !found {
		_, err := fmt.Fprintf(dataFile, "%s,%s,%s\n", info.Name(), info.ModTime().Format("2006-01-02 15:04:05"), fileChecksum)
		if err != nil {
			fmt.Printf("Cannot append %s to data.txt\n", info.Name())
			return false, err, false
		}
		newInfo := infoTable{
			Name:     info.Name(),
			TimeModif: info.ModTime(),
			CheckSum: fileChecksum,
		}
		storedFileInfo = append(storedFileInfo, newInfo)
		return true, nil, false
	}
	
	if storedFileInfo[index].TimeModif == info.ModTime() {
		return false, nil, false
	} else if info.IsDir(){
		return true, nil, false
	} else{
		if storedFileInfo[index].CheckSum == fileChecksum {
			return false, nil, false
		} else {
			storedFileInfo[index].TimeModif = info.ModTime()
			storedFileInfo[index].CheckSum = fileChecksum
			return true, nil, true
		}
	} 
}


func copySingleFile(src, dest string, log *os.File) error{
	source, err := os.Open(src)
	if err != nil {
		fmt.Printf("Cannot open %s\n", src)
		return err
	}
	defer source.Close()

	destination, err := os.Create(dest + "\\" + filepath.Base(src))
	if err != nil {
		fmt.Printf("Cannot create %s\n", destination.Name())
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)

	if err != nil {
		return err
	} else {
		fmt.Fprintf(log, "File %s is copied successfully\n", src)
		fmt.Printf("File %s is copied successfully\n", src)
	}
	return nil
}

func copyEntireDir(src, dest string, log *os.File) error {
	srcInfo, err := os.Stat(src)
	if(err != nil) {
		fmt.Println(err)
		fmt.Printf("Cannot find file info for %s\n", src)
		return err
	}
	var checksum string
	if !srcInfo.IsDir() {
		checksum, err = calculateChecksum(src)
		if(err != nil) {
			fmt.Println(err)
			fmt.Printf("Cannot calculate checksum for %s\n", src)
			return err
		}
	}

	var dataFile *os.File
    var err1 error
    dataFile, err1 = os.OpenFile(dest + "\\data.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
    if err1 != nil {
		fmt.Printf("Cannot open data.txt in copying\n")
        return err1
	}
	var modify bool
	var proceed bool
	if proceed, err, modify = changesMade(srcInfo, dataFile, checksum); err != nil{
		fmt.Println(err)
		fmt.Printf("Cannot figure out if %s was modified\n", srcInfo.Name())
		return err
	} else if !proceed {
		return nil
	}

	if !srcInfo.IsDir() {
		err := copySingleFile(src, dest, log)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("Cannot copy %s to %s\n", src, dest)
			return err
		} else if modify{
			fmt.Fprintf(log,"File %s is modified successfully\n", src)
			fmt.Printf("File %s is modified successfully\n", src)
		} else {
			fmt.Fprintf(log,"File %s is copied successfully\n", src)
			fmt.Printf("File %s is copied successfully\n", src)
		}
		return nil
	}

	files, err := os.ReadDir(src)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Cannot read %s\n", src)
		return err
	}

	for _, file := range files {
		fullPathSrc := filepath.Join(src,"\\", file.Name())
		fullPathDest := filepath.Join(dest,"\\", file.Name())

		fileInfo, err := os.Stat(fullPathSrc)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("Cannot find file info for %s\n", fullPathSrc)
			return err
		}
		var checksum string
		if !fileInfo.IsDir() {
			checksum, err = calculateChecksum(fullPathSrc)
			if(err != nil) {
				fmt.Printf("Cannot calculate checksum for %s\n", fullPathSrc)
				fmt.Println(err)
				return err
			}
		}
		var proceed bool
		if proceed, err, modify = changesMade(fileInfo, dataFile, checksum); err != nil{
			fmt.Println(err)
			fmt.Printf("Cannot figure out if %s was modified\n", fileInfo.Name())
			return err
		} else if !proceed {
			continue
		}
		if fileInfo.IsDir() {
			err := os.MkdirAll(fullPathDest, 0755)
			if err != nil {
				fmt.Println(err)
				fmt.Printf("Cannot create directory %s\n", fullPathDest)
				return err
			}
			err = copyEntireDir(fullPathSrc, fullPathDest, log)
			if err != nil {
				fmt.Println(err)
				fmt.Printf("Cannot copy %s to %s\n", fullPathSrc, fullPathDest)
				return err
			} else if modify {
				fmt.Fprintf(log, "Directory %s is modified successfully\n", src)
				fmt.Printf("Directory %s is modified successfully\n", src)
			} else {
				fmt.Fprintf(log, "Directory %s is copied successfully\n", src)
				fmt.Printf("Directory %s is copied successfully\n", src)
			}
		} else {
			if err:= copySingleFile(fullPathSrc, dest, log); err != nil {
				fmt.Println(err)
				fmt.Printf("Cannot copy %s to %s\n", fullPathSrc, dest)
				return err
			} else if modify {
				fmt.Fprintf(log, "Directory %s is modified successfully\n", src)
				fmt.Printf("Directory %s is modified successfully\n", src)
			} else {
				fmt.Fprintf(log, "Directory %s is copied successfully\n", src)
				fmt.Printf("Directory %s is copied successfully\n", src)
			}
		}
	}
	return nil
}

func main() {
	cmd.Execute()
	src := config.SourcePath
	dest := config.DestinPath + "\\backup"
	os.MkdirAll(dest, 0777)
	log, err := os.Create(dest + "\\log.txt")
	
	if err != nil {
		fmt.Printf("Cannot create log.txt\n")
	    panic(err)
	}
	defer log.Close()
    readPrevData(dest)
	if err:= copyEntireDir(src, dest, log); err != nil {
		fmt.Printf("Cannot copy %s to %s\n", src, dest)
		panic(err)
	} else {
		fmt.Fprintf(log, "Directory %s is backed up successfully\n", src)
		fmt.Printf("Directory %s is backed up successfully\n", src)
	}
}
