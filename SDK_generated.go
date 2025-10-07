package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// --- 配置常量 ---
const (
	mspID        = "Org1MSP"
	cryptoPath   = "/home/zzc/SDK/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore/"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
	channelName  = "mychannel"
	chaincodeName = "basic" // 假设链码名称为 'basic'
)

// SDK 结构体封装了与 Fabric 网络的连接
type SDK struct {
	contract *client.Contract
	gateway  *client.Gateway
}

// --- SDK 初始化 ---

// NewSDK 创建并初始化一个新的 SDK 实例
func NewSDK() (*SDK, error) {
	clientConnection := newGrpcConnection()
	id := newIdentity()
	sign := newSign()

	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("无法连接到 Gateway: %w", err)
	}

	network := gateway.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	return &SDK{
		contract: contract,
		gateway:  gateway,
	}, nil
}

// Close 关闭 Gateway 连接
func (sdk *SDK) Close() {
	if sdk.gateway != nil {
		sdk.gateway.Close()
	}
}

// --- 辅助函数 ---

func newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("创建 gRPC 连接失败: %w", err))
	}

	return connection
}

func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("读取密钥目录失败: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))
	if err != nil {
		panic(fmt.Errorf("读取私钥文件失败: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func loadCertificate(path string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取证书文件失败: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// formatJSON 格式化并打印 JSON 字节
func formatJSON(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return string(data)
	}
	return prettyJSON.String()
}

// --- 链码函数封装 ---

// AddStudent 提交一个添加新学生的申请
func (sdk *SDK) AddStudent(school, major, id, name string) error {
	fmt.Printf("--> 提交交易: AddStudent, 参数: %s, %s, %s, %s\n", school, major, id, name)
	_, err := sdk.contract.SubmitTransaction("addStudent", school, major, id, name)
	if err != nil {
		return fmt.Errorf("提交 AddStudent 交易失败: %w", err)
	}
	fmt.Println("<-- 交易成功")
	return nil
}

// AddGrade 提交一个添加成绩的申请
func (sdk *SDK) AddGrade(courseName, courseId, teacher, school, studentId, year, score, semester string) error {
	fmt.Printf("--> 提交交易: AddGrade, 参数: %s, %s, %s, %s, %s, %s, %s, %s\n", courseName, courseId, teacher, school, studentId, year, score, semester)
	_, err := sdk.contract.SubmitTransaction("addGrade", courseName, courseId, teacher, school, studentId, year, score, semester)
	if err != nil {
		return fmt.Errorf("提交 AddGrade 交易失败: %w", err)
	}
	fmt.Println("<-- 交易成功")
	return nil
}

// AddPrice 提交一个添加奖项的申请
func (sdk *SDK) AddPrice(school, studentId, prizeName, prizeId, year, level, institution string) error {
	fmt.Printf("--> 提交交易: AddPrice, 参数: %s, %s, %s, %s, %s, %s, %s\n", school, studentId, prizeName, prizeId, year, level, institution)
	_, err := sdk.contract.SubmitTransaction("addPrice", school, studentId, prizeName, prizeId, year, level, institution)
	if err != nil {
		return fmt.Errorf("提交 AddPrice 交易失败: %w", err)
	}
	fmt.Println("<-- 交易成功")
	return nil
}

// ValidateStudent 审批一个学生身份
func (sdk *SDK) ValidateStudent(school, studentId, newStatus string) error {
	fmt.Printf("--> 提交交易: ValidateStudent, 参数: %s, %s, %s\n", school, studentId, newStatus)
	_, err := sdk.contract.SubmitTransaction("validateStudent", school, studentId, newStatus)
	if err != nil {
		return fmt.Errorf("提交 ValidateStudent 交易失败: %w", err)
	}
	fmt.Println("<-- 交易成功")
	return nil
}

// ValidateGrade 审批一个成绩
func (sdk *SDK) ValidateGrade(school, studentId, courseId, year, semester, newStatus string) error {
	fmt.Printf("--> 提交交易: ValidateGrade, 参数: %s, %s, %s, %s, %s, %s\n", school, studentId, courseId, year, semester, newStatus)
	_, err := sdk.contract.SubmitTransaction("validateGrade", school, studentId, courseId, year, semester, newStatus)
	if err != nil {
		return fmt.Errorf("提交 ValidateGrade 交易失败: %w", err)
	}
	fmt.Println("<-- 交易成功")
	return nil
}

// ValidatePrice 审批一个奖项
func (sdk *SDK) ValidatePrice(priceId, newStatus string) error {
	fmt.Printf("--> 提交交易: ValidatePrice, 参数: %s, %s\n", priceId, newStatus)
	_, err := sdk.contract.SubmitTransaction("validatePrice", priceId, newStatus)
	if err != nil {
		return fmt.Errorf("提交 ValidatePrice 交易失败: %w", err)
	}
	fmt.Println("<-- 交易成功")
	return nil
}

// QueryStudent 查询已批准的学生信息
func (sdk *SDK) QueryStudent(school, studentId string) (string, error) {
	fmt.Printf("--> 查询: QueryStudent, 参数: %s, %s\n", school, studentId)
	evaluateResult, err := sdk.contract.EvaluateTransaction("queryStudent", school, studentId)
	if err != nil {
		return "", fmt.Errorf("查询 QueryStudent 失败: %w", err)
	}
	result := formatJSON(evaluateResult)
	fmt.Printf("<-- 查询结果: %s\n", result)
	return result, nil
}

// QueryGrade 查询已批准的成绩信息
func (sdk *SDK) QueryGrade(school, studentId, courseId, year, semester string) (string, error) {
	fmt.Printf("--> 查询: QueryGrade, 参数: %s, %s, %s, %s, %s\n", school, studentId, courseId, year, semester)
	evaluateResult, err := sdk.contract.EvaluateTransaction("queryGrade", school, studentId, courseId, year, semester)
	if err != nil {
		return "", fmt.Errorf("查询 QueryGrade 失败: %w", err)
	}
	result := formatJSON(evaluateResult)
	fmt.Printf("<-- 查询结果: %s\n", result)
	return result, nil
}

// QueryPrice 查询已批准的奖项信息
func (sdk *SDK) QueryPrice(priceId string) (string, error) {
	fmt.Printf("--> 查询: QueryPrice, 参数: %s\n", priceId)
	evaluateResult, err := sdk.contract.EvaluateTransaction("queryPrice", priceId)
	if err != nil {
		return "", fmt.Errorf("查询 QueryPrice 失败: %w", err)
	}
	result := formatJSON(evaluateResult)
	fmt.Printf("<-- 查询结果: %s\n", result)
	return result, nil
}

// --- Main 函数 (示例用法) ---
func main() {
	fmt.Println("============ SDK aplication starting ============")

	sdk, err := NewSDK()
	if err != nil {
		fmt.Printf("创建 SDK 失败: %v\n", err)
		return
	}
	defer sdk.Close()

	fmt.Println("\n**************************************************")
	fmt.Println("1. 学生 'Tom' 申请学生身份")
	// 注意：这里的学号 '1001' 应该是字符串
	if err := sdk.AddStudent("PKU", "CS", "1001", "Tom"); err != nil {
		fmt.Printf("错误: %v\n", err)
	}

	fmt.Println("\n**************************************************")
	fmt.Println("2. 查询 'Tom' 的信息 (此时应失败，因为未审批)")
	if _, err := sdk.QueryStudent("PKU", "1001"); err != nil {
		fmt.Printf("预期中的错误: %v\n", err)
	}

	fmt.Println("\n**************************************************")
	fmt.Println("3. 验证者审批 'Tom' 的学生身份")
	if err := sdk.ValidateStudent("PKU", "1001", "Approved"); err != nil {
		fmt.Printf("错误: %v\n", err)
	}

	fmt.Println("\n**************************************************")
	fmt.Println("4. 再次查询 'Tom' 的信息 (此时应成功)")
	if _, err := sdk.QueryStudent("PKU", "1001"); err != nil {
		fmt.Printf("错误: %v\n", err)
	}

	fmt.Println("\n**************************************************")
	fmt.Println("5. 'Tom' 为自己添加成绩")
	if err := sdk.AddGrade("Blockchain", "CS101", "Dr. Nakamoto", "PKU", "1001", "2025", "95.5", "1"); err != nil {
		fmt.Printf("错误: %v\n", err)
	}

	fmt.Println("\n**************************************************")
	fmt.Println("6. 验证者审批 'Tom' 的成绩")
	if err := sdk.ValidateGrade("PKU", "1001", "CS101", "2025", "1", "Approved"); err != nil {
		fmt.Printf("错误: %v\n", err)
	}

	fmt.Println("\n**************************************************")
	fmt.Println("7. 查询 'Tom' 的成绩")
	if _, err := sdk.QueryGrade("PKU", "1001", "CS101", "2025", "1"); err != nil {
		fmt.Printf("错误: %v\n", err)
	}

	fmt.Println("\n============ SDK aplication finished ============")
}
