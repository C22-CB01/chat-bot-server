from flask import Flask, request
from flask.json import jsonify
from tensorflow import keras

# Machine learning model
model = keras.models.load_model("my_model.h5")

app = Flask(__name__)


# TODO: Make it possible to predict
def predict(data):
    # Removing punctuation and converting to lowercase
    # text_p = []
    # message = [letters.lower()
    #            for letters in data if letters not in string.punctuation]
    # message = ''.join(message)
    # text_p.append(message)

    prediction = model.predict([data])
    return prediction


@app.route("/", methods=["GET", "POST"])
def index():
    if request.method == "POST":
        data = request.json
        if not data:
            return jsonify({"error": "No message"})
        try:
            prediction = predict(data['message'])
            print(prediction)
            return data
        except Exception as e:
            return jsonify({"error": str(e)})

    return "OK"


if __name__ == "__main__":
    app.run(debug=True)
