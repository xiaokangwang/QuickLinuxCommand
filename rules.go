package main

type WebBotPT struct {
	Include []string
	Contain [](ReplyC)
}

type ReplyC struct {
	MatchingCond    []string
	MatchingCondTag string
	OutPut          string
	Tags            string
}

func (rc *ReplyC) AddCond(cc string) *ReplyC {
	rc.MatchingCond = append(rc.MatchingCond, cc)
	return rc
}
func (rc *ReplyC) SetCondT(cc string) *ReplyC {
	rc.MatchingCondTag = cc
	return rc
}
func (rc *ReplyC) SetOutputT(cc string) *ReplyC {
	rc.Tags = cc
	return rc
}
func (rc *ReplyC) SetOutput(cc string) *ReplyC {
	rc.OutPut = cc
	return rc
}
func NewChat() *ReplyC {
	return &ReplyC{}
}
