import json
day = []
m15=[]
with open('Resultday.json') as f:
    day = set(json.load(f))
with open('Result15m.json') as f:
    m15 = set(json.load(f))



print(day&m15)