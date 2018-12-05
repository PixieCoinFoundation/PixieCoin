package dao

import (
	"appcfg"
	"constants"
	"encoding/json"
	"fmt"
	"github.com/go-xorm"
	. "model"
	. "pixie_contract/api_specification"
)

var engine *xorm.Engine

func init() {
	if appcfg.GetServerType() != constants.SERVER_TYPE_GM {
		return
	}

	var err error
	dbUsername := appcfg.GetString("db_username", "")
	dbPassword := appcfg.GetString("db_password", "")
	dbIPPort := appcfg.GetString("db_ipport", "")
	dbName := appcfg.GetString("db_name", "")
	if engine, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?readTimeout=15s&writeTimeout=15s&timeout=10s&charset=utf8mb4", dbUsername, dbPassword, dbIPPort, dbName)); err != nil {
		panic(err)
	}
	engine.SetMaxIdleConns(10)
}

func GetEngine() *xorm.Engine {
	return engine
}

func GetPlayerByUsername(username string) (err error, status Status) {
	status.Username = username
	var exist bool

	if exist, err = engine.Get(&status); err != nil {
		return
	} else if !exist {
		err = constants.PlayerNotExist
		return
	}
	return
}

func GetPlayerByNickname(nickname string) (err error, status Status) {
	status.Nickname = nickname
	var exist bool
	if exist, err = engine.Get(&status); err != nil {
		return
	} else if !exist {
		err = constants.PlayerNotExist
		return
	}
	return
}

func InsertOfficialPaper(opaper GMOfficialPaper) (err error) {
	if _, err = engine.Table("pixie_official_paper").Insert(opaper); err != nil {
		fmt.Println(err)
		return err
	}
	return
}

func GetOfficialPaper(offset, pageSize int) (err error, count int64, Res []GMOfficialPaper) {
	if err = engine.Table("pixie_official_paper").Where("part_type!=?", PAPER_TYPE_BEIJING).Desc("id").Limit(pageSize, offset).Find(&Res); err != nil {
		return
	}
	count, _ = engine.Table("pixie_official_paper").Where("part_type!=?", PAPER_TYPE_BEIJING).Count()
	return
}

func GetOfficialBGPaper(offset, pageSize int) (err error, count int64, Res []GMOfficialPaper) {
	if err = engine.Table("pixie_official_paper").Where("part_type=?", PAPER_TYPE_BEIJING).Desc("id").Limit(pageSize, offset).Find(&Res); err != nil {
		return
	}
	count, _ = engine.Table("pixie_official_paper").Where("part_type=?", PAPER_TYPE_BEIJING).Count()
	return
}

func GetOneOfficialPaper(id int64) (err error, paper GMOfficialPaper) {
	if _, err = engine.Table("pixie_official_paper").Id(id).Get(&paper); err != nil {
		return
	}
	return
}

func GetPlayerClothes(status int) (err error, list []DesignPaper) {
	if err = engine.Table("pixie_paper").Where("part_type!=? and status=?", PAPER_TYPE_BEIJING, status).Find(&list); err != nil {
		return
	}
	return
}

func GetPlayerBG(status int) (err error, list []DesignPaper) {
	if err = engine.Table("pixie_paper").Where("part_type=? and status=?", PAPER_TYPE_BEIJING, status).Find(&list); err != nil {
		return
	}
	return
}

func GetPlayerClothesExist(clothesIDList []string) (err error, list []DesignPaperGM) {
	if err = engine.Table("pixie_paper").In("id", clothesIDList).Find(&list); err != nil {
		return
	}
	return
}

func GetSaleClothes(offset, size int) (err error, count int64, list []SalePaperGM) {
	if err = engine.Table("pixie_occupy").Limit(size, offset).Find(&list); err != nil {
		return
	}
	count, _ = engine.Table("pixie_occupy").Count()
	return
}

func GetOnePlayerPaper(paperID int64) (err error, paper DesignPaperGM) {
	paper.PaperID = paperID
	if _, err = engine.Table("pixie_paper").Get(&paper); err != nil {
		return
	}
	return
}

func DeleteOfficialPaper(id int64) (err error) {
	var paper GMOfficialPaper
	paper.ID = id
	if _, err = engine.Table("pixie_official_paper").Id(id).Delete(&paper); err != nil {
		return
	}
	return
}

func GetVersionInfo() (err error, v Version) {
	if _, err = engine.Table("version").Get(&v); err != nil {
		return
	}
	return
}

func UpdateVersion(v Version) (err error) {
	if _, err = engine.Table("version").Update(v); err != nil {
		return
	}
	return
}

func GetMaintenanceList() (err error, Res []Maintenance) {
	if err = engine.Table("gf_maintain_job").Find(&Res); err != nil {
		return
	}
	return
}

func AddMaintenance(m Maintenance) (effect int64, err error) {
	if effect, err = engine.Table("gf_maintain_job").Insert(m); err != nil {
		return
	}
	return
}

func GetValByKeyName(keyName string) (exist bool, w ConfigInfo, err error) {
	w.Key = keyName
	if exist, err = engine.Table("gf_config").Get(&w); err != nil {
		return
	}
	return
}

func AddKeyVal(s string, existKey bool, keyName string) (effect int64, err error) {
	var w ConfigInfo
	var keyInsert string
	w.Value = s
	w.Key = keyName
	keyInsert = keyName

	if existKey {
		if effect, err = engine.Table("gf_config").Where("`key`=?", keyInsert).Update(w); err != nil {
			return
		}
	} else {
		if effect, err = engine.Table("gf_config").Insert(w); err != nil {
			return
		}
	}
	return
}

func DelWhiteAccount(list string) (effect int64, err error) {
	var w ConfigInfo
	w.Key = constants.MAINTAIN_WHITE_LIST_KEY
	w.Value = list
	if effect, err = engine.Table("gf_config").Update(w); err != nil {
		return
	}
	return
}

func AddAdmin(agm AdminGM) (err error) {
	if _, err = engine.Table("gf_admin").Insert(agm); err != nil {
		return
	}
	return
}

func ChangeRole(id int64, roleID int64) (err error) {
	var agm AdminGM
	agm.ID = id
	agm.Role = roleID
	if _, err = engine.Table("gf_admin").Where("id=?", id).Cols("role").Update(&agm); err != nil {
		return
	}
	return
}

func GetOneAdmin(id int64) (err error, agm AdminGM) {
	if _, err = engine.Table("gf_admin").Where("id=?", id).Get(&agm); err != nil {
		return
	}
	return
}

func AdminExist(username, password string) (exist bool, agm AdminGM) {
	agm.Username = username
	agm.Password = password
	var err error
	if exist, err = engine.Table("gf_admin").Get(&agm); err != nil {
		return
	}
	return
}

func GetAdminList() (err error, list []AdminGM) {
	if err = engine.Table("gf_admin").Find(&list); err != nil {
		return
	}
	return
}

func DeleteAdmin(id int64) (err error) {
	var temp AdminGM
	temp.ID = id
	if _, err = engine.Table("gf_admin").Id(id).Delete(temp); err != nil {
		return
	}
	return
}

func AddRole(roleName string, desc string) (id int64, err error) {
	input := new(AdminRole)
	input.RoleName = roleName
	input.Desc = desc
	if _, err = engine.Table("gf_admin_role").InsertOne(input); err != nil {
		return
	}
	output := new(AdminRole)
	if _, err = engine.Table("gf_admin_role").Where("role_name=?", roleName).Get(output); err != nil {
		return
	}
	return output.ID, nil
}

func AddRoleApiRriv(exist bool, roleID int64, apiID []int64) (err error) {
	input := new(RoleApiPriv)
	input.RoleID = roleID
	data, _ := json.Marshal(apiID)
	input.MenuIDList = string(data)
	if !exist {
		if _, err = engine.Table("gf_menu_role_priv").Insert(input); err != nil {
			return
		}
	} else {
		if _, err = engine.Table("gf_menu_role_priv").Where("role_id=?", roleID).Cols("menu_id").Update(input); err != nil {
			return
		}
	}

	return
}

func GetMenuRolePriv(roldID int64) (exist bool, err error, res RoleApiPriv) {
	res.RoleID = roldID
	if exist, err = engine.Table("gf_menu_role_priv").Where("role_id=?", roldID).Get(&res); err != nil {
		return
	}
	return
}

func GetRoleList() (list []AdminRole, err error) {
	if err = engine.Table("gf_admin_role").Find(&list); err != nil {
		return
	}
	return
}

func GetAdminRolePriv(adminID int64) (err error, temp AdminRolePriv) {
	temp.AdminID = adminID
	if _, err = engine.Table("gf_admin_role_priv").Get(&temp); err != nil {
		return
	}
	return
}

//预审列表
func GetVerifyList(status PaperStatus, isClothes bool) (err error, list []DesignPaperGM) {
	if isClothes {
		if err = engine.Table("pixie_paper").Where("status=? AND part_type!=?", status, PAPER_TYPE_BEIJING).Find(&list); err != nil {
			return
		}
	} else {
		if err = engine.Table("pixie_paper").Where("status=? AND part_type=?", status, PAPER_TYPE_BEIJING).Find(&list); err != nil {
			return
		}
	}
	return
}

//获取一个审核作品
func GetVerifyOne(id int64) (error, DesignPaperGM) {
	var err error
	var exist bool
	var dp DesignPaperGM
	if exist, err = engine.Table("pixie_paper").Where("id=?", id).Get(&dp); err != nil {
		return err, DesignPaperGM{}
	} else if !exist {
		return nil, DesignPaperGM{}
	}
	return nil, dp
}

func PassVerify(id int64, hiddenTag string) (err error) {
	var dp DesignPaperGM
	dp.Status = int(PAPER_STATUS_QUEUE)
	dp.STag = hiddenTag
	if _, err = engine.Table("pixie_paper").Cols("status", "stag").Where("id=? AND status=?", id, PAPER_STATUS_ADMIN_QUEUE).Update(&dp); err != nil {
		return err
	}
	return
}

func RejectVerify(id int64, rejectReason string) (err error) {
	var dp DesignPaperGM
	dp.Status = int(PAPER_STATUS_FAIL)
	dp.Reason = rejectReason
	if _, err = engine.Table("pixie_paper").Cols("status", "reason").Where("id=? AND status=?", id, PAPER_STATUS_ADMIN_QUEUE).Update(&dp); err != nil {
		return err
	}
	return
}

func GetPoolStatusClothes(status int, limit, size int64) (err error, list []SubjectPool) {
	if err = engine.Table("pixie_recommend_pool").Where("status=?", status).Limit(int(size), int(limit)).Find(&list); err != nil {
		return
	}
	return
}

func FindSalePaper(clothesIDList []string) (err error, list []SalePaperGM) {
	if err = engine.Table("pixie_occupy").In("paper_id", clothesIDList).Find(&list); err != nil {
		return
	}
	return
}

func FindOffPaper(clothesIDList []string) (err error, list []OffPaperGM) {
	if err = engine.Table("pixie_official_paper").In("id", clothesIDList).Find(&list); err != nil {
		return
	}
	return
}
