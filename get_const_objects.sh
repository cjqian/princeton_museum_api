#start server
#go run api.go root helloworld puamapi

#get constituent in question
const=$1
echo $const

path="http://localhost:8080/api/apiconobjxrefs?constituentid=$const"
echo $path

#curl $path
curl -s $path | python -mjson.tool
