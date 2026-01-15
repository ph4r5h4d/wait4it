+++
date = '2026-01-15T20:35:00+01:00'
draft = false
title = 'DNS Check'
+++
DNS Check validates that a hostname can be resolved via DNS. Supports all record types (A, AAAA, CNAME, MX, TXT, SRV, NS, PTR), custom DNS servers, and expected value matching.

## Environment Variables

| Variable        | Description                                              | Default   |
|-----------------|----------------------------------------------------------|-----------|
| W4IT_TYPE       | The type of check (set to `dns`).                        | -         |
| W4IT_TIMEOUT    | Timeout in seconds.                                      | 30        |
| W4IT_HOST       | The hostname to resolve.                                 | 127.0.0.1 |
| W4IT_DNS_TYPE   | DNS record type (A, AAAA, CNAME, MX, TXT, SRV, NS, PTR). | A         |
| W4IT_DNS_EXPECT | Expected value to find in DNS records (substring match). | -         |
| W4IT_DNS_SERVER | Custom DNS server to query (e.g., `8.8.8.8:53`).         | -         |

## Command-Line Arguments

| Argument     | Description                                              | Default   |
|--------------|----------------------------------------------------------|-----------|
| -type        | The type of check (set to `dns`).                        | -         |
| -t           | Timeout in seconds.                                      | 30        |
| -h           | The hostname to resolve.                                 | 127.0.0.1 |
| -dns-type    | DNS record type (A, AAAA, CNAME, MX, TXT, SRV, NS, PTR). | A         |
| -dns-expect  | Expected value to find in DNS records (substring match). | -         |
| -dns-server  | Custom DNS server to query (e.g., `8.8.8.8:53`).         | -         |

## Supported Record Types

| Type  | Description                                    |
|-------|------------------------------------------------|
| A     | IPv4 address records                           |
| AAAA  | IPv6 address records                           |
| CNAME | Canonical name records                         |
| MX    | Mail exchange records                          |
| TXT   | Text records (SPF, DKIM, domain verification)  |
| SRV   | Service records (Kubernetes, service discovery)|
| NS    | Nameserver records                             |
| PTR   | Reverse DNS (pointer) records                  |

## Notes
{{< callout type="info" >}}
- If `-dns-expect` is not specified, the check succeeds if any records are found.
- The expected value uses substring matching, so partial matches work.
- Custom DNS server format is `host:port` (e.g., `8.8.8.8:53`).
- For SRV records, use the full format: `_service._proto.name`
{{< /callout >}}

## Exit Codes
| Code | Meaning                           |
|------|-----------------------------------|
| 0    | DNS resolution successful.        |
| 1    | Timed out.                        |
| 2    | Validation error or invalid input.|


## Usage Examples
### Basic DNS Resolution
```bash
# Check if hostname resolves (A record by default)
./wait4it -type=dns -h=example.com -t=60

# Simple check with short timeout
./wait4it -type=dns -h=myservice.local -t=10
```

### Record Type Examples

#### A Record (IPv4)
```bash
# Check A record resolves
./wait4it -type=dns -h=example.com -dns-type=A -t=60

# Check A record matches specific IP prefix
./wait4it -type=dns -h=example.com -dns-type=A -dns-expect="93.184" -t=60

# Check A record matches exact IP
./wait4it -type=dns -h=example.com -dns-type=A -dns-expect="93.184.216.34" -t=60
```

#### AAAA Record (IPv6)
```bash
# Check IPv6 record exists
./wait4it -type=dns -h=google.com -dns-type=AAAA -t=60

# Check IPv6 record matches prefix
./wait4it -type=dns -h=google.com -dns-type=AAAA -dns-expect="2607:f8b0" -t=60
```

#### CNAME Record
```bash
# Check CNAME record exists
./wait4it -type=dns -h=www.github.com -dns-type=CNAME -t=60

# Check CNAME points to expected target
./wait4it -type=dns -h=www.github.com -dns-type=CNAME -dns-expect="github.com" -t=60
```

#### MX Record (Mail)
```bash
# Check MX record exists
./wait4it -type=dns -h=gmail.com -dns-type=MX -t=60

# Check MX record points to Google mail
./wait4it -type=dns -h=gmail.com -dns-type=MX -dns-expect="google.com" -t=60

# Wait for new mail server to propagate
./wait4it -type=dns -h=mydomain.com -dns-type=MX -dns-expect="mail.mydomain.com" -t=300
```

#### TXT Record
```bash
# Check TXT record exists
./wait4it -type=dns -h=google.com -dns-type=TXT -t=60

# Check SPF record exists
./wait4it -type=dns -h=google.com -dns-type=TXT -dns-expect="v=spf1" -t=60

# Check DKIM record
./wait4it -type=dns -h=selector._domainkey.example.com -dns-type=TXT -dns-expect="v=DKIM1" -t=60

# Wait for domain verification (Google, AWS, etc.)
./wait4it -type=dns -h=example.com -dns-type=TXT -dns-expect="google-site-verification" -t=300
```

#### SRV Record (Service Discovery)
```bash
# Check SRV record for Kubernetes service
./wait4it -type=dns -h=_mongodb._tcp.default.svc.cluster.local -dns-type=SRV -t=60

# Check SRV record for SIP service
./wait4it -type=dns -h=_sip._tcp.example.com -dns-type=SRV -t=60

# Check SRV record for LDAP
./wait4it -type=dns -h=_ldap._tcp.example.com -dns-type=SRV -t=60
```

#### NS Record (Nameservers)
```bash
# Check nameserver records
./wait4it -type=dns -h=example.com -dns-type=NS -t=60

# Check for specific nameserver
./wait4it -type=dns -h=example.com -dns-type=NS -dns-expect="cloudflare" -t=60
```

#### PTR Record (Reverse DNS)
```bash
# Check reverse DNS for an IP
./wait4it -type=dns -h=8.8.8.8 -dns-type=PTR -t=60

# Check reverse DNS matches expected hostname
./wait4it -type=dns -h=8.8.8.8 -dns-type=PTR -dns-expect="google" -t=60
```

### Custom DNS Server Examples
```bash
# Query Google Public DNS
./wait4it -type=dns -h=example.com -dns-server=8.8.8.8:53 -t=60

# Query Cloudflare DNS
./wait4it -type=dns -h=example.com -dns-server=1.1.1.1:53 -t=60

# Query internal corporate DNS
./wait4it -type=dns -h=myservice.internal -dns-server=10.0.0.1:53 -t=60

# Check DNS propagation on specific server
./wait4it -type=dns -h=newsite.com -dns-type=A -dns-server=ns1.mydns.com:53 -dns-expect="192.168.1.100" -t=300
```

### Docker Examples
```bash
# Basic DNS check
docker run ph4r5h4d/wait4it -type=dns -h=example.com -t=60

# Check MX record
docker run ph4r5h4d/wait4it -type=dns -h=gmail.com -dns-type=MX -t=60

# Check with custom DNS server
docker run ph4r5h4d/wait4it -type=dns -h=example.com -dns-server=8.8.8.8:53 -t=60

# Check TXT record for domain verification
docker run ph4r5h4d/wait4it -type=dns -h=example.com -dns-type=TXT -dns-expect="v=spf1" -t=60
```

### Real-World Use Cases

#### Wait for DNS Propagation After Record Change
```bash
# Wait up to 10 minutes for DNS propagation
./wait4it -type=dns -h=newsite.example.com -dns-type=A -dns-expect="203.0.113.50" -t=600
```

#### Wait for Kubernetes Service DNS
```bash
# Wait for service to be discoverable
./wait4it -type=dns -h=myservice.default.svc.cluster.local -t=120
```

#### Wait for Internal Service Discovery
```bash
# Wait for Consul service DNS
./wait4it -type=dns -h=redis.service.consul -t=60

# Wait for CoreDNS entry
./wait4it -type=dns -h=cache.internal -dns-server=10.96.0.10:53 -t=60
```

#### Pre-flight Check Before Deployment
```bash
# Verify DNS is configured correctly before app starts
./wait4it -type=dns -h=api.example.com -dns-type=CNAME -dns-expect="loadbalancer.example.com" -t=30
```