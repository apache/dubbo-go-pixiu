# benchmark

This dir isn't guaranteed to be run successfully

benchmark for dubbo-go-pixiu

1.根据开发手册进行部署环境
https://dubbogoproxy.yuque.com/docs/share/e7aca6d0-1957-4f85-af4e-306104c2bbe4?#

以下测试基于http->dubbo进行协议转换

2. 配置lua脚本 (可以使用script/wrk_post.lua)


```
local counter = 1
local threads  = {}

--local json = require("json")

local req  = {

    --group = "dubbo-test",
    --version = "1.0.0",
    --method= "GetUserByName",
    types = "types",
    values = "tcccccccc"
}

function setup(thread)
-- 给每个线程设置一个id参数
thread:set("id",counter)
table.insert(threads,thread)
counter = counter + 1

end

function init(args)
-- 初始化两个参数 每个线程都有独立的 requests, response 参数
requests = 0
response = 0
-- 打印线程被创建的消息，打印完后，线程正式开始运行
local msg = " thread %d created"
print(msg.format(counter))
end

function request()

    wrk.headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
    wrk.headers["User-Agent"] = "wrk"
    wrk.headers["Connection"] = "keep-alive"
    wrk.body = "{\"types\":\"string\",\"values\":\"tc\" }"

    print("本次请求包体为: %s",wrk.body)
    return wrk.format("POST",wrk.path,wrk.headers,wrk.body)
end

function response(status,headers,body)
if status ~= 200 then
print(body)
return
end

    -- 打印响应体
    local resp = json.decode(body)
    print(json.encode(body)..'-->'..body)

end


--function done(summary,latency,requests)
--    print("99 latency:"..latency:percentile(99.0))
--    for index,thread in ipairs(threads) do
--        local id = thread:get("id")
--        local requests = thread:get("requests")
--        local responses = thread:get("responses")
--        -- 打印每个线程发起了多少请求 得到了多少响应
--        print(wrk:format(id,requests,responses))
--    end
--    print(latency)
--end
```





- 性能测试部分: (consumer请求是否经过pixiu网关)


不经过pixiu网关的测试方式:  dubbogo文件夹下的pixiu_test方法
经过pixu网关的测试方式:  目前仅支持http->dubbo的方式，执行pixiu_test




- 压力测试部分:
分别在不同运行时间，不同连接数，不同线程数，不同包体的情况下做测试

● 对照组1: 在不同连接数下的性能测试

```
case1:
运行时间: 3s   连接数: 10  线程数: 1 包体: 33字节
Thread Stats   Avg      Stdev     Max   +/- Stdev
Latency     9.41ms   26.34ms 214.80ms   95.53%
Req/Sec     2.76k   684.49     3.55k    79.31%
Latency Distribution
50%    2.88ms
75%    6.28ms
90%   11.55ms
99%  166.03ms
8038 requests in 3.01s, 1.53MB read
Requests/sec:   2671.24
Transfer/sec:    521.73KB
结果： 每个线程的平均延迟在166.03ms, 每秒请求数 2.76k个/s，延迟分布在t99的 主要延迟时长是166.03ms

接口对总体平均每秒数处理请求数为 2671.24QPS
```

```
case2:

运行时间: 3s   连接数: 20  线程数: 1 包体: 33字节
Thread Stats   Avg      Stdev     Max   +/- Stdev
(平均值)    (标准差)  (最大值)    (正负一个标准差所占比例)
Latency   285.15ms  152.23us 285.28ms   66.67%
(延迟)
Req/Sec     9.00      0.00     9.00    100.00%
(每秒请求数 tps)
Latency Distribution
(延迟分布)
50%  285.18ms
75%  285.28ms
90%  285.28ms
99%  285.28ms
3 requests in 3.09s, 600.00B read  (3.09s内处理了3个请求，耗费流量600.00B)
Requests/sec:      0.97     (QPS 0.97 即平均每秒数处理请求数0.97  )
Transfer/sec:     194.20B   (平均每秒流量194.20B )
```

结果： 
```每个线程的平均延迟在285.15ms, 每秒请求数 9个/s，延迟分布在t99的 主要延迟时长是285.28ms```

接口对总体请求的平均每秒数处理请求数为 ```0.97QPS```


对照结论: 在连接数增加的情况下平均延迟增大,每秒接受的请求数变少，平均处理时长降低

● 对照组2: 在不同线程下的性能测试:

case1:

运行时间: 3s   连接数: 20  线程数: 1 包体: 33字节

```
Thread Stats   Avg      Stdev     Max   +/- Stdev
(平均值)    (标准差)  (最大值)    (正负一个标准差所占比例)
Latency   285.15ms  152.23us 285.28ms   66.67%
(延迟)
Req/Sec     9.00      0.00     9.00    100.00%
(每秒请求数 tps)
Latency Distribution
(延迟分布)
50%  285.18ms
75%  285.28ms
90%  285.28ms
99%  285.28ms
3 requests in 3.09s, 600.00B read  (3.09s内处理了3个请求，耗费流量600.00B)
Requests/sec:      0.97     (QPS 0.97 即平均每秒数处理请求数0.97  )
Transfer/sec:     194.20B   (平均每秒流量194.20B )
```

结果： 
```每个线程的平均延迟在285.15ms, 每秒请求数 9个/s，延迟分布在t99的 主要延迟时长是285.28ms```

接口对总体平均每秒数处理请求数为 ```0.97QPS```



case2:
运行时间: 3s   连接数: 20  线程数: 5 包体: 33字节 (不稳定 pixiu网关很长一段时间连不到注册中心)


```
Thread Stats   Avg     Stdev     Max   +/- Stdev
(平均值)  (标准差)  (最大值)   (正负一个标准差所占比例)
Latency     6.81ms    5.90ms  58.48ms   83.17%
(延迟)
Req/Sec   582.34    127.02     0.95k    74.12%
(每秒请求数 tps)
Latency Distribution
(延迟分布)
50%    4.91ms
75%    9.68ms
90%   14.63ms
99%   27.15ms
4941 requests in 3.05s, 0.94MB read  (3.05s内处理了4941个请求，耗费流量0.94MB)
Requests/sec:   1621.19 (QPS 163.76即平均每秒数处理请求数163.76 )
Transfer/sec:    316.64KB (平均每秒流量316.64KB)
```
结果： 
```每个线程的平均延迟在6.81ms,  每秒平均请求数 582.34个/s，延迟分布在t99的 主要延迟时长是27.15ms```

接口对总体请求的平均每秒数处理请求数为 ```163.76QPS```

对照结论：在仅提升处理线程数的情况下，每个线程的平均延迟降低，每秒平均请求数提升.t99的延迟时长降低

● 对照组3: 在不同运行时间的情况下测试

case1:
运行时间: 3s   连接数: 20  线程数: 5 包体: 33字节 (不稳定 pixiu网关很长一段时间连不到注册中心)



Thread Stats   Avg     Stdev     Max   +/- Stdev
(平均值)  (标准差)  (最大值)   (正负一个标准差所占比例)
Latency     6.81ms    5.90ms  58.48ms   83.17%
(延迟)
Req/Sec   582.34    127.02     0.95k    74.12%
(每秒请求数 tps)
Latency Distribution
(延迟分布)
50%    4.91ms
75%    9.68ms
90%   14.63ms
99%   27.15ms
4941 requests in 3.05s, 0.94MB read  (3.05s内处理了4941个请求，耗费流量0.94MB)
Requests/sec:   1621.19 (QPS 163.76即平均每秒数处理请求数163.76 )
Transfer/sec:    316.64KB (平均每秒流量316.64KB)

结果： 
```每个线程的平均延迟在6.81ms,  每秒平均请求数 582.34个/s，延迟分布在t99的 主要延迟时长是27.15ms```

接口对总体请求的平均每秒数处理请求数为 ```163.76QPS```


case2:
运行时间: 60s   连接数: 20  线程数: 5 包体: 33字节 (不稳定 pixiu网关很长一段时间连不到注册中心)

wrk -d60s -c20 -t5 --latency  -s test.lua "http://localhost:8883/api/v1/test-dubbo/UserProvider/com.dubbogo.pixiu.UserService?group=dubbo-test&version=1.0.0&method=GetUserByName"


Thread Stats   Avg      Stdev     Max   +/- Stdev
(平均值)   (标准差)  (最大值)   (正负一个标准差所占比例)
Latency     6.29ms    5.38ms  53.31ms   81.72%
(延迟)
Req/Sec   618.23    142.44     0.91k    70.00%
(每秒请求数 tps)
Latency Distribution
(延迟分布)
50%    4.45ms
75%    8.94ms
90%   13.46ms
99%   23.92ms
6172 requests in 1.00m, 1.18MB read (一分钟内处理了6172个请求，耗费流量1.18MB)
Requests/sec:    102.74  (QPS 102.74即平均每秒数处理请求数102.74 )
Transfer/sec:     20.07KB  (平均每秒流量20.07KB )

结果： 
```每个线程的平均延迟在6.29ms,  每秒平均请求数 618.23个/s，延迟分布在t99的 主要延迟时长是23.92ms```

接口对总体请求的平均每秒数处理请求数为 ```102.74QPS```


对照结论：在不同运行时长的情况下，运行时间越久，平均处理的请求数越多，接口对总体请求的平均处理时长越短，其余参数没有变化



