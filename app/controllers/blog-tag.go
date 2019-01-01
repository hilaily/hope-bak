package controllers

import (
	"hope/app/models"
	"strconv"

	"github.com/revel/revel"
)

// BlogTag controller
type BlogTag struct {
	*revel.Controller
}

//
func (b *BlogTag) Index(ident string) revel.Result {
	tag := new(models.Tag)
	tag, err := tag.GetByIdent(ident)
	if err != nil {
		revel.ERROR.Panic("wrong")
	}
	blogs := tag.FindBlogByTag("")
	b.ViewArgs["flag"] = "tag"
	b.ViewArgs["tag"] = tag
	b.ViewArgs["blogs"] = blogs
	return b.RenderTemplate("Main/Blog4Search.html")
}

// GetAllTags to find all tags
// 获取所有的标签
func (b *BlogTag) GetAllTags() revel.Result {
	tagModel := new(models.Tag)
	tags, err := tagModel.ListAll()
	if err != nil {
		revel.ERROR.Println("find all tags error: ", err)
	}
	return b.RenderJSON(tags)
}

// QueryTags to Search for tag
// 根据用户输入的单词匹配 tag
func (b *BlogTag) QueryTags(t string) revel.Result {
	tag := new(models.Tag)
	res, err := tag.QueryTags(t)
	if err != nil {
		return b.RenderJSON(&ResultJson{Success: false, Msg: err.Error(), Data: ""})
	}
	resMap := make(map[int64]string, 0)
	for _, v := range res {
		id, err := strconv.Atoi(string(v["id"]))
		if err == nil {
			resMap[int64(id)] = string(v["name"])
		}
	}
	return b.RenderJSON(&ResultJson{Success: true, Msg: "", Data: resMap})
}
