package main

import (
	"Backend/internal/model/local"
	"Backend/internal/router"
)

func main() {
	// 引用连接数据库
	local.InitDb()
	router.InitRouter()
	//errcode, err := Basic.UpdateRDS("1190073115559253")
	//fmt.Println(errcode, " ", err)
	//matched, err := regexp.MatchString("^[0-9a-zA-Z_-]{20}$", "rm-uf61h9w1satmik2zq")
	//fmt.Println(matched, err) //true <nil>
	/*
		input := &local.Rules{
		UUID:           "8cc48595-b0fd-473d-9681-8a50ff13a1f0",
		RuleName:       "mobile",
		PreferColumn:   "mobile,phone",
		RegExpress:     "^(13[0-9]|14[579]|15[0-3,5-9]|16[6]|17[0135678]|18[0-9]|19[89])\\d{8}$",
		CategoryUUID:   "237b0bee-90cf-465c-99b9-b76e350e002b",
		SensitiveGrade: "S1",
		Comments:       "",
		}
		_, _ = Basic.UpdateRules(*input)
		errcode, err := Basic.UpdateDBDetails("pgm-uf64xc9s753fcgz6")
		errcode, err, RDSClient := aliyunSDK.CreateRDSClient("1709821521553093")
		Basic.InitRDSAccount("rm-uf6m4dn81fm6kz30h", "MySQL", "UJH*xyr2nhg7bhp9cam", RDSClient)
		aliyunSDK.GrantAccountPrivilege("rm-uf6m4dn81fm6kz30h", "icontact", RDSClient)
		aliyunSDK.DeleteAccount("rm-uf6m4dn81fm6kz30h", RDSClient)
		errcode, err = aliyunSDK.ModifySecurityIps("rm-uf6m4dn81fm6kz30h", "Delete", RDSClient)
		errcode, err = aliyunSDK.ReleaseInstancePublicConnection("pgm-uf64xc9s753fcgz6", "v4w44tnl7p.pg.rds.aliyuncs.com", RDSClient)

		matched, err := regexp.MatchString("^1(3\\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\\d|9[0-35-9])\\d{8}$", "asdfasdfas13524983768asdf")
		fmt.Println(matched, err) //true <nil>
		matched, err = regexp.MatchString("^1(3\\d|4[5-9]|5[0-35-9]|6[2567]|7[0-8]|8\\d|9[0-35-9])\\d{8}$", "35524983268")
		fmt.Println(matched, err) //false <nil>
	*/
}
