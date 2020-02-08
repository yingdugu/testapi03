package controllers

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testapi03/models"

	"github.com/astaxie/beego"
)

// ProjectImageController operations for ProjectImage
type ProjectImageController struct {
	beego.Controller
}

// URLMapping ...
func (c *ProjectImageController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("Upload", c.Upload)
	c.Mapping("Getimage", c.Getimage)
}

// Post ...
// @Title Post
// @Description create ProjectImage
// @Param	body		body 	models.ProjectImage	true		"body for ProjectImage content"
// @Success 201 {int} models.ProjectImage
// @Failure 403 body is empty
// @router / [post]
func (c *ProjectImageController) Post() {
	var v models.ProjectImage
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddProjectImage(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = v
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get ProjectImage by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.ProjectImage
// @Failure 403 :id is empty
// @router /:id [get]
func (c *ProjectImageController) GetOne() {
	//idStr := c.Ctx.Input.Param(":id")
	idStr := c.GetString("name")
	//id, _ := strconv.Atoi(idStr)
	name := idStr
	v, err := models.GetProjectImageByName(name)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get ProjectImage
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.ProjectImage
// @Failure 403
// @router / [get]
func (c *ProjectImageController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllProjectImage(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the ProjectImage
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.ProjectImage	true		"body for ProjectImage content"
// @Success 200 {object} models.ProjectImage
// @Failure 403 :id is not int
// @router /:id [put]
func (c *ProjectImageController) Put() {
	//idStr := c.Ctx.Input.Param(":id")
	idStr := c.GetString("name")
	//id, _ := strconv.Atoi(idStr)
	id := idStr
	v := models.ProjectImage{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateProjectImageById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the ProjectImage
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *ProjectImageController) Delete() {
	//idStr := c.Ctx.Input.Param(":id")
	idStr := c.GetString("name")
	//id, _ := strconv.Atoi(idStr)
	id := idStr
	if err := models.DeleteProjectImage(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Upload ...
// @Title Upload
// @Description Upload the ProjectImage
// @Param	id		path 	string	true		"The id you want to Upload"
// @Success 200 {string} Upload success!
// @Failure 403 id is empty
// @router / :id [post]
func (c *ProjectImageController) Upload() {
	fmt.Println("Hello, World!")
	f, h, err := c.GetFile("image")
	result := make(map[string]interface{})
	img := ""
	id := UniqueId()
	exStrArr := strings.Split(h.Filename, ".")
	nameStr := exStrArr[len(exStrArr)-2]
	if err == nil {
		exStr := strings.ToLower(exStrArr[len(exStrArr)-1])
		if exStr != "jpg" && exStr != "png" && exStr != "gif" {
			result["code"] = 1
			result["message"] = "上传只能.jpg 或者png格式"
		}
		img = "static/upload/" + id + "." + exStr
		c.SaveToFile("image", img) // 保存位置在 static/upload, 没有文件夹要先创建
		result["code"] = 0
		result["message"] = img
	} else {
		result["code"] = 2
		result["message"] = "上传异常" + err.Error()
	}
	//defer f.Close()
	//c.Data["json"] = result
	// 定义数据库地址数据存储
	data := models.ProjectImage{}
	data.Name = nameStr
	data.Id = id
	data.Imagepath = img
	data.ImageUrl = c.GetUrl(id)
	fmt.Println(data)

	fmt.Println("图片读取成功")
	// 将数据保存在数据库中
	if _, err := models.AddProjectImage(&data); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = data
	} else {
		c.Data["json"] = err.Error()
	}
	defer f.Close()
	c.Data["json"] = result
	c.ServeJSON()
}

// 通过图片名称找到图片并下载
func (c *ProjectImageController) Getimage() {
	fmt.Println("方法已经调用")
	//idStr := c.Ctx.Input.Param("name")
	idStr := c.GetString("name")
	//imageurl, _ := strconv.Atoi(idStr)
	name := idStr
	fmt.Println(name)
	v, err := models.GetProjectImageByName(name)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	//通过url获取图片
	fmt.Println("运行到这里")
	//c.Ctx.Output.Download("static/upload/04b4926d2b63758d6dfa66a09bc0f6eb.jpg", "04b4926d2b63758d6dfa66a09bc0f6eb.jpg")
	imagepath := v.Imagepath
	imageid := v.Id
	imagepathArr := strings.Split(imagepath, ".")
	fmt.Println(imageid + "." + strings.ToLower(imagepathArr[len(imagepathArr)-1]))
	imagename := "static/upload/" + imageid + "." + strings.ToLower(imagepathArr[len(imagepathArr)-1])
	//c.Ctx.Output.Download(imagepath, imageid + "." + strings.ToLower(imagepathArr[len(imagepathArr)-1]))
	c.Ctx.Output.Header("Content-Type", "image"+"/"+strings.ToLower(imagepathArr[len(imagepathArr)-1]))
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", imagename))
	file, err := ioutil.ReadFile(imagename)
	c.Ctx.WriteString(string(file))
	c.ServeJSON()
}

// 生成图片URL地址
func (c *ProjectImageController) GetUrl(s string) (url string) {
	//data, _ := models.GetConnInfoById(1)
	//ip := data.IpAddress
	//port := data.PortAddress
	ip := "127.0.0.1"
	port := "8082"
	var str strings.Builder
	str.WriteString("http://")
	str.WriteString(ip + ":")
	str.WriteString(port + "/")
	str.WriteString(s)
	url = str.String()
	return
}

func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// 自动生成id号码
func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return Md5(base64.URLEncoding.EncodeToString(b))
}
