package sqlParser

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
	//dbUrl := "mysql-puamapi.catfcdegqces.us-west-2.rds.amazonaws.com:3306"
	fmt.Println("attempting to connect")
	//db, err := sqlx.Connect("mysql", username+":"+password+"@tcp(mysql-puamapi.catfcdegqces.us-west-2.rds.amazonaws.com:3306)/"+environment)
	db, err := sqlx.Connect("mysql", username+":"+password+"@tcp(localhost:3306)/"+environment)
	check(err)
	fmt.Println("connected")
	globalDB = *db

	colTypeMap = GetColTypeMap()
	fmt.Println("Connected")
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
	tableRawBytes := make([]byte, 1)
	tableInterface := make([]interface{}, 1)

	tableInterface[0] = &tableRawBytes

	queryStr := "select distinct " + columnName + " from " + tableName
	fmt.Println(queryStr)
	rows, err := globalDB.Query(queryStr)
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
	fmt.Println("Calling " + queryStr)
	rows, err := globalDB.Queryx(queryStr)
	fmt.Printf(queryStr + "\n")
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

func GetObjQueryStr(tableName string, whereStr string, limStr string) string {
	queryStr := "select * from " + tableName + " as apiobjects " + whereStr + limStr
	return queryStr
}

//returns data from the apiobjconxrefs data
func GetObjBibliography(whereStr string, limStr string) []map[string]interface{} {
	return QueryRows(GetObjQueryStr("apiobjbibxrefs", whereStr, limStr))
	//queryStr := "select ReferenceID, ObjCitation, apibibobjxrefs.SysTimeStamp, apiobjects.ObjectID from apibibobjxrefs INNER JOIN apiobjects ON apibibobjxrefs.ObjectID = apiobjects.ObjectID " + whereStr + limStr

	//return QueryRows(queryStr)
}

func GetObjConstituents(whereStr string, limStr string) []map[string]interface{} {
	/*
		params := "apiconstituents.Active, apiconstituents.AlphaSort, apiconstituents.Approved, apiconstituents.BeginDate, apiconstituents.BeginDateISO, apiconstituents.ConstituentID, apiobjconxrefs.DisplayOrder, apiobjconxrefs.Displayed, apiobjconxrefs.Prefix, apiobjconxrefs.Remarks, apiobjconxrefs.Role, apiobjconxrefs.Suffix, apiconstituents.SysTimeStamp, apiobjects.ObjectID"

		//queryStr := "select " + params + " from apiobjconxrefs " + "INNER JOIN apiobjects ON apiobjconxrefs.ObjectID = apiobjects.ObjectID " + "INNER JOIN apiconstituents ON apiobjconxrefs.ConstituentID = apiconstituents.ConstituentID " + whereStr + limStr + " ORDER BY apiobjconxrefs.DisplayOrder"

		//TODO: get order by query right
		queryStr := "select " + params + " from apiobjconxrefs " + "INNER JOIN apiobjects ON apiobjconxrefs.ObjectID = apiobjects.ObjectID " + "INNER JOIN apiconstituents ON apiobjconxrefs.ConstituentID = apiconstituents.ConstituentID " + whereStr + limStr
	*/
	return QueryRows(GetObjQueryStr("apiobjconxrefs", whereStr, limStr))
}

func GetObjDimElements(whereStr string, limStr string) []map[string]interface{} {

	return QueryRows(GetObjQueryStr("apiobjdimxrefs", whereStr, limStr))
	/*
		queryStr := "select apiobjdimxrefs.* from apiobjdimxrefs INNER JOIN apiobjects ON apiobjdimxrefs.ObjectID = apiobjects.ObjectID " + whereStr + limStr
		//queryStr := `select apidimelements.*, apiobjdimxrefs.* from apiobjdimxrefs
		//INNER JOIN apiobjects ON apiobjdimxrefs.ObjectID = apiobjects.ObjectID
		//INNER JOIN apidimelements ON apiobjdimxrefs.DimItemElemXrefID = apidimelements.DimItemElemXrefID ` + whereStr + limStr

		return QueryRows(queryStr)
	*/
}

//returns data from the apiobjconxrefs data
func GetObjExhibitions(whereStr string, limStr string) []map[string]interface{} {

	return QueryRows(GetObjQueryStr("apiobjexhxrefs", whereStr, limStr))
	/*
		queryStr := "select ExhibitionID, RunningCaption, apiexhobjxrefs.SysTimeStamp, apiobjects.ObjectID from apiexhobjxrefs INNER JOIN apiobjects on apiobjects.ObjectID = apiexhobjxrefs.ObjectID " + whereStr + limStr

		return QueryRows(queryStr)
	*/
}

//returns data from the apiobjconxrefs data
func GetObjGeography(whereStr string, limStr string) []map[string]interface{} {

	return QueryRows(GetObjQueryStr("apiobjgeography", whereStr, limStr))
	/*
		ueryStr := `select ObjGeographyID, apiobjgeography.ObjectID, GeoCode, PrimaryDisplay, Continent, SubContinent, Country,
			Region, State, City, Country, SubRegion, Locale, Locus, River, Excavation,
			Latitude, Longitude, GeoNames, apiobjgeography.SysTimeStamp from apiobjgeography INNER JOIN apiobjects ON apiobjects.ObjectID = apiobjgeography.ObjectID` + whereStr + limStr

		return QueryRows(queryStr)
	*/
}

func GetObjMedia(whereStr string, limStr string) []map[string]interface{} {

	return QueryRows(GetObjQueryStr("apiobjmediaxrefs", whereStr, limStr))
	/*
		queryStr := `select apimedia.* from apiobjmediaxrefs
			INNER JOIN apiobjects ON apiobjmediaxrefs.ObjectID = apiobjects.ObjectID
			INNER JOIN apimedia ON apiobjmediaxrefs.MediaMasterID = apimedia.MediaMasterID ` + whereStr + limStr

		return QueryRows(queryStr)
	*/
}

//returns data from the apiobjconxrefs data
func GetObjTerms(whereStr string, limStr string) []map[string]interface{} {

	return QueryRows(GetObjQueryStr("apiobjtermsxrefs", whereStr, limStr))
	/*
		queryStr := "select apiobjtermsxrefs.* from apiobjtermsxrefs INNER JOIN apiobjects ON apiobjects.ObjectID = apiobjtermsxrefs.ObjectID " + whereStr + limStr

		return QueryRows(queryStr)*/
}

func GetObjTitles(whereStr string, limStr string) []map[string]interface{} {

	return QueryRows(GetObjQueryStr("apiobjtitlexrefs", whereStr, limStr))
	/*
		//change it back to apititleobjxrefs when new data is pulled NOTE
		queryStr := "select apititleobjxrefs.* from apititleobjxrefs INNER JOIN apiobjects ON apiobjects.ObjectID = apititleobjxrefs.ObjectID " + whereStr + limStr
		fmt.Println(queryStr)
		return QueryRows(queryStr)
	*/
}

//returns data from the apiobjconxrefs data
func GetConstAltNames(whereStr string, limStr string) []map[string]interface{} {
	queryStr := "select * from apiconaltnames as apiconstituents " + whereStr

	return QueryRows(queryStr)
}

//returns data from the apiobjconxrefs data
func GetConstGeography(whereStr string, limStr string) []map[string]interface{} {
	queryStr := "select * from apicongeography as apiconstituents " + whereStr

	return QueryRows(queryStr)
}

//returns data from the apiobjconxrefs data
func GetConstObjects(whereStr string, limStr string) []map[string]interface{} {
	queryStr := "select objectid from apiconobjxrefs as apiconstituents " + whereStr

	return QueryRows(queryStr)
}

//returns data from the apiobjconxrefs data
func GetConstUri(whereStr string, limStr string) []map[string]interface{} {
	//hack because wherestr recognises constieutn
	queryStr := "select * from apiconuris as apiconstituents " + whereStr

	return QueryRows(queryStr)
}

func AddSubObjects(curResult map[string]interface{}, tableName string, subResults []map[string]interface{}, subIdx int) int {
	subArray := make([]interface{}, 0)

	for subIdx < len(subResults) && curResult["ObjectID"] == subResults[subIdx]["ObjectID"] {
		delete(subResults[subIdx], "ObjectID")
		subArray = append(subArray, subResults[subIdx])
		subIdx++
	}

	curResult[tableName] = subArray

	return subIdx
}

/*********************************************************************************
 * GET TABLES
 ********************************************************************************/
func QueryConstituents(whereStr string, limStr string, rowCount int, results []map[string]interface{}) {
	fmt.Println(whereStr)
	altResults := GetConstAltNames(whereStr, limStr)
	geoResults := GetConstGeography(whereStr, limStr)
	//objResults := GetConstObjects(whereStr, limStr)
	uriResults := GetConstUri(whereStr, limStr)

	altIdx := 0
	geoIdx := 0
	//objIdx := 0
	uriIdx := 0

	for i := 0; i < len(results); i++ {
		altIdx = AddSubObjects(results[i], "AltNames", altResults, altIdx)
		geoIdx = AddSubObjects(results[i], "Geography", geoResults, geoIdx)
		//objIdx = AddSubObjects(results[i], "Objects", objResults, objIdx)
		uriIdx = AddSubObjects(results[i], "URIs", uriResults, uriIdx)
	}
}

func QueryObjects(whereStr string, limStr string, rowCount int, results []map[string]interface{}) {
	fmt.Println("Bibliography")
	bibliographyResults := GetObjBibliography(whereStr, limStr)
	fmt.Println("Constituents")
	constituentResults := GetObjConstituents(whereStr, limStr)

	//fmt.Println("DimElements")
	dimensionResults := GetObjDimElements(whereStr, limStr)

	fmt.Println("Exhibitions")
	exhibitionResults := GetObjExhibitions(whereStr, limStr)
	fmt.Println("Geography")
	geographyResults := GetObjGeography(whereStr, limStr)
	//TODO: speed up apidimelements?
	//TODO: fix geography, terms, titles

	fmt.Println("Media")
	mediaResults := GetObjMedia(whereStr, limStr)

	//fmt.Println("Terms")
	termResults := GetObjTerms(whereStr, limStr)

	//fmt.Println("Titles")
	titleResults := GetObjTitles(whereStr, limStr)

	bibIdx := 0
	constituentIdx := 0
	dimIdx := 0
	exhIdx := 0
	geoIdx := 0
	mediaIdx := 0
	termIdx := 0
	titleIdx := 0
	for i := 0; i < len(results); i++ {

		bibIdx = AddSubObjects(results[i], "Bibliography", bibliographyResults, bibIdx)
		constituentIdx = AddSubObjects(results[i], "Constituents", constituentResults, constituentIdx)
		dimIdx = AddSubObjects(results[i], "Dimensions", dimensionResults, dimIdx)
		//results[i]["Constituents"] = constituentResults[i]
		//results[i]["Dimensions"] = dimensionResults[i]
		exhIdx = AddSubObjects(results[i], "Exhibitions", exhibitionResults, exhIdx)
		geoIdx = AddSubObjects(results[i], "Geography", geographyResults, geoIdx)

		mediaIdx = AddSubObjects(results[i], "Media", mediaResults, mediaIdx)
		//results[i]["Media"] = mediaResults[i]
		termIdx = AddSubObjects(results[i], "Terms", termResults, termIdx)
		titleIdx = AddSubObjects(results[i], "Titles", titleResults, titleIdx)
	}
}

/*********************************************************************************
 * GET FUNCTIONALITY
 ********************************************************************************/
func GetNumRows(tableName string, whereStr string, limStr string) int {
	var rowCount int
	queryStr := "SELECT count(*) FROM " + tableName + " " + whereStr + limStr

	err := globalDB.Get(&rowCount, queryStr)
	if err != nil {
		panic(err)
	}

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

//makes a view with name "queryView", no return value
func MakeView(tableName string, whereStr string, limitStr string) {
	selectStatement := "select * from " + tableName + " " + whereStr + " " + limitStr
	query := "create view queryView as " + selectStatement
	_, err := globalDB.Query(query)
	if err != nil {
		panic(err)
	}

	//now, there should be a view with the name "queryView"
}

func Get(tableName string, tableParameters []string, specialParameters map[string]int) ([]map[string]interface{}, error) {
	//pagination
	regStr := ""
	//joinStr := ""

	cols := GetColumnNames(tableName)
	for _, col := range cols {
		regStr += col + ","
	}

	if len(cols) > 0 {
		regStr = regStr[:len(regStr)-1]
	}

	//if joinStr != "" {
	//joinStr = ", " + joinStr[:len(joinStr)-1]
	//}

	whereStr := GetWhereString(tableParameters)
	limStr := GetLimString(specialParameters)

	fmt.Println("Parsed all the little things")

	//make the view
	queryStr := "select * from " + tableName + " " + whereStr + " " + limStr

	fmt.Println(queryStr)
	//get number of rows
	rowCount := GetNumRows(tableName, whereStr, limStr)

	fmt.Println("We will have \n", rowCount)

	//map into an array of type map[colName]value
	rowArray := QueryRows(queryStr)

	//query the special tables
	if tableName == "apiobjects" {
		QueryObjects(whereStr, limStr, rowCount, rowArray)
	} else if tableName == "apiconstituents" {
		QueryConstituents(whereStr, limStr, rowCount, rowArray)
	}

	//then, remove the view
	return rowArray, nil
}
