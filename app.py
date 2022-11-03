import os
import flask
from flask import Flask, request, url_for
from werkzeug.utils import secure_filename
from markupsafe import escape

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
    return os.listdir(repoDir)

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
