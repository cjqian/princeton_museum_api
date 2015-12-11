#Object

Contains information on objects in the collections. This is one of the tables you can explore in the [api](www.github.com).

##Get apiobjects

`GET /api/apiobjects` will return all objects.

Here are parameters you can use to filter: (for more information, use `/info`)

| Parameter | Value 
| :--------- | ----- 
| ObjectID | ID number of the object 
| ObjectNumber | String field of the object number 
| SortNumber | Description
| Dated | Description
| DateBegin | Description
| DateEnd | Description
| Medium | Description
| DimensionsLabel | Description
| CreditLine | Description
| Restrictions | Description
| CatRais | Description
| Edition | Description
| Department | Description
| Classification | Description
| ObjectStatus | Description
| CuratorApproved | Description
| NoWebUse | Description
| SysTimeStamp | Description

##Queries
For character fields, the search functionality is case insensitive and is space-separated by underscores.

Multiple queries are separated by `&`. 
The parser does accept `<` and `>` signs for numerical comparisons!

Here are some example queries:

> http://localhost:8080/apiobjects?objectID=2497
> This returns an object with objectID 2497.
>
> http://localhost:8080/apiobjects?department=american_art
> This returns all the objects in the American Art department. 
>
> http://localhost:8080/apiobjects?department=american_art&objectID<7000
> This returns all the objects in the American Art department with objectID < 7000.


Special case: you can also query ID (in this case, objectID) as follows:


> http://localhost:8080/apiobjects/2497
> This also returns an object with objectID 2497.

###Responses

Here's an example response.

```json
{
    request: {
        Type: "api",
        TableName: "apiobjects",
          Parameters: [
          "objectid=2497"   
          ]
     },
    numresults: 1,
    results: [
        {
            CatRais: null,
            Classification: "Ceramic",
            CreditLine: "Museum purchase, Fowler McCormick, Class of 1921, Fund",
            CuratorApproved: "0",
            DateBegin: -2000,
            DateEnd: -2000,
            Dated: "ca. 2000 B.C.",
            Department: "Asian Art",
            Dimensions: [
                {
                    DimItemElemXrefID: 51058,
                    DimUnit: "centimeters",
                    Dimension: "20.5000000000",
                    DimensionType: "Height",
                    Rank: 1
                },
                {
                    DimItemElemXrefID: 51058,
                    DimUnit: "centimeters",
                    Dimension: "17.0000000000",
                    DimensionType: "Width",
                    Rank: 3
                },
                {
                    DimItemElemXrefID: 51058,
                    DimUnit: "centimeters",
                    Dimension: "12.7000000000",
                    DimensionType: "Depth",
                    Rank: 4
                }
            ],
            DimensionsLabel: "h. 20.5 cm., w. 17.0 cm., d. 12.7 cm. (8 1/16 x 6 11/16 x 5 in.)",
            Edition: null,
            Medium: "Earthenware",
            NoWebUse: 0,
            ObjectID: 2497,
            ObjectNumber: "2000-345",
            ObjectStatus: "Accessioned Object",
            Restrictions: null,
            SortNumber: " 2000 345 ",
            SysTimeStamp: "" 
        }
 ]
}
```
