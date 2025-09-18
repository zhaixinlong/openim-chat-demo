package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	APIAddress = "http://127.0.0.1:10002" // 替换成你的 OpenIM API 地址
	AdminID    = "imAdmin"
	Secret     = "openIM123"
)

// 请求结构
type TokenRequest struct {
	UserID     string `json:"userID" binding:"required"`
	PlatformID int    `json:"platformID" binding:"required"`
}

// 响应结构
type TokenResponse struct {
	UserID string `json:"userID"`
	Token  string `json:"token"`
}

func main() {
	r := gin.Default()
	// ✅ 启用 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允许所有来源（生产环境建议改成前端域名）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 签发 Token
	r.POST("/token", func(c *gin.Context) {
		var req TokenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Printf("req: %+v", req)

		adminToken, err := getAdminToken()
		if err != nil {
			fmt.Println("获取管理员 token 失败:", err)
			return
		}
		fmt.Println("管理员 token:", adminToken)

		userToken, err := getUserToken(adminToken, req.UserID, req.PlatformID)
		if err != nil {
			fmt.Println("获取用户 token 失败:", err)
			return
		}
		fmt.Println("用户 token:", userToken)

		c.JSON(http.StatusOK, TokenResponse{
			UserID: req.UserID,
			Token:  userToken,
		})
	})

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.Run(":8081") // 启动服务，监听 8081 端口
}

// 管理员 token 响应结构
type AdminTokenResp struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    struct {
		Token             string `json:"token"`
		ExpireTimeSeconds int64  `json:"expireTimeSeconds"`
	} `json:"data"`
}

// 用户 token 响应结构
type UserTokenResp struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	Data    struct {
		Token             string `json:"token"`
		ExpireTimeSeconds int64  `json:"expireTimeSeconds"`
	} `json:"data"`
}

// 获取管理员 token
func getAdminToken() (string, error) {
	url := fmt.Sprintf("%s/auth/get_admin_token", APIAddress)

	reqBody := map[string]string{
		"secret": Secret,
		"userID": AdminID,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	fmt.Printf("getAdminToken reqBody: %+v", reqBody)

	req, _ := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("operationID", fmt.Sprintf("%d", time.Now().UnixMilli()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, _ := ioutil.ReadAll(resp.Body)
	var respData AdminTokenResp
	if err := json.Unmarshal(respBytes, &respData); err != nil {
		return "", err
	}

	if respData.ErrCode != 0 {
		return "", fmt.Errorf("获取管理员 token 失败: %d %s", respData.ErrCode, respData.ErrMsg)
	}

	return respData.Data.Token, nil
}

// 获取指定用户 token
func getUserToken(adminToken, userID string, platformID int) (string, error) {
	url := fmt.Sprintf("%s/auth/get_user_token", APIAddress)

	reqBody := map[string]interface{}{
		"userID":     userID,
		"platformID": platformID,
	}
	bodyBytes, _ := json.Marshal(reqBody)
	fmt.Printf("getUserToken reqBody: %+v", reqBody)

	req, _ := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("operationID", fmt.Sprintf("%d", time.Now().UnixMilli()))
	req.Header.Set("token", adminToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, _ := ioutil.ReadAll(resp.Body)
	var respData UserTokenResp
	if err := json.Unmarshal(respBytes, &respData); err != nil {
		return "", err
	}

	if respData.ErrCode != 0 {
		return "", fmt.Errorf("获取用户 token 失败: %d %s", respData.ErrCode, respData.ErrMsg)
	}

	return respData.Data.Token, nil
}
