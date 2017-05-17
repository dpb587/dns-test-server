package handler

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

//
// Generate a response from the requested domain.
//
// Supported hosts...
//
// * hdr-aa - invert the authoritative header flag
// * hdr-qr - invert the response header flag
// * hdr-ra - invert the recursion available header flag
// * hdr-rd - invert the recursion desired header flag
// * hdr-tc - invert the truncated header flag
// * compress - enable response compression
// * ttl-{seconds} - set TTL for generated answers (default 0)
// * rcode-{string} - set the response code (default noerror; example servfail, notimpl, nxdomain)
// * delay-{duration} - delay the response by a golang duration (example 2s)
// * size-{count} - add {count} dummy answers
// * size-{bytes}b - add dummy answers until the response size is <= {bytes} bytes; remove answers, if necessary
// * answer - add a dummy answer
//
func Generate(res *dns.Msg, domain string) error {
	var ttl uint32

	question := res.Question[0]
	answers := 0

	hosts := strings.Split(domain, ".")
	sort.Sort(sort.Reverse(sort.StringSlice(hosts)))

	for _, host := range hosts {
		if host == "hdr-aa" {
			res.Authoritative = !res.Authoritative
		} else if host == "hdr-qr" {
			res.Response = !res.Response
		} else if host == "hdr-ra" {
			res.RecursionAvailable = !res.RecursionAvailable
		} else if host == "hdr-rd" {
			res.RecursionDesired = !res.RecursionDesired
		} else if host == "hdr-tc" {
			res.Truncated = !res.Truncated
		} else if host == "compress" {
			res.Compress = !res.Compress
		} else if strings.HasPrefix(host, "ttl-") {
			hostTTL := strings.TrimPrefix(host, "ttl-")

			ttlInt, err := strconv.Atoi(hostTTL)
			if err != nil {
				return fmt.Errorf("dns-fuzz: invalid ttl: %s", host)
			}

			ttl = uint32(ttlInt)
		} else if strings.HasPrefix(host, "rcode-") {
			hostRcode := strings.ToUpper(strings.TrimPrefix(host, "rcode-"))
			hostRcodeInt := -1

			for rcode, rcodeString := range dns.RcodeToString {
				if hostRcode != rcodeString {
					continue
				}

				hostRcodeInt = rcode

				break
			}

			if hostRcodeInt < 0 {
				return fmt.Errorf("dns-fuzz: invalid rcode: %s", hostRcode)
			}

			res.Rcode = hostRcodeInt
		} else if strings.HasPrefix(host, "delay-") {
			hostDelay := strings.TrimPrefix(host, "delay-")

			delay, err := time.ParseDuration(hostDelay)
			if err != nil {
				return fmt.Errorf("dns-fuzz: invalid delay: %s", hostDelay)
			}

			time.Sleep(delay)
		} else if strings.HasPrefix(host, "size-") {
			hostSize := strings.TrimPrefix(host, "size-")

			if strings.HasSuffix(hostSize, "b") {
				hostSizeInt, err := strconv.Atoi(strings.TrimSuffix(hostSize, "b"))
				if err != nil {
					return fmt.Errorf("dns-fuzz: invalid size bytes: %s", hostSize)
				}

				for i := 0; i < 102400 && res.Len() < hostSizeInt; i++ {
					res.Answer = append(res.Answer, &dns.A{
						Hdr: dns.RR_Header{
							Name:   question.Name,
							Rrtype: question.Qtype,
							Class:  question.Qclass,
							Ttl:    ttl,
						},
						A: net.ParseIP(fmt.Sprintf("198.51.100.%d", i%255)),
					})
				}

				for res.Len() > hostSizeInt {
					res.Answer = res.Answer[:len(res.Answer)-1]
				}
			} else {
				hostSizeInt, err := strconv.Atoi(hostSize)
				if err != nil {
					return fmt.Errorf("dns-fuzz: invalid size: %s", hostSize)
				}

				for i := 0; i < hostSizeInt; i++ {
					res.Answer = append(res.Answer, &dns.A{
						Hdr: dns.RR_Header{
							Name:   question.Name,
							Rrtype: question.Qtype,
							Class:  question.Qclass,
							Ttl:    ttl,
						},
						A: net.ParseIP(fmt.Sprintf("203.0.113.%d", i%255)),
					})
				}
			}
		} else if host == "answer" {
			res.Answer = append(res.Answer, &dns.A{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: question.Qtype,
					Class:  question.Qclass,
					Ttl:    ttl,
				},
				A: net.ParseIP(fmt.Sprintf("192.0.2.%d", answers)),
			})

			answers++
		} else {
			return fmt.Errorf("dns-fuzz: invalid host: %s", host)
		}
	}

	return nil
}
