/* this is the main server script
 * for the api
 * crystal qian 2015 */

package main

import (
	"./sqlParser"
	"./urlParser"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

var (
	addr     = flag.Bool("addr", false, "find open address and print to final-port.txt")
	username = os.Args[1]
	password = os.Args[2]
	database = os.Args[3]

	//initializing the database connects and writes a column type map
	//(see sqlParser for more details)
	db = sqlParser.InitializeDatabase(username, password, database)
)

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	path := r.URL.Path[1:]
	request := urlParser.ParseURL(path)

	tableName := request.TableName
	parameters := request.Parameters

	if tableName == "" {
		resp := sqlParser.GetTableNames()
		enc := json.NewEncoder(w)
		enc.Encode(resp)
	} else if len(parameters) <= 0 {
		resp := sqlParser.GetColumnNames(tableName)
		enc := json.NewEncoder(w)
		enc.Encode(resp)
	} else {
		resp := make(map[string]interface{}, 0)
		for _, columnName := range parameters {
			resp[columnName] = sqlParser.GetColumnValues(tableName, columnName)
		}

		enc := json.NewEncoder(w)
		enc.Encode(resp)
	}
}

//handles all calls to the API
func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//url of type "/table?parameterA=valueA&parameterB=valueB/id
	path := r.URL.Path[1:]
	fmt.Println(path)
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}

	request := urlParser.ParseURL(path)
	fmt.Println(request)

	//note: tableName could also refer to a view
	tableName := request.TableName
	tableParameters := request.Parameters
	var rows []map[string]interface{}

	//GETS the request
	if tableName != "" {
		rows, _ = sqlParser.Get(tableName, tableParameters)
		/*if err != nil {*/
		/*errString = err.Error()*/
		/*}*/
	} else {
		rows = nil
	}

	enc := json.NewEncoder(w)
	enc.Encode(rows)
}

func main() {
	fmt.Println("Starting server.")
	flag.Parse()

	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/info/", infoHandler)

	if *addr {
		//runs on home
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

	http.ListenAndServe(":8080", nil)
}
