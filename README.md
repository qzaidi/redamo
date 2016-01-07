This is an experiment

What if we speak the redis protocol, but use groupcache internally, and use dynamo as 
a backing store.

So on write, we update groupcache and back it up in dynamo. On read, we fetch from groupcache, if not found
from dynamodb. 


