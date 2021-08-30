# Api

Introduction to API model, recommended reading before customizing API for api_config.yaml.

## Api Gateway

API is the key feature of the dubbo-go-pixiu. With this feature, you can expose your dubbo service as an HTTP service.

### Configuration

Sample:
``` yaml
name: api name
description: api description
resources:
  - path: '/'
    type: restful
    description: resource documentation
    methods:
      - httpVerb: GET
        enable: true
        inboundRequest:
          requestType: http
          queryStrings:
            - name: id
              required: false
        integrationRequest:
          requestType: dubbo
          mappingParams:
            - name: queryStrings.id
              mapTo: 1
          applicationName: "BDTService"
          interface: "com.ikurento.user.UserProvider"
          method: "GetUser"
          clusterName: "test_dubbo"

definitions:
  - name: modelDefinition
    schema: >-
      {
        "type" : "object",
        "properties" : {
          "id" : {
            "type" : "integer"
          },
          "type" : {
            "type" : "string"
          },
          "price" : {
            "type" : "number"
          }
        }
      }
```
name
:
> The name of the Gateway.


description
:
> The description of the Gateway

resources
:
> Defines the API resources.

path
:
> Defines the path of the API. Parameterized path is acceptable. For instance, you can specify /users/:id, the :id will be considered as parameter in query string. You can also specify /users/:id/transactions/:transactionId, which :id, and :transactionId will be used as query parameters.

type
:
> The type of the API, possible values: restful. Type dubbo will be support later to provide dubbo-to-http gateway function.

filters
:
> The filters in resource level. Use comma to separate the filters.

methods
:
> Defines the methods within the resource. 

httpVerb
:
> The http method, accept GET/POST/PUT/DELETE/OPTIONS/PATCH/HEAD/HEAD/ANY(Not supported yet)

enable
:
> Defines the API is online/offline. true -> online, false -> offline. When the API is offline, it returns 406.

inboundRequest
:
> Defines the way to expose the API.

requestType
: 
> the inbound request type. Values: http. Type dubbo will be support later to provide dubbo-to-http gateway function.

queryString
:
> the query string defines the parameters that the API exposes. The parameters you defined in the path will be merged into the configuration here.

    name: the name of the query parameter
    required: if the parameter is required.

requestBody
:
> The parameters in request body. It works with the definition config. 
    
    definitionName: The name of the defined definition maps to this request body.

integrationRequest
:
> Defines the detail of the actual backend service.

requestType
:
> The request type in integrationRequest defines the type of the backend service. 'dubbo' is the value we support now. 'http' will be support later with the dubbo-to-http pixiu feature.

mappingParams
:
> It is the field that describe how to map the inboundRequest parameters to backend service parameters.

    name: The name of the inboundRequest parameter name. use queryStrings.* for the parameters defines in queryStrings, 
    mapTo: The index of the target parameters. Starts from 0.

application
:
> The name of the target dubbo application.

interface
:
> The interface represents the interface in dubbo.

method
:
> The method in integrationRequest represents the method in dubbo

clusterName
:
> The clusterName defines which dubbo cluster to call(Will release later)

definitions
:
> the definitions of the models. They will be used as request body.

    name: The name of the definitions.
    schema: The detail definition of the definition. Use json.schema.

