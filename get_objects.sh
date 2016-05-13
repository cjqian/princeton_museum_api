declare -a arr=("11540" "2340")

url="localhost:8080/api/apiobjects/"

for i in "${arr[@]}"
do
    path="$url$i"
    outpath="examples/$i.json"
    echo "$path"
    echo "$outpath"
done
