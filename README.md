# Info
Contains information on the layout of the database .

## GET info
`GET /info` will return a list of all tables.

Each query will return 1 page of 10 objects as a default. To parse pages or increase the page size, you can also change the following parameters:

| Parameter | Value 
| :------ | -----
| Page  | The page you wish to access
| Size | The number of objects per page

Note: if you increase size by too much, performance may suffer.

## Queries
For character fields, the search functionality is case insensitive and is space-separated by underscores.

Multiple queries are separated by `&`.
The parser does accept `<` and `>` signs for numerical comparisons!

Here are the types of queries:
> http://localhost:8080/info
> This returns a  list of all the tables in the database.
>
> http://localhost:8080/info/apiobjects
> This returns all the columns in the apiobjects table. 
>
> http://localhost:8080/info/apiobjects/department
> This returns all the unique departments in the apiobjects table.
>
> http://localhost:8080/info/apiobjects/department&objectID
> This returns all the unique departments and the unique object IDs in the apiobjects table.

###Responses
```json
{
request: {
Type: "info",
TableName: "",
Parameters: [ ],
SpecialParameters: { }
},
metadata: {
recordsperquery: 0,
numrecords: 8,
numpages: 8,
curpage: 0
},
records: [
"apiconobjxrefs",
"apiconstituents",
"apidimelements",
"apidimobjxrefs",
"apimedia",
"apimediaxrefs",
"apiobjects",
"apititleobjxrefs"
]
}
```



