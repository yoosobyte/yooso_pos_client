//go:build !windows

package service

import "fmt"

func DoUpdate(downloadURL, newVersion string) error {
	return fmt.Errorf("当前系统不支持自动更新")
}

func DownloadUpdater() {
	// 非 Windows 平台不支持自动更新，不执行任何操作
}
