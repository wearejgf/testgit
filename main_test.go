package main

import (
	"apiserver/api/apibonus"
	"apiserver/api/apicococloud"
	"apiserver/api/apipartner"
	"apiserver/api/apiperformance"
	"apiserver/api/apirule"
	"apiserver/conf"
	"apiserver/feature/upgrade"
	"apiserver/lib/logger"
	"apiserver/sdk"
	"apiserver/types"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/BurntSushi/toml"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"os"
	"path"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	var config conf.RootConfig
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		log.Fatal("Error parse file config.toml: ", err)
	}

	sdk.Initialize(&config)

	exitCode := m.Run()

	sdk.Shutdown()

	os.Exit(exitCode)
}

func TestFinanceFeatureDownloadContract(t *testing.T) {
	secret := sdk.Pbkdf2Sha256(sdk.Config.Finance.FtDownloadContract.AppId, sdk.Config.Auth.JWT.Secret)
	assert.Equal(t, secret, sdk.Config.Finance.FtDownloadContract.AppSecret)
}

// 获取所有裂变团长
func TestGetFissionPartners(t *testing.T) {
	_, sitePartners := apipartner.GetFissionPartnerInfo(nil, nil)

	f := excelize.NewFile()
	sheet1 := f.NewSheet("Sheet1")
	f.SetActiveSheet(sheet1)

	f.SetCellValue("Sheet1", "A1", "序号")
	f.SetCellValue("Sheet1", "B1", "主站ID")
	f.SetCellValue("Sheet1", "C1", "子站ID")
	f.SetCellValue("Sheet1", "D1", "团长ID")
	f.SetCellValue("Sheet1", "E1", "团长级别")

	// 设置单元格的值
	i := 0
	for _, userIds := range sitePartners {
		for _, id := range userIds {
			n := sdk.PTree.ReferTree.GetNode(id)
			p := n.Data.(*types.Partner)

			// Create a new sheet.
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), i+1)
			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), p.MasterSiteId)
			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), p.SlaveSiteId)
			f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), p.UserId)
			f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), p.Level)

			i += 1
		}

	}

	if err := f.SaveAs("output.xlsx"); err != nil {
		log.Println(err)
	}

}

func TestPossibleB1(t *testing.T) {
	partners := make([]*types.Partner, 0)

	f := excelize.NewFile()
	sheet1 := f.NewSheet("Sheet1")
	f.SetActiveSheet(sheet1)

	f.SetCellValue("Sheet1", "A1", "序号")
	f.SetCellValue("Sheet1", "B1", "主站ID")
	f.SetCellValue("Sheet1", "C1", "子站ID")
	f.SetCellValue("Sheet1", "D1", "合格团长ID")

	// 设置单元格的值
	i := 0

	for _, p := range partners {
		// Create a new sheet.
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), i+1)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), p.MasterSiteId)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), p.SlaveSiteId)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), p.UserId)

		i += 1
	}

	if err := f.SaveAs("test.xlsx"); err != nil {
		log.Println(err)
	}
}

func TestUpgrade(t *testing.T) {
	upgrade.Upgrade()
}

func TestFamily(t *testing.T) {
	ids := sdk.PTree.GetFamilyMemberIds(96792338, false, true)
	t.Log(ids)
}

func TestSurpassed(t *testing.T) {
	allSurpassed := sdk.PTree.GetSurpassedIds(91198875, false)
	t.Log("all surpassed:", allSurpassed)

	directSurpassed := sdk.PTree.GetSurpassedIds(91198875, true)
	t.Log("direct surpassed:", directSurpassed)

	teamIds := sdk.PTree.GetSurpassedTeamIds(91198875, false)
	t.Log("teamIds:", teamIds)
}

func TestBQ(t *testing.T) {
	ids := sdk.PTree.GetMyBQIds(96199630, false, false)
	t.Log("All BQs::", ids)

	ignoreSurpassedIds := sdk.PTree.GetMyBQIds(96199630, false, true)
	t.Log("All BQs(ignore surpassed):", ignoreSurpassedIds)

	directIds := sdk.PTree.GetMyBQIds(96199630, true, false)
	t.Log("Direct BQs:", directIds)

	ignoreSurpassedDirectIds := sdk.PTree.GetMyBQIds(96199630, true, true)
	t.Log("Direct BQs(ignore surpassed):", ignoreSurpassedDirectIds)
}

func TestUpgradeBQ(t *testing.T) {
	pt := sdk.CreatePartnerTree()

	p := pt.GetPartner(81087842)
	p.IsQualified = 0

	c := apirule.NewBQChecker(pt)
	c.Process()
}

func TestUpgradeB1(t *testing.T) {
	pt := sdk.CreatePartnerTree()

	c := apirule.NewB1Checker(pt)
	c.Process()
}

func TestBred(t *testing.T) {
	bredByIds := sdk.PTree.GetBredByIds(88993627)
	t.Log("Bred by:", bredByIds)
}

func TestYYY(t *testing.T) {
	amount, _ := apibonus.GetReceivedBredAmountFrom(121612, 1588262400, 1590940799,
		[]int64{121613, 121614, 121615, 121613, 121616, 121623, 121617, 121628, 121624, 121619, 121629, 121614, 121618, 121621, 121615, 121622, 121620, 121630})
	t.Log("aaa:", amount)
}

func TestGetB2TeamGMV(t *testing.T) {
	fmt.Println("包括顾问自身销售额的GMV:")
	memberIds := sdk.PTree.GetTeamMemberIds(81504401, false, false)
	t.Log(apiperformance.GetGMVsByUserIds("2020-05-01", "2020-05-31", memberIds))

	fmt.Println("不包括顾问自身销售额的GMV:")
	memberIds = sdk.PTree.GetTeamMemberIds(81504401, false, true)
	t.Log(apiperformance.GetGMVsByUserIds("2020-05-01", "2020-05-31", memberIds))
}

func TestGetDirectMemberUpgradeTimes(t *testing.T) {
	allB1s := sdk.PTree.GetMyB1s(98318240, true, true)
	for _, b1 := range allB1s {
		t.Logf("UserId: %d, bq upgrade time: %s", b1.UserId, sdk.FromUnixTime(b1.UpgradeTime.ToBQ))
		t.Logf("UserId: %d, b1 upgrade time: %s", b1.UserId, sdk.FromUnixTime(b1.UpgradeTime.ToB1))
		t.Logf("UserId: %d, b2 upgrade time: %s", b1.UserId, sdk.FromUnixTime(b1.UpgradeTime.ToB2))
	}
}

func TestCheckUpgradeTimes(t *testing.T) {
	toChecks := []int64{81504401, 90107996}
	for _, uid := range toChecks {
		p := sdk.PTree.GetPartner(uid)
		if p == nil {
			continue
		}

		t.Logf("UserId: %d, bq upgrade time: %s", p.UserId, sdk.FromUnixTime(p.UpgradeTime.ToBQ))
		t.Logf("UserId: %d, b1 upgrade time: %s", p.UserId, sdk.FromUnixTime(p.UpgradeTime.ToB1))
		t.Logf("UserId: %d, b2 upgrade time: %s", p.UserId, sdk.FromUnixTime(p.UpgradeTime.ToB2))
	}
}

func TestNewAddedPartners(t *testing.T) {
	siteInfos, newPartnerInfo := apipartner.GetNewFissionPartners("2020-06-05", "2020-07-05", nil, nil)
	for siteId, siteName := range siteInfos {
		t.Logf("Site: %s, NewAdded: %d", siteName, len(newPartnerInfo[siteId]))
	}
}

func TestNewAddedGMV(t *testing.T) {
	siteInfos, siteGMVs := apiperformance.GetGMVsBySite("2020-06-05", "2020-07-05", nil, nil)
	_, siteNewGMVs := apiperformance.GetNewGMVs("2020-06-05", "2020-07-05", nil, nil)
	for siteId, siteName := range siteInfos {
		t.Logf("%s, %d, %d", siteName, siteNewGMVs[siteId], siteGMVs[siteId])
	}
}

func TestAllBonus(t *testing.T) {
	siteInfos, allPartnerInfo := apipartner.GetFissionPartnerInfo(nil, nil)
	_, newPartnerInfo := apipartner.GetNewFissionPartners("2020-06-05", "2020-07-05", nil, nil)
	for siteId, siteName := range siteInfos {
		var totalTraingBonus, totalTeamBonus, totalBredBonus int

		allUserIds := allPartnerInfo[siteId]
		recvAmounts, err := apibonus.BatchGetReceivedAmount(allUserIds, "2020-06-05", "2020-07-05")
		if err != nil {
			continue
		}
		totalTraingBonus = recvAmounts[apibonus.TypeTrainingBonus]
		totalTeamBonus = recvAmounts[apibonus.TypeTeamBonus]
		totalBredBonus = recvAmounts[apibonus.TypeBredBonus]

		t.Logf("SiteName: %s, Total: %d", siteName, len(allUserIds))
		t.Logf("%s, %d, %d, %d", siteName, totalTraingBonus, totalTeamBonus, totalBredBonus)

		var newTraingBonus, newTeamBonus, newBredBonus int
		newUserIds := newPartnerInfo[siteId]
		newRecvAmounts, err := apibonus.BatchGetReceivedAmount(newUserIds, "2020-06-05", "2020-07-05")
		if err != nil {
			continue
		}
		newTraingBonus = newRecvAmounts[apibonus.TypeTrainingBonus]
		newTeamBonus = newRecvAmounts[apibonus.TypeTeamBonus]
		newBredBonus = newRecvAmounts[apibonus.TypeBredBonus]

		t.Logf("SiteName: %s, New: %d", siteName, len(newUserIds))
		t.Logf("%s, %d, %d, %d", siteName, newTraingBonus, newTeamBonus, newBredBonus)
	}
}

func TestTimes(t *testing.T) {
	beginTime := sdk.GetMonthBeginTime(1)
	endTime := sdk.GetMonthEndTime(1)

	lastBeginTime := sdk.GetMonthBeginTimeSince(beginTime.Unix(), -1)
	lastEndTime := sdk.GetMonthEndTimeSince(endTime.Unix(), -1)

	t.Log("Begin Time: ", beginTime)
	t.Log("End Time: ", endTime)
	t.Log("Last Begin Time: ", lastBeginTime)
	t.Log("Last End Time: ", lastEndTime)
}

func TestCocoCloudSignUrl(t *testing.T) {
	// 椰云众包实名认证成功的回调地址
	// IMPORTANT: 因为服务器上根据protect_prefix来判断是否转发请求到api server,
	// 所以这里必须加上protect_prefix
	signRedirectUrl := path.Join(sdk.Config.Api.ProtectPrefix + sdk.Config.Api.PublicPrefix + sdk.Config.Wxmp.CocoCloudSignReturn)
	signReturnUrl := getReturnUrl(signRedirectUrl, 116342)

	logger.Debug("cocoCloudSignUrl", "sign return url:", signReturnUrl)
	signUrl, err := apicococloud.GetPersonSignUrl(signReturnUrl)

	if err != nil {
		logger.Error("cocoCloudSignUrl", "Error get cocoCloud sign url: ", err)
		return
	}

	fmt.Println("signUrl==>", signUrl)

}

func getReturnUrl(redirectUrl string, userId int64) string {
	var wxmpHost string
	switch sdk.Config.Env {
	case "prod":
		wxmpHost = sdk.Config.Wxmp.Prod.Host
	default:
		wxmpHost = sdk.Config.Wxmp.Test.Host
	}

	// user_id只能以path parameter的方式传过去
	returnUrl := fmt.Sprintf("%s%s/%v", wxmpHost, redirectUrl, userId)
	return returnUrl
}

func TestCallBack(t *testing.T) {
	cocoCloudSignReturn := apicococloud.CocoCloudSignReturn{
		BankCardNum:       "bankcardnum111111",
		IdCard:            "idcard11111111",
		Name:              "name111111",
		ResultDescription: "完成",
		Sign:              "DB431A84B810786050053604C2BC0F1BC5CC637DF5EE9A7F809BA80C7256881B",
		SignResult:        "2",
		SignTime:          "2020-08-05 15:34:00",
	}

	apicococloud.CheckSign(cocoCloudSignReturn)

	// 将签约信息插入fission_financial_info表
	signIngo := apicococloud.SignInfo{
		UserId:      121973,
		Name:        cocoCloudSignReturn.Name,
		IdNumber:    cocoCloudSignReturn.IdCard,
		BankAccount: cocoCloudSignReturn.BankCardNum,
		IsVerified:  1,
		IsSigned:    1,
		Created:     time.Now().Unix(),
	}

	fmt.Println(signIngo)

	fmt.Println("InsertOrUpdateSignInfo==>", apicococloud.InsertOrUpdateSignInfo(&signIngo))
}

func TestSignStatus(t *testing.T) {

	signInfo := apicococloud.SignInfo{
		Id:              0,
		UserId:          0,
		Name:            "杨博",
		IdNumber:        "43122819960101001x",
		BankName:        "",
		BankAccount:     "",
		TrdPartyId:      "",
		TrdPartyTransId: "",
		IsVerified:      1,
		IsSigned:        1,
		DownloadUrl:     "",
	}

	// 根据signinfo生成签名
	statusSign := apicococloud.GetStatusSign(signInfo)

	// 初始化请求参数
	statusSignParam := apicococloud.SignStatus{
		Name:   signInfo.Name,
		IdCard: signInfo.IdNumber,
		Sign:   statusSign,
	}

	logger.Info("statusSignParam==>", statusSignParam)

	// 请求接口
	statusSignResult, err := apicococloud.GetSignStatus(statusSignParam)
	if err != nil {
		logger.Error("cocoCloudSignUrl", "Error request user sign info")
		return
	}

	logger.Info("statusSignResult中文==>", statusSignResult)
}

// 测试数据库隔离设置
func testDBIsolation() {
	cocoCloudSignReturn := apicococloud.CocoCloudSignReturn{
		BankCardNum:       "666666666666666666",
		IdCard:            "888888888888888888",
		Name:              "黄金帅",
		ResultDescription: "完成",
		Sign:              "DB431A84B810786050053604C2BC0F1BC5CC637DF5EE9A7F809BA80C7256881B",
		SignResult:        "2",
		SignTime:          "2020-08-03 15:25:27",
	}

	// 将签约信息插入fission_financial_info表
	signIngo := apicococloud.SignInfo{
		UserId:      66666,
		Name:        cocoCloudSignReturn.Name,
		IdNumber:    cocoCloudSignReturn.IdCard,
		BankAccount: cocoCloudSignReturn.BankCardNum,
		IsVerified:  1,
		IsSigned:    1,
	}
	insertFlag := apicococloud.InsertOrUpdateSignInfo(&signIngo)
	if !insertFlag {
		logger.Error("fddSignCallback", "error insert sign info")
		return
	}
}

func TestDBIsolation(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			testDBIsolation()
		}()
	}
	time.Sleep(time.Second * 5)
}

func TestPTree(t *testing.T) {
	jsonBytes, err := json.Marshal(sdk.PTree.GetPartner(121973))
	if err != nil {
		fmt.Println(err)
	}
	logger.Info(string(jsonBytes))

}

func TestConfig(t *testing.T) {
	fmt.Println(sdk.Config.CocoCloud.CocoCloudApi.Test.RequestUrl)
	fmt.Println(sdk.Config.CocoCloud.CocoCloudApi.Test.AppSecret)
	fmt.Println(sdk.Config.CocoCloud.CocoCloudApi.Test.SignUrlPrefix)
	fmt.Println(sdk.Config.CocoCloud.CocoCloudApi.Test.MerchantId)
	fmt.Println(sdk.Config.CocoCloud.CocoCloudApi.Test.AppId)
	fmt.Println(sdk.Config.CocoCloud.CocoCloudApi.Test.AccountName)
}

func TestZxj(t *testing.T) {
	fmt.Println("demo updated by Zxj!")
}

func TestZxjStep(t *testing.T) {
	fmt.Println("Step3:demo updated by Zxj!")
}

func TestDevUpdate(t *testing.T) {
	fmt.Println("TestDevUpdate:dev branch updated by Zxj!")
	fmt.Println("dev branch second time update!")
}
      

func TestZxjBranchUpdate(t *testing.T) {"TestZxjBranchUpdate:zxj branch updated by zxj!"}
222222
111111
