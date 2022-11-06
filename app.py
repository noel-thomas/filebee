import os
import flask
import shutil
from flask import Flask, request, url_for
from werkzeug.utils import secure_filename
from markupsafe import escape
from hashlib import md5
from mmap import ACCESS_READ, mmap

#UPLOAD_FOLDER = '/tmp'
#ALLOWED_EXTENSIONS = {'txt'}

app = Flask(__name__)
#app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER

# setting max file size to 16 MB
#app.config['MAX_CONTENT_LENGTH'] = 16 * 1000 * 1000
#app.config['MAX_CONTENT_LENGTH'] = 2 * 1000
repoDir = "/tmp/filebee/"
# reply for root
@app.route('/')
def index():
    return 'Hello!'

# response to api /add - to upload files
@app.post('/add')
def add_files():
    f = request.files['file'] # the_file identification of uploaded file
    f.save(repoDir + f"{secure_filename(f.filename)}")
    return "uploaded!"

# response for the api ls - to list the files
@app.route('/ls')
def list_files(): 
    if not os.listdir(repoDir):
        return ["Empty file-store!"]
    else:
        return os.listdir(repoDir)

@app.route('/wc')
def word_count():
    if os.listdir(repoDir) == [] :
        return "Empty file-store!"
    else:
        count = 0
        for i in os.listdir(repoDir):
        # path to the file in file-store
            path = repoDir + i
            file = open(path, "rt")
            data = file.read()
            words = data.split()
            count = count + len(words)

    return f"{count}"


# file hash verification
@app.post('/hash')
def hash_files():
    remoteHash = request.get_json()
    # REMOVE
    
    #return remoteHash
    #return '[{"Name":"data.txt","Hash":"764efa883dda1e11db47671c4a3bbd9e"},{"Name":"opera.txt","Hash":"7c501a0514172db8a0cad8b627de5f98"}]'

    returnContent = []
    ########### REMOVE
    
    #for j in range(len(remoteHash)):
    #    returnContent.append(remoteHash[j])
    #    #j = j+1
    #return returnContent

    # if file-store is empty reply all files are absent
    if os.listdir(repoDir) == [] :
        for k in range(len(remoteHash)):
            value = {'Name': remoteHash[k]['Name'], 'State': 'absent'}
            returnContent.append(value)
        return returnContent

    # iterate through the files in the file-store
    for i in os.listdir(repoDir):
        # path to the file in file-store
        path = repoDir + i

        # check if file-store is empty
        if os.path.getsize(path) != 0:
            # open file
            with open(path) as file, mmap(file.fileno(), 0, access=ACCESS_READ) as file:
            
                #calculate the checksum
                data = md5(file).hexdigest()
                
                # COMPARE THE HASH
                for j in range(len(remoteHash)):
                    # check if remote file hash and local file hash are equal
                    if remoteHash[j]['Hash'] == data:
                        
                        # check if both file-names are equal 
                        if remoteHash[j]['Name'] == i:
                            state = 'present'
                            value = {'Name': remoteHash[j]['Name'], 'State': 'present'}
                            
                        else:
                            # if names are different replicate the file in the file-store
                            state = 'replicate'
                            value = {'Name': remoteHash[j]['Name'], 'State': 'replicate'}
                            
                    else:
                        # if hash does'nt match then file is absent
                        state = 'absent'
                        value = {'Name': remoteHash[j]['Name'], 'State': 'absent'}
                    
                    # if checking remote content for the first time
                    # we have to add that info into the returnContent
                    if j > (len(returnContent) - 1):

                        # if it's a new record add into the return content
                        returnContent.append(value)
                        
                        # even if the record is new and the state is replicate then repicate the file in the file-store
                        if state == 'replicate':
                            fpath = repoDir + remoteHash[j]['Name']
                            
                            if os.path.exists(fpath):
                                with open(fpath) as cfile, mmap(cfile.fileno(), 0, access=ACCESS_READ) as cfile:
                                    cdata = md5(cfile).hexdigest()
                                    if cdata != data:
                                        shutil.copy(path, fpath)
                                        print(path, fpath)
                                cfile.close()
                            else:
                                shutil.copy2(path, fpath)
                            
                    else:
                        # if both current state and the recorded state are equal then do nothing
                        if returnContent[j]['State'] == state:
                            break
                        # if the recorded state in the returnContent is absent and the current state is not absent
                        # or if the recorded state is replicate but if the file is present
                        # then record the current state in the return content
                        elif (returnContent[j]['State'] == 'absent' and state != 'absent') or (returnContent[j]['State'] == 'replicate' and state == 'present') :
                            fpath = repoDir + remoteHash[j]['Name']
                            
                            # if the file exists but with different content then replicate
                            if os.path.exists(fpath):
                                with open(fpath) as cfile, mmap(cfile.fileno(), 0, access=ACCESS_READ) as cfile:
                                    cdata = md5(cfile).hexdigest()
                                    if cdata != data:
                                        shutil.copy(path, fpath)
                                        print(path, fpath)
                                cfile.close()
                            else:
                                # if file does'nt exist then replicate
                                shutil.copy2(path, fpath)
                            returnContent[j]['State'] = state
                    
            file.close()


    return returnContent

# remove requested file from store
@app.post('/rm')
def remove_files():
   filenames = request.get_json()
   returnContent = []
   for i in filenames:
        if os.path.isfile(repoDir + str(i)) == False:
            returnContent.append(f"{i} file not found!")
        else:
            try:
                os.remove(repoDir + str(i))
                returnContent.append(f"{i} deleted")
            except OSError:
                returnContent.append("Error: unable to remove the file")
   return returnContent

# sample test for custom route
@app.route('/<string:datas>')
def list_print(datas):
    return f"Hi, {escape(datas)}"


if __name__ == "__main__":
    app.run(debug=True)
