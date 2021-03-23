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
collector = SDK("agent:7000")

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

    face = {}
    if len(face_landmarks_list) != 0:
        face = face_landmarks_list[0]
        face["off_center"] = abs(face["nose_tip"][0][0] - 640) / 1280 + abs(face["nose_tip"][0][1] - 360) / 720
    collector.Image(
        now.strftime("%H:%M:%S.%f.jpeg"),
        imgByteArr.getvalue(),
        face
    )
    return jsonify(face)
