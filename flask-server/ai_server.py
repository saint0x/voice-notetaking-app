from flask import Flask, request, jsonify
import os
from openai import OpenAI

app = Flask(__name__)

# Initialize OpenAI client
client = OpenAI(api_key=os.environ.get("OPENAI_API_KEY"))

# Define system prompt
SYSTEM_PROMPT = "System prompt: Please provide your input. Concepts: {}"


@app.route('/ai-call', methods=['POST'])
def ai_call():
  # Extract concepts and prompt from the request
  data = request.json
  concepts = data.get('concepts', [])
  prompt = data.get('prompt', '')

  # Create chat completion request with system prompt
  chat_completion = client.chat.completions.create(
      messages=[{
          "role": "system",
          "content": SYSTEM_PROMPT.format(', '.join(concepts)),
      }, {
          "role": "user",
          "content": prompt,
      }],
      model="gpt-3.5-turbo",
  )

  # Process response and update the graph
  messages = chat_completion.get('choices', [])
  if not messages:
    return jsonify({"error": "No messages received in response"}), 500

  message_content = None
  for message in messages:
    if message.get('role') == 'assistant':
      message_content = message.get('message', {}).get('content')
      # Parse the message content and update the graph accordingly
      break  # Exit the loop once we find the assistant message

  if message_content is None:
    return jsonify({"error": "No assistant message found in response"}), 500

  return jsonify({"success": True}), 200


if __name__ == '__main__':
  app.run(debug=True)
