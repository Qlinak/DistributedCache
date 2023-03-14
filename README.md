# DistributedCache

A simplified version of memCached   <br>  

An in-memory distributed caching solution supporting resource control, concurrent read-write, smart replacement, inter-process communication, etc.<br> 

Apply Least Recently Used (LRU) algorithm for cache replacement in a single node<br>

Use mutex lock to avoid inter-process concurrency issue and avoid cache penetration<br>

Apply consistent hash to pick peers upon single node cache miss, load balancing amongst multiple nodes and avoid cache avalanche upon change in cache node cardinality<br> 

Ensure concurrent requests on the same key only invoke once<br>

[pending to be added] Use Protobuf as RPC data exchange type to decrease data transmission size, ensuring high performance
