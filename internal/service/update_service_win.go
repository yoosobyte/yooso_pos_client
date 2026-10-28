//go:build windows

package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime" // 添加这个导入
	"syscall"
	"time"

	"yoosoPosV3/internal/config"
)

const updaterURL = "https://qnres.yoosobyte.com/updater_20260705225655A011.exe"

// DownloadUpdater 应用启动时下载更新器到数据目录
func DownloadUpdater() {
	appDataDir := config.GetAppDataDir()
	updaterPath := filepath.Join(appDataDir, "updater.exe")

	// 检查文件是否存在且大小正常
	if info, err := os.Stat(updaterPath); err == nil {
		if info.Size() > 0 {
			log.Printf("[更新器] 已存在: %s (大小: %d 字节)", updaterPath, info.Size())
			return
		}
		log.Printf("[更新器] 文件大小为0，重新下载")
		os.Remove(updaterPath)
	}

	log.Printf("[更新器] 下载：%s -> %s", updaterURL, updaterPath)
	resp, err := http.Get(updaterURL)
	if err != nil {
		log.Printf("[更新器] 下载失败：%v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[更新器] 下载失败 HTTP %d", resp.StatusCode)
		return
	}

	// 先下载到临时文件，成功后再重命名
	tmpPath := updaterPath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		log.Printf("[更新器] 创建临时文件失败: %v", err)
		return
	}

	written, err := io.Copy(f, resp.Body)
	f.Close()

	if err != nil || written == 0 {
		os.Remove(tmpPath)
		log.Printf("[更新器] 写入失败: %v, 字节数: %d", err, written)
		return
	}

	// 删除旧文件，重命名新文件
	os.Remove(updaterPath)
	if err := os.Rename(tmpPath, updaterPath); err != nil {
		// 重命名失败，尝试复制
		log.Printf("[更新器] 重命名失败，尝试复制: %v", err)
		copyFile(tmpPath, updaterPath)
		os.Remove(tmpPath)
	}

	log.Printf("[更新器] 下载完成 (大小: %d 字节)", written)
}

// DoUpdate 下载新版本并启动更新器
func DoUpdate(downloadURL, newVersion string) error {
	if downloadURL == "" {
		return fmt.Errorf("下载地址为空")
	}

	log.Printf("[更新] === 开始更新流程 ===")

	// 确保更新器存在
	DownloadUpdater()

	// 准备路径
	appDataDir := config.GetAppDataDir()
	fileName := fmt.Sprintf("yooso_pos_%s.exe", newVersion)
	newPath := filepath.Join(appDataDir, fileName)
	updaterPath := filepath.Join(appDataDir, "updater.exe")

	// 获取当前 exe 路径
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取当前进程路径失败：%w", err)
	}

	// 诊断信息
	log.Printf("[更新] 诊断信息:")
	log.Printf("[更新]   - 更新器路径: %s", updaterPath)
	log.Printf("[更新]   - 当前程序路径: %s", currentExe)
	log.Printf("[更新]   - 新程序路径: %s", newPath)

	// 验证更新器存在
	if info, err := os.Stat(updaterPath); err != nil {
		return fmt.Errorf("更新器不存在: %s", updaterPath)
	} else {
		log.Printf("[更新]   - 更新器大小: %d 字节", info.Size())
	}

	// 下载新版本
	log.Printf("[更新] 开始下载新版本: %s", downloadURL)

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("下载失败：%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，HTTP状态码：%d", resp.StatusCode)
	}

	// 下载到临时文件
	tmpPath := newPath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("创建临时文件失败：%w", err)
	}
	written, err := io.Copy(f, resp.Body)
	f.Close()

	if err != nil || written == 0 {
		os.Remove(tmpPath)
		return fmt.Errorf("写入失败，%d 字节：%v", written, err)
	}

	// 重命名为正式文件
	os.Remove(newPath)
	if err := os.Rename(tmpPath, newPath); err != nil {
		log.Printf("[更新] 重命名失败，尝试复制: %v", err)
		if err := copyFile(tmpPath, newPath); err != nil {
			os.Remove(tmpPath)
			return fmt.Errorf("复制新版本失败: %w", err)
		}
		os.Remove(tmpPath)
	}

	log.Printf("[更新] 新版本下载完成: %s (%d 字节)", newPath, written)

	// 获取当前进程PID
	pid := os.Getpid()

	log.Printf("[更新] 启动更新器...")
	log.Printf("[更新]   - 命令: %s", updaterPath)
	log.Printf("[更新]   - 参数: new=%s, target=%s, pid=%d", newPath, currentExe, pid)

	// 启动更新器
	cmd := exec.Command(updaterPath, newPath, currentExe, fmt.Sprintf("%d", pid))

	// 只在 Windows 系统上设置特定的进程属性
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
		}
	} else {
		// macOS/Linux 不需要特殊设置，或者设置一个空的 SysProcAttr
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}

	// 不共享标准输入输出，完全独立
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		log.Printf("[更新] 启动失败，尝试备用方式: %v", err)

		// 备用方式：通过 cmd /c start 启动（仅 Windows）
		if runtime.GOOS == "windows" {
			cmd2 := exec.Command("cmd", "/c", "start", "",
				updaterPath, newPath, currentExe, fmt.Sprintf("%d", pid))
			cmd2.SysProcAttr = &syscall.SysProcAttr{
				CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
			}
			if err := cmd2.Start(); err != nil {
				return fmt.Errorf("启动更新器失败（主方式和备用方式都失败）: %w", err)
			}
			log.Printf("[更新] 备用方式启动成功")
		} else {
			// macOS/Linux 尝试直接启动
			if err := cmd.Start(); err != nil {
				return fmt.Errorf("启动更新器失败: %w", err)
			}
		}
	} else {
		log.Printf("[更新] 更新器已启动，PID: %d", cmd.Process.Pid)
	}

	log.Printf("[更新] 等待更新器接管...")

	// 给更新器时间启动和初始化
	time.Sleep(2 * time.Second)

	log.Printf("[更新] 主程序退出，等待被更新器关闭")

	// 主程序退出
	os.Exit(0)

	return nil
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}
