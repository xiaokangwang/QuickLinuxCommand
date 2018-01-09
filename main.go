package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func main() {
	//res := FindAllPackage()
	//p := Debpackage{Url: "/unstable/zutils/index.html", Name: "coreutils"}
	//r2 := GetManpagestructFromManurl("/unstable/zutils/zcat.1.en.html", p)
	//spew.Dump(r2)
	Elab()
}

func CrawMans() {

	f, err := os.Create("manpages.json")
	if err != nil {
		panic(err)
	}

	outp := json.NewEncoder(f)

	//result := make([]DebpackageManEntry, 0, 2000)
	pkgs := FindAllPackage()
	for ei, e := range pkgs {
		fmt.Printf("Fetching package %v, %v/%v\n", e.Name, ei, len(pkgs))
		//Check if skip
		switch e.Name {
		case "coreutils":
		case "manpages":
		case "manpages-zh":
		case "curl":
		case "wget":
		case "lftp":
		case "net-tools":
		case "gnupg2":
		case "tar":
		case "gzip":
		case "openssl":
		default:
			if strings.HasPrefix(e.Name, "systemd") {
				goto shouldfetch
			}
			if strings.HasPrefix(e.Name, "openssh") {
				goto shouldfetch
			}
			continue
		}
	shouldfetch:
		r2 := GetManEntryByDebPackage(e)
		for crurli, crurl := range r2 {
			fmt.Printf("Fetching man %v, %v/%v\n", crurl, crurli, len(r2))
			s := GetManpagestructFromManurl(crurl, e)
			outp.Encode(s)
		}
	}

	f.Close()

}

func Elab() {
	f, err := os.Open("manpages.json")
	if err != nil {
		panic(err)
	}
	var root map[string]Leaf
	root = make(map[string]Leaf)
	for i := 0; i <= 3; i++ {
		jsond := json.NewDecoder(f)
		for jsond.More() {
			var et DebpackageManEntry
			jsond.Decode(&et)
			FindEntery(i, et, root)
		}
		f.Seek(0, 0)
	}
	f.Close()
	r := Construct(root)
	js := new(WebBotPT)
	js.Contain = r
	js.Include = make([]string, 0)
	f, err = os.Create("jsonout.json")
	if err != nil {
		panic(err)
	}
	e := json.NewEncoder(f)
	e.Encode(js)
	f.Close()

}
