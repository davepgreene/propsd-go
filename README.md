# Propsd

Propsd does dynamic property management at scale, across thousands of servers
and changes from hundreds of developers.

We built Propsd with lessons learned from years of running [Conqueso][] on
large scale systems. High availability is achieved by leveraging [Amazon S3][]
to deliver properties and [Consul][] to handle service discovery. Composable
layering lets you set properties for an organization, a single server, and
everything in between. Plus, flat file storage makes backups and audits
a breeze.

So if your [Conqueso][] server's starting to heat up or you just want an audit
trail when things change, give Propsd a try.

## Features
Propsd allows the user to supply a [index](https://github.com/rapid7/propsd/blob/master/docs/getting-started/usage.md#index-files) file which defines the layering of configuration sources.  Propsd expects these configuration sources to be in json format.

Propsd will serve up this layering of configuration (last in, first out) in combination of other services (including a given [Consul][] catalog) via HTTP to any requesting services.

For example, if you have a configuration layering schema like this:

~~~text
global
|_global.json (foo = global)
|_regional
  |_region-1
    |_region.json (foo = region-1)
    |_service
      |_service_name.json (foo = service_specific)
~~~
Propsd would return flattened configuration where the `service_name.json` value would win (netting a `foo` value of `service_specific`).

Propsd can also consume from various source locations.  This ranges from a Consul catalog, to local files, to remote S3 buckets.  This feature helps the user package once and let the configuration source location contain the differences between environments.  Said another way - with Propsd, you won't find yourself repackaging your software in order to move it through your environments.

Propsd will regularly inspect and reload these configuration settings, include new layers, etc. - no restarts of Propsd required.

[Conqueso]: https://github.com/rapid7/conqueso
[Amazon S3]: https://aws.amazon.com/s3/
[Consul]: https://www.consul.io/