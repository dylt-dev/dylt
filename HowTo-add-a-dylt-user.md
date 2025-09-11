### How to add a dylt whitebelt

Adding a whitebelt is pretty straightforward
- Get a public key for a user
- Confirm public key is valid
- Create an incus VM
    - @note don't worry about resource contraints for now
- Create a proxy on a port to funnel request to 22 (unclear if I can do this by user except with ssh_config ForceCommand
- Add the user
- Add the public key
- Try the login
- Setup DNS for the user
- Setup etcd config info

Next steps
- Cloudflare API REST for command-line DNS add
- incus profile for ssh proxy
- add install-incus function to /opt/daylight.sh (using zabbly)

Big swings
- Add SRV-lookup to dylt ssh
- ForceCommand for ssh proxying