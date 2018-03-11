#!/usr/bin/env python

import json
import requests

baseURL = "https://dockstore.org:8443/api/ga4gh/v1/tools/"

res = requests.get(baseURL)

nout = open("dockstore.nodes", "w")
eout = open("dockstore.edges", "w")

def write_node(n):
    nout.write(json.dumps(n) + "\n")

def write_edge(n):
    eout.write(json.dumps(n) + "\n")

id_map = {}

for tool in res.json():
    print json.dumps(tool)
    for v in tool['versions']:
        vm = {
            "gid": v["id"],
            "label" : "ToolVersion",
            "data" : v
        }
        write_node(vm)
        tve = {
            "from": tool['id'],
            "to" : v["id"],
            "label" : "version"
        }
        write_edge(tve)

        if v["image"] is not None:
            if v["image"] not in id_map:
                di = {
                    "gid" : v["image"],
                    "label" : "DockerImage"
                }
                write_node(di)
                id_map[v["image"]] = True
            vie = {
                "from" : v["id"],
                "to" : v["image"],
                "label": "image"
            }
            write_edge(vie)


    td = {
        "gid" : tool['id'],
        "label" : "Tool",
        "data" : tool
    }
    del td["data"]["versions"]
    write_node(td)

nout.close()
eout.close()
