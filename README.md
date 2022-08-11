# Go HTTP/3 Proxy

A minimal simple HTTP/3 Proxy written in Go.
The proxy receives requests in HTTP 1.1 but proxy them using HTTP/3.
This proxy is optimized to be used as a way to expose a Cloud Run endpoint on different region using the minimum possible round trips per request.   


## Environment Variables:

- `PORT` the port on which the proxy listens (default `8080`). 
- `TARGET_HOST` the target host to proxy to (default `172.17.0.1`). 

## Usage

```bash
docker run -e TARGET_HOST=goolge.com -e PORT=80 -p 8080:80 http3-proxy
```

## Build from source

```bash
docker build -t http3-proxy .
```
