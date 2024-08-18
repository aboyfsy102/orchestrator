
```bash
yyin@fedora:~/Desktop/workspace/8_local/j5v3/infra$ curl --verbose --insecure https://j5v3-alb-1882735532.ap-southeast-1.elb.amazonaws.com/order?flag=`date +"%H%M"`
* Host j5v3-alb-1882735532.ap-southeast-1.elb.amazonaws.com:443 was resolved.
* IPv6: (none)
* IPv4: 3.1.1.35, 52.221.90.28, 13.215.62.52
*   Trying 3.1.1.35:443...
* Connected to j5v3-alb-1882735532.ap-southeast-1.elb.amazonaws.com (3.1.1.35) port 443
* ALPN: curl offers h2,http/1.1
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.2 (IN), TLS handshake, Certificate (11):
* TLSv1.2 (IN), TLS handshake, Server key exchange (12):
* TLSv1.2 (IN), TLS handshake, Server finished (14):
* TLSv1.2 (OUT), TLS handshake, Client key exchange (16):
* TLSv1.2 (OUT), TLS change cipher, Change cipher spec (1):
* TLSv1.2 (OUT), TLS handshake, Finished (20):
* TLSv1.2 (IN), TLS handshake, Finished (20):
* SSL connection using TLSv1.2 / ECDHE-RSA-AES128-GCM-SHA256 / secp256r1 / rsaEncryption
* ALPN: server accepted h2
* Server certificate:
*  subject: O=ACME Examples, Inc; CN=example.com
*  start date: Aug 18 08:53:48 2024 GMT
*  expire date: Aug 18 20:53:48 2024 GMT
*  issuer: O=ACME Examples, Inc; CN=example.com
*  SSL certificate verify result: self-signed certificate (18), continuing anyway.
*   Certificate level 0: Public key type RSA (2048/112 Bits/secBits), signed using sha256WithRSAEncryption
*   Certificate level 1: Public key type RSA (2048/112 Bits/secBits), signed using sha256WithRSAEncryption
* using HTTP/2
* [HTTP/2] [1] OPENED stream for https://j5v3-alb-1882735532.ap-southeast-1.elb.amazonaws.com/order?flag=1814
* [HTTP/2] [1] [:method: GET]
* [HTTP/2] [1] [:scheme: https]
* [HTTP/2] [1] [:authority: j5v3-alb-1882735532.ap-southeast-1.elb.amazonaws.com]
* [HTTP/2] [1] [:path: /order?flag=1814]
* [HTTP/2] [1] [user-agent: curl/8.6.0]
* [HTTP/2] [1] [accept: */*]
> GET /order?flag=1814 HTTP/2
> Host: j5v3-alb-1882735532.ap-southeast-1.elb.amazonaws.com
> User-Agent: curl/8.6.0
> Accept: */*
> 
< HTTP/2 200 
< server: awselb/2.0
< date: Sun, 18 Aug 2024 10:14:14 GMT
< content-type: application/json
< content-length: 88
< 
* Connection #0 to host j5v3-alb-1882735532.ap-southeast-1.elb.amazonaws.com left intact
{"message": "GET request processed", "path": "/order", "query_params": {"flag": "1814"}}
```