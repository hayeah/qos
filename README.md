# QOS

Monitor remote host's quality of service, and record data in InfluxDB.

Starts

#

Starts a server counterpart:

```
qos-server port
```

Start qos monitoring on client:

```
qos host port
```

The client will measure:

1. A steady transfer rate.
2. Round trip time.
