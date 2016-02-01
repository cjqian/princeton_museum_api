package sqlParser

import (
	/*"fmt"*/
	"strconv"
	"strings"
)

/*
func GetTableIDMap() map[string]string {
	tableIDMap := make(map[string]string, 0)

	tableIDMap["apiconobjxrefs"] = "ObjectID"
	tableIDMap["apiconstituents"] = "ConstituentID"
	tableIDMap["apidimelements"] = "DimItemElemXrefID"
	tableIDMap["apidimobjxrefs"] = "ObjectID"
	tableIDMap["apimedia"] = "MediaMasterID"
	tableIDMap["apimediaxrefs"] = "ID"
	tableIDMap["apiobjects"] = "ObjectID"
	tableIDMap["apititleobjxrefs"] = "TitleID"

	return tableIDMap
}

func GetTableXRefMap() map[string]string {
	tableXrefMap := make(map[string]string, 0)

	tableXrefMap["apiconstituents"] = "apiconobjxrefs"
	tableXrefMap["apidimelements"] = "apidimobjxrefs"
	tableXrefMap["apimedia"] = "apimediaxrefs"

	return tableXrefMap
}*/

//returns a map of each column name in table to its appropriate GoLang tpye (name string)
func GetColTypeMap() map[string]string {
	colMap := make(map[string]string, 0)

	cols, err := globalDB.Queryx("SELECT DISTINCT COLUMN_NAME, COLUMN_TYPE FROM information_schema.columns")
	check(err)

	for cols.Next() {
		var colName string
		var colType string

		err = cols.Scan(&colName, &colType)
		//split because SQL type returns are sometimes ex. int(11)
		colMap[colName] = strings.Split(colType, "(")[0]
	}

	return colMap
}

//given a []byte b and type t, return b in t form
func StringToType(b []byte, t string) (interface{}, error) {
	//all unregistered types (datetime for now, etc) are type string
	s := string(b)

	if t == "bigint" || t == "int" || t == "integer" || t == "tinyint" {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return i, nil
	} else if t == "double" {
		float, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		return float, nil
	} else if t == "varchar" {
		return s, nil
	} else {
		return string(b), nil
	}

}
