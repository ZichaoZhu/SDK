package main

import (
	"fmt"
    "log"
	"strconv"

    "github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
    "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
    "github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/gin-gonic/gin"
)

var (
	SDK        *fabsdk.FabricSDK
	channelClient *channel.Client
	channelName = "mychannel"
	chaincodeName = "fabcar"
	orgName  = "Org1"
	orgAdmin = "Admin"
	org1Peer0 = "peer0.org1.example.com"
	org2Peer0 = "peer0.org2.example.com"
)

// student 结构体定义
type Student struct {
    School string `json:"school"`
    Major  string `json:"major"`
    Id     int    `json:"id"`
    Name   string `json:"name"`
    Owner  string `json:"owner"` // 创建者的唯一ID
    Status string `json:"status"`// 状态: Pending, Approved, Rejected
}

// ValidateGradeReq 结构体定义
type ValidateGradeReq struct {
    School    string `json:"school" binding:"required"`
    StudentID int    `json:"studentId" binding:"required"`
    NewStatus string `json:"newStatus" binding:"required"`
}

// AddGradeReq 结构体定义
type AddPriceReq struct {
    School      string `json:"school" binding:"required"`
    StudentID   int    `json:"studentId" binding:"required"`
    PrizeName   string `json:"prizeName" binding:"required"`
    PrizeID     string `json:"prizeId" binding:"required"`
    Year        int    `json:"year" binding:"required"`
    Level       string `json:"level" binding:"required"`
    Institution string `json:"institution" binding:"required"`
}

// ValidatePriceReq 结构体定义
type ValidatePriceReq struct {
    PriceID   string `json:"priceId" binding:"required"`
    NewStatus string `json:"newStatus" binding:"required"`
}

// Invoke 是对 ChannelExecute 的简单封装，接受字符串参数切片
func Invoke(funcName string, strArgs []string) (channel.Response, error) {
    var byteArgs [][]byte
    for _, a := range strArgs {
        byteArgs = append(byteArgs, []byte(a))
    }
    return ChannelExecute(funcName, byteArgs)
}

func ChannelExecute(funcName string, args [][]byte)(channel.Response,error){
	var err error
	configPath := "./config.yaml"
	configProvider := config.FromFile(configPath)
	SDK,err = fabsdk.New(configProvider)
	if err != nil{
		log.Fatalf("Failed to create new SDK: %s", err)
	} 
	ctx := SDK.ChannelContext(channelName,fabsdk.WithOrg(orgName),fabsdk.WithUser(orgAdmin))
	channelClient,err = channel.New(ctx)
	response,err := channelClient.Execute(channel.Request{
		ChaincodeID : chaincodeName,
		Fcn : funcName,
		Args: args,
	})
	if err != nil{
		return response,err
	}
	SDK.Close()
	return response,nil
}

func main() {
	r := gin.Default()

	// addStudent
	r.POST("/addStudent",func(c *gin.Context){
		var student Student
		c.BindJSON(&student)
		var result channel.Response
		result, err := Invoke("addStudent", []string{
			student.School,
			student.Major,
			strconv.Itoa(student.Id),
			student.Name,
		})
		fmt.Println(result)
		if err != nil{
			log.Fatalf("Failed to evaluate transaction: %s\n", err)
		}
		c.JSON(200,gin.H{
			"code" : "200",
			"message" : "Add Student Success",
			"result" : string(result.Payload),
		})
	})

	// validateStudent
	r.POST("/validateStudent", func(c *gin.Context) {
		var req ValidateStudentReq
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"code":400,"error":err.Error()}); return
		}
		resp, err := Invoke("validateStudent", []string{
			req.School,
			strconv.Itoa(req.StudentID),
			req.NewStatus,
		})
		if err != nil {
			c.JSON(400, gin.H{"code":400,"error":err.Error()}); return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "validateStudent success",
			"result":  string(resp.Payload),
    	})
	})







	r.Run(":9099")
}