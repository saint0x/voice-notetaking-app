import os
import logging
import base64
import requests
import time

# Set up logging
logging.basicConfig(level=logging.INFO)


def upload_file_and_get_upload_url(file_path):
  # Set API key from environment variable
  api_key = os.getenv("ASSEMBLY_AI_KEY")
  headers = {
      "Authorization": api_key,
      "Content-Type": "application/octet-stream"
  }
  # Read the file and encode it to base64
  with open(file_path, "rb") as file:
    file_content = file.read()
    base64_content = base64.b64encode(file_content).decode("utf-8")

  # Prepare the data for the request
  data = {"data": f"data:application/octet-stream;base64,{base64_content}"}

  # Make the POST request to get the upload URL
  response = requests.post("https://api.assemblyai.com/v2/upload",
                           headers=headers,
                           json=data)
  if response.status_code == 200:
    return response.json()["upload_url"]
  else:
    logging.error(f"Failed to get upload URL: {response.text}")
    return None


def create_audio_transcript(upload_url):
  # Set API key from environment variable
  api_key = os.getenv("ASSEMBLY_AI_KEY")
  headers = {"Authorization": api_key, "Content-Type": "application/json"}
  # Data for creating the audio transcript
  data = {
      "audio_url": upload_url,
      "language_code": "en_us",
      "auto_highlights": True,
      "redact_pii": False,
      "redact_pii_policies": [],  # Empty list if no PII to redact
      "summarization": True
      # Add other parameters as needed
  }

  # Make the POST request to create the audio transcript
  response = requests.post("https://api.assemblyai.com/v2/transcript",
                           headers=headers,
                           json=data)
  if response.status_code == 200:
    return response.json()["id"]
  else:
    logging.error(f"Failed to create audio transcript: {response.text}")
    return None


def get_transcription(transcript_id):
    # Set API key from environment variable
    api_key = os.getenv("ASSEMBLY_AI_KEY")
    headers = {
        "Authorization": api_key
    }

    # Loop to check the status of the transcription
    while True:
        # Make the GET request to get the transcription
        response = requests.get(f"https://api.assemblyai.com/v2/transcript/{transcript_id}", headers=headers)
        if response.status_code == 200:
            response_json = response.json()
            status = response_json.get("status")
            if status == "completed":
                transcription = response_json.get("text")
                if transcription:
                    logging.info("Transcription retrieved successfully!")
                    print("Transcription:", transcription)
                    return transcription
                else:
                    logging.error("Transcription text is missing in the response JSON.")
                    return None
            elif status == "failed":
                logging.error("Transcription failed.")
                return None
            else:
                logging.info(f"Transcription status: {status}. Waiting for completion...")
                time.sleep(10)  # Wait for 10 seconds before checking again
        else:
            logging.error(f"Failed to get transcription. Status Code: {response.status_code}. Response Text: {response.text}")
            return None





if __name__ == "__main__":
  # File path of the audio file to transcribe
  file_path = "./testaudio.mp3"

  # Step 1: Upload file and get upload URL
  upload_url = upload_file_and_get_upload_url(file_path)
  if upload_url:
    logging.info("Upload URL obtained successfully!")

    # Step 2: Create audio transcript
    transcript_id = create_audio_transcript(upload_url)
    if transcript_id:
      logging.info("Audio transcript created successfully!")
      print("Audio transcript created successfully!"
            )  # Add this line for debugging

      # Step 3: Get transcription
      transcription = get_transcription(transcript_id)
      if transcription:
        logging.info("Transcription retrieved successfully!")
        print("Transcription:",
              transcription)  # Add this line to print the transcription
