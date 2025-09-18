Read the YAML configuration file to get the number
of servers

```yaml
servers:
	server: name1
		weight: 0
	server: name2
		weight: 50
	server: name3
		weight: 100
```

Add these servers with there weight information in a minimum heap,
that orders itself according to the number of requests that a single
server is handling.

Regular health checks of the servers, if a server is down move it
to the `down` queue and checking it on the intervals to see if they are
back up or not.
