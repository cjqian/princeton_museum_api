/* this is the main server script
 * for the api
 * crystal qian 2015 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cjqian/princeton_museum_api/outputFormatter"
	"github.com/cjqian/princeton_museum_api/sqlParser"
	"github.com/cjqian/princeton_museum_api/urlParser"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

var (
	addr = flag.Bool("addr", false, "find open address and print to final-port.txt")
	/*
		username = os.Args[1]
		password = os.Args[2]
		database = os.Args[3]
	*/
	username = "root"
	password = "helloworld"
	database = "puamapi"

	//initializing the database connects and writes a column type map
	//(see sqlParser for more details)
	db = sqlParser.InitializeDatabase(username, password, database)
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println(r.URL.Path)
	path := r.URL.Path[1:]
	request := urlParser.ParseURL(path)
	fmt.Println(request)
	tableName := request.TableName
	tableParameters := request.Parameters
	specialParameters := request.SpecialParameters

	var resp interface{}
	fmt.Println(tableName)
	fmt.Println(tableParameters)
	fmt.Println(specialParameters)

	if tableName == "" {
		records := sqlParser.GetTableNames()
		numResults := len(records)
		resp = outputFormatter.MakeApiWrapper(request, records, numResults, specialParameters)

	} else if len(tableParameters) <= 0 {
		records := sqlParser.GetColumnNames(tableName)
		numResults := len(records)
		resp = outputFormatter.MakeApiWrapper(request, records, numResults, specialParameters)

	} else {
		records := make(map[string][]string, 0)
		rowWrappers := make(map[string]interface{}, 0)
		for _, columnName := range tableParameters {
			records[columnName] = sqlParser.GetColumnValues(tableName, columnName)

			curNumResults := len(records[columnName])
			rowWrappers[columnName] = outputFormatter.MakeInfoWrapper(curNumResults, records[columnName])
		}
		numResults := len(records)
		resp = outputFormatter.MakeApiWrapper(request, rowWrappers, numResults, specialParameters)

	}

	fmt.Println(resp)

	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

//handles all calls to the API
func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//url of type "/table?parameterA=valueA&parameterB=valueB/id
	path := r.URL.Path[1:]
	//fmt.Println(path)
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}

	request := urlParser.ParseURL(path)
	//fmt.Println(request)

	//note: tableName could also refer to a view
	tableName := request.TableName
	tableParameters := request.Parameters
	specialParameters := request.SpecialParameters

	var records []map[string]interface{}

	//GETS the request
	if tableName != "" {
		fmt.Println("Getting request")
		records, _ = sqlParser.Get(tableName, tableParameters, specialParameters)
		/*if err != nil {*/
		/*errString = err.Error()*/
		/*}*/
	} else {
		records = nil
	}

	fmt.Println("Getting numbers, lim String")
	numResults := sqlParser.GetNumRows(tableName, sqlParser.GetWhereString(tableParameters), sqlParser.GetLimString(specialParameters))

	fmt.Println("Making wrapper")
	resp := outputFormatter.MakeApiWrapper(request, records, numResults, specialParameters)
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

func main() {
	//fmt.Println("Starting server.")
	//fmt.Println(os.Getenv("PORT"))
	flag.Parse()

	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/info/", infoHandler)

	if *addr {
		//runs on home
		//l, err := net.Listen("tcp", os.Getenv("PORT"))
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
		if err != nil {
			panic(err)
		}
		s := &http.Server{}
		s.Serve(l)
		return
	}
	fmt.Println(":" + os.Getenv("PORT"))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
