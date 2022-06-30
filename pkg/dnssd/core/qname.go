package core

import "strings"

type QNamePart = []byte

type FullQName struct {
	data      []byte
	nameCount uint8
}

func ParseFullQName(args ...string) (*FullQName, error) {
	data := make([]byte, 0)
	var i uint8 = 0
	for _, arg := range args {
		if arg == "" || len(arg) > 63 {
			continue
		}
		data = append(data, byte(len(arg)))
		data = append(data, []byte(arg)...)
		i++
	}
	data = append(data, 0)
	return &FullQName{
		data:      data,
		nameCount: i + 1,
	}, nil
}

func (n *FullQName) Bytes() []byte {
	return n.data
}

func (n *FullQName) String() (domainName string) {
	domainName = ""
	qname := n.data
	for i := 0; qname[i] != 0; {
		domainLen := int(qname[i])
		// since length of domain name also occupies an octet
		// for example, "google.com":
		// i++ make j begin at domain name 'g' or 'c', instead of length '0x06' or '0x03'
		i++
		for j := 0; j < domainLen; j++ {
			domainName += string(qname[i+j])
		}
		// it has to be NOTICED that "google.com" will be translated into "google.com."
		domainName += "."
		i += domainLen
	}
	// trim the last '.'
	return strings.Trim(domainName, ".")
}
