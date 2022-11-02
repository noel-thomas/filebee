import os
from flask import Flask, request, url_for
from werkzeug.utils import secure_filename
from markupsafe import escape

UPLOAD_FOLDER = '/tmp'
ALLOWED_EXTENSIONS = {'txt'}

app = Flask(__name__)
app.config['UPLOAD_FOLDER'] = UPLOAD_FOLDER

# reply for root
@app.route('/')
def index():
    return 'Hello!'

# response to api /add - to upload files
@app.post('/add')
def add_files():
    f = request.files['the_file'] # the_file identification of uploaded file
    f.save(f"/tmp/{secure_filename(f.filename)}")
    return {"msg": "added"}

# response for the api ls - to list the files
@app.route('/ls')
def list_files():
    return f"{request(files)}"

# sample test for custom route
@app.route('/<string:datas>')
def list_print(datas):
    return f"Hi, {escape(datas)}"
