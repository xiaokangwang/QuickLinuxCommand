package main

import (
	"fmt"
	"regexp"
)

func FindEntery(scan int, scanning DebpackageManEntry, root map[string]Leaf) {
	switch scan {
	case 0:
		if scanning.Mantype != "1" {
			return
		}
		if scanning.Lang != "en" {
			return
		}
	case 1:
		if scanning.Mantype == "1" {
			return
		}
		if scanning.Lang != "en" {
			return
		}
	case 2:
		if scanning.Mantype != "1" {
			return
		}
		if scanning.Lang == "en" {
			return
		}
	case 3:
		if scanning.Mantype == "1" {
			return
		}
		if scanning.Lang == "en" {
			return
		}
	}
	//Is there a key?
	rootLeaf, ok := root[scanning.Name]
	if !ok {
		if scan >= 2 {
			return
		}
		root[scanning.Name] = Leaf{Entry: scanning}
	} else {
		switch scan {
		case 1:
			if rootLeaf.Type == nil {
				rootLeaf.Type = make(map[string]Leaf)
			}
			rootLeaf.Type[scanning.Mantype] = Leaf{Entry: scanning}
		case 2:
			if rootLeaf.Lang == nil {
				rootLeaf.Lang = make(map[string]Leaf)
			}
			if scanning.Lang == "zh_CN" {
				rootLeaf.Lang[scanning.Lang] = Leaf{Entry: scanning}
			}
		case 3:
			typ, ok3 := rootLeaf.Type[scanning.Mantype]
			if !ok3 {
				return
			}
			if scanning.Lang == "zh_CN" {
				if typ.Lang == nil {
					typ.Lang = make(map[string]Leaf)
				}
				typ.Lang[scanning.Lang] = Leaf{Entry: scanning}
			}
		}
	}

}

type Leaf struct {
	Entry DebpackageManEntry
	Lang  map[string]Leaf
	Type  map[string]Leaf
}

func RegxEscape(in string) string {
	var desret string
	exp := regexp.MustCompile(`[-[\]{}()*+?.,\\^$|#]`)
	desret = exp.ReplaceAllString(in, "\\$0")
	return desret
}

func Construct(root map[string]Leaf) [](ReplyC) {
	desret := make([](ReplyC), 0, 5000)
	for _, rtele := range root {
		briefTag := fmt.Sprintf("MAN%v00", rtele.Entry.Name)

		if rtele.Lang != nil {
			zh := rtele.Lang["zh_CN"]

			typBriefTag := fmt.Sprintf("MAN%vLANG%v00", zh.Entry.Name, zh.Entry.Lang)
			typec := NewChat().
				AddCond("(?i)zh").
				SetOutput(RegxEscape(fmt.Sprintf("%v", zh.Entry.Brief))).
				SetOutputT(typBriefTag)
			typecm := GenMore(zh.Entry, typBriefTag, fmt.Sprintf("MAN%vLANG%vMORE00", zh.Entry.Name, zh.Entry.Lang))
			desret = append(desret, *typec, typecm)
		}

		if rtele.Type != nil {
			for _, currtyp := range rtele.Type {
				typBriefTag := fmt.Sprintf("MAN%vTYPE%v00", currtyp.Entry.Name, currtyp.Type)

				if currtyp.Lang != nil {
					zh := currtyp.Lang["zh_CN"]

					typBriefTagL := fmt.Sprintf("MAN%vTYPE%vLANG%v00", zh.Entry.Name, zh.Entry.Mantype, zh.Entry.Lang)
					typec := NewChat().
						AddCond("(?i)zh").
						SetOutput(RegxEscape(FormartBrief(zh))).
						SetOutputT(typBriefTagL)
					typecm := GenMore(zh.Entry, typBriefTagL, fmt.Sprintf("MAN%vTYPE%vLANG%vMORE00", zh.Entry.Name, zh.Entry.Mantype, zh.Entry.Lang))
					desret = append(desret, *typec, typecm)
				}

				typec := NewChat().
					AddCond(RegxEscape(currtyp.Entry.Mantype)).
					SetOutput(RegxEscape(FormartBrief(currtyp))).
					SetOutputT(typBriefTag)
				typecm := GenMore(currtyp.Entry, typBriefTag, fmt.Sprintf("MAN%vTYPE%vMORE00", rtele.Entry.Name, currtyp.Type))
				desret = append(desret, *typec, typecm)
			}
		}

		brief := NewChat().
			AddCond("man").
			AddCond(" " + RegxEscape(rtele.Entry.Name) + "$").
			SetOutput(RegxEscape(FormartBrief(rtele))).
			SetOutputT(briefTag)
		more :=
			GenMore(rtele.Entry, briefTag, fmt.Sprintf("MAN%vMORE00", rtele.Entry.Name))
		desret = append(desret, *brief, more)
	}
	return desret
}

func GenMore(Entry DebpackageManEntry, RespondTag string, moreTag string) ReplyC {
	return *NewChat().
		AddCond("(?i)more").
		SetCondT(RegxEscape(RespondTag)).
		SetOutput(RegxEscape(fmt.Sprintf("%v", Entry.More))).
		SetOutputT(moreTag)
}

func FormartBrief(lea Leaf) string {
	ent := lea.Entry
	var ret string
	ret += fmt.Sprintf("下面是来自包 %v 的 %v 的手册的描述。\n", ent.Pkg.Name, ent.Name)
	ret += ent.Brief
	ret += `您可以回复 more 或者访问下放链接获取来获取完整的手册` + "\n"
	ret += "https://manpages.debian.org" + ent.Url + "\n"
	ret += fmt.Sprintf("这个内容来自 %v 章节\n", GetSecName(ent.Mantype))
	if lea.Lang != nil {
		ret += `您可以回复 zh 获取来获取本内容的中文版本\n`
	}
	if lea.Type != nil {
		for _, ctx := range lea.Type {
			ret += fmt.Sprintf("您也可以回复 %v 获取其在 %v 章节的内容。", ctx.Entry.Mantype, GetSecName(ctx.Entry.Mantype))
		}
	}
	return ret
}

func GetSecName(s string) string {
	return s
}
