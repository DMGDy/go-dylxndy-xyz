# go dylxndy.xyz
A dedicated HTTPS webserver for my personal website [dylxndy.xyz](https://dylxndy.xyz) 
to replace [Apache/httpd](https://httpd.apache.org/). This is also a rewite of
[https-server](https://github.com/DMGDy/https-server) I did in C and attempted to
implement the Linux [epoll](https://man7.org/linux/man-pages/man7/epoll.7.html) API. 


Go's coroutines (goroutines) make that sort of issue trivial. Go's standard library make
a lot of boiler plate for an https server simple, especially with is [crypto/tls](https://pkg.go.dev/crypto/tls) package.


# Building
This doesn't depend on anything else so:
```
go build -o dylxndy_server main.go
```
Note: this assumes you have the propper X509 Key Pairs in locally in `cert/` as
`cert/cert.pem` and `cert/key.pem`
