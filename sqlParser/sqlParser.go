package sqlParser

import (
	_ "./mysql"
	"./sqlx"
	"fmt"
	"strconv"
	//"encoding/json"
)

var tableNameToId = map[string]string{
	"apiobjects":      "ObjectID",
	"apiconstituents": "ConstituentID",
}

var (
	globalDB sqlx.DB

	colTypeMap map[string]string
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*********************************************************************************
 * DB INITIALIZE: Connects given DB creds, creates ColMap FOR SESSION
 ********************************************************************************/
func InitializeDatabase(username string, password string, environment string) sqlx.DB {
	db, err := sqlx.Connect("mysql", username+":"+password+"@tcp(localhost:3306)/"+environment)
	check(err)

	globalDB = *db

	colTypeMap = GetColTypeMap()
	return *db
}

/*********************************************************************************
 * HELPER FUNCTIONS
 ********************************************************************************/
//returns array of table name strings from queried database
func GetTableNames() []string {
	var tableNames []string

	tableRawBytes := make([]byte, 1)
	tableInterface := make([]interface{}, 1)

	tableInterface[0] = &tableRawBytes

	rows, err := globalDB.Query("show tables")
	check(err)

	for rows.Next() {
		err := rows.Scan(tableInterface...)
		check(err)

		tableNames = append(tableNames, string(tableRawBytes))
	}

	return tableNames
}

//returns array of column names from table in database
func GetColumnNames(tableName string) []string {
	tableRawBytes := make([]byte, 1)
	tableInterface := make([]interface{}, 1)

	tableInterface[0] = &tableRawBytes

	rows, err := globalDB.Query("SELECT distinct column_name from information_schema.columns where table_name='" + tableName + "' ORDER BY column_name asc")
	check(err)

	colMap := make([]string, 0)

	for rows.Next() {
		err = rows.Scan(tableInterface...)
		check(err)

		colMap = append(colMap, string(tableRawBytes))
	}

	return colMap
}

func GetColumnValues(tableName string, columnName string) []string {
	//fmt.Println(tableName + ", " + columnName)
	tableRawBytes := make([]byte, 1)
	tableInterface := make([]interface{}, 1)

	tableInterface[0] = &tableRawBytes

	rows, err := globalDB.Query("select distinct " + columnName + " from " + tableName)
	check(err)

	colMap := make([]string, 0)

	for rows.Next() {
		err = rows.Scan(tableInterface...)
		check(err)

		colMap = append(colMap, string(tableRawBytes))
	}

	return colMap
}

/*********************************************************************************
 * APIOBJECTS HELPER FUNCTIONS
 ********************************************************************************/
func QueryRows(queryStr string) []map[string]interface{} {
	fmt.Println(queryStr)

	rows, err := globalDB.Queryx(queryStr)
	if err != nil {
		panic(err)
	}

	rowArray := make([]map[string]interface{}, 0)

	for rows.Next() {
		//map the column name to its value
		results := make(map[string]interface{}, 0)
		err = rows.MapScan(results)
		if err != nil {
			panic(err)
		}

		for k, v := range results {
			if b, ok := v.([]byte); ok {
				results[k], err = StringToType(b, colTypeMap[k])
				if err != nil {
					panic(err)
				}
			}
		}

		rowArray = append(rowArray, results)
	}

	return rowArray
}

//returns data from the apiconobjxrefs data
func GetBibliography(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := "select ReferenceID, ObjCitation, apibibobjxrefs.SysTimeStamp from apibibobjxrefs INNER JOIN apiobjects " + whereStr + limStr

	return QueryRows(queryStr)
}

func GetConstituentsTrunc(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	params := "apiconstituents.Active, apiconstituents.AlphaSort, apiconstituents.Approved, apiconstituents.BeginDate, apiconstituents.BeginDateISO, apiconstituents.ConstituentID, apiconobjxrefs.DisplayOrder, apiconobjxrefs.Displayed, apiconobjxrefs.Prefix, apiconobjxrefs.Remarks, apiconobjxrefs.Role, apiconobjxrefs.Suffix, apiconstituents.SysTimeStamp"

	queryStr := "select " + params + " from apiconobjxrefs " + "INNER JOIN apiobjects ON apiconobjxrefs.ObjectID = apiobjects.ObjectID " + "INNER JOIN apiconstituents ON apiconobjxrefs.ConstituentID = apiconstituents.ConstituentID " + whereStr + " ORDER BY apiconobjxrefs.DisplayOrder" + limStr

	fmt.Println(queryStr)

	return QueryRows(queryStr)
}

func GetDimElements(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := `select apidimelements.*, apidimobjxrefs.* from apidimobjxrefs 
		INNER JOIN apiobjects ON apidimobjxrefs.ObjectID = apiobjects.ObjectID 
		INNER JOIN apidimelements ON apidimobjxrefs.DimItemElemXrefID = apidimelements.DimItemElemXrefID `
	queryStr += whereStr + limStr

	return QueryRows(queryStr)
}

//returns data from the apiconobjxrefs data
func GetExhibitions(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := "select ExhibitionID, RunningCaption, apiexhobjxrefs.SysTimeStamp from apiexhobjxrefs INNER JOIN apiobjects " + whereStr + limStr

	return QueryRows(queryStr)
}

//returns data from the apiconobjxrefs data
func GetGeography(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := `select ObjGeographyID, apiobjgeography.ObjectID, GeoCode, PrimaryDisplay, Continent, SubContinent, Country, 
		Region, State, City, Country, SubRegion, Locale, Locus, River, Excavation, 
		Latitude, Longitude, GeoNames, apiobjgeography.SysTimeStamp from apiobjgeography INNER JOIN apiobjects `
	queryStr += whereStr + limStr

	return QueryRows(queryStr)
}

func GetMedia(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := `select apimedia.* from apimediaxrefs 
		INNER JOIN apiobjects ON apimediaxrefs.ID = apiobjects.ObjectID 
		INNER JOIN apimedia ON apimediaxrefs.MediaMasterID = apimedia.MediaMasterID `

	queryStr += whereStr + limStr

	return QueryRows(queryStr)
}

//returns data from the apiconobjxrefs data
func GetTerms(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := "select * from apitermsobjxrefs INNER JOIN apiobjects " + whereStr + limStr

	return QueryRows(queryStr)
}

func GetTitles(whereStr string, limStr string, rowCount int) []map[string]interface{} {
	queryStr := "select * from apititleobjxrefs INNER JOIN apiobjects " + whereStr + limStr

	return QueryRows(queryStr)
}

//returns data from the apiconobjxrefs data
func GetUri(constituentID int, channel chan interface{}) {
	queryStr := "select * from apiconuris where ConstituentID = " + strconv.Itoa(constituentID)
	rows, err := globalDB.Queryx(queryStr)
	if err != nil {
		panic(err)
	}

	rowArray := make([]map[string]interface{}, 0)

	//for each row
	for rows.Next() {
		//map the column name to its value
		results := make(map[string]interface{}, 0)
		err = rows.MapScan(results)
		if err != nil {
			panic(err)
		}

		for k, v := range results {
			if b, ok := v.([]byte); ok {
				results[k], err = StringToType(b, colTypeMap[k])
				if err != nil {
					panic(err)
				}
			}
		}

		rowArray = append(rowArray, results)
	}

	channel <- rowArray
}

/*********************************************************************************
 * GET TABLES
 ********************************************************************************/
//func QueryConstituents(tableName string, idVal int, results map[string]interface{}) {
//uriChan := make(chan interface{})
//go GetUri(idVal, uriChan)
//results["URIs"] = <-uriChan
//}

func QueryObjects(whereStr string, limStr string, rowCount int, results []map[string]interface{}) {
	bibliographyResults := GetBibliography(whereStr, limStr, rowCount)
	constituentResults := GetConstituentsTrunc(whereStr, limStr, rowCount)
	dimensionResults := GetDimElements(whereStr, limStr, rowCount)
	exhibitionResults := GetExhibitions(whereStr, limStr, rowCount)
	geographyResults := GetGeography(whereStr, limStr, rowCount)
	mediaResults := GetMedia(whereStr, limStr, rowCount)
	termResults := GetTerms(whereStr, limStr, rowCount)
	titleResults := GetTitles(whereStr, limStr, rowCount)

	for i := 0; i < rowCount; i++ {
		results[i]["Bibliography"] = bibliographyResults[i]
		results[i]["Constituents"] = constituentResults[i]
		results[i]["Dimensions"] = dimensionResults[i]
		results[i]["Exhibitions"] = exhibitionResults[i]
		results[i]["Geography"] = geographyResults[i]
		results[i]["Media"] = mediaResults[i]
		results[i]["Terms"] = termResults[i]
		results[i]["Titles"] = titleResults[i]
	}
}

/*********************************************************************************
 * GET FUNCTIONALITY
 ********************************************************************************/
func GetNumRows(tableName string, whereStr string, limStr string) int {
	var rowCount int
	err := globalDB.Get(&rowCount, "SELECT count(*) FROM "+tableName+" "+whereStr+limStr)
	if err != nil {
		panic(err)
	}

	fmt.Println(rowCount)

	return rowCount
}

func GetWhereString(tableParameters []string) string {
	whereStr := ""
	//where
	if len(tableParameters) > 0 {
		whereStr += " where "

		for _, v := range tableParameters {
			whereStr += v + " and "
		}

		whereStr = whereStr[:len(whereStr)-4]
	}

	return whereStr
}

func GetLimString(specialParameters map[string]int) string {
	//pagination
	size := specialParameters["size"]
	page := specialParameters["page"]

	//limitation
	startNum := (page - 1.0) * size
	limStr := " LIMIT " + strconv.Itoa(startNum) + ", " + strconv.Itoa(size)

	return limStr
}

func Get(tableName string, tableParameters []string, specialParameters map[string]int) ([]map[string]interface{}, error) {

	regStr := ""
	joinStr := ""

	cols := GetColumnNames(tableName)
	for _, col := range cols {
		regStr += tableName + "." + col + ","
	}

	regStr = regStr[:len(regStr)-1]

	if joinStr != "" {
		joinStr = ", " + joinStr[:len(joinStr)-1]
	}

	whereStr := GetWhereString(tableParameters)
	limStr := GetLimString(specialParameters)
	queryStr := "select " + regStr + joinStr + " from " + tableName + " " + whereStr + limStr

	//do the query
	rows, err := globalDB.Queryx(queryStr)
	if err != nil {
		return nil, err
	}

	//get number of rows
	rowCount := GetNumRows(tableName, whereStr, limStr)
	fmt.Println(rowCount)

	//map into an array of type map[colName]value
	rowArray := make([]map[string]interface{}, rowCount)

	//query the special tables
	if tableName == "apiobjects" {
		QueryObjects(whereStr, limStr, rowCount, rowArray)
	} else if tableName == "apiconstituents" {
		//QueryConstituents(whereStr, size)
	}

	//for each row
	i := 0
	for rows.Next() {
		//map the column name to its value
		results := make(map[string]interface{}, 0)
		err = rows.MapScan(results)
		if err != nil {
			return nil, err
		}

		for k, v := range results {
			if b, ok := v.([]byte); ok {
				rowArray[i][k], err = StringToType(b, colTypeMap[k])
				if err != nil {
					return nil, err
				}
			}

		}

		i++
	}

	return rowArray, nil
}
