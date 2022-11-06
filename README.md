# filebee

Filebee is a file-store management app. The client side is created in Golang and server on Python(Flask). Using instructions are below:

1. container created from the app.py is available at `quay.io/noeltredhat/filebee`
   Pull the image using the command `podman pull quay.io/noeltredhat/filebee` (you can use your preferred runtime)
   The application is listening to port 8000, so map this port to any host port when running.

2. To run the container `podman run -d --name <optional: container-name> -p 8000:<host-port> quay.io/noeltredhat/filebee`
   It will pull and execute the container in you local machine.
   
3. open the store.go and replace the variable `Url = localhost:<host-port>`

4. Run `go build store.go` to build the binary

5. Then execute `./store' to see help

### Usage:

a. To add files to file store. Ensure the file extension should be `.txt` for all the files
    `$ store add file1.txt file2.txt file3.txt`
    You can pass multiple files are argument. If the file content is available on the server it duplicates the file instead of sending again.
    
b. To list files from file store
    `$ store ls`
    
c. To remote file from file store
    `$ store rm file.txt file2.txt file3.txt`
    You can remove multiple files from the server
   
d. Update file, It implements same functionality like add, later verions might add PUT method to send data
    `$ store update file.txt file2.txt`
    
e. To count the total words from files in file store
    `$ store wc`
