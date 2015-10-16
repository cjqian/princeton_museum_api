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

func requestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resp := sqlParser.GetTableNames()
	enc := json.NewEncoder(w)
	enc.Encode(resp)
}

//handles all calls to the API
func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request:")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, X-Requested-With, Content-Type")

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
	//for error p urposes

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

	//encoder writes the resultant "Response" struct (see outputFormatter) to writer
	enc := json.NewEncoder(w)
	enc.Encode(rows)

}

func main() {
	fmt.Println("Starting server.")
	flag.Parse()

	http.HandleFunc("/api/", apiHandler)
	http.HandleFunc("/request/", requestHandler)

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
