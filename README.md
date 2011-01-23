go-redir-ws
==========
This is a simple webserver which simply sends redirects for a
predefined list of URLs and logs them and keeps a counter of hits.

Overview
----------

This is a simple webserver which will send redirects and keep logs and
counters about the hits.
It will read the "path to URL" mappings from a file and for each hit,
it will log a record in MongoDB (containing path, time, referer, IP
and userg agent) and will increment the counter "hit:path" (where path
is the actual path) in Redis.

This is the Go implementation of the server, and it's part of a
project which aims at implementing the same server in multiple
languages/frameworks/backends for benchmarking using more realistic
workloads.

Installation
----------

It depends on the following go modules:

*   github.com/mikejs/gomongo/mongo
*   my fork of [redis.go](https://github.com/wladh/redis.go)

Also you need, of course, MongoDB and Redis.

Configuration
----------

You need to invoke the webserver as:
**ws** map\_file [pool\_size]

**map\_file** contains the URLs to redirect. Its format is one entry
  per line, local path and destination URL are separated by spaces,
  like this:  
  `test1 http://www.example.com/file1.html`  
  It will redirect _http://localhost/test1_ to
  _http://www.example.com/file1.html_  
**pool\_size** is the maximum number of connections to the databases.

You also need to create the database "logs" in MongoDB.
