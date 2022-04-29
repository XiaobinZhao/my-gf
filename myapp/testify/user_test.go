package testify

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/stretchr/testify/suite"
)

var (
	ctx          = context.TODO()
	client       = g.Client()
	resultStruct = &Result{}
)

type Result struct {
	Code    int
	Message string
	Data    interface{}
}

type MyTestSuit struct {
	suite.Suite
}

func (s *MyTestSuit) SetupSuite() {
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%s", "8199"))
	fmt.Println("【SetupSuite】config http client and get token before all test")
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
	// list user
	listContentStr := client.GetContent(ctx, "/users")
	json.Unmarshal([]byte(listContentStr), resultStruct)
	s.Assert().NotEmpty(resultStruct.Data)
	s.Assert().Equal(resultStruct.Code, 0)
	s.Assert().Greater(resultStruct.Data.(map[string]interface{})["total"], float64(0))
}

func TestExample(t *testing.T) {
	suite.Run(t, new(MyTestSuit))
}
