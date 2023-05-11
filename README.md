# certcli
certcli 是一个基于Let's Encrypt 命令行https 证书申请工具, 同时支持更新上传至腾讯云的ssl证书及其关联资源

# 1. Usage
```Bash
$ certcli -h

certcli certcli 是一个基于Let's Encrypt 命令行https 证书申请工具, 同时支持更新上传至腾讯云的ssl证书及其关联资源

Usage:
  certcli [command]

Available Commands:
  dnshelp     Shows additional help for the '--dns' global option
  help        Help about any command
  run         生成基于Let's Encrypt 颁发的ssl 证书(Generate ssl certificates based on Let's Encrypt)
  sync        检查上传至云平台的ssl证书是否过期,并自动更新证书及关联资源
  upload      上传基于Let's Encrypt 颁发的ssl证书至云平台(Upload the ssl certificate issued based on Let's Encrypt to the cloud platform)

Flags:
  -y, --accept-tos                 By setting this flag to true you indicate that you accept the current Let's Encrypt terms of service.
      --cert.timeout int           Set the certificate timeout value to a specific value in seconds. Only used when obtaining certificates. (default 30)
      --csr string                 Certificate signing request filename, if an external CSR is to be used.
      --debug                      Enable debug mode
      --dns string                 Solve a DNS-01 challenge using the specified provider. Can be mixed with other types of challenges. Run 'certcli dnshelp' for help on usage.
      --dns-timeout int            Set the DNS timeout value to a specific value in seconds. Used only when performing authoritative name server queries. (default 10)
      --dns.disable-cp             By setting this flag to true, disables the need to await propagation of the TXT record to all authoritative name servers.
      --dns.resolvers strings      Set the resolvers to use for performing (recursive) CNAME resolving and apex domain determination.
  -d, --domains strings            指定域名,可指定多个以逗号分割(Specify the domain name, you can specify multiple separated by commas)
      --eab                        Use External Account Binding for account registration. Requires --kid and --hmac.
  -m, --email string               Email used for registration and recovery contact.
  -h, --help                       help for certcli
      --hmac string                MAC key from External CA. Should be in Base64 URL Encoding without padding format. Used for External Account Binding.
      --http                       Use the HTTP-01 challenge to solve challenges. Can be mixed with other types of challenges.
      --http-timeout int           Set the HTTP timeout value to a specific value in seconds.
      --http.port string           Set the port and interface to use for HTTP-01 based challenges to listen on. Supported: interface:port or :port. (default ":80")
      --http.proxy-header string   Validate against this HTTP header when solving HTTP-01 based challenges behind a reverse proxy. (default "host")
      --http.webroot string        Set the webroot folder to use for HTTP-01 based challenges to write directly to the .well-known/acme-challenge file.
  -k, --key_type string            Key type to use for private keys. Supported: rsa2048, rsa4096, rsa8192, ec256, ec384. (default "rsa2048")
      --kid string                 Key identifier from External CA. Used for External Account Binding.
      --pem                        Generate an additional .pem (base64) file by concatenating the .key and .crt files together.
      --pfx                        Generate an additional .pfx (PKCS#12) file by concatenating the .key and .crt and issuer .crt files together.
      --pfx-pass string            The password used to encrypt the .pfx (PCKS#12) file. (default "changeit")
  -s, --server string              CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client. (default "https://acme-v02.api.letsencrypt.org/directory")
      --tls                        Use the TLS-ALPN-01 challenge to solve challenges. Can be mixed with other types of challenges.
      --tls.port string            Set the port and interface to use for TLS-ALPN-01 based challenges to listen on. Supported: interface:port or :port. (default ":443")
      --user-agent string          Add to the user-agent sent to the CA to identify an application embedding lego-cli
  -v, --v Level                    number for the log level verbosity

Use "certcli [command] --help" for more information about a command.
```

## 1.1 场景一：生成一个基于Let's Encrypt的证书

### 1.1.1 介绍
该功能基于 https://github.com/go-acme/lego

用法和lego run 一致 


```Bash
$ certcli run -h

生成基于Let's Encrypt 颁发的ssl 证书(Generate ssl certificates based on Let's Encrypt)

Usage:
  certcli run [flags]

Flags:
      --always-deactivate-authorizations   Force the authorizations to be relinquished even if the certificate request was successful.
  -h, --help                               help for run
      --must-staple                        Include the OCSP must staple TLS extension in the CSR and generated certificate.
      --no-bundle                          Do not create a certificate bundle by adding the issuers certificate to the new certificate.
      --path string                        Directory to use for storing the data. (default "/Users/wangruihua/sensorsdata/project/certcli/.certcli")
      --preferred-chain string             If the CA offers multiple certificate chains, prefer the chain with an issuer matching this Subject Common Name.
      --run-hook string                    Define a hook. The hook is executed when the certificates are effectively created.

Global Flags:
  -y, --accept-tos                 By setting this flag to true you indicate that you accept the current Let's Encrypt terms of service.
      --cert.timeout int           Set the certificate timeout value to a specific value in seconds. Only used when obtaining certificates. (default 30)
      --csr string                 Certificate signing request filename, if an external CSR is to be used.
      --debug                      Enable debug mode
      --dns string                 Solve a DNS-01 challenge using the specified provider. Can be mixed with other types of challenges. Run 'certcli dnshelp' for help on usage.
      --dns-timeout int            Set the DNS timeout value to a specific value in seconds. Used only when performing authoritative name server queries. (default 10)
      --dns.disable-cp             By setting this flag to true, disables the need to await propagation of the TXT record to all authoritative name servers.
      --dns.resolvers strings      Set the resolvers to use for performing (recursive) CNAME resolving and apex domain determination.
  -d, --domains strings            指定域名,可指定多个以逗号分割(Specify the domain name, you can specify multiple separated by commas)
      --eab                        Use External Account Binding for account registration. Requires --kid and --hmac.
  -m, --email string               Email used for registration and recovery contact.
      --hmac string                MAC key from External CA. Should be in Base64 URL Encoding without padding format. Used for External Account Binding.
      --http                       Use the HTTP-01 challenge to solve challenges. Can be mixed with other types of challenges.
      --http-timeout int           Set the HTTP timeout value to a specific value in seconds.
      --http.port string           Set the port and interface to use for HTTP-01 based challenges to listen on. Supported: interface:port or :port. (default ":80")
      --http.proxy-header string   Validate against this HTTP header when solving HTTP-01 based challenges behind a reverse proxy. (default "host")
      --http.webroot string        Set the webroot folder to use for HTTP-01 based challenges to write directly to the .well-known/acme-challenge file.
  -k, --key_type string            Key type to use for private keys. Supported: rsa2048, rsa4096, rsa8192, ec256, ec384. (default "rsa2048")
      --kid string                 Key identifier from External CA. Used for External Account Binding.
      --pem                        Generate an additional .pem (base64) file by concatenating the .key and .crt files together.
      --pfx                        Generate an additional .pfx (PKCS#12) file by concatenating the .key and .crt and issuer .crt files together.
      --pfx-pass string            The password used to encrypt the .pfx (PCKS#12) file. (default "changeit")
  -s, --server string              CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client. (default "https://acme-v02.api.letsencrypt.org/directory")
      --tls                        Use the TLS-ALPN-01 challenge to solve challenges. Can be mixed with other types of challenges.
      --tls.port string            Set the port and interface to use for TLS-ALPN-01 based challenges to listen on. Supported: interface:port or :port. (default ":443")
      --user-agent string          Add to the user-agent sent to the CA to identify an application embedding lego-cli
  -v, --v Level                    number for the log level verbosity
```

### 1.1.2 Demo
推荐使用dns进行challenges

如下采用aldns，需要先设置aliyun apikey

> `NOTICE`
>
> 需设置环境变量
>
> ```BASH
>export ALICLOUD_ACCESS_KEY="[YOUR ACCESS KEY]"
>export ALICLOUD_SECRET_KEY="[YOUR SECRET KEY]"
>```


```Bash
$ certcli run -d *.test.example.cn -m example@example.cn --dns alidns -y

```


## 1.2 场景二：生成一个基于Let's Encrypt的证书 并上传至云平台进行托管

### 1.2.1 介绍
支持生成证书并直接上传至云平台进行托管
> `目前只支持腾讯云`
>

```Bash
$ certcli upload -h

上传基于Let's Encrypt 颁发的ssl证书至云平台(Upload the ssl certificate issued based on Let's Encrypt to the cloud platform)

Usage:
  certcli upload [flags]

Flags:
  -f, --filename string   通过文件指定域名,一行一个域名(Specify the domain name by file, one domain name per line)
  -h, --help              help for upload
  -p, --provider string   指定云厂商,暂时只支持(txcloud) (default "txcloud")
  
Global Flags:
  -y, --accept-tos                 By setting this flag to true you indicate that you accept the current Let's Encrypt terms of service.
      --cert.timeout int           Set the certificate timeout value to a specific value in seconds. Only used when obtaining certificates. (default 30)
      --csr string                 Certificate signing request filename, if an external CSR is to be used.
      --debug                      Enable debug mode
      --dns string                 Solve a DNS-01 challenge using the specified provider. Can be mixed with other types of challenges. Run 'certcli dnshelp' for help on usage.
      --dns-timeout int            Set the DNS timeout value to a specific value in seconds. Used only when performing authoritative name server queries. (default 10)
      --dns.disable-cp             By setting this flag to true, disables the need to await propagation of the TXT record to all authoritative name servers.
      --dns.resolvers strings      Set the resolvers to use for performing (recursive) CNAME resolving and apex domain determination.
  -d, --domains strings            指定域名,可指定多个以逗号分割(Specify the domain name, you can specify multiple separated by commas)
      --eab                        Use External Account Binding for account registration. Requires --kid and --hmac.
  -m, --email string               Email used for registration and recovery contact.
      --hmac string                MAC key from External CA. Should be in Base64 URL Encoding without padding format. Used for External Account Binding.
      --http                       Use the HTTP-01 challenge to solve challenges. Can be mixed with other types of challenges.
      --http-timeout int           Set the HTTP timeout value to a specific value in seconds.
      --http.port string           Set the port and interface to use for HTTP-01 based challenges to listen on. Supported: interface:port or :port. (default ":80")
      --http.proxy-header string   Validate against this HTTP header when solving HTTP-01 based challenges behind a reverse proxy. (default "host")
      --http.webroot string        Set the webroot folder to use for HTTP-01 based challenges to write directly to the .well-known/acme-challenge file.
  -k, --key_type string            Key type to use for private keys. Supported: rsa2048, rsa4096, rsa8192, ec256, ec384. (default "rsa2048")
      --kid string                 Key identifier from External CA. Used for External Account Binding.
      --pem                        Generate an additional .pem (base64) file by concatenating the .key and .crt files together.
      --pfx                        Generate an additional .pfx (PKCS#12) file by concatenating the .key and .crt and issuer .crt files together.
      --pfx-pass string            The password used to encrypt the .pfx (PCKS#12) file. (default "changeit")
  -s, --server string              CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client. (default "https://acme-v02.api.letsencrypt.org/directory")
      --tls                        Use the TLS-ALPN-01 challenge to solve challenges. Can be mixed with other types of challenges.
      --tls.port string            Set the port and interface to use for TLS-ALPN-01 based challenges to listen on. Supported: interface:port or :port. (default ":443")
      --user-agent string          Add to the user-agent sent to the CA to identify an application embedding lego-cli
  -v, --v Level                    number for the log level verbosity
```

### 1.2.2 Demo

如下采用aldns，需要先设置aliyun apikey
同时设置腾讯云 apikey

> `NOTICE`
>
> 需设置环境变量
>
> ```BASH
>#DNS 托管key
>export ALICLOUD_ACCESS_KEY="[YOUR ACCESS KEY]"
>export ALICLOUD_SECRET_KEY="[YOUR SECRET KEY]"
>
>#云厂商key
>export TENCENTCLOUD_SECRET_ID="[YOUR ACCESS KEY]"
>export TENCENTCLOUD_SECRET_KE="[YOUR SECRET KEY]"
>```


```Bash
$ certcli upload -d *.test.example.cn --dns alidns -p txcloud

```

## 1.3 场景三：支持检测托管至云平台的证书,进行自动更新

### 1.3.1 介绍
支持检测托管至云平台的证书，设置过期时间阈值，超过，即自动生成基于Let's Encrypt的证书进行更新，并更新关联资源

> `目前只支持腾讯云`
> 
>`更新关联资源, 目前只支持 [clb、cdn]`


```Bash
$ certcli sync -h

检查上传至云平台的ssl证书是否过期,并自动更新证书及关联资源

Usage:
  certcli sync [flags]

Flags:
  -A, --all-cert                指定此选项,即会监听所有证书进行定时更新,谨慎使用！！！
      --days int                指定证书还剩下多少天可以更新(The number of days left on a certificate to renew it) (default 14)
  -f, --filename string         通过文件指定域名,一行一个域名(Specify the domain name by file, one domain name per line)
  -h, --help                    help for sync
      --new-id string           指定要更新的新证书ID
      --old-id string           指定要更新的旧证书ID
  -p, --provider string         指定云厂商,暂时只支持(txcloud) (default "txcloud")
  -r, --resource-type strings   指定需要更新的资源, 可指定多个,以逗号分割(目前只支持资源类型:clb,cdn)

Global Flags:
  -y, --accept-tos                 By setting this flag to true you indicate that you accept the current Let's Encrypt terms of service.
      --cert.timeout int           Set the certificate timeout value to a specific value in seconds. Only used when obtaining certificates. (default 30)
      --csr string                 Certificate signing request filename, if an external CSR is to be used.
      --debug                      Enable debug mode
      --dns string                 Solve a DNS-01 challenge using the specified provider. Can be mixed with other types of challenges. Run 'certcli dnshelp' for help on usage.
      --dns-timeout int            Set the DNS timeout value to a specific value in seconds. Used only when performing authoritative name server queries. (default 10)
      --dns.disable-cp             By setting this flag to true, disables the need to await propagation of the TXT record to all authoritative name servers.
      --dns.resolvers strings      Set the resolvers to use for performing (recursive) CNAME resolving and apex domain determination.
  -d, --domains strings            指定域名,可指定多个以逗号分割(Specify the domain name, you can specify multiple separated by commas)
      --eab                        Use External Account Binding for account registration. Requires --kid and --hmac.
  -m, --email string               Email used for registration and recovery contact.
      --hmac string                MAC key from External CA. Should be in Base64 URL Encoding without padding format. Used for External Account Binding.
      --http                       Use the HTTP-01 challenge to solve challenges. Can be mixed with other types of challenges.
      --http-timeout int           Set the HTTP timeout value to a specific value in seconds.
      --http.port string           Set the port and interface to use for HTTP-01 based challenges to listen on. Supported: interface:port or :port. (default ":80")
      --http.proxy-header string   Validate against this HTTP header when solving HTTP-01 based challenges behind a reverse proxy. (default "host")
      --http.webroot string        Set the webroot folder to use for HTTP-01 based challenges to write directly to the .well-known/acme-challenge file.
  -k, --key_type string            Key type to use for private keys. Supported: rsa2048, rsa4096, rsa8192, ec256, ec384. (default "rsa2048")
      --kid string                 Key identifier from External CA. Used for External Account Binding.
      --pem                        Generate an additional .pem (base64) file by concatenating the .key and .crt files together.
      --pfx                        Generate an additional .pfx (PKCS#12) file by concatenating the .key and .crt and issuer .crt files together.
      --pfx-pass string            The password used to encrypt the .pfx (PCKS#12) file. (default "changeit")
  -s, --server string              CA hostname (and optionally :port). The server certificate must be trusted in order to avoid further modifications to the client. (default "https://acme-v02.api.letsencrypt.org/directory")
      --tls                        Use the TLS-ALPN-01 challenge to solve challenges. Can be mixed with other types of challenges.
      --tls.port string            Set the port and interface to use for TLS-ALPN-01 based challenges to listen on. Supported: interface:port or :port. (default ":443")
      --user-agent string          Add to the user-agent sent to the CA to identify an application embedding lego-cli
  -v, --v Level                    number for the log level verbosity
```

### 1.3.2 Demo

如下采用aldns，需要先设置aliyun apikey
同时设置腾讯云 apikey

> `NOTICE`
>
> 需设置环境变量
>
> ```BASH
>#DNS 托管key
>export ALICLOUD_ACCESS_KEY="[YOUR ACCESS KEY]"
>export ALICLOUD_SECRET_KEY="[YOUR SECRET KEY]"
>
>#云厂商key
>export TENCENTCLOUD_SECRET_ID="[YOUR ACCESS KEY]"
>export TENCENTCLOUD_SECRET_KE="[YOUR SECRET KEY]"
>```

+ 支持指定多个域名来更新 (通过指定 -d 参数)
```Bash
$ certcli sync -d *.test.example.cn,*.dev.example.cn --dns alidns -r cdn,clb --days 15

```

+ 支持将域名写入文件(一行一个域名)，然后指定文件名来更新 （通过指定-f 参数）
```Bash
$ certcli sync -f domains.txt --dns alidns -r cdn,clb --days 15

```

+ 支持检测所有上传至云平台的证书（通过指定 -A 参数 `谨慎使用`）

```Bash
$ certcli sync -A --dns alidns -r cdn,clb --days 15

```

## 1.4 场景四：支持通过指定新老证书ID，来更新云平台上证书关联资源

### 1.4.1 介绍
支持通过指定新老证书ID，来更新云平台上证书关联资源。
适用于, 不想使用基于Let's Encrypt的证书，手动上传至新证书至云平台上，然后去更新证书关联资源的场景

### 1.4.2 Demo

设置腾讯云 apikey

> `NOTICE`
>
> 需设置环境变量
>
> ```BASH
>#云厂商key
>export TENCENTCLOUD_SECRET_ID="[YOUR ACCESS KEY]"
>export TENCENTCLOUD_SECRET_KE="[YOUR SECRET KEY]"
>```


```Bash
$ certcli sync  --new-id "新的证书ID" --old-id "老的证书ID" -r clb,cdn

```