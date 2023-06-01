# Module - "net"

```golang
net := import("net")
```

## Functions

### is_ip
```golang
is_ip(input string) => bool
```
Checks if the provided input is an IP address. Runs the `ParseIP()` Go function from the `net` module.

#### Example
```golang
fmt := import("fmt")
net := import("net")

fmt.println(net.is_ip("127.0.0.1"))
fmt.println(net.is_ip("example.com"))
fmt.println(net.is_ip("2001:db8::68"))
```
```
Output:
true
false
true
```
