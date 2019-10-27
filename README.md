# Web Service Specification

## Implimented endpoints

```
/repocheck/v1/commits{?limit=[0-9]+{&auth=<access-token>}}
/repocheck/v1/languages{?limit=[0-9]+{&auth=<access-token>}}
/repocheck/v1/issues{?type=(users|labels){&auth=<access-token>}}
/repocheck/v1/status
/repocheck/v1/webhooks
```

## Todos:

Delete authorisation cache after a certain amout of hours/days.
Delete method for webhooks.
Write test.
Figure out what a get payload is.

# Commits

This endpoint will return the repositories with the highest numbers of commits, with the `{?limit}` parameter indicating the number of returned repositories. 

## Request

The endpoint accepts GET requests with an empty payload.

If not specified, the parameter `limit` returns 5 repositories, 

If `auth` parameter is not specified, it will return publicly avalable repositories

## Response

Calls to the endpoint produce output according to the following JSON schema specification:

```
{ 
  "repos": [
    {
         // repository name of repo with most commits
      "repository": "path_with_namespace",
         // commit count of the repo with the most commits
      "commits": count
    }, 
    {
         // the repo with second most commits
      "repository": "path_with_namespace",
         // commit count of second repo
      "commits": count
    },
    ...
   ],
   "auth": true // true or false indicating whether the call has been made with our without authentication
}
```
# Languages

This endpoint returns the languages used in given projects by distribution in descending order.

## Request

The endpoint accept GET requests. If the payload is empty, it returns top-ranking languages across all accessible repositories irrespective of the distribution value.

### Todo:
If a payload is specified, only the listed repositories be considered when identifying the top-ranking languages. The payload format is as an array of project names, i.e., :

```
{
    "projects": [ "project1", "project2", ... ]
}
```
(project names are the 'name' field in https://git.gvk.idi.ntnu.no/api/v4/projects)
## Response

Calls to the endpoint produce output according to the following JSON schema specification and list the most frequently ranked languages (based on returned ranking for the individual projects):

```
{ 
  "languages": ["Go", "C#", ... ],
  "auth": true // true or false indicating whether the call has been made with our without authentication
}
```
# Issues 
This endpoint will return the name of the users or labels (see parameters) for the attached to the most issues for a given project.  

## Request (Todo)

The endpoint should accept GET requests with the following payload:

```
{
  "project": "project name" // name of the project whose issues are analysed
}
```

The parameter `type` indicates whether users with the most posted issues (value `users`), or the most frequently referred labels (value `labels`) should be returned. If not specified, a corresponding error and status code should be returned.

If not specified in the `auth` parameter, the request should occur without authentication.

## Response

Calls to the endpoint produce output according to the following JSON schema specification:

If `type=users` is queried, the return format should look as follows:

```
{ 
   "users": [
              {
                 "username": "username1",
                 "count": count // count of issues by user
               },
               ... ],
   "auth": true // true or false indicating whether the call has been made with our without authentication
}
```

If `type=labels` is queried, the return format should look as follows:

```
{ 
   "labels": [
               {
                 "label": "label1",
                 "count": count // count of issues with label
               },
               ... ],
   "auth": true // true or false indicating whether the call has been made with our without authentication
}
```
# Status

The status endpoint returns information about availability of invoked service and database connectivity.

## Request

The endpoint accept a GET request with an empty body.

## Response

The response body should look as shown in the following:

```
{
  "gitlab": statusCode, // indicates whether gitlab service is available based on HTTP status code
  "database": statusCode, // similar as above
  "uptime": seconds, // seconds since service start
  "version": "v1"
}
```

# Webhook Registration

The system support the registration of one or more webhooks which is activated upon invocation of any of the endpoints. The registration of webhooks is persistent

## Registration

### Request

The registration occur by sending a POST request with the following payload:

```
{
  "event": "(commits|languages|issues|status)", // any of the types of requests to the service (expressed as regex) 
  "url": "webhook url" // URL invoked upon event trigger
}
```

### Response

The response body contain the newly created resource. 

Example: 
```
{
    "id": 1,
    "event": "issue",
    "url": "example.com",
    "time": "2019-10-27 00:32:45.6611177 +0200 CEST m=+37.440569401"

}
```

## Viewing registered webhooks

A GET request to `/repocheck/v1/webhooks` list all the registered webhooks in the following format:

```
[ 
  { 
    id: 1, 
    event: "eventtype", // event the webhook is subscribed to
    url: "webhook URL", // URL of registered webhook
    time: "timestamp" // timestamp of registration
  },
  ...
] 
```

Individual webhooks are viewable using the respective identifier returned during registration (e.g., `/repocheck/v1/webhooks/1`). The response body is an individual webhook object, e.g., 

```
{ 
    id: 1, 
    event: "eventtype", // event the webhook is subscribed to
    url: "webhook URL", // URL of registered webhook
    time: "timestamp" // timestamp of registration
}
```

## Invocation

Upon being triggered in the endpoint associated with the event the webhook is registered with (i.e., commits, languages, issues, or status), the service send a GET request to the registered webhook URL with the following payload:

```
{
  "event": "eventtype", // event type as per specification for registration
  "params": ["param1", "param2"] // parameters passed along in triggering request
  "time": "timestamp" // timestamp of request
}
```

## Deletion (Todo)

To delete a webhook, send a DELETE request to the URL identifying the webhook (resource id) to be deleted.

# Deviation/Changes from assignment requirements

Specified Languages payload to be
```
{
    "projects": [ "project1", "project2", ... ]
}
```
and that project names assums 'name' field in https://git.gvk.idi.ntnu.no/api/v4/projects

instead of just
```
 [ "project1", "project2", ... ]
```

The invocation gives a string array of the diffrent parameters, instead of a string

Instead of responding with just the ID when you create a webook, I respone with the entire date for the webhook as a JSON format. This is done in purpouse, becuase I belive this makes more sense

In issues:
```
The parameter `type` indicates whether users with the most labels
```
changed to:
The parameter `type` indicates whether users with the most posted issues, as that is how I interpreted it, because the original did not make sense to me

