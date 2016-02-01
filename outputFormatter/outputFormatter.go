package outputFormatter

/******************************************************************
outputFormatter contains:
* Wrapper struct, which is written to the stream in server.go
  * * Resp interface{}, which is the response to the user query
    * Version, which is the version of the API
* MakeWrapper(r interface{}), which wraps r into a struct to encode
	*****************************************************************/

//import "fmt"

type ApiWrapper struct {
	Request  interface{} `json:"request"`
	Metadata interface{} `json:"metadata"`
	Records  interface{} `json:"records"`
}

type MetadataWrapper struct {
	RecordsPerQuery int `json:"recordsperquery"`
	NumRecords      int `json:"numrecords"`
	Pages           int `json:"numpages"`
	Page            int `json:"curpage"`
}

type InfoWrapper struct {
	NumRecords int      `json:"numrecords"`
	Records    []string `json:"records"`
}

//wraps the given interface r into a returned Wrapper
//prepped for encoding to stream
func MakeApiWrapper(req interface{}, records interface{}, numRecords int, specialParameters map[string]int) interface{} {
	//fmt.Println(req)
	//fmt.Println(records)
	//fmt.Println(numRecords)
	//fmt.Println(specialParameters)

	//version is hard coded to "1.1"
	//all of this is variable
	metadataWrapper := MakeMetadataWrapper(numRecords, specialParameters)
	w := ApiWrapper{req, metadataWrapper, records}

	return w
}

func MakeInfoWrapper(numRecords int, records []string) interface{} {
	w := InfoWrapper{numRecords, records}

	return w
}

func MakeMetadataWrapper(numRecords int, specialParameters map[string]int) interface{} {
	pages := numRecords / (specialParameters["size"] + 1)

	metadataWrapper := MetadataWrapper{specialParameters["size"], numRecords, pages, specialParameters["page"]}

	return metadataWrapper
}
