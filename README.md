Gethashes
=========

This is something I made for me to learn Golang and to get a feel for how to make Goland applications and combine it with containers. Currently this is made so that it has to be used with Docker. But I want to make this work in a rootless Podman setup at a later time.

This will setup two containers, the first being a pure nginx-reverse proxy and the second being the gethashes app that is set up to only communicate with the nginx proxy.

Requirements
------------

Docker
Docker-compose

The nginx reverse proxy virtual host is set up to listen for gethashes.localhost, meaning your /etc/hosts file or your DNS server must be configured to give the correct IP address for gethashes.localhost
