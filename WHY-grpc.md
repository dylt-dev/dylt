## Why `grpc`?

Why have a `grpc` server at all?

Specifically, why have a grpc server whose primary responsibility is to play with etcd, and then return to the client? Couldn't the client just connect directly to etcd and do what they want?

It's a valid question. Either the client uses a CLI that executes business logic against a remote etcd store ... or the client uses a CLI that uses grpc to connect to a server to do the same thing.
This is probably unclear, so let's review the players. We have ...
- A human user
- a CLI
- a chunk of business logic
- a remote etcd server

We can do this ...

human user - CLI - bizlogic <=========== Internet (via etcd client) ============> etcd

or we can do this ...

human user - CLI <============= Internet (via grpc client) ===============> bizlogic etcd

What's the difference?

### aesthetics

Some will think the first approach is cleaner, more elegant, or just right.
And some will think the second approach is cleaner, more elegant or just right.
Those camps will argue forever and make no progress.
Leave them to it.

### technology aka objective truth

Let's avoid the neurological temptation of making this a question of right vs wrong
Let's instead talk about circumstantial realities, that might make one approach better vs the other

### server availability / maintenance

Let's say our resource server is mature, rock solid, actively supported ... total bedrock.
And let's say a grpc server is something we'd be making from scratch, possibly because we're interested in the technology.
Is it really better to create a whole grpc service, to interact with the resource service, just so our CLI can use gRPC?
Or should our CLI just hit the RS directly. Remove aesthetics and it's hard to argue against this.

### Backend simplicity / homogeneity

On the other hand what if the backend has to connect to a bunch of different services, possibly hit some local files, and basically has a hard job.
In that case you might want your CLI to have an easy job, and the back end to have a hard job. Among other things, it's very likely that it's easier for the back end to avoid or deal with problems, than wherever your workstation happens to be.

### Where is the server, anyway

On the other hand, maybe the 'server' is on your Intranet or otherwise is almost as good as local. In that case, it might make a lot more sense to use the server. The server can know stuff like what load-balancer to switch over to, who to email in case of a persistent issue, etc. That's a lot of config that would need to live on the client. If the client only has to talk to the server, all the better. Of course all that config _could_ be kept in the server, if for some reason that was desirable. But just because something _can_ happen on the server doesn't mean that it should.

### Are we getting somewhere?

It's sounding like for trivial examples, it doesn't matter what the CLI does -- if the CLI does the work directly or delegates it to a remote piece.
But there are scenarios, many of which correlate with scaling, where it makes sense for a CLI to delegate the hard stuff to a remote service, via gRPC.
Or put another way, when the CLI just needs to talk directly to One True Resource Server, it should go ahead and do that. And the role of a gRPC service is to be that One True Resource Server.


