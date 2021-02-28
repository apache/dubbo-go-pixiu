# Filter

Filter is the composition of filter chain, make use more control.

## Default filter

### Remote filter

call downstream service, only for call, not process the response. 

#### field

> level: mockLevel 
 
- 0:open Global mock is open, api need config `mock=true` in `api_config.yaml` will mock response. If some api need mock, you need use this. 
- 1:close Not mock anyone.
- 2:all Global mock setting, all api auto mock.

result
```json
{
    "message": "success"
}
```

#### pixiu log 
```bash
2020-11-17T11:31:05.718+0800    DEBUG   remote/call.go:92       [dubbo-go-pixiu] client call resp:map[age:88 iD:3213 name:tiecheng time:<nil>]
```

### Timeout filter

api timeout control, independent config for each interface.

#### Basic response

[reference](../user/response.md#timeout)

### Response filter

Response result or err.

#### Common result

[reference](../sample/dubbo-body.md)

### Host filter

 