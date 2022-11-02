package main

import (
//	"net/http"
	"os"
//	"bytes"
//	"path"
//	"path/filepath"
//	"mime/multipart"
//	"io"
	"fmt"
)

func main() {

// testing commandline args
	fmt.Println(len(os.Args[1:]), os.Args[1:])
// Identify the current working directory
	out, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

}
