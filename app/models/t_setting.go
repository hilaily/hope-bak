package models

import (
	"hope/app/support"

	"github.com/revel/revel"
)

var SiteInfos *SiteInfo

//Setting model
type Setting struct {
	Key   string `xorm:"not null pk VARCHAR(20)"`
	Value string `xorm:"not null TEXT"`
}

//SiteInfo model
type SiteInfo struct {
	Title      string
	SubTitle   string
	Url        string
	Seo        string
	Reg        string
	Foot       string
	Statistics string
	Status     string
	Comment    string
}

//Loaded setting info to cache.
func LoadCache() {
	set := new(Setting)
	res, err := set.FindAll()
	if err != "" {
		revel.ERROR.Printf("Loaded setting info to cache error: %v", err)
		return
	}
	if len(res) > 0 {
		for i := 0; i < len(res); i++ {
			s := res[0]
			support.Cache.Set(s.Key, s.Value, 0)
		}
	}
}

//find all setting info.
func (s *Setting) FindAll() ([]Setting, string) {
	set := make([]Setting, 0)
	err := support.Xorm.Find(&set)
	if err != nil {
		return set, err.Error()
	}
	return set, ""
}

//Get setting value for key.
func (s *Setting) Get() (string, error) {
	res, err := support.Cache.Get(s.Key).Result()
	if res == "" {
		set := new(Setting)
		has, err := support.Xorm.Where("`key` = ?", s.Key).Get(set)
		if has && err == nil {
			return set.Value, err
		}
	}
	return res, err
}

//Put setting K,V in cache and db.
func (s *Setting) Put() (bool, error) {
	has, err := support.Xorm.InsertOne(s)
	if err == nil && has > 0 {
		support.Cache.Set(s.Key, s.Value, 0)
	}
	return has > 0, err
}

//Update value for key.
func (s *Setting) Update() (bool, error) {
	set := new(Setting)
	set.Value = s.Value
	has, err := support.Xorm.Where("`key` = ?", s.Key).Update(set)
	if err == nil && has > 0 {
		support.Cache.Del(s.Key)
		support.Cache.Set(s.Key, s.Value, 0)
	}
	return has > 0, err
}

//Query site setting info.
func (s *Setting) GetSiteInfo() (*SiteInfo, string) {
	site := new(SiteInfo)
	data, err := s.FindAll()
	if err != "" {
		return site, err
	}
	if len(data) > 0 {
		for _, tmp := range data {
			switch tmp.Key {
			case "site-foot":
				site.Foot = tmp.Value
			case "site-reg":
				site.Reg = tmp.Value
			case "site-seo":
				site.Seo = tmp.Value
			case "site-status":
				site.Status = tmp.Value
			case "site-subtitle":
				site.SubTitle = tmp.Value
			case "site-title":
				site.Title = tmp.Value
			case "site-url":
				site.Url = tmp.Value
			case "site-statistics":
				site.Statistics = tmp.Value
			case "site-comment":
				site.Comment = tmp.Value
			}
		}
	}
	SiteInfos = site
	return site, err
}

//Add new setting info or update
func (s *Setting) InsertAndModify(key, value string) error {
	set := new(Setting)
	set.Key = key
	res, err := set.Get()
	if err == nil && res != "" {
		set.Key = key
		set.Value = value
		has, err := set.Update()
		revel.INFO.Printf("NewSiteInfo::Put has: %v,error: %v", has, err)
		if err != nil {
			return err
		}
	} else {
		set.Key = key
		set.Value = value
		has, err := set.Put()
		revel.INFO.Printf("NewSiteInfo::Put has: %v,error: %v", has, err)
		if err != nil {
			return err
		}
	}
	return err
}

//Add new site info
func (s *Setting) NewSiteInfo(title, subtitle, url, seo, reg, foot,
	statistics, status, comment string) error {

	keyArr := []string{"site-title", "site-subtitle", "site-url", "site-seo", "site-reg", "site-foot",
		"site-statistics", "site-status", "site-comment"}
	valueArr := []string{title, subtitle, url, seo, reg, foot, statistics, status, comment}

	var err error
	for k, v := range valueArr {
		if v != "" {
			err = s.InsertAndModify(keyArr[k], v)
			if err != nil {
				revel.ERROR.Printf("更新站点设置表错误，错误字段[%s],内容：%v\n", keyArr[k], err)
				return err
			}
		}
	}
	return err
}
