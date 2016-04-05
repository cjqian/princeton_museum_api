import pprint
import sys
import re
reload(sys)
sys.setdefaultencoding('utf-8')

import string
import json
import requests
import argparse

url = "http://localhost:8080/"

#flags
C_FLAG = 0; #query constituents
O_FLAG = 1; #query objects
CO_FLAG = 2; #query constituent objects

#set up the parser
parser = argparse.ArgumentParser()
parser.add_argument("id", type=int, help="id of what you are seeking")
parser.add_argument('-o', '--objects', action='store_true', help='query an object')
parser.add_argument('-c', '--constituents', action='store_true', help='query a constituent')

#gets a request from the api and returns a JSON object with the necessary strings
def getObject(objectID):
    base = url + "api/apiobjects/" + str(objectID)

    request = requests.get(base)
    if (request is None) or (request == '') or (request.json() is None):
        print("get failed")
        return -1

    return (request.json()["records"][0])

def getConst(constID):
    base = url + "api/apiconstituents/" + str(constID)

    request = requests.get(base)
    if (request is None) or (request == '') or (request.json() is None):
        print("get failed")
        return -1

    return (request.json()["records"][0])
     
#returns an array of objects
def getConstObjects(constID):
    jsonArray = []

    base = url + "api/apiconobjxrefs?constituentid=" + str(constID)

    request = requests.get(base)
    objectArray = request.json()["records"]
    for object in objectArray:
        curID = object["ObjectID"]
        curObject = getObject(curID)
        jsonArray.append(curObject)

    return jsonArray

#python make_results.py maxNumGames
def writeData(id, flag): 
    if flag == CO_FLAG:
        data = getConstObjects(id)
        fileName = "const_object_" + str(id) + ".json"
    elif flag == C_FLAG:
        data = getConst(id)
        fileName = "const_" + str(id) + ".json"
    elif flag == O_FLAG:
        data = getObject(id)
        fileName = "object_" + str(id) + ".json"

    #write to file "results_123542.json"
    f = open("output/" + fileName, 'w')
    f.write(json.JSONEncoder().encode(data))
    f.close()

def main():
        args = parser.parse_args()
        print args
        #make sure id exists
        if not args.id:
                print ("ERROR: NO ID SPECIFIED")
                return
        #cases
        if not args.objects and not args.constituents:
                print ("ERROR: NO QUERY TYPE SPECIFIED")
                return
        elif args.objects and args.constituents:
                writeData(args.id, CO_FLAG)
        elif args.objects:
                writeData(args.id, O_FLAG)
        elif args.constituents:
                writeData(args.id, C_FLAG)
main()
