package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

type VersionInfo struct {
	Current       string `json:"current"`
	IsNew         int    `json:"isNew"`
	NewVersion    string `json:"newVersion"`
	DownloadURL   string `json:"downloadUrl"`
	UpdateContent string `json:"updateContent"`
	ForceUpdate   int    `json:"forceUpdate"`
	CheckAt       string `json:"checkAt"`
}

const Version = "1.0.2"

// const remoteCheckURL = "http://localhost:8080/pos/system/bypass/getVersion"
const remoteCheckURL = "https://admin.yoosobyte.com/api/pos/system/bypass/getVersion"

func GetVersion() string {
	return Version
}

func CheckVersion() *VersionInfo {
	info := &VersionInfo{
		Current: Version,
		CheckAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	clientID := GetSetting("client_id")
	storeName := GetSetting("store_name")
	secondaryPassword := GetSetting("secondary_password")
	if storeName == "" {
		storeName = "悦速用户"
	}
	if secondaryPassword == "" {
		secondaryPassword = "123456"
	}

	params := url.Values{}
	params.Set("clientId", clientID)
	params.Set("versionCode", Version)
	params.Set("storeName", storeName)
	params.Set("secondaryPassword", secondaryPassword)

	reqURL := remoteCheckURL + "?" + params.Encode()
	log.Printf("[版本检查] 请求：%s", reqURL)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		log.Printf("[版本检查] 请求失败：%v", err)
		return info
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8192))
	if err != nil {
		log.Printf("[版本检查] 读取响应失败：%v", err)
		return info
	}
	log.Printf("[版本检查] 响应：%s", string(body))

	var result struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data VersionInfo `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("[版本检查] 解析失败：%v", err)
		return info
	}
	if result.Code != 200 {
		log.Printf("[版本检查] 接口返回错误：code=%d msg=%s", result.Code, result.Msg)
		return info
	}

	if runtime.GOOS == "windows" {
		info.IsNew = 0
		info.NewVersion = ""
		info.DownloadURL = ""
		info.UpdateContent = ""
		info.ForceUpdate = 0
		return info
	}

	info.IsNew = result.Data.IsNew
	info.NewVersion = result.Data.NewVersion
	info.DownloadURL = result.Data.DownloadURL
	info.UpdateContent = result.Data.UpdateContent
	info.ForceUpdate = result.Data.ForceUpdate

	log.Printf("[版本检查] 结果：isNew=%d newVersion=%s forceUpdate=%d", info.IsNew, info.NewVersion, info.ForceUpdate)
	return info
}
