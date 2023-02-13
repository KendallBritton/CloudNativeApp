#! /bin/sh

curl "http://localhost:8000/list"

curl "http://localhost:8000/create?item=TV&price=999.99"

curl "http://localhost:8000/create?item=Soundbar&price=450" 

curl "http://localhost:8000/list" 

curl "http://localhost:8000/delete?item=shoes" 

curl "http://localhost:8000/create?item=Pie&price=e" 

curl "http://localhost:8000/create?item=Pie&price=2.50" 

curl "http://localhost:8000/list" 

curl "http://localhost:8000/update?item=Soundbar&price=250" 

curl "http://localhost:8000/update?item=TV&price=599.99" 

curl "http://localhost:8000/read" 

curl "http://localhost:8000/price?item=Pie" 

curl "http://localhost:8000/price?item=Soundbar" 

curl "http://localhost:8000/delete?item=Pie" 