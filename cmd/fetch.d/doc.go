// Dataflow kit - main
//
// Copyright © 2017-2018 Slotix s.r.o. <dm@slotix.sk>
//
//
// All rights reserved. Use of this source code is governed
// by the BSD 3-Clause License license.

// Fetcher service of the Dataflow kit downloads html content from web pages to feed Dataflow kit scrapers.
//
// Currently two types of fetcher are available : Headless Chrome Fetcher and Base Fetcher.
//
// Base fetcher is used for downloading html web page using Go standard library's http.
//
// Chrome Fetcher connects to Headless Chrome which renders JavaScript pages.
//
// Accessing Fetcher endpoints
//
// Examples
//		fetch a web page using CDP
//		curl -XPOST  localhost:8000/fetch -d '{"type":"chrome", "url":"http://example.com"}'
//		fetch a web page with base fetcher. For base fetcher type parameter may be omitted.
//		curl -XPOST  localhost:8000/fetch -d '{"url":"http://example.com"}'
//
// Flags and configuration settings
//
//General settings
//		DFK_FETCH: HTTP listen address of Fetch service (defaults to "127.0.0.1:8000")
//		CHROME: Headless Chrome address. It is used for fetching JS driven web pages (defaults to http://127.0.0.1:9222)
//		PROXY: Proxy address http://username:password@proxy-host:port . (defaults to "")
//Storage settings
//		STORAGE_TYPE: Storage type may be Diskv or Cassandra. (defaults to "Diskv")
//		Storage stores auxiliary information generated by fetcher.
//		DISKV_BASE_DIR: diskv base directory for Diskv Storage type (defaults to "diskv").
//		Find more information about Diskv storage at https://github.com/peterbourgon/diskv
//		CASSANDRA: Cassandra host address (defaults to 127.0.0.1)
//
package main

// EOF
