curl -H "Content-Type: application/json" -d '{"id":"0000000003","platform":"web","namespace":"test","userId":"testUser","sessionId":"testSession","date":"0001-01-01T00:00:00Z","msgType":"testType","msg":[{"key":"testKey", "type":"string", "value":"testval"}]}' http://localhost:8080/v1/log