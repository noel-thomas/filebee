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
	"encoding/hex"
	"crypto/md5"
)

// initialize the server/file-store URL here
var Url string = "http://127.0.0.1:8000"

var exitCode int = 0 

// struct for holding file names and hashes to convert into json array
type hashinfo struct {
	Name string //`json: "name"`
	Hash string //`json: "hash"`
}

type replyinfo struct {
	Name string
	State string
}

func verifyFiles() int {
	for _, element := range os.Args[2:]{
		// get file extension
		fileExt := filepath.Ext(element)

		// check whether file extension is ".txt"
		if fileExt != ".txt" {
			// exit code 1 for invalid extension
			exitCode = 1
			fmt.Fprintf(os.Stderr, "Invalid filename %s \nPlease use ONLY plain-text files with '.txt' extension", element)
	
		}
	}
	return exitCode
}


func hashFiles(responseSlice *[]replyinfo) {
	// initialize string slice to hold file hashes
	//h := make([]string, len(os.Args[2:]))
	fileInfo := []hashinfo{}

	// calculate the hashes for individual files
	for _, element := range os.Args[2:]{
		file, err := os.Open(element)
		if err != nil {
			// exit code 2 for invalid filename
			exitCode = 2
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}
		defer file.Close()
	
		hash := md5.New()

		if _, err := io.Copy(hash, file); err != nil {
			// exit code 8 for unable to calculate file hash
			exitCode = 8
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}

		// add file hash into the slice h
		//h = append(h, string(hash.Sum(nil)))
		data := hashinfo{Name: element, Hash: hex.EncodeToString(hash.Sum(nil))}
		//data := hashinfo{element, hex.EncodeToString(hash.Sum(nil))}
		fileInfo = append(fileInfo, data)
		//fmt.Println(fileInfo)
	}
	
	// convert to json
	payload, _ := json.Marshal(fileInfo)

	response, responseErr := http.Post(Url + "/hash", "application/json", bytes.NewBuffer(payload))
	if responseErr != nil {
		// exit 7 unable to sent data to remove files
		exitCode = 7
		fmt.Fprintf(os.Stderr, "%s\n", responseErr)
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
	//var responseSlice []replyinfo
	//_ = json.Unmarshal([]byte(responseBody), &responseSlice)
	_ = json.Unmarshal([]byte(responseBody), responseSlice)

	//fmt.Println(responseSlice)
	
	// processed data from reply
	//for _, element := range *responseSlice {
	//fmt.Println(element.Name, element.State)
	//}	

	defer response.Body.Close()
	return

}


func addFiles(){
		
		// Identify the current working directory
		dirPath, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		
		// verify the file extension
		if verifyFiles() == 1 {
			return
		}

		var rinfo []replyinfo
		// iterate through each files
		//for _, element := range os.Args[2:]{
		hashFiles(&rinfo)
		
		for _, element := range rinfo{

			// verify the file extension
			//if verifyFiles() == 1 {
			//	return
			//}
			
			if element.State == "absent" {


				// absolute path of the files to be uploaded
				fileRoute := path.Join(dirPath, element.Name)
        
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
				part, _ := writer.CreateFormFile("file", element.Name)
				io.Copy(part, openFile)
				writer.Close()
        
				req, _ := http.NewRequest("POST", Url + "/add", body)
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
        
				// use the below line for a single update after upload
                		//_, responseErr = ioutil.ReadAll(response.Body)
                		if responseErr != nil {
                        		fmt.Fprintf(os.Stderr, "%s\n", responseErr)
                        		// exit code 5 for file not upload error
                        		exitCode = 5
                        		return
                		}
				// print for each file uploaded - comment of not required
				fmt.Printf("%s %v\n", element.Name, string(responseBody))
				defer response.Body.Close()
			}else if element.State == "replicate"{
				fmt.Printf("%s %v\n", element.Name, "replicated")
			}else {
				fmt.Fprintf(os.Stderr, "%s already exists\n", element.Name)
			}
	}
	// to get a single update for multiple file upload uncomment the below line
	// fmt.Println("Uploaded!")
}

	
func listFiles(){

	// sending get request to http api to fetch file list
	response, responseErr := http.Get(Url + "/ls")
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
	
	defer response.Body.Close()
	}
}

func wordCount(){

	// sending get request to http api to fetch word count
	response, responseErr := http.Get(Url + "/wc")
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
	//var responseOut int
	//_ = json.Unmarshal([]byte(responseBody), responseOut)
	
	fmt.Println("Total words in file-store:", string(responseBody))
	
	defer response.Body.Close()
	
}


func freqWords() {
    if os.Args[2] =="-n" || os.Args[2] == "--limit" {
	

	// convert to json
	payload, _ := json.Marshal(os.Args)

	response, responseErr := http.Post(Url + "/freq", "application/json", bytes.NewBuffer(payload))
	if responseErr != nil {
		// exit 7 unable to sent data to remove files
		exitCode = 7
		fmt.Fprintf(os.Stderr, "%s\n", responseErr)
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
	//_ = json.Unmarshal([]byte(responseBody), &responseSlice)
	_ = json.Unmarshal([]byte(responseBody), &responseSlice)

	//fmt.Println(responseSlice)
	
	// processed data from reply
	for _, element := range responseSlice {
		fmt.Println(element)
	}	

	defer response.Body.Close()
	return

	}
}

func removeFiles() {
	// verify the files 
	if verifyFiles() == 1 {
		return
	}
	payload, _ := json.Marshal(os.Args[2:])
	response, responseErr := http.Post(Url + "/rm", "application/json", bytes.NewBuffer(payload))
	if responseErr != nil {
		// exit 7 unable to sent data to remove files
		exitCode = 7
		fmt.Fprintf(os.Stderr, "%s\n", responseErr)
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
	
	defer response.Body.Close()
	return
}


func main() {
	// initial value of exit code
	defer func(){
		os.Exit(exitCode)
	}()
// testing commandline args
//	fmt.Println(len(os.Args[2:]), os.Args[2:])
	//hashFiles()

	// if no cmdline args
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: store ['ls', 'add', 'rm', 'update', 'wc'] FILE\n")
		exitCode = 6
		return
	}
	
	// each cmdline options
	if os.Args[1] == "add" {
		addFiles()
	}else if os.Args[1] == "freq-words"{
		freqWords()
	}else if os.Args[1] == "ls" {
		listFiles()
	}else if os.Args[1] == "rm" {
		removeFiles()
	}else if os.Args[1] == "update" {
		addFiles()
	}else if os.Args[1] == "wc" {
		wordCount()
	}else {
		fmt.Fprintf(os.Stderr, "Invalid option!\nUsage: store ['ls', 'add', 'rm', 'update', 'wc'] FILE\n")
		exitCode = 6
		return
	}



}
