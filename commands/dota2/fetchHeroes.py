#!/usr/bin/env python3

import urllib.request, json 
with urllib.request.urlopen("https://raw.githubusercontent.com/odota/dotaconstants/master/build/heroes.json") as url:
    data = json.loads(url.read().decode())
    
    f = open("heroes.go", "a")

    txt = "package dota2\n\nvar heroes = []string{\n"

    for x in range(1, 139):
        if data.get(str(x), None):
            name = data[str(x)]["localized_name"]
        else:
            name = ""
        
        txt += f"\t\"{name}\",\n"


    txt += "}\n"

    f.write(txt)
    f.close()
