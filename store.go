package main

import (
	"net/http"
	"os"
	"bytes"
	"path"
	"path/filepath"
	"mime/multipart"
	"io"
	"io/ioutil"
	"fmt"
	"time"
	"encoding/json"
)


var exitCode int = 0 

func fileAdd(){
		
		// Identify the current working directory
		dirPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		
		// iterate through each files
		for _, element := range os.Args[2:]{
			file_ext := filepath.Ext(element)
			// verify the file extension
			if file_ext != ".txt" {
				// exit code 1 for invalid extension
				err := fmt.Errorf("Please use ONLY plain-text files with '.txt' extension")
				exitCode = 1
				fmt.Fprintf(os.Stderr, "%s\n", err)
				return

			}

			// absolute path of the files to be uploaded
			fileRoute := path.Join(dirPath, element)

//-------------------			fmt.Println(fileRoute)  //##### REMOVE ######
			// open files to be uploaded
			openFile, openErr := os.Open(fileRoute)	
			if openErr != nil{
				fmt.Fprintf(os.Stderr, "%s\n", openErr)
				// exit code 2 for invalid filename
				exitCode = 2
				return
			}
//--------------------			fmt.Println(filepath.Base(openFile.Name()))
//---------------			fmt.Println(element)
			defer openFile.Close()
// http://127.0.0.1:5000
// Read file contents	
//			result, _ := ioutil.ReadAll(openFile)
//			fmt.Println(string(result))
// upload files
			
			// Initialize new empty buffer
			body := &bytes.Buffer{}
			// creating multipart writer 
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", element)
			io.Copy(part, openFile)
			writer.Close()

			req, _ := http.NewRequest("POST", "http://127.0.0.1:5000/add", body)
			req.Header.Add("Content-Type", writer.FormDataContentType())
			client := &http.Client{Timeout: 120 * time.Second}
			response, responseErr := client.Do(req)
			if responseErr != nil {
				fmt.Fprintf(os.Stderr, "%s\n", responseErr)
				// exit code 3 for POST request failure
				exitCode = 3
				return
			}

			// read reply after updloading
        		responseBody, responseErr := ioutil.ReadAll(response.Body)
        		if responseErr != nil {
                		fmt.Fprintf(os.Stderr, "%s\n", responseErr)
                		// exit code 5 for file not upload error
                		exitCode = 5
                		return
        		}

			fmt.Printf("%s %v\n", element, string(responseBody))
			defer response.Body.Close()
	}
}

	
func listFiles(){

	// sending get request to http api to fetch file list
	response, responseErr := http.Get("http://127.0.0.1:5000/ls")
	if responseErr != nil {
		fmt.Fprintf(os.Stderr, "%s\n", responseErr)
		// exit code 4 for GET request failure
		exitCode = 4
		return
	}

	// read body from the get request
	responseBody, responseErr := ioutil.ReadAll(response.Body)
	if responseErr != nil {
                fmt.Fprintf(os.Stderr, "%s\n", responseErr)
                // exit code 5 for no response body
                exitCode = 5
                return
        }

	// convert response body to string and print
	var responseSlice []string
	_ = json.Unmarshal([]byte(responseBody), &responseSlice)
	
	for _, element := range responseSlice{
	fmt.Println(element)
	}
}


func main() {
	// initial value of exit code
	defer func(){
		os.Exit(exitCode)
	}()
// testing commandline args
//	fmt.Println(len(os.Args[2:]), os.Args[2:])

	// if the cmdline option is add
	if len(os.Args) < 3 && os.Args[1] != "ls"{
		fmt.Println("Usage: store ['ls', 'add', 'rm', 'update'] FILE")
		return
	}else if os.Args[1] == "add" {
		fileAdd()
	} else if os.Args[1] == "ls" {
		listFiles()
	}



}
