package controllers

import (
	"encoding/json"
	"errors"
	"pet/models"
	"strconv"
	"strings"
	"time"
	"web"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

// oprations for Likes
type LikesController struct {
	beego.Controller
}

func (this *LikesController) URLMapping() {
	this.Mapping("Post", this.Post)
	this.Mapping("GetOne", this.GetOne)
	this.Mapping("GetAll", this.GetAll)
	this.Mapping("Put", this.Put)
	this.Mapping("Delete", this.Delete)
}

// @Title Post
// @Description create Likes
// @Param	body		body 	models.Likes	true		"body for Likes content"
// @Success 200 {int} models.Likes.Id
// @Failure 403 body is empty
// @router / [post]
func (this *LikesController) Post() {
	var v models.Likes
	var err error

	valid := validation.Validation{}
	this.ParseForm(&v)
	photoIdStr := this.GetString("photo_id")
	photoId, _ := strconv.Atoi(photoIdStr)

	passed, _ := valid.Valid(&v)
	if !passed {
		outPut := helper.Reponse(1, nil, valid.Errors[0].Key+" "+valid.Errors[0].Message)
		this.Data["json"] = outPut
	} else {
		v.TargetId, err = models.GetPhotosById(photoId)
		if err != nil {
			outPut := helper.Reponse(1, nil, err.Error())
			this.Data["json"] = outPut
			this.ServeJson()
			return
		}
		v.CreatedAt = time.Now()
		v.UpdatedAt = time.Now()
		userSession := this.GetSession("user").(models.Users)
		v.UserId = &userSession
		v.TargetId.UserId = &userSession

		if id, err := models.AddLikes(&v); err == nil {
			v.Id = int(id)
			outPut := helper.Reponse(0, v, "创建成功")
			this.Data["json"] = outPut
		} else {
			outPut := helper.Reponse(1, nil, err.Error())
			this.Data["json"] = outPut
		}
	}

	this.ServeJson()
}

// @Title Get
// @Description get Likes by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Likes
// @Failure 403 :id is empty
// @router /:id [get]
func (this *LikesController) GetOne() {
	idStr := this.Ctx.Input.Params[":id"]
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetLikesById(id)
	if err != nil {
		this.Data["json"] = err.Error()
	} else {
		this.Data["json"] = v
	}
	this.ServeJson()
}

// @Title Get All
// @Description get Likes
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Likes
// @Failure 403
// @router / [get]
func (this *LikesController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query map[string]string = make(map[string]string)
	var limit int64 = 10
	var offset int64 = 0

	// fields: col1,col2,entity.col3
	if v := this.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := this.GetInt("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := this.GetInt("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := this.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := this.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := this.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.Split(cond, ":")
			if len(kv) != 2 {
				this.Data["json"] = errors.New("Error: invalid query key/value pair")
				this.ServeJson()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllLikes(query, fields, sortby, order, offset, limit)
	if err != nil {
		this.Data["json"] = err.Error()
	} else {
		this.Data["json"] = l
	}
	this.ServeJson()
}

// @Title Update
// @Description update the Likes
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Likes	true		"body for Likes content"
// @Success 200 {object} models.Likes
// @Failure 403 :id is not int
// @router /:id [put]
func (this *LikesController) Put() {
	idStr := this.Ctx.Input.Params[":id"]
	id, _ := strconv.Atoi(idStr)
	v := models.Likes{Id: id}
	json.Unmarshal(this.Ctx.Input.RequestBody, &v)
	if err := models.UpdateLikesById(&v); err == nil {
		this.Data["json"] = "OK"
	} else {
		this.Data["json"] = err.Error()
	}
	this.ServeJson()
}

// @Title Delete
// @Description delete the Likes
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (this *LikesController) Delete() {
	var v models.Likes
	idStr := this.Ctx.Input.Params[":id"]
	id, _ := strconv.Atoi(idStr)

	valid := validation.Validation{}
	this.ParseForm(&v)
	passed, _ := valid.Valid(&v)
	if !passed {
		outPut := helper.Reponse(1, nil, valid.Errors[0].Key+" "+valid.Errors[0].Message)
		this.Data["json"] = outPut
	} else {
		userSession := this.GetSession("user").(models.Users)
		if num, err := models.DeleteLikedPhoto(userSession.Id, id); err == nil {
			outPut := helper.Reponse(0, num, "删除成功")
			this.Data["json"] = outPut
		} else {
			outPut := helper.Reponse(1, nil, err.Error())
			this.Data["json"] = outPut
		}
	}
	this.ServeJson()
}
