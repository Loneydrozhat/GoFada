# GoFada
[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://github.com/dudsan/GoFada/blob/master/LICENSE)

A light web crawler/directory brute force scanner written in Golang

## Brief
I have written this mainly as an exercise on Golang, yet it is a fully functional web crawler/directory brute force scanner to be used for mapping standard web applications. It's only 5 MB large and fully text-based, so it might be useful in pivoting scenarios.

This tool is under constant enhancement. Contributions are most welcome. :-)

You might want to take a look at GoFada's board on Trello: https://trello.com/b/j5Hmw6PS/gofada.
## Getting Started
In order to compile GoFada you're supposed to have Go installed. Please refer to https://golang.org/doc/install for detail. Go packages may be installed using go get, e.g.
```
$ go get packagex
```
will install Go's packagex on your system. Make sure all packages imported by GoFada are installed prior to compiling it. Compile by issuing
```
$ go build GoFada.go
```
## Usage
After compiling GoFada, launch it by issuing
```
$ ./GoFada
```
and you'll see a puking stickman, followed by a menu. It is extremely easy to use GoFada since it comes with a text-based menu. For instance, to set a base URL (for crawling/scanning) pick option "1", followed by option "1" as shown below.
```
@> ./GoFada

     (}
    /Y\`;,
    /^\  ;:,
:::GoFada:::

:::pick an option:::
(1) set params
(2) crawl
(3) discover
(4) show crawled
(5) show brute forced
(0) quit
(*) show puking stickman

gfd>1

:::current params:::
(1) base_url: localhost
(2) scope_filter: localhost
(3) depth: 2
(4) wordlist: none
(5) throttle: 0 ms
(6) jitter: at most 0 ms
(*) back

gfd|params>1

:::set new base_url:::
gfd|set params|base_url>http://www.foo.bar/
```
## Acknowledgements
* For information on SVNDigger wordlists (which I uploaded to this project) please refer to https://www.netsparker.com/blog/web-security/svn-digger-better-lists-for-forced-browsing/ - I did not compile these myself. 
## Authors
* Eduardo Vasconcelos - esmev@protonmail.ch - dudsan (main author)
* Diogo Brandão - diogo.brandao@perallis.com (contributor)
## License
This project is lincensed under the terms of the GPL License. Please refer to https://www.gnu.org/licenses/gpl.html for detail.
