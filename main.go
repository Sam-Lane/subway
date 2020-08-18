package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sam-lane/golor"
)

const (
	BANNER = `
====================================
   ______  _____                  
  / __/ / / / _ )_    _____ ___ __
 _\ \/ /_/ / _  | |/|/ / _ '/ // /
/___/\____/____/|__,__/\_,_/\_, / 
                           /___/  
====================================
`
)

func readWordList(wordlist string) []string {
	fileBytes, err := ioutil.ReadFile(wordlist)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sliceData := strings.Split(string(fileBytes), "\n")
	return sliceData
}

func argsToSlice(args string) []string {
	return strings.Split(args, ",") // please work
}

func colourStatusCode(statusCode string) string {
	status := string(statusCode[0])
	switch status {
	case "2":
		return golor.Green(statusCode)
	case "3":
		return golor.Yellow(statusCode)
	case "4":
		return golor.Red(statusCode)
	case "5":
		return golor.Blue(statusCode)
	default:
		return golor.White(statusCode)
	}
}

func dnsMode(hostvar string, subdomainList []string) {
	for index, sub := range subdomainList {
		fmt.Printf("\r%d:%d", index+1, len(subdomainList))
		subdomain := fmt.Sprintf("%s.%s", sub, hostvar)

		_, err := net.LookupHost(subdomain)
		if err != nil {
			continue
		}

		fmt.Printf("%c[2K\r%s\n", 27, subdomain)
	}
}

func isInList(x string, ignoreList []string) bool {
	for _, element := range ignoreList {
		if x == element {
			return true
		}
	}
	return false
}

func virtMode(hostvar string, subdomainList []string, ip string, hs string, wc string) {
	hideStatusList := argsToSlice(hs)
	wordCountList := argsToSlice(wc)

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	for index, sub := range subdomainList {
		fmt.Printf("\r%d:%d", index+1, len(subdomainList))
		subdomain := fmt.Sprintf("%s.%s", sub, hostvar)
		request, err := http.NewRequest("GET", fmt.Sprintf("http://%s", ip), bytes.NewBuffer([]byte("")))
		request.Host = subdomain

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		resp, err := client.Do(request)

		if err != nil {
			continue
		}
		bodyLength, _ := io.Copy(ioutil.Discard, resp.Body)

		if resp.StatusCode == 404 ||
			isInList(strconv.Itoa(resp.StatusCode), hideStatusList) ||
			isInList(strconv.Itoa(int(bodyLength)), wordCountList) {
			continue
		}

		fmt.Printf("%c[2K\r%s\t%d\t%s\n", 27, colourStatusCode(strconv.Itoa(resp.StatusCode)), bodyLength, subdomain)
	}
}

func main() {
	var hostvar string   //host to look up
	var wordlist string  //wordlist to use
	var virtualmode bool //use virtual hosting (eg: hackthebox.eu)
	var ipvar string     //IP to use for virtual hosting
	var hsvar string
	var hcvar string

	flag.StringVar(&hostvar, "H", "", "base apex domain to enumerate agains, eg: example.com")
	flag.StringVar(&wordlist, "w", "", "path to wordlist to use, eg: /usr/share/subdomains.txt")
	flag.BoolVar(&virtualmode, "v", false, "enumerate subdomains by virtual hosts instead of DNS lookup")
	flag.StringVar(&ipvar, "ip", "", "IP used for virtual hosting (REQUIRED), eg: 127.0.0.1")
	flag.StringVar(&hsvar, "hs", "", "Ignore status codes, eg: 302,301,401. Status 404 is always ignored")
	flag.StringVar(&hcvar, "hc", "", "Hide responses that contain character length, eg: 2000,20043") //TODO

	flag.Parse()

	if len(hostvar) == 0 || len(wordlist) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	// print banner and echo back settings used
	println(BANNER)
	fmt.Println("host:", hostvar)
	fmt.Println("wordlist:", wordlist)
	if len(hsvar) > 0 {
		fmt.Println("ignore status:", hsvar)
	}
	if len(hcvar) > 0 {
		fmt.Println("ignore content length:", hcvar)
	}
	fmt.Println("====================================")

	subdomainList := readWordList(wordlist)

	if virtualmode {
		if len(ipvar) == 0 {
			os.Exit(1)
		}
		virtMode(hostvar, subdomainList, ipvar, hsvar, hcvar)
	} else {
		dnsMode(hostvar, subdomainList)
	}

	fmt.Println("\ndone")
}
