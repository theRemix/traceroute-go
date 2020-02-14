<h1 align="center">🚧 Traceroute go</h1 >

<p align="center">
  <strong>Traceroute program implemented in go</strong>
</p>

## Quick Start

```sh
go build -o traceroute-go main.go
sudo setcap cap_net_raw=ep traceroute-go
./traceroute-go
```

## Traceroute

See [./main.go](./main.go)


## Example

_unfinished_

```
2020/02/13 17:47:50 sending 1
2020/02/13 17:47:50 type time exceeded
2020/02/13 17:47:50 hop ttl=254 src=10.33.104.1 dst=10.33.109.176 ifindex=2
10.33.104.1
2020/02/13 17:47:50 sending 2
2020/02/13 17:47:50 type time exceeded
2020/02/13 17:47:50 hop ttl=254 src=157.130.196.213 dst=10.33.109.176 ifindex=2
157.130.196.213
2020/02/13 17:47:50 sending 3
2020/02/13 17:47:50 type time exceeded
2020/02/13 17:47:50 hop ttl=250 src=140.222.225.243 dst=10.33.109.176 ifindex=2
140.222.225.243
2020/02/13 17:47:50 sending 4
2020/02/13 17:47:50 type time exceeded
2020/02/13 17:47:50 hop ttl=250 src=129.250.9.249 dst=10.33.109.176 ifindex=2
129.250.9.249
2020/02/13 17:47:50 sending 5
2020/02/13 17:47:50 type time exceeded
2020/02/13 17:47:50 hop ttl=250 src=129.250.2.2 dst=10.33.109.176 ifindex=2
129.250.2.2
2020/02/13 17:47:50 sending 6
2020/02/13 17:47:50 type time exceeded
2020/02/13 17:47:50 hop ttl=249 src=131.103.117.82 dst=10.33.109.176 ifindex=2
131.103.117.82
2020/02/13 17:47:50 sending 7
timeout 7
```
