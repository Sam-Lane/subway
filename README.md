<h1 align="center">🥪 SUBway </h1>
<p align="center">Another subdomain enumberation tool</p>


Enumerate subdomains by either using DNS lookup or by virtual hosting HTTP requests, useful for things like [Hack The Box](https://www.hackthebox.eu) or [Try Hack Me](https://tryhackme.com/). SUBway requires a wordlist to use for subdomain discovery, [SecLists](https://github.com/danielmiessler/SecLists) is the recomended pairing for use with this tool.


## Usage
```
DNS lookup
subway -h example.com -w subdomains.txt

Virtual hosts:
subway -h example.com -w subdomains.txt -v -ip 127.0.0.1

-h string
        base apex domain to enumerate against, eg: example.com
-hc string
        Hide responses that contain character length, eg: 2000,20043
-hs string
        Ignore status codes, eg: 302,301,401. Status 404 is always ignored
-ip string
        IP used for virtual hosting (REQUIRED), eg: 127.0.0.1
-v    enumerate subdomains by virtual hosts instead of DNS lookup
-w string
        path to wordlist to use, eg: /usr/share/subdomains.txt
```

## Examples

### DNS look up
```
subway -h example.com -w subdomains.txt
====================================
   ______  _____                  
  / __/ / / / _ )_    _____ ___ __
 _\ \/ /_/ / _  | |/|/ / _ '/ // /
/___/\____/____/|__,__/\_,_/\_, / 
                           /___/  
====================================

host: example.com
wordlist: subdomains.txt
====================================
www.example.com
wild.example.com
mail.example.com
server.example.com
dev1.example.com
100:100
done
```

### Virtual hosting look up
When looking by virtual hosts you get additional information such as the status code of the response and the content length.
```
subway -h example.com -w subdomains.txt -v -ip 127.0.0.1
====================================
   ______  _____                  
  / __/ / / / _ )_    _____ ___ __
 _\ \/ /_/ / _  | |/|/ / _ '/ // /
/___/\____/____/|__,__/\_,_/\_, / 
                           /___/  
====================================

host: example.com
wordlist: subdomains.txt
====================================
200 388178  www.example.com
200 388178  wild.example.com
200 11720  mail.example.com
302 186015  server.example.com
401 4438  dev1.example.com
100:100
done
```

##### Filter out status codes
Lets filter out any responses that contain a 302 and a 401
```
subway -h example.com -w subdomains.txt -v -ip 127.0.0.1 -hs 302,401
====================================
   ______  _____                  
  / __/ / / / _ )_    _____ ___ __
 _\ \/ /_/ / _  | |/|/ / _ '/ // /
/___/\____/____/|__,__/\_,_/\_, / 
                           /___/  
====================================

host: example.com
wordlist: subdomains.txt
ignore status: 401,302
====================================
200 388178  www.example.com
200 388178  wild.example.com
200 11720  mail.example.com
100:100
done
```

##### Filter out content lengths
Lets filter out any responses that contain a content-length. This is useful for where a wild card subdomain has been setup to redirect to the base domain.
Previously www.example.com and wild.example.com both direct to the same site and have the same content-length. Lets filter this out any repsonse with content length `388178`.
```
subway -h example.com -w subdomains.txt -v -ip 127.0.0.1 -wc 388178
====================================
   ______  _____                  
  / __/ / / / _ )_    _____ ___ __
 _\ \/ /_/ / _  | |/|/ / _ '/ // /
/___/\____/____/|__,__/\_,_/\_, / 
                           /___/  
====================================

host: example.com
wordlist: subdomains.txt
ignore content length: 388178
====================================
200 11720  mail.example.com
302 186015  server.example.com
401 4438  dev1.example.com
100:100
done
```