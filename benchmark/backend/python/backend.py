from flask import Flask, request, jsonify
from hashlib import sha256
from waitress import serve
from time import time
from jwt import PyJWT
import json
import math
import re

PORT = 1234
EMAIL_RE = re.compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$")
SECRET_KEY = "temporary_key_i_will_change_for_the_real_project_inshallah"
EXPIRATION = 60 * 30 # 30 minutes

app = Flask(__name__)

def good_password(password: str):
    return re.findall("[0-9]", password) \
        and re.findall("[a-z]", password) \
        and re.findall("[A-Z]", password) \
        and re.findall("[@$\\/?<>#:+!-]", password)

def calculate(t0, t1, it):
    start = time()
    angle: None | float = None
    x = 0
    print("bar")
    for i in range(0, it):
        x += i
        if (t1[2] == t0[2]) \
            or (t0[2] > 0 and t1[2] < 0) \
            or (t0[2] < 0 and t1[2] > 0):
            angle = None
            continue
        vel = [
            t1[0] - t0[0],
            t1[1] - t0[1],
            t1[2] - t0[2],
        ]
        k = math.sqrt(vel[0] ** 2 + vel[1] ** 2 + vel[2] ** 2)
        angle = math.asin(vel[2] / k)
        angle = math.floor(angle * -180 / math.pi * 100) / 100
    print(x)
    end = time()
    return angle, (end - start) * 1000

@app.route("/hello/<username>")
def hello(username):
    try:
        x = json.loads(request.args.get("t0"))
        y = json.loads(request.args.get("t1"))
        angle, ms = calculate(
            x,
            y,
            1_000_000
        )
    except (json.JSONDecodeError, TypeError) as e:
        return str(e), 400
    if angle is None:
        return "Hello {}! your ball wont reach, computed in {:.2f}ms".format(username, ms)
    else:
        return "Hello {}! your incidence angle is {:.2f}, computed in  {:.2f}ms".format(username, angle, ms)

def new_token(email):
    payload = {
        'email': email,
        'exp': time.time() + EXPIRATION
    }
    return PyJWT().encode(payload=payload, key=SECRET_KEY, algorithm='HS256')

def untoken_your_token(tok):
    return PyJWT().decode(jwt=tok, key=SECRET_KEY, algorithms=['HS256'])

def get_email_password():
    data = request.get_json() 
    if not data:
        return False, jsonify({"message": "no data given"}), 400
    if not "password" in data or not "email" in data:
        return False, jsonify({"message": "no email or password"}), 400
    if not EMAIL_RE.match(data["email"]):
        return False, jsonify({"message": "invalid email"}), 400
    if not good_password(data["password"]):
        return False, jsonify({"message": "invalid password"}), 400
    if not type(data["password"]) == str:
        return False, jsonify({"message": "data is not a string ?"}), 400
    password: str = data["password"]
    hashedPassword = sha256(password.encode("utf-8")).hexdigest()
    email = data["email"]
    return True, email, hashedPassword

@app.route('/api/register', methods=['POST'])
def handle_register():
    success, x, y = get_email_password()
    if not success:
        return x, y
    print(f"Email: {x}")
    print(f"Password: {y}")
    token = new_token(x)
    # TODO: insert email and hashedpassword in db
    return jsonify({"token": token}), 200

@app.route('/api/login', methods=['POST'])
def handle_login():
    success, x, y = get_email_password()
    if not success:
        return x, y
    print(f"Email: {x}")
    print(f"Password: {y}")
    # TODO: check if it exists in db and get the token
    token = new_token(x)
    return jsonify({"token": token}), 200

if __name__ == "__main__":
    print(f"=> server listens on port {PORT}")
    serve(app, port=PORT)
