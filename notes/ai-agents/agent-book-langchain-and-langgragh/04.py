from dotenv import load_dotenv

# from langchain_openai import ChatOpenAI
# from langchain_core.messages import SystemMessage, HumanMessage, AIMessage
from langsmith import Client

load_dotenv()
# model = ChatOpenAI(model="gpt-4o-mini", temperature=0)

#     SystemMessage(content="Your are a helpful assistant"),
#     HumanMessage(content="こんにちは。私はJhonといいます"),
#     #     AIMessage(content="こんにちは、ジョンさん！どのようにお手伝いできますか？"),
#     #     HumanMessage(content="私の名前がわかりますか？"),
# ]
#
# # ai_message = model.invoke(messages)
# # print(ai_message.content)
#
# for chunk in model.stream(messages):
#     print(chunk.content, end="", flush=True)

client = Client()
prompt = client.pull_prompt("oshima/recipe")

prompt_value = prompt.invoke({"dish": "カレー"})
print(prompt_value)
