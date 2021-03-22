# dns

A simple DNS server in pure go. Currently using the DNS implementation from: `golang.org/x/net/dns/dnsmessage`.

This project is able to answer a `dig` query.

```
go run *.go
```

```
dig @localhost -p 1053 www.github.com
```
