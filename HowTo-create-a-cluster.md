# How to Create a Cluster

## _'3 VMs and the Truth'_

A cluster is ...
- A collection of Hosts that share config through etcd
- 3 Host, each in their own provider: OVH, IONOS, Liquid Web
- Each Host runs Incus, for VMs
- Each host has a cloud config, which can either be usws on startup or post-startup since that might be easier
- Each cluster can have multiple admins via RBAC
- Each user can have multiple clusters through multitenancy
- ... and each has `dylt`, of course!

One possible approach ...
- Install dylt early in the cloud config process
- Use dylt to install the cluster
- Use a cluster script, on cloud-init-startup, to configure the host

This would make restarting a host a hard option: it would include installing all a user's services from the cluster. But this might be the sort of hard option that clarifies everything.

There's another factor here. Services are installed in VMs. VMs have their own lifecycle and configuration management through Incus. Is there value in supporting Incus features here? Does restart really need to kill the Incus VM, start a new VM, and restart all services? This seems like it would make more sense as an option, rather than as unavoidable behavior.

Also ... is tunneling into a VM even all that helpful, vs running services right on the Host? Answer: yes, if we would like to support multiple projects or multitenancy via Incus, which might make sense. So I need to do more Incus VM tunneling.

? When dylt invokes an init script, does that script invoke dylt? I guess it probably should.
(actually dylt doesn't have to invoke it initially. just download & install it)

> Note: this does not accomodate the idea of a "last-mile" on-prem replacement, where I'm retiring an on-prem VM in favor of a hosted Incus VM, which can be stopped, snapshooted, etc. That's a whole other use case, with a whole other type of customer. It might be pretty simple to get up and going, which is worth thinking about, but generalizing it with the service stuff could be tough.Then again ... maybe I have services, and VMs. And a VM is something a client is hands-on with, with zero IaaS except for what Incus provides. And services make things easier. Just no hybrid stuff. Unless the customer wants hybrid, which maybe is ok. Oy. This is exactly the mess I wanted to avoid.




