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
	// 腾晖
	//info.CorpId = "ww4a3407dd4c27e725"
	//info.CorpSecret = "cOO2IqQXZGwSLYHdaRoMuwo0Bhk4bvrqBH4httj_Vv8"
	//info.AgentId = 1000003
	// 个人
	info.CorpId = "ww06bd2f666a354c94"
	info.CorpSecret = "UINmPVLShl4xDGs1kWfX8dzipbSf45SE2GyVDHWf2ZY"
	info.AgentId = 1000002
}

func Tenghui(info *CorpWeChatInfo) *CorpWeChatInfo {
	info.TokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	info.FileURL = "https://qyapi.weixin.qq.com/cgi-bin/media/upload"
	info.StringURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	info.CorpId = "ww4a3407dd4c27e725"
	info.CorpSecret = "cOO2IqQXZGwSLYHdaRoMuwo0Bhk4bvrqBH4httj_Vv8"
	info.AgentId = 1000003

	return info
}

func Acke(info *CorpWeChatInfo) *CorpWeChatInfo {
	info.TokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	info.FileURL = "https://qyapi.weixin.qq.com/cgi-bin/media/upload"
	info.StringURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send"
	info.CorpId = "ww06bd2f666a354c94"
	info.CorpSecret = "UINmPVLShl4xDGs1kWfX8dzipbSf45SE2GyVDHWf2ZY"
	info.AgentId = 1000002

	return info
}

/*
企业信息工厂
  */
func CreateCompany(corpName string) *CorpWeChatInfo {
	switch corpName {
	case "tenghui":
		info := &CorpWeChatInfo{}
		return Tenghui(info)
	default:
		info := &CorpWeChatInfo{}
		return Acke(info)

	}
}
