curl localhost:8000/hello/123 -i
curl localhost:8000/hello/abc -i
curl localhost:8000/hello/james/bond -i
curl localhost:8000/hello/james -i  
curl localhost:8000/hello/1 -i
curl "localhost:8000/query?name=james bond" -i
curl "localhost:8000/query?name=james bond&surname=smith" -i
curl "localhost:8000/query?name='james bond'&surname=smith" -i
curl "localhost:8000/query2?name=james bond&id=007" -i
curl "localhost:8000/wildcard/hello/world" -i 
curl "localhost:8000/error" -i
curl "localhost:8000/hello" -i -X POST
curl "localhost:8000/v1/hello" -i
curl "localhost:8000/v2/hello" -i
curl "localhost:8000/user/login" -i
curl "localhost:8000/server" -i
curl "localhost:8000/env" -i
curl "localhost:8000/env" | jq
curl "localhost:8000/body" -d 'hello' -i
curl "localhost:8000/body" -d '{"name":"James Bond"}' -i
curl "localhost:8000/body" -d '{"name":"James Bond"}' -i -H content-type:application/json
curl "localhost:8000/body" -d '{"id": 1,"name":"James Bond"}' -i -H content-type:application/json
curl "localhost:8000/body2" -d '{"id": 1,"name":"Bond"}' -i -H content-type:application/json
curl "localhost:8000/body" -d 'id=1&name=bond' -i -H content-type:application/x-www-form-urlencoded

curl "localhost:8000/signup" -H content-type:application/json -d '{"username":"bond", "password":"bond"}' -i
curl "localhost:8000/signup" -H content-type:application/json -d '{"username":"bond2", "password":"bond2"}' -i
curl "localhost:8000/signup" -H content-type:application/json -d '{"username":"bond3", "password":"bond3"}' -i
curl "localhost:8000/login" -H content-type:application/json -d '{"username":"bond", "password":"bond"}' -i

curl "localhost:8000/login" -H content-type:application/json -d '{"username":"bond", "password":"bond"}' | jq
curl "localhost:8000/hello" -i -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ1NTA2NDAsImlzcyI6IjEifQ.gskOMhDG4C0G8nwYhlwYqtpJEztpIauVgyqQJPAJi7c"
curl "localhost:8000/hello" -i -H "Authorization:Bearer ${jwtToken}"