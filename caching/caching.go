package caching

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		return globals.Error, nil
	}
	mtime := info.ModTime()
	timenow := time.Now()
	fmt.Println("CHecking dates on", path)
	fmt.Println("Is", timenow.Sub(mtime).Hours(), "larger than", 24, "?")
	//Does not care for timezones btw
	if timenow.Sub(mtime).Hours() > 24 {
		fmt.Println("Cache is old, Run Update")
		file.Close() //A new file will be created
		status = DeleteFile(filename, dir)
		if status == globals.Deleted {
			return simpleCreateFile(filename, dir)
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
				return globals.DirFail, nil
			}
		}
		//Create the file
		fmt.Println("trying to create", path)
		file, err = os.Create(path)
		if err != nil {

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
func FileExist(filename string, dir string) *os.File {
	path := dir + "/" + filename
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
func DeleteFile(filename string, dir string) globals.FileMsg {
	path := dir + "/" + filename
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
		return globals.Error
	}
	return globals.Deleted
}

//Does not to all the cheking that CreateFile does
func simpleCreateFile(filename string, dir string) (globals.FileMsg, *os.File) {
	path := dir + "/" + filename
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return globals.Error, nil
	}
	return globals.OldRenew, file
}
