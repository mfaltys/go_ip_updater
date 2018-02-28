# Go_ip_updater
[![Build Status (Travis)](https://travis-ci.org/unixvoid/go_ip_updater.svg)](https://travis-ci.org/unixvoid/go_ip_updater)  
Go_ip_updater is a tool written in golang designed to monitor and update changing IP's.
This tool is designed to combat an ISP's dynamic IP assignment. Every time
an exteral IP changes go_ip_updater can update route53 almost immediately.


## Running go_ip_updater
There are 2 main ways to run go_ip_updater:  

1. **From Source**: After obtaining the binary file you can simply run it with a few
parameters (shown below).  
`make dependencies`  
`make run`

4. **ACI/rkt**: We have publically hosted ACI images available for use, check them
out [here](https://cryo.unixvoid.com/bin/rkt/go_ip_updater/) or build for yourself
with:  
`make dependencies`  
`make build_aci`  
If you are not interested in building the image yourself you can grab it with:  
`rkt fetch unixvoid.com/go_ip_updater`.  
To run go_ip_updater we recommend using systemd.  You can find our service file
[here](https://github.com/unixvoid/go_ip_updater/blob/master/deps/go_ip_updater.service)  
To set up to run all the time copy `go_ip_updater.service` to `/etc/systemd/system/`
and issue the following commands:  
`sudo systemctl daemon-reload` to get the service file loaded  
`sudo systemctl start go_ip_updater` to start the service  
`sudo systemctl enable go_ip_updater` to run the service on boot  
before getting started with systemd make sure to follow our configuration guide below!


## Configuration
go_ip_updater can take a couple command line arguments to get going quickly, but uses
a configuration file for the rest.  

**Command line arguments**:
  `-loglevel`: will set the loglevel for the updater (defaults to info)  
    possible choices are `debug`, `info`, `error`.  
  `-config`: path the the configuration (defaults to `goip.list`)

**Config file**:  
```
# go_ip_updater config
configKey:     <key>
configSecret:  <secret>
configZoneId:  <zone_id>
configTTL:     30
configURL:     http://checkip.amazonaws.com

# domain list
testa.unixvoid.com
testb.unixvoid.com
testc.unixvoid.com

```
The configuration file is pretty straightforward but will need to be filled out
with specific AWS details.  `configKey` and `configSecret` will come straight from
user invoking the service while `configZoneId` is the zone ID for your route53 domain.  
`configTTL` is not only the TTL for the configured A records, but will be how often
(in seconds) that go_ip_updater polls to see if an IP has changed.  
`configURL` is the URL that will be used to check the external IP. This doesn't have
to be set to amazons service, but we find it to have the best response times.  

The domain list is the list of A names that go_ip_updater will update, there is no
limit to the ammount of domain names to control.
