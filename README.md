# Princeton Art Museum API

Partially derived from my work on Comcast's GoTO API prototype over the summer (see my other repositories), 
this API allows for interaction with data from Princeton's art museum.

## Install

1. First, fork a copy of this sick repo. Go to a directory of your choice and type in

```
git clone https://github.com/cjqian/princeon_museum_api.git

```
2. Then, make a `.dbInfo` file that follows this syntax, 
  replacing the content in brackets with your own data:
  ```
  USERNAME="[databaseUsername]"
  PASSWORD="[databasePassword]"
  DATABASE="[databaseName]"
  ```
  For example, if you want to work with the `foo` database with username `johndoe` and password `password`, 
  your `.dbInfo file should look like this:
  ```
  USERNAME="johndoe"
  PASSWORD="password"
  DATABASE="foo"
  ```

3. Now, you can run the server by typing this into your terminal:
  ```
  ./run
  ```
    Now, you can submit requests to the server.

## Debugging
  If you're getting errors in the Install process or you happen to be Mark, make sure you can answer "yes" to
  the following questions. If you're still having issues, that really sucks.
  * Do you have the most recent version of Go [installed](https://golang.org/doc/install)? Try uninstalling/reinstalling.
  * Did you make a `.dbInfo` file? (See step two of the [Install](http://github.com/cjqian/GoTO#install) notes.)

  See `./run` for execution examples. Also, are your database credentials correct?
  * Is your `mysql` up and running? Type `mysql` into your terminal to verify.
  * Do you have the latest version of this code? Run `git pull` to get an update. 
  * Also, make sure you've checked out `master` branch and not a development branch.

## Syntax
    We're still developing the objects that you can get from this API, so this part is coming soon!

##Packages
###Local
  * sqlParser processes all interactions with the database. It contains `sqlParser.go`, which contains most of the CRUD methods, and `sqlTypeMap`, which has functions mapping values of type interface{} to string and vice-versa.
  * urlParser parses the url into a Request.

  There are more details in the comments of each of these packages.
###Other
  * I'm also using AngularJS, jQuery, Bootstrap.
  * `jmoiron/sqlx` has been super useful. Thanks!
Note: you're going to need to add these dependencies to the main api.go code. I'll add more details on this next week.
