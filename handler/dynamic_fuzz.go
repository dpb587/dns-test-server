package handler

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

type DynamicFuzz struct {
	base string
}

func NewDynamicFuzz(base string) DynamicFuzz {
	return DynamicFuzz{base: base}
}

func (h DynamicFuzz) ServeDNS(w dns.ResponseWriter, req *dns.Msg) {
	if len(req.Question) == 0 {
		fmt.Println(fmt.Errorf("dns-fuzz: query missing question"))
	}

	question := req.Question[0]

	if !strings.HasSuffix(question.Name, h.base) {
		fmt.Println(fmt.Errorf("dns-fuzz: missing question"))

		return
	}

	res := &dns.Msg{}
	res.SetReply(req)

	err := Generate(res, strings.TrimSuffix(question.Name, "."+h.base))
	if err != nil {
		fmt.Println(err)

		return
	}

	err = w.WriteMsg(res)
	if err != nil {
		fmt.Println(fmt.Errorf("dns-fuzz: failed to write response: %s", err.Error()))
	}
}
