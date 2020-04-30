# sr-cache
A distributed cache system

# 2012-04-28
## Features:
Basic LRU cache replacement policy(LRU-1).
Data partition and peer picking according to consistent hash.
HTTP server.

## Considered Todos:
A tool for deploying the system on docker or k8s easily.
GRPC server.
LRU-k cache replacement policy.
Expiration time of cache.
Memory limit.
Cache and Peer visit statistics.
Prevent cache penetration, cache breakdown and cache avalanche.

