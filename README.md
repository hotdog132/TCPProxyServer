# TCPProxyServer
TCP Server to communicate between external API and client using Go.

## To-Do
- [x] Add TCP server and mock external API server prototypes
- [x] Add request limiter mechanism
- [x] Add test cases
- [x] Handle the case when external API is unreachable or unavailable
- [x] Add http endpoint for showing current connection count, current request rate, processed request count, remaining jobs

## Addition
- [ ] Handle case when the external api is available but not response
- [ ] Clear jobs in the pending queue when peer disconnect
- [ ] Save jobs to DB to reduce memory usage (Refactoring using interface)

## Setup & launch

### Initialization
```
go mod init
```

### Launch mock external API
```
go run externalAPI/main.go
```

Default config of mock external API:
- Port:8888
- It takes 2 seconds to handle the request from TCP server


### Launch TCP server
```
go run tcpServer/main.go
```

Default config of TCP server:
- Port: 8000
- 30 requests/second
- Peer connection timeout: 120 seconds


### Http endpoint:
- http://localhost:7000/statistics
- Refresh the page every 2 seconds


### Use clients to test TCP server
```
nc localhost 8000
```

Type any query strings and send it to TCP server.
Type 'quit' or force stop using cmd to disconnect from TCP server.

Open http://localhost:7000/statistics to check current status of TCP server.

