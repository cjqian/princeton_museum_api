package outputFormatter

/******************************************************************
outputFormatter contains:
* Wrapper struct, which is written to the stream in server.go
  * * Resp interface{}, which is the response to the user query
    * Version, which is the version of the API
* MakeWrapper(r interface{}), which wraps r into a struct to encode
	*****************************************************************/
/*import "fmt"*/

type ApiWrapper struct {
	Resp        interface{}     `json:"response"`
	Cols        []Column        `json:"columns"`
	ColWrappers []ColumnWrapper `json:"colWrappers"`
	Error       string          `json:"error"`
	IsTable     bool            `json:"isTable"`
	Version     float64         `json:"version"`
}

//wraps the given interface r into a returned Wrapper
//prepped for encoding to stream
func MakeApiWrapper(r interface{}, c []string, err string, isTable bool) ApiWrapper {
	//version is hard coded to "1.1"
	//all of this is variable
	w := ApiWrapper{r, MakeColumns(c), MakeColumnWrappers(c), err, isTable, 1.1}
	return w
}

type ColumnWrapper struct {
	Field        string `json:"field"`
	DisplayName  string `json:"displayName"`
	ColumnFilter bool   `json:"columnFilter"`
}

type Column struct {
	Name             string                 `json:"colName"`
	ForeignKey       bool                   `json:"isForeignKey"`
	ForeignKeyValues map[string]interface{} `json:"foreignKeyValues"`
}

func MakeColumns(columns []string) []Column {
	c := make([]Column, 0)

	for _, column := range columns {
		var w Column
		w = Column{column, false, nil}
		c = append(c, w)
	}
	return c
}

func MakeColumnWrappers(columns []string) []ColumnWrapper {
	cw := make([]ColumnWrapper, 0)
	for _, column := range columns {
		w := ColumnWrapper{column, column, true}
		cw = append(cw, w)
	}

	return cw
}
