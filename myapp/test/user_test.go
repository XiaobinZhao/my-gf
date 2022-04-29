package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/stretchr/testify/suite"
)

type Result struct {
	Code    int
	Message string
	Data    interface{}
}

type User struct {
	Uuid        string
	UserName    string
	DisplayName string
	Email       string
	Phone       string
	Desc        string
}

var (
	ctx    = context.TODO()
	client = g.Client()
	admin  = g.Map{
		"userName": "admin",
		"password": "password",
	}
	userData = g.Map{
		"userName":    "zhangsanzhangsan",
		"displayName": "张三",
		"email":       "san.zhang@gmail.com",
		"phone":       "17628272827",
		"password":    "123qwe.",
		"desc":        "我是zhang三",
	}
	userStruct   = &User{}
	resultStruct = &Result{}
	uerUuid      = ""
	token        = ""
)

type MyTestSuit struct {
	suite.Suite
}

func (s *MyTestSuit) SetupSuite() {
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%s", "8199"))
	s.login()
	client.SetHeader("Authorization", "Bearer "+token)
	fmt.Println("【SetupSuite】config http client and get token before all test")
}

func (s *MyTestSuit) login() {
	getContentStr := client.PostContent(ctx, "/login", admin)
	json.Unmarshal([]byte(getContentStr), resultStruct)
	s.Assert().Equal(resultStruct.Code, 0)
	s.Assert().Equal(resultStruct.Data.(map[string]interface{})["user"].(map[string]interface{})["userName"], admin["userName"])
	token = fmt.Sprint(resultStruct.Data.(map[string]interface{})["token"])
}

func (s *MyTestSuit) logout() {
	getContentStr := client.DeleteContent(ctx, "/logout/"+uerUuid)
	json.Unmarshal([]byte(getContentStr), resultStruct)
	s.Assert().Equal(resultStruct.Code, 0)
}

func (s *MyTestSuit) TearDownSuite() {
	fmt.Println("【TearDownSuite】delete token after all test")
}

func (s *MyTestSuit) SetupTest() {
}

func (s *MyTestSuit) TearDownTest() {
}

func (s *MyTestSuit) BeforeTest(suiteName, testName string) {
}

func (s *MyTestSuit) AfterTest(suiteName, testName string) {
}

func (s *MyTestSuit) TestUserCRUD() {
	//create user
	createContentStr := client.PostContent(ctx, "/users", userData)
	json.Unmarshal([]byte(createContentStr), resultStruct)
	s.Assert().Equal(resultStruct.Code, 0)
	s.Assert().Equal(resultStruct.Data.(map[string]interface{})["userName"], userData["userName"])
	uerUuid = resultStruct.Data.(map[string]interface{})["uuid"].(string)
	// GET user
	getContentStr := client.GetContent(ctx, "/users/"+uerUuid)
	json.Unmarshal([]byte(getContentStr), resultStruct)
	s.Assert().Equal(resultStruct.Code, 0)
	s.Assert().Equal(resultStruct.Data.(map[string]interface{})["userName"], userData["userName"])
	// list user
	listContentStr := client.GetContent(ctx, "/users")
	json.Unmarshal([]byte(listContentStr), resultStruct)
	s.Assert().Equal(resultStruct.Code, 0)
	s.Assert().Greater(resultStruct.Data.(map[string]interface{})["total"], float64(0))
	// update user
	updateContentStr := client.PatchContent(ctx, "/users/"+uerUuid, g.Map{"displayName": "wangmazi"})
	json.Unmarshal([]byte(updateContentStr), resultStruct)
	s.Assert().Equal(resultStruct.Code, 0)
	s.Assert().Equal(resultStruct.Data.(map[string]interface{})["displayName"], "wangmazi")
	// delete user
	deleteContentStr := client.DeleteContent(ctx, "/users/"+uerUuid)
	json.Unmarshal([]byte(deleteContentStr), resultStruct)
	s.Assert().Equal(resultStruct.Code, 0)

}

func TestExample(t *testing.T) {
	suite.Run(t, new(MyTestSuit))
}
