# Precision Analytics

An analytics backend built with Go

### Introduction

This analytics backend provides a simple REST API interface for logging messages from applications or services. All messages recieved by the server contain the following parameters in JSON format:

```javascript
{
	// JWT Token
	"token": "",
	// A UUID for this message
	"id": "",

	// Platform name (Web, Android, iOS, MacOS, Windows, Linux)
	"platform": "Android",

	// App or Service name
	"namespace": "com.companyname.appname",

	// UUID for user
	"userId": "",

	// Date message was logged
	"date": "",

	// Name of message type
	"msgType": "",

	// Collection of key-type-values
	"msg":[
		{
			"name": "user-search",
			"type": "bool",
			"value": "Ferrari 355"
		}
	]
}
```

User-defined parameters are contained in the `msg` parameter of the JSON object. Authentication is provided via login using a valid API key, after which the client is provided with an authorization token for posting data to the server.

Messages logged by the server are stored in database (currently SQL) which can be queried through the API using a seperate API key. API keys are stored in the database with provisions for accessing different services of the API.

This backend provides no front-end to view analytics data, allowing you to choose which frontend is most suitable, such as [Google Analytics Data Studio](https://www.google.com/analytics/data-studio/).

### Authentication

To authenticate with the server and receive an access token, POST to: `[address]:[port]/[version]/auth with the following JSON:

```javascript
{
	// Valid API key with permission to use service
	"apiKey": "",
	// UUID of user
	"userId": ""
}
```

If authentication is successful, the client will recieve the response:

```javascript
{
	// JWT Token
	"token": ""
}
```

The token is then used in all future POST requests to the server.

## Versioning

Many API requests are prepended with a version. This ensures when a user performs an action like as logging a message, the format will not change for that version of the API. If the format of a request is changed, the version number must be changed to maintain backwards-compatibility.

### Routes

#### Log Messages

`/[version]/log`


### Testing with Curl

#### Get Auth

`curl -H "Content-Type: application/json" -d '{"apiKey":"61529673-6c86-4f54-9bdd-838bf12360a6", "userId":"johnsmith"}' http://localhost:8080/v1/auth`

Reply:

```javascript
{
	"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJ0b213YWxrZXIiLCJpc3MiOiJwcmVjaXNpb24tYW5hbHl0aWNzIn0.ipEkn7bHv53h72R-KBhErlQtLyThsYfbP6dz41tLiuM"
}
```

#### Log with Auth

curl -H "Content-Type: application/json" -d '{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJleHAiOjE0OTIxMDU1ODYsImlhdCI6MTQ5MjEwMTk4NiwiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyMTAxOTg2fQ.hwbhzsZxnGsKtSGnLzrmefQ6w7-H03CUXHGBqL7pVNE","id":"0000000004","platform":"web","namespace":"test","version":"1.0.0","userId":"testUser","sessionId":"testSession","date":"0001-01-01T00:00:00Z","msgType":"testType","msg":[{"key":"testKey", "type":"string", "value":"testval"},{"key":"testKey2", "type":"string", "value":"testval2"}]}' http://localhost:8080/v1/log




