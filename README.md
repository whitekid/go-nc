# Simple Bastian Jumper

## Setup

build

    make linux

copy to host

    scp bin/gn <bastian>:.

## ssh config

    Host <hostname>
      HostName <hostname>
      User <login>
      ServerAliveInterval 20
      ServerAliveCountMax 10
      GSSAPIAuthentication yes
      GSSAPIDelegateCredentials yes
      ProxyCommand ssh <bastian-host> ./gn %h %p
