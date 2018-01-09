package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Debpackage struct {
	Name string
	Url  string
}

func FindAllPackage() []Debpackage {
	dres := make([]Debpackage, 0, 4000)
	doc, err := goquery.NewDocument("https://manpages.debian.org/contents-unstable.html")
	if err != nil {
		panic(err)
	}
	elem := doc.Find("div#content ul li a")
	elem.Each(func(arg0 int, arg1 *goquery.Selection) {
		var pkg Debpackage
		PackagePageLink, exists := arg1.Attr("href")
		if !exists {
			return
		}
		pkg.Url = PackagePageLink
		PackageName := arg1.Text()
		pkg.Name = PackageName
		dres = append(dres, pkg)
	})
	return dres
}

type DebpackageManEntry struct {
	Name    string
	Mantype string
	Lang    string
	Url     string
	Brief   string
	Pkg     Debpackage
	More    string
}

func GetManEntryByDebPackage(pkg Debpackage) []string {
	dres := make([]string, 0, 40)
	doc, err := goquery.NewDocument("https://manpages.debian.org" + pkg.Url)
	if err != nil {
		panic(err)
	}
	elem := doc.Find("div#content ul li a")
	elem.Each(func(arg0 int, arg1 *goquery.Selection) {
		url, exists := arg1.Attr("href")
		if !exists {
			return
		}
		dres = append(dres, url)
	})
	return dres
}

func GetManpagestructFromManurl(manpageurl string, pkg Debpackage) DebpackageManEntry {
	var dres DebpackageManEntry
	dres.Pkg = pkg
	dres.Url = manpageurl

	//url to field
	urls := strings.Split(manpageurl, "/")

	dataf := urls[3]

	pl := strings.Split(dataf, ".")

	dres.Name = pl[0]
	dres.Mantype = pl[1]
	dres.Lang = pl[2]
	GetManpageContent(&dres, manpageurl)
	return dres
}

func GetManpageContent(etr *DebpackageManEntry, url string) {
	doc, err := goquery.NewDocument("https://manpages.debian.org" + url)
	if err != nil {
		log.Println(err)
		return
	}
	ele := doc.Find(".manual-text").First().Contents()
	morecutoff := false
	h1counter := 0
	ele.Each(func(arg0 int, arg1 *goquery.Selection) {
		nodename := goquery.NodeName(arg1)
		switch nodename {
		case "h1":
			h1counter++
			if h1counter >= 4 {
				morecutoff = true
			}
		case "dl":
			morecutoff = true
		}
		stradd := arg1.Text()
		if !morecutoff {
			etr.Brief += stradd
		}
		etr.More += stradd
	})

}
