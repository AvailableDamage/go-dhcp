# go-dhcp
A simple dhcpv4 server written in Golang.

**Work in Progress**

## Build
```
$ git clone https://github.com/AvailableDamage/go-dhcp
$ cd go-dhcp
$ make build
```
The executable file is placed in ./bin/.
The executable needs CAP_NET_RAW and CAP_NET_BIND_SERVICE, because it uses a privileged Port (67/UDP) and uses raw sockets (to manually set the clients macaddress in the ethernet frame).
```
$ setcap 'cap_net_bind_service=+ep' /path/to/go-dhcp
$ setcap 'cap_net_raw=pe' /path/to/go-dhcp
```
The directories /etc/goDHCP and /var/run/goDHCP have to be created.
In /etc/goDHCP directory is the configuration stored.
In /var/run/goDHCP is a csv-file for the leases.

## Configuration 

Example:
```
{
  "Interface": "lo",
  "Nameserver": "10.0.200.250",
  "Gateway": "192.168.123.1",
  "LeaseTime": 3600,
  "Server ID": "10.0.20.13",
  "Pools": [
    {
      "Name": "pool1",
      "Interface": "enp34s0",
      "StartIP": "10.0.202.20",
      "EndIP": "10.0.202.90",
      "Options": {
        "Subnetmask": "255.255.255.0",
        "Nameserver": "10.0.200.250",
        "Gateway": "10.0.202.1",
        "Broadcast Address": "10.0.202.255",
        "Lease Time": "600",
        "Domain": "example.com"
      }
    },
    {
      "Name": "pool2",
      "Interface": "enp42s0",
      "StartIP": "10.0.20.50",
      "EndIP": "10.0.20.250",
      "Options": {
        "Subnetmask": "255.255.255.0",
        "Nameserver": "10.0.200.250",
        "Gateway": "10.0.20.1",
        "Broadcast Address": "10.0.20.255",
        "Lease Time": "3600",
        "Domain": "example.com"
      }
    }
  ]
}
```
