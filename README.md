[![GoDoc](https://godoc.org/github.com/qzaidi/redamo?status.svg)](https://godoc.org/github.com/qzaidi/redamo)

This is an experiment.

What if we speak the redis protocol, but use dynamo as a backing store.

So on write, we back it up in dynamo. On read, we fetch from dynamodb (or a cache). There is no ttl.

TODO
----

1. Use a TTL cache to make GETs faster, if they haven't changed. This will mean
   we can't run multiple servers, so it would be a config choice.
2. Fix keymapping logic (maybe via config file)
3. Support INCR and INCRBY
