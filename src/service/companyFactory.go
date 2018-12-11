package service

type CorpWeChatInfo struct {
	TokenURL   string
	StringURL  string
	FileURL    string
	CorpId     string
	CorpSecret string
	AgentId    int
}

type CompanyFactory interface {
	GetCompany(corpName string)
	InitWeChatInfo(info *CorpWeChatInfo)
}

func InitWeChatInfo(info *CorpWeChatInfo) {
	info.TokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	info.FileURL = "https://qyapi.weixin.qq.com/cgi-bin/media/upload"
	info.StringURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	//info.CorpId = "ww4a3407dd4c27e725"
	//info.CorpSecret = "cOO2IqQXZGwSLYHdaRoMuwo0Bhk4bvrqBH4httj_Vv8"
	//info.AgentId = 1000003
	info.CorpId = "ww06bd2f666a354c94"
	info.CorpSecret = "UINmPVLShl4xDGs1kWfX8dzipbSf45SE2GyVDHWf2ZY"
	info.AgentId = 1000002
}

func CreateCompany(companyFactory *CompanyFactory, corpName string) *CorpWeChatInfo {
	// TODO 根据生产的工厂决定CorpWeChatInfo的属性
	switch corpName {
	case "tenghui":
		info := &CorpWeChatInfo{}
		InitWeChatInfo(info)
		return info
	default:
		info := &CorpWeChatInfo{}
		InitWeChatInfo(info)
		return info
	}
}
