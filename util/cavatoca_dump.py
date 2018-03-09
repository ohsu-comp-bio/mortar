#!/usr/bin/env python

import os
import json
import sevenbridges as sbg
from sevenbridges.models.file import File

api = sbg.Api(url='https://cavatica-api.sbgenomics.com/v2', token=os.environ['SB_AUTH_TOKEN'])

project_id = "yuankun/kf-gvcf-benchmarking"

nodes = open("cavatica.nodes", "w")
edges = open("cavatica.edges", "w")

n = {
    "gid" : project_id,
    "label" : "Project"
}
nodes.write(json.dumps(n) + "\n")

#List tasks
task_list = api.tasks.query(project=project_id).all()
for i in task_list:
    n = {
        "gid" : i.id,
        "label" : "Task",
        "data" : {
            "name" : i.name,
            "status" : i.status,
            "app" : i.app
        }
    }
    nodes.write(json.dumps(n) + "\n")
    for name, f in i.inputs.items():
        if isinstance(f, list):
            for s in f:
                e = {
                    "to" : s.id,
                    "from" : i.id,
                    "label" : "input",
                    "data" : {
                        "name" : name
                    }
                }
                edges.write(json.dumps(e) + "\n")

        else:
            e = {
                "to" : f.id,
                "from" : i.id,
                "label" : "input",
                "data" : {
                    "name" : name
                }
            }
            edges.write(json.dumps(e) + "\n")
    if i.outputs is not None:
        for name, f in i.outputs.items():
            if isinstance(f, list):
                """
                for s in f:
                    e = {
                        "to" : s.id,
                        "from" : i.id,
                        "label" : "input",
                        "data" : {
                            "name" : name
                        }
                    }
                    edges.write(json.dumps(e) + "\n")
                """
            elif f is not None:
                if isinstance(f, File):
                    e = {
                        "to" : f.id,
                        "from" : i.id,
                        "label" : "input",
                        "data" : {
                            "name" : name
                        }
                    }
                    edges.write(json.dumps(e) + "\n")
        
    #print i.id, i.name, i.status, i.app, i.type, i.inputs, i.outputs

#list files
file_list = api.files.query(project=project_id).all()
# http://sevenbridges-python.readthedocs.io/en/latest/quickstart/#file-properties
for i in file_list:
    n = {
        "gid" : i.id,
        "label" : "File",
        "data" : {
            "size" : i.size,
            "metadata" : i.metadata,
            "tags" : i.tags
        }
    }
    nodes.write(json.dumps(n) + "\n")
    e = {
        "from" : i.id,
        "to" : i.project,
        "label" : "project"
    }
    edges.write(json.dumps(e) + "\n")
    
    if i.origin.task is not None:
        e = {
            "from" : i.id,
            "to" : i.origin.task,
            "label" : "origin"
        }
        edges.write(json.dumps(e) + "\n")


nodes.close()
edges.close()
