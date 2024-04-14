# stats - statistics plugin for CoreDNS

stats is a CoreDNS plugin which can be used to store request statics into a database backend. The plugin is designed to be used with a SQL database backend, but can be extended to support other databases as well.

Currently the plugin supports PostgreSQL.


## Usage
### CoreDNS Binary

You can include _stats_ in your CoreDNS just as you would include any other CoreDNS plugin.

```bash
# Clone coredns to a local location
$ git clone git@github.com:coredns/coredns.git ~/dns-server/coredns

# Clone stats plugin to a close location
$ git clone git@github.com:GustavoKatel/stats.git ~/dns-server/stats

# Symlink stats location into coredns/plugin/stats
$ cd ~/dns-server/coredns/plugin
$ ln -s ../stats ./stats

# Update plugin.cfg and put the line "stats:stats" before the "forward:forward" line

# Build CoreDNS
$ cd ~/dns-server/coredns
$ go generate
$ make
$ ./coredns -conf Corefile
```

### Corefile

The stats directive inside Corefile requires one argument, which is the connection string of the backend. The connection string should be in the format of `driver://user:password@host:port/database`.

Optionally, you can specify a block with the following options:

- `workers` [integer] The number of workers to use for processing the requests. The default value is 3.
- `queryTimeout` [Go duration] The timeout for the database queries. The default value is 5 seconds.
- `statsPrefix` [string] The prefix to be used for the statistics table. The default value is `coredns`.


```
stats postgresql://user:pass@postgres.postgres.svc.cluster.local:5432/coredns
```

Or

```
stats postgresql://user:pass@postgres.postgres.svc.cluster.local:5432/coredns {
    workers 3
    queryTimeout 5s
    statsPrefix coredns
}
```

This is a sample Corefile including the stats directive. It will store statistics in a PostgreSQL database and forward the requests to Cloudflare DNS.

```Corefile
.:53 {
    metadata
    prometheus 0.0.0.0:9153

    bind 0.0.0.0

    stats postgresql://user:pass@postgres.postgres.svc.cluster.local:5432/coredns {
        workers 3
        queryTimeout 5s
        statsPrefix coredns
    }

    forward . tls://1.1.1.1 tls://1.0.0.1 {
        tls_servername cloudflare-dns.com
    }
    log
    errors
}
```

### plugin.cfg

This is a sample middleware configuration file. The order of plugins here is important. This is the order in which plugins will be executed for incoming requests.

metadata:metadata
prometheus:metrics
log:log
forward:forward
stats:stats

### Interaction with Other CoreDNS Plugins
#### `metadata`
The plugin will store all metadata stored by the `metadata` plugin for each request. The metadata will be stored in a JSON format.
