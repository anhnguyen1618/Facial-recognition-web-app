from flask import Flask, render_template, request
import json, time, string, random, os, requests
from imgurpython import ImgurClient

app = Flask(__name__)

register_pic = None

@app.route('/')
def home():
    return render_template('index.html')


def get_file_path(name):
	'''
	input: filename
	random new name for files
	output : current working derectory + new filename
	'''
    name_parts = name.split('.')
    spec_characters = ""
    for i in range(4):
        spec_characters += random.choice(string.ascii_letters)
    extension = name_parts[1]
    fileName = str(int(time.time())) + spec_characters + '.' + extension
    path = 'C:/Users/zozos/Desktop/ws/untitled/static/'
    return path + fileName

def upload(file_path):
	'''
	input: directory of file
	Upload to imgur
	output: link to that image in imgur 
	'''
    client_id = '32e70e60c837b9d'
    client_secret = '58a3cc5096cda09e38afa837ebcb4038f943e543'
    client = ImgurClient(client_id, client_secret)
    config = {'album': None, 'name': 'Catastrophe!', 'title': 'Catastrophe!',
              'description': 'Cute kitten being cute on'}
    image = client.upload_from_path(file_path, config=config, anon=False)
    link = image['link']
    return link


@app.route('/registerperson', methods=['POST'])
def register_person():
    global register_pic
    REGISTER_NAME_URL = "https://api.projectoxford.ai/face/v1.0/persongroups/authorized_people_id/persons"
    TRAIN_BOT_URL= "https://api.projectoxford.ai/face/v1.0/persongroups/authorized_people_id/train"
    name = request.form['name']
    filePath = register_pic
    if not name:
        return ""

    if not filePath:
        try:
            file = request.files['idPic']
            filePath = get_file_path(file.filename)
            file.save(filePath)
        except:
            return ""

    link = upload(filePath)
    register_pic = None

    body = {"name" : name}
    person_data = call_api(REGISTER_NAME_URL, body)
    person_id = person_data["personId"]
    ADD_IMAGE_URL = "https://api.projectoxford.ai/face/v1.0/persongroups/authorized_people_id/persons/"+person_id+"/persistedFaces"
    call_api(ADD_IMAGE_URL, {"url": link})
    call_api(TRAIN_BOT_URL,{})
    return "ok"

def call_api(url, body, head = None):
    headers = {'Content-Type': 'application/json', 'Ocp-Apim-Subscription-Key': '57f524c61e8a4f95bdcbd3ffa7e18cdd'}
    if head:
        headers = {'Content-Type': 'application/json', 'Ocp-Apim-Subscription-Key': 'ffe87b798fd64705945dcec50ec647a5'}
    data_callback = requests.post(url, data=json.dumps(body), headers=headers)
    try:
        return data_callback.json()
    except:
        return {}

@app.route('/authenticate', methods=['POST'])
def authenticate_person():
    GO_URL_OF_PARTNER = "http://551c2397.ngrok.io/identifyperson"
    file = request.files['webcam']
    filePath = get_file_path(file.filename)
    file.save(filePath)
    link = upload(filePath)
    call_url =  "https://api.projectoxford.ai/face/v1.0/detect?returnFaceId=true&returnFaceLandmarks=false"
    data_callback = call_api(call_url, {"url":link})
    faceId=""

    try:
        faceId = data_callback[0]['faceId']
    except:
        faceId = ""

    status = requests.post(kha_url, data=json.dumps({"faceId": faceId}))
    try:
        status = status.json()["name"]
    except:
        status = None
    if status:
        try:
            emotion = call_api('https://api.projectoxford.ai/emotion/v1.0/recognize', {"url": link}, True)[0]['scores']
            emo=""
            max = 0
            for key in emotion:
                if emotion[key] > max:
                    emo = key
                    max = emotion[key]
            return json.dumps({'emotion':emo, 'index': max, 'name': status})
        except:
            return json.dumps({'emotion': "undetected", 'index': 0, 'name': status})
    return json.dumps({'emotion':None, 'index': None, 'name': ""})

@app.route('/uploadpic', methods=['POST'])
def upload_pic():
    global register_pic
    file= request.files['webcam']
    filePath = get_file_path(file.filename)
    file.save(filePath)
    register_pic= filePath
    return "ok"

if __name__ == '__main__':
    app.run()
