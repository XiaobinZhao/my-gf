package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/test/gtest"
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
	ctx      = context.TODO()
	userData = g.Map{
		"userName":    "qwesssssssss",
		"displayName": "张三",
		"email":       "san.zhang@gmail.com",
		"phone":       "17628272827",
		"password":    "123qwe.",
		"desc":        "我是zhang三",
	}
	userStruct   = &User{}
	resultStruct = &Result{}
	uerUuid      = "6c2fi104x40cj8wr4vpd32og00yv0pi2"
)

func Test_User_CRUD(t *testing.T) {
	client := g.Client()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%s", "8199"))
	gtest.C(t, func(t *gtest.T) {
		//create user
		createContentStr := client.PostContent(ctx, "/users", userData)
		json.Unmarshal([]byte(createContentStr), resultStruct)
		fmt.Printf("create user: %+v \n", resultStruct)
		t.Assert(resultStruct.Code, 0)
		t.Assert(resultStruct.Data.(map[string]interface{})["userName"], userData["userName"])
		uerUuid = resultStruct.Data.(map[string]interface{})["uuid"].(string)
		// GET user
		getContentStr := client.GetContent(ctx, "/users/"+uerUuid)
		json.Unmarshal([]byte(getContentStr), resultStruct)
		fmt.Printf("get user: %+v \n", resultStruct)
		t.Assert(resultStruct.Code, 0)
		t.Assert(resultStruct.Data.(map[string]interface{})["userName"], userData["userName"])
		// list user
		listContentStr := client.GetContent(ctx, "/users")
		json.Unmarshal([]byte(listContentStr), resultStruct)
		fmt.Printf("list user: %+v \n", resultStruct)
		t.Assert(resultStruct.Code, 0)
		t.AssertGT(resultStruct.Data.(map[string]interface{})["total"], 0)
		// update user
		updateContentStr := client.PatchContent(ctx, "/users/"+uerUuid, g.Map{"displayName": "wangmazi"})
		json.Unmarshal([]byte(updateContentStr), resultStruct)
		fmt.Printf("update user: %+v \n", resultStruct)
		t.Assert(resultStruct.Code, 0)
		t.Assert(resultStruct.Data.(map[string]interface{})["displayName"], "wangmazi")
		// delete user
		deleteContentStr := client.DeleteContent(ctx, "/users/"+uerUuid)
		json.Unmarshal([]byte(deleteContentStr), resultStruct)
		fmt.Printf("delete user: %+v \n", resultStruct)
		t.Assert(resultStruct.Code, 0)
	})
}
