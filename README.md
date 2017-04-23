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

## Access Control

### Groups

#### Introduction

Groups restrict API key access to specific routes in the API. When an API Key is created, it is a member of a single group. The routes that the API Key can access are specified in the group description.

Two groups are created by default when the database is created, the 'root' group and the 'client' group. API Keys are created for these groups and are ('created at startup' - Not yet implemented, manually insert keys to use) and stored in the pa.config file. 

The 'root' group has access to every route in the API and is used for initial setup and configuration. It is highly recommended to create a new 'admin' group for yourself and refrain from using the root group, as it cannot be modified or removed later, making it difficult to deactivate in-case it is disclosed to a 3rd party.

The 'client' group only has access to the Log route, for logging analytics messages from an application. This is just an example and can be removed or modified later.

Except for creating a new Admin group to retreive analytics data, its possible that you will never need to interact with the API's group management. Access Control allows finely-grained permissions, for example, a group could be created that can ONLY retreive log messages, but has no control over group management or creation of API keys.

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

#### List of Permissions

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

#### Listing all Groups

#### Creating a Group

#### Modifying a Group

#### Removing a Group

### API Keys

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


#### Listing all Keys

--- | ---
** Route ** 		| /[version]/key/get
** Method ** 		| GET
** Requires Auth ** | yes

** Request Body: **

```javascript```
{

}
```

#### Creating a Key

#### Modifying a Key

#### Removing a Key


### Testing with Curl

#### Get Auth

`curl -H "Content-Type: application/json" -d '{"apiKey":"abcdefgh-1234-5678-stuvwxyz", "userId":"johnsmith"}' http://localhost:8080/v1/auth`

Reply:

```javascript
{
	"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJncm91cCI6InJvb3QiLCJleHAiOjE0OTI5ODY0MTcsImlhdCI6MTQ5Mjk4MjgxNywiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyOTgyODE3fQ.uJWde7gfvC6HkeHko4PIs3fHUmlyHW4QX1AF9maEta8"
}
```

#### Log with Auth

curl -H "Content-Type: application/json" -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJncm91cCI6InJvb3QiLCJleHAiOjE0OTI5ODY0MTcsImlhdCI6MTQ5Mjk4MjgxNywiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyOTgyODE3fQ.uJWde7gfvC6HkeHko4PIs3fHUmlyHW4QX1AF9maEta8" -d '{"token":"","id":"0000000004","platform":"web","namespace":"test","version":"1.0.0","userId":"testUser","sessionId":"testSession","date":"0001-01-01T00:00:00Z","msgType":"testType","msg":[{"key":"testKey", "type":"string", "value":"testval"},{"key":"testKey2", "type":"string", "value":"testval2"}]}' http://localhost:8080/v1/log

#### Add Group

curl -H "Content-Type: application/json" -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJncm91cCI6InJvb3QiLCJleHAiOjE0OTI5ODY0MTcsImlhdCI6MTQ5Mjk4MjgxNywiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyOTgyODE3fQ.uJWde7gfvC6HkeHko4PIs3fHUmlyHW4QX1AF9maEta8" -d '{"groups": [{"name":"admin","perms":"Log,ShowKeys,SetKeys,RemoveKeys,ShowGroups,SetGroups,RemoveGroups"}]}' http://localhost:8080/v1/group/set

#### Add Key

curl -H "Content-Type: application/json" -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJqb2huc21pdGgiLCJncm91cCI6InJvb3QiLCJleHAiOjE0OTI5ODY0MTcsImlhdCI6MTQ5Mjk4MjgxNywiaXNzIjoicHJlY2lzaW9uIiwibmJmIjoxNDkyOTgyODE3fQ.uJWde7gfvC6HkeHko4PIs3fHUmlyHW4QX1AF9maEta8" -d '{"keys": [{"key":"mynewkey", "expires":"false","expDate":"0001-01-01T00:00:00Z","active":"true","group":"admin"]}' http://localhost:8080/v1/key/set








