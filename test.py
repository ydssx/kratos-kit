from pydantic_ai.models.ollama import OllamaModel
from pydantic_ai import Agent
import asyncio
import requests
model = OllamaModel(model_name="llama3.2:3b")

agent = Agent(model=model, system_prompt="用中文回答所有问题")


async def main():
    async with agent.run_stream("帮我做一份圣诞节礼物清单") as response:
        async for chunk in response.stream(debounce_by=0.01):
            print(chunk, end="", flush=False)


if __name__ == "__main__":
    # asyncio.run(main())
    # from openai import OpenAI

    # chat = OpenAI(api_key="715028322e509853f0b60a74a3ce8a26", base_url="http://localhost:8000").chat.completions.create(
    #     model="jimeng-2.1",
    #     messages=[{"role": "user", "content": "少女"}],
    # )
    # print(chat.choices[0].message.content)
    body = {
        "model": "jimeng-2.1",
        "messages": [{"role": "user", "content": "少女"}],
        "stream": False,
    }
    headers = {
        "Authorization": "Bearer 715028322e509853f0b60a74a3ce8a26",
        "Content-Type": "application/json",
    }
    response = requests.post("http://localhost:8000/v1/chat/completions", json=body, headers=headers)
    print(response.json())
