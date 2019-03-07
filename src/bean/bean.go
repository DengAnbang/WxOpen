package bean

import (

	"net/http"
	"encoding/json"
	"gitee.com/DengAnbang/goutils/utils"

)
const (
	OK        = 0  //成功
	NormalErr = -1 //普通错误
)
type ResultData struct {
	Code         int         `json:"code"`
	Type         int         `json:"type"`
	Message      string      `json:"message"`
	DebugMessage string      `json:"debug_message"`
	Data         interface{} `json:"data"`
}
type RequestData struct {
	Code         int               `json:"code"`
	Type         int               `json:"type"`
	Message      string            `json:"message"`
	DebugMessage string            `json:"debug_message"`
	Data         map[string]string `json:"data"`
}

func (r *ResultData) Error() string {
	return r.Message
}
func NewSucceedMessage(data interface{}) *ResultData {
	return &ResultData{Code: OK, Message: "", Data: data}
}
func NewErrorMessage(message string) *ResultData {
	return &ResultData{Code: NormalErr, Message: message}
}
func (r *ResultData) SetDeBugMessage(message string) *ResultData {
	r.DebugMessage = message
	return r
}
func (r *ResultData) WriterResponse(w http.ResponseWriter) {
	bytes, err := json.Marshal(r)
	if err != nil {
		NewErrorMessage("编码错误").WriterResponse(w)
		return
	}
	w.Write(bytes)
}
func (r *ResultData) GetJson() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		return NewErrorMessage("编码错误").GetJson()
	}
	return string(bytes)
}

type UserBean struct {
	UserName  string `json:"user_name"`
	LoginName string `json:"login_name"`
	UserId    string `json:"user_id"`
	SessionId string `json:"-"`
	Pwd       string `json:"-"`
}
type TreeBean struct {
	Id           string `json:"-"`
	NodeId       string `json:"node_id"`
	TreeId       string `json:"tree_id"`
	NodeName     string `json:"node_name"`
	NodeParentId string `json:"node_parent_id"`
	Depth        string `json:"depth"`
	Code         string `json:"code"`
	Sort         string `json:"sort"`
	ModelID      string `json:"model_id"`
	No           string `json:"no"`
}

func (t TreeBean) SetData(data map[string]string) TreeBean {
	t.NodeId = data["node_id"]
	t.NodeName = data["node_name"]
	t.TreeId = data["tree_id"]
	t.NodeParentId = data["node_parent_id"]
	t.Depth = data["depth"]
	t.Code = data["code"]
	t.Sort = data["sort"]
	t.ModelID = data["model_id"]
	t.ModelID = data["model_id"]
	t.No = data["no"]
	t.Id = data["id"]
	return t
}

type TemplateBean struct {
	TemplateId string `json:"template_id"`
	TreeNodeId string `json:"tree_node_id"`
	No         string `json:"no"`
	Name       string `json:"name"`
	Unit       string `json:"unit"`
	DesignNum  string `json:"design_num"`
	Price      string `json:"price"`
	Remarks    string `json:"remarks"`
	Id         string `json:"-"`
}

func (t TemplateBean) SetData(data map[string]string) TemplateBean {
	t.TemplateId = data["template_id"]
	t.TreeNodeId = data["tree_node_id"]
	t.No = data["no"]
	t.Name = data["name"]
	t.Unit = data["unit"]
	t.DesignNum = data["design_num"]
	t.Price = data["price"]
	t.Remarks = data["remarks"]
	t.Id = data["id"]
	return t
}

type ModelBean struct {
	ModelId   string `json:"model_id"`
	NodeId    string `json:"node_id"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	CreatedId string `json:"created_id"`
	Versions  int    `json:"versions"`
}

func (m ModelBean) SetData(data map[string]string) ModelBean {
	m.NodeId = data["node_id"]
	m.ModelId = data["model_id"]
	m.CreatedId = data["created_id"]
	m.Name = data["name"]
	m.Path = data["path"]
	m.Versions = utils.String2int(data["versions"], 0)
	return m
}

type CreatedTreeBean struct {
	TreeBean
	CreatedId string `json:"created_id"`
}
type ModelTree struct {
	Id   string `json:"id"`
	Pid  string `json:"pid"`
	Name string `json:"name"`
	Flag string `json:"flag"`
	Tid  string `json:"tid"`
}

type ModelNodeAddition struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Unit          string `json:"unit"`
	PropertyGroup string `json:"property_group"`
	Description   string `json:"description"`
}
