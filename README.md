Load balancers are crucial in modern software development. If you've ever wondered how requests are distributed across multiple servers, or why certain websites feel faster even during heavy traffic, the answer often lies in efficient load balancing.

### What is a Load Balancer?
A load balancer is a system that distributes incoming network traffic across multiple servers. It ensures that no single server bears too much load, preventing bottlenecks and improving the overall user experience. Load balancing approach also ensure that if one server fails, then the traffic can be automatically re-routed to another available server, thus reducing the impact of the failure and increasing availability.

### Why do we use Load Balancers?
- High availability: By distributing traffic, load balancers ensure that even if one server fails, traffic can be routed to other healthy servers, making the application more resilient.
- Scalability: Load balancers allow you to scale your system horizontally by adding more servers as traffic increases.
- Efficiency: It maximizes resource utilization by ensuring all servers share the workload equally.

### Load balancing algorithms
There are different algorithms and strategies to distribute the traffic:
- Round Robin: One of the simplest methods available. It distributes requests sequentially among the available servers. Once it reaches the last server, it starts again from the beginning.
- Weighted Round Robin: Similar to round robin algorithm except each server is assigned some fixed numerical weighting. This given weight is used to determine the server for routing traffic.
- Least Connections: Routes traffic to the server with the least active connections.
- IP Hashing: Select the server based on the client's IP address.

---------------
In this repo, I'll focus on implementing a Round Robin load balancer.

### What is a Round Robin algorithm?
A round robin algorithm sends each incoming request to the next available server in a circular manner. If server A handles the first request, server B will handle the second, and server C will handle the third. Once all servers have received a request, it starts again from server A.