import os
from openai import OpenAI


class LanguageModel():
  OPENAI_MODEL = 'gpt-4'
  client = OpenAI(api_key=os.environ['OPENAI_API_KEY'])

  def complete(self, context: str) -> str | None:
    messages = [
        dict(role='system', content=context),
    ]
    chat_completion = self.client.chat.completions.create(
        model=self.OPENAI_MODEL, messages=messages)
    response = chat_completion.choices[0].message.content
    return response 
