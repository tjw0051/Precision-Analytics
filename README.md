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

## Quickstart

To get started straight away, delete the existing database in the Precision-Analytics directory called 'PA.db'. Open pa.config and change the root-user API key. You could use [This UUID Generator](https://www.uuidgenerator.net) generator: 

'''
// Enter API Key for the root group (access to everything)
// NOTE: Once the server is running, access server management using
// the root API key. Create a new user ASAP with the required
// permissions and refrain from using the root group.
root-key= abcdefgh-1234-5678-stuvwxyz
'''

Run the server from the Precision-Analytics/ directory by entering:

```
./Precision-Analytics
```

Lets create an API Key for logging analytics messages from an application. First use the root key set up in pa.config to authenticate with the server. This will give us an authentication token for interacting with the server (replace YOUR_API_KEY_HERE with the key you entered in pa.config and use any userId you wish):

```
`curl -H "Content-Type: application/json" -d '{"apiKey":"YOUR_API_KEY_HERE", "userId":"johnsmith"}' http://localhost:8080/v1/auth`
```

The server should respond with an authentication token in JSON format:

"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJncm91cCI6InJvb3QiLCJleHAiOjE0OTI5ODY0MTcsImlhdCI6MTQ5Mjk4MjgxNywiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyOTgyODE3fQ.uJWde7gfvC6HkeHko4PIs3fHUmlyHW4QX1AF9maEta8"

**NOTE:**: If the token expires you will need to request a new one using the method above.

We can now use the token we received to create a new API Key. Helpfully, an access-group with only the permission to log messages to the server has been created by default, the 'client' group. This limits the ability of our key to only log messages, so that a malicious user can't de-compile our application, steal the key, and cause all sorts of shenanigans on our server:

'''
curl -H "Content-Type: application/json" -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJncm91cCI6InJvb3QiLCJleHAiOjE0OTI5ODY0MTcsImlhdCI6MTQ5Mjk4MjgxNywiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyOTgyODE3fQ.uJWde7gfvC6HkeHko4PIs3fHUmlyHW4QX1AF9maEta8" -d '{"keys": [{"key":"", "expires":"false","expDate":"0001-01-01T00:00:00Z","active":"true","group":"client"}]}' http://localhost:8080/v1/key/set
'''

In the above JSON, notice the "key" property has been left blank. The key can be provided by the user, otherwise the server will autogenerate the key in UUID version 4 format if the property is left empty. The server will then reply with the created key:

```
{
	"keys": [
	{
		"key":"81716f18-6ce3-4136-920d-80734e7fa3f7", 
		"expires":"false",
		"expDate":"0001-01-01T00:00:00Z",
		"active":"true",
		"group":"client"
	]
}
```

We can now use this newly-created API Key to retrieve a token, and log messages from our application.

Get Auth Token:

'''
`curl -H "Content-Type: application/json" -d '{"apiKey":"81716f18-6ce3-4136-920d-80734e7fa3f7", "userId":"johnsmith"}' http://localhost:8080/v1/auth`
'''

Log Message:

'''
curl -H "Content-Type: application/json" -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJncm91cCI6InJvb3QiLCJleHAiOjE0OTI5ODY0MTcsImlhdCI6MTQ5Mjk4MjgxNywiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyOTgyODE3fQ.uJWde7gfvC6HkeHko4PIs3fHUmlyHW4QX1AF9maEta8" -d '{"token":"","id":"0000000004","platform":"web","namespace":"test","version":"1.0.0","userId":"testUser","sessionId":"testSession","date":"0001-01-01T00:00:00Z","msgType":"testType","msg":[{"key":"testKey", "type":"string", "value":"testval"},{"key":"testKey2", "type":"string", "value":"testval2"}]}' http://localhost:8080/v1/log
'''

## Authentication

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

The token is then used to access API routes that require authentication. The token is sent to the server in the HTTP header using the Bearer format:

```
Authentication: BEARER [token]
```

## Versioning

Many API requests are prepended with a version. This ensures when a user performs an action like as logging a message, the format will not change for that version of the API. If the format of a request is changed, the version number must be changed to maintain backwards-compatibility.

## Groups

### Introduction

Groups restrict API key access to specific routes in the API. When an API Key is created, it is a member of a single group. The routes that the API Key can access are specified in the groups 'perms' property.

Two groups are created by default when the database is created, the 'root' group and the 'client' group. API Keys are created for these groups and are ('created at startup' - Not yet implemented, manually insert keys to use) and stored in the pa.config file. 

The 'root' group has access to every route in the API and is used for initial setup and configuration. It is highly recommended to create a new 'admin' group for yourself and refrain from using the root group, as it cannot be modified or removed later, making it difficult to deactivate if it is disclosed to a 3rd party.

The 'client' group only has access to the Log route, for logging analytics messages from an application. This is just an example and can be removed or modified later.

Except for creating a new Admin group, it is possible that you will never need to interact with the API's group management. Access Control allows finely-grained permissions, for example, a group could be created that can ONLY retreive log messages, but has no control over group management or creation of API keys. In most circumstances however, there will be many API keys attached to the 'client' group for logging messages from applications, and 1 Key attached to the Admin group for viewing the analytics data received by the server.

### Group Structure

The structure of a group is as follows:

```javascript
{
	"name": "admin",
	"perms": "Log,ShowKeys,SetKeys,RemoveKeys,ShowGroups,SetGroups,RemoveGroups"
}
```

The properties are:
- **name** - The name of the group (case-sensitive)
- **perms** - Comma-delimited list of permissions

### List of Permissions

**NOTE:** All routes are prepended with an API version number. E.g. For version 1, the log route would be *example.com/v1/log*

| Permission   | Route          | description                  |
|--------------|----------------|------------------------------|
| Log          | /log           | Log analytics messages       |
| ShowKeys     | /key/get       | List all API Keys            |
| SetKeys      | /key/set       | Create/Modify a Key          |
| RemoveKeys   | /key/remove    | Remove a key                 |
| ShowGroups   | /group/get     | List all groups              |
| SetGroups    | /group/set     | Create/Modify groups         |
| RemoveGroups | /group/remove  | Remove groups                |

### Listing all Groups

### Creating a Group

### Modifying a Group

### Removing a Group

## API Keys

API Keys are used to access routes in the API. As outlined in the Groups section, every API key is a member of a group, which defines what routes the key can access. The structure of a key is as follows:

```javascript
{
	"key": "",
	"expires": "true",
	"expDate": "",
	"active": "true",
	"group": "admin"
}
```
The purpose of each property is:
- **key** - The API Key, typically in UUID version 4 format.
- **expires** - Does the key expire or not
- **expDate** - if *expires* is true, expDate is the date the key expires
- **active** - If false, the key is deactivated and cannot be used
- **group** - The group this key is associated with.


### Listing all Keys

--- | ---
** Route ** 		| /[version]/key/get
** Method ** 		| GET
** Requires Auth ** | yes

** Request Body: **

```javascript```
{

}
```

### Creating a Key

### Modifying a Key

### Removing a Key


## Testing with Curl

### Get Auth

`curl -H "Content-Type: application/json" -d '{"apiKey":"abcdefgh-1234-5678-stuvwxyz", "userId":"johnsmith"}' http://localhost:8080/v1/auth`

Reply:

```javascript
{
	"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJncm91cCI6InJvb3QiLCJleHAiOjE0OTI5ODY0MTcsImlhdCI6MTQ5Mjk4MjgxNywiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyOTgyODE3fQ.uJWde7gfvC6HkeHko4PIs3fHUmlyHW4QX1AF9maEta8"
}
```

### Log with Auth

curl -H "Content-Type: application/json" -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJncm91cCI6InJvb3QiLCJleHAiOjE0OTI5ODY0MTcsImlhdCI6MTQ5Mjk4MjgxNywiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyOTgyODE3fQ.uJWde7gfvC6HkeHko4PIs3fHUmlyHW4QX1AF9maEta8" -d '{"token":"","id":"0000000004","platform":"web","namespace":"test","version":"1.0.0","userId":"testUser","sessionId":"testSession","date":"0001-01-01T00:00:00Z","msgType":"testType","msg":[{"key":"testKey", "type":"string", "value":"testval"},{"key":"testKey2", "type":"string", "value":"testval2"}]}' http://localhost:8080/v1/log

### Add Group

curl -H "Content-Type: application/json" -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJncm91cCI6InJvb3QiLCJleHAiOjE0OTI5ODY0MTcsImlhdCI6MTQ5Mjk4MjgxNywiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyOTgyODE3fQ.uJWde7gfvC6HkeHko4PIs3fHUmlyHW4QX1AF9maEta8" -d '{"groups": [{"name":"admin","perms":"Log,ShowKeys,SetKeys,RemoveKeys,ShowGroups,SetGroups,RemoveGroups"}]}' http://localhost:8080/v1/group/set

### Add Key

curl -H "Content-Type: application/json" -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJncm91cCI6InJvb3QiLCJleHAiOjE0OTI5ODY0MTcsImlhdCI6MTQ5Mjk4MjgxNywiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyOTgyODE3fQ.uJWde7gfvC6HkeHko4PIs3fHUmlyHW4QX1AF9maEta8" -d '{"keys": [{"key":"mynewkey", "expires":"false","expDate":"0001-01-01T00:00:00Z","active":"true","group":"admin"]}' http://localhost:8080/v1/key/set








