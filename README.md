http test tool
==========

# install
```
go get -u github.com/hidu/curl-flow
```


# params
```
  -c int
        concurrency Number of multiple requests to make (default 1)
  -n int
        Number of each requests to perform (default 1)
  -t uint
        Timeout of request (default 10 second) 
  -ui
        use termui
  -url string
       request url
```

# usage

## 1. stream:
```
cat _example/requests.txt|curl-flow  -c 100
```

stream contents format:
```
line 1: request 1 ( JSON string, not indented )
line 2: request 2 ( JSON string, not indented )
line 2: request 3 ( JSON string, not indented )
```
see the [_example/requests.txt](./_example/requests.txt)


one request (JSON string, after indented)：
```
{
    "url": "http://127.0.0.1:8088/test.php?i=0",
    "method": "post",
    "header": {
        "Content-Type": "application/x-www-form-urlencoded",
        "head-a": "head-v"
    },
    "payload": "id=0&now=1611202835"
}
```

real time stream with script:
```
php _example/gen_requests.php|curl-flow
```

## 2. by args:
```
curl-flow -url "http://example.com/xxx" -c 10 -n 10000
```
send `10000` GET requests, concurrency is `10` .

