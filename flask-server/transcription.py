from flask import Flask, request, jsonify
import os
import assemblyai as aai

app = Flask(__name__)

# Endpoint for file upload and transcription
@app.route("/transcribe", methods=["POST"])
def transcribe_audio():
    # Receive file from the request
    file = request.files['file']
    file.save("./temp_audio.m4a")  # Save the file temporarily

    # Upload file and get URL
    file_url = UploadFileAndGetURL("./temp_audio.m4a")

    if isinstance(file_url, str):
        # Get transcription using the URL
        transcription = GetTranscriptionFromURL(file_url)
        os.remove("./temp_audio.m4a")  # Remove temporary file
        return jsonify({"transcription": transcription})
    else:
        os.remove("./temp_audio.m4a")  # Remove temporary file
        return jsonify({"error": file_url}), 500

def UploadFileAndGetURL(file_path):
    # Replace with your API key
    api_key = os.getenv("ASSEMBLY_AI_KEY")
    aai.settings.api_key = api_key

    # Upload the file and get the URL
    uploader = aai.Uploader()
    upload_response = uploader.upload(file_path)

    if upload_response.status == aai.UploadStatus.success:
        return upload_response.upload_url
    else:
        return upload_response.error

def GetTranscriptionFromURL(file_url):
    # Replace with your API key
    api_key = os.getenv("ASSEMBLY_AI_KEY")
    aai.settings.api_key = api_key

    # Get transcription using the URL
    transcriber = aai.Transcriber()
    transcript = transcriber.transcribe(file_url)

    if transcript.status == aai.TranscriptStatus.error:
        return transcript.error
    else:
        return transcript.text

if __name__ == "__main__":
    app.run(debug=True)
