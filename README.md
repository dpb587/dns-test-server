# dns-fuzz

A DNS server to generate fake responses for testing client behavior.

    $ go run main.go &
    $ dig -p 35053 +tcp @127.0.0.1 rcode-servfail.delay-8s.size-4.ttl-7.hdr-tc.hdr-aa.test

    ; <<>> DiG 9.8.3-P1 <<>> -p 35053 +tcp @127.0.0.1 rcode-servfail.delay-8s.size-4.ttl-7.hdr-tc.hdr-aa.test
    ; (1 server found)
    ;; global options: +cmd
    ;; Got answer:
    ;; ->>HEADER<<- opcode: QUERY, status: SERVFAIL, id: 37250
    ;; flags: qr aa tc rd; QUERY: 1, ANSWER: 4, AUTHORITY: 0, ADDITIONAL: 0
    ;; WARNING: recursion requested but not available

    ;; QUESTION SECTION:
    ;rcode-servfail.delay-8s.size-4.ttl-7.hdr-tc.hdr-aa.test. IN A

    ;; ANSWER SECTION:
    rcode-servfail.delay-8s.size-4.ttl-7.hdr-tc.hdr-aa.test. 7 IN A	203.0.113.0
    rcode-servfail.delay-8s.size-4.ttl-7.hdr-tc.hdr-aa.test. 7 IN A	203.0.113.1
    rcode-servfail.delay-8s.size-4.ttl-7.hdr-tc.hdr-aa.test. 7 IN A	203.0.113.2
    rcode-servfail.delay-8s.size-4.ttl-7.hdr-tc.hdr-aa.test. 7 IN A	203.0.113.3

    ;; Query time: 8003 msec
    ;; SERVER: 127.0.0.1#35053(127.0.0.1)
    ;; WHEN: Tue May 16 23:03:29 2017
    ;; MSG SIZE  rcvd: 357
