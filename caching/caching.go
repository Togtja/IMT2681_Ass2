package caching

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"RESTGvkGitLab/globals"
)

//ShouldFileCache Check if a file should be cached
//(it will create a file if it does not exist from before)
func ShouldFileCache(filename string, dir string) (globals.FileMsg, *os.File) {
	status, file := CreateFile(filename, dir)
	//If it doesn't exist we have either created it or
	//returned a nil
	path := dir + "/" + filename
	if status != globals.Exist {
		return status, file
	}

	//The file exist and we need to see how old it is
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		file.Close()
		return globals.Error, nil
	}

	if fileAge(path, info, globals.DeleteAge) {
		fmt.Println("Cache is old, Run Update")
		file.Close() //A new file will be created
		status := DeleteFile(path)
		if status == globals.Deleted {
			return simpleCreateFile(path)
		}
		return status, nil
	}
	fmt.Println("Cache is recent, No need to update")
	return globals.Exist, file
}

//CreateFile creates file in directory a file
func CreateFile(filename string, dir string) (globals.FileMsg, *os.File) {
	path := dir + "/" + filename
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Can not open file")
		if dir != "" {
			fmt.Println("Trying to create directory")
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Println("Failed to create directory")
				file.Close()
				return globals.DirFail, nil
			}
		}
		//Create the file
		file.Close()
		fmt.Println("trying to create", path)
		file, err = os.Create(path)
		if err != nil {
			file.Close()
			return globals.Error, nil
		}
		fmt.Println(path, "created")
		return globals.Created, file
	}
	return globals.Exist, file
}

//CacheStruct creates a file in the dir with the struct as a json
func CacheStruct(file *os.File, v interface{}) {
	vBytes, _ := json.Marshal(v)
	file.Write(vBytes)
	fmt.Println("We are done")
	file.Close()
}

//FileExist Sees if file exist, if it does return it
func FileExist(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	return file
}

//ReadFile read from file to struct
func ReadFile(file *os.File, v interface{}) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	json.Unmarshal(data, &v)
	file.Close()
	return err
}

//DeleteFile Given a filename and dir
func DeleteFile(path string) globals.FileMsg {
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return globals.Error
	}
	return globals.Deleted
}

//Does not to all the cheking that CreateFile does
func simpleCreateFile(path string) (globals.FileMsg, *os.File) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return globals.Error, nil
	}
	return globals.OldRenew, file
}

//CleanUp Goes through all files in dir and delete files that are older than time
//Also include semi optional publicFile which should be the public file, and will not be cleaned
func CleanUp(dir string, age float64, publicFile string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		path := dir + "/" + f.Name()
		if f.Name() == publicFile {
			continue
		}
		if fileAge(path, f, age) {
			DeleteFile(path)
		}
	}
}

//CleanUpPrivateAll Cleans Everything older than age except public files
func CleanUpPrivateAll(age float64) {
	fmt.Println("Running Cleanup!")
	CleanUp(globals.COMMITDIR, age, globals.PUBLIC+globals.COMMITFILE)
	CleanUp(globals.LANGDIR, age, globals.PUBLIC+globals.LANGFILE)
	CleanUp(globals.PROJIDDIR, age, globals.PUBLIC+globals.PROJIDFILE)
}

//CleanUpAll Cleans Everything older than age including public files
func CleanUpAll(age float64) {
	fmt.Println("Running Cleanup!")
	CleanUp(globals.COMMITDIR, age, "")
	CleanUp(globals.LANGDIR, age, "")
	CleanUp(globals.PROJIDDIR, age, "")
}

//FileAge checks the age of the file(info) is hours older that age
func fileAge(path string, info os.FileInfo, age float64) bool {
	//No need to check if older
	if age == -1 {
		return false
	}
	timenow := time.Now()
	mtime := info.ModTime()
	fmt.Println("Checking dates on", path)
	fmt.Println("Is", timenow.Sub(mtime).Hours(), "larger than", age, "?")
	//Does not care for timezones btw
	if timenow.Sub(mtime).Hours() > age {
		return true
	}
	return false
}

//CleanUpInterval cleans file ever age hour, checking if file is fileage old
func CleanUpInterval(age float64, fileage float64) {
	ticker := time.NewTicker(time.Duration(age) * time.Hour)
	//can be closed by running close(quit)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				CleanUpPrivateAll(fileage)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
