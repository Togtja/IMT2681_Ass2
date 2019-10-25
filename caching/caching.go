package caching

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"../globals"
)

//ShouldFileCache Check if a file should be cached
//(it will create a file if it does not exist from before)
func ShouldFileCache(filename string, dir string) (globals.FileMsg, *os.File) {
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
		fmt.Println("trying to create file")
		file, err = os.Create(path)
		if err != nil {

			return globals.Error, nil
		}
		fmt.Println("file created")
		return globals.Created, file
	}

	//The file exist and we need to see how old it is
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		return globals.Error, nil
	}
	mtime := info.ModTime()
	fmt.Println("Last changed:", mtime)
	timenow := time.Now()
	fmt.Println("Time now:", timenow)
	fmt.Println("Is", timenow.Sub(mtime).Hours(), "larger than", 24, "?")
	//Does not care for timezones btw
	if timenow.Sub(mtime).Hours() > 24 {
		fmt.Println("Cache is old, Run Update")
		err := os.Remove(path)
		if err != nil {
			fmt.Println(err)
			return globals.Error, nil
		}
		file, err = os.Create(path)
		if err != nil {
			fmt.Println(err)
			return globals.Error, nil
		}
		return globals.OldRenew, file
	}
	fmt.Println("Cache is recent, No need to update")
	return globals.Exist, file
}

//CacheStruct creates a file in the dir with the struct as a json
func CacheStruct(filename string, dir string, v interface{}) {
	should, file := ShouldFileCache(filename, dir)
	if should == globals.Error || should == globals.DirFail {
		fmt.Println("Failed to find or create file")
		return
	}
	if should == globals.Exist {
		//No need to cache the files
		return
	}
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
