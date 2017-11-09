http test tool
==========

# install
```
go get -u github.com/hidu/curl-flow
```

# useage

```
php _example/gen_requests.php|curl-flow
cat requests.txt|curl-flow  -c 100
```

# params
```
  -c int
        concurrency Number of multiple requests to make (default 1)
  -n int
        Number of requests to perform (default 1)
  -t uint
        Timeout of request (default 10)
  -ui
        use termui
  -url string
        test url
```