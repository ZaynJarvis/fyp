import json
import io
from PIL import Image
import face_recognition
from datetime import datetime
from flask import Flask, jsonify, request, redirect
from flask_cors import CORS
from sdk import SDK

app = Flask(__name__)
CORS(app)
collector = SDK("host.docker.internal:7000")

@app.route('/', methods=['POST'])
def root():
    if 'image' not in request.files:
        return "405"

    file = request.files['image']

    if file.filename == '':
        return "405"

    imgArr = face_recognition.load_image_file(file)
    face_landmarks_list = face_recognition.face_landmarks(imgArr, model="small")

    now = datetime.now()
    image = Image.fromarray(imgArr)
    imgByteArr = io.BytesIO()
    image.save(imgByteArr, format='JPEG')
    
    collector.Image(
        now.strftime("%H:%M:%S.%f.jpeg"),
        imgByteArr.getvalue(),
        {} if len(face_landmarks_list) == 0 else face_landmarks_list[0]
    )
    return jsonify(face_landmarks_list)
