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
        return ["Empty!"]
    else:
        return os.listdir(repoDir)


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
    if os.listdir(repoDir) == [] :
        for k in range(len(remoteHash)):
            value = {'Name': remoteHash[k]['Name'], 'State': 'absent'}
            returnContent.append(value)
        return returnContent

    for i in os.listdir(repoDir):
        path = repoDir + i

        if os.path.getsize(path) != 0:
            with open(path) as file, mmap(file.fileno(), 0, access=ACCESS_READ) as file:
                data = md5(file).hexdigest()
                # COMPARE THE HASH
                #j = 0
                #for k in remoteHash[j].keys():
                for j in range(len(remoteHash)):
                    if remoteHash[j]['Hash'] == data:
                        if remoteHash[j]['Name'] == i:
                            # returnContent.file_exist = yes
                            state = 'present'
                            value = {'Name': remoteHash[j]['Name'], 'State': 'present'}
                            # returnContent.append(value)
                        else:
                            # replcate i with the name j['name']
                            state = 'replicate'
                            value = {'Name': remoteHash[j]['Name'], 'State': 'replicate'}
                            # returnContent.append(value)
                    else:
                        # returnContent.file_name = j['name']
                        # returnContent.file_exist = no
                        state = 'absent'
                        value = {'Name': remoteHash[j]['Name'], 'State': 'absent'}
                    #returnContent.append(value)
                    #LOOK HERE
                    if j > (len(returnContent) - 1):
                        returnContent.append(value)
                    else:
                        if returnContent[j]['State'] == state:
                            break
                        elif (returnContent[j]['State'] == 'absent' and state != 'absent') or (returnContent[j]['State'] == 'replicate' and state == 'present') :
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
                            returnContent[j]['State'] = state
                    #j = j+1
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
