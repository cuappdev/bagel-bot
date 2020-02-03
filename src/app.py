import json
from flask import Flask, request

app = Flask(__name__)

@app.route('/event/', methods=['POST'])
def receive_slack_event():
    body = json.loads(request.data)
    if body['type'] == 'url_verificaton':
        return body['challenge']

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
