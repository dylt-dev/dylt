## HowTo - Create a Two Node `etcd` cluster

In a sense there is no such thing as a 2 node `etcd` cluster. etcd uses a consensus protocol. Consensus protocols tend not to work with 2 nodes. If 2 nodes disagree, it's hard to break the tie. For this reason and others, etcd only works with 2 nodes if both nodes are active. This is according to etcd's own docs. And if both nodes have to be active, you don't really have any benefits of a cluster. You just have 2 nodes to worry about instead of 1.

But can't 2 nodes be better than one? They can. You can run a 1 node cluster, and plan on failing over to the 2nd node if necessary. There are a few ways to achieve that in etcd. We're going to look at one way: Create a 1 node etcd-cluster on Host #1, and startup etcd on Host #2 in Learner Mode.

### What's Learner Mode?

Details are here -> https://etcd.io/docs/v3.3/learning/learner/

Learner Mode was created to let a new node gently join an existing etcd cluster, without overwhelming the cluster and threatening its stability. New nodes are needy - they need to be updated with the entire state of the cluster, including history, all while the cluster remains online serving users. In a Learner Mode, a new node does not vote, cannot serve as leader, and is not part of quorum. It just receives updates.

In typical usage, an etcd node that is started in Learner Mode will be restarted in normal mode once it has been synched to the cluster. But this isn't necessary. A node can be left in Learner Mode permanently instead. If something were to happen to the 'real' node, the new joiner can be started in normal mode and take over as the new cluster.

### Steps

1. Get IPs / hostnames and SSH keys for the main node and backup node.
2. Create a startup script to run etcd on the main node
3. scp the startup script to the main node
4. Run the startup script manually on the main node
5. Test that the startup script works
6. Create a systemd service for etcd on the main node
7. Create a startup script to run etcd on the backup node
8. Run the etcd script manually on the backup node
9. Test the startup script
10. use `etcd member add --learn` to add the new node in Learner Mode
11. Create a systemd service for the backup node
12. Test the cluster
13. Explicitly test the backup node

### Failing over to Backup Node

1. Shut down main node
2. Remove backup node from cluster
3. Re add backup node to cluster, without Learner Mode
4. Test former backup node as new main node