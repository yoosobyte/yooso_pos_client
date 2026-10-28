package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	// 同时输出到控制台和文件
	setupLogging()

	printHeader()

	if len(os.Args) < 4 {
		printError("参数错误")
		fmt.Printf("用法: updater.exe <newExe> <targetExe> <pid>\n")
		fmt.Printf("收到: %v\n", os.Args)
		fmt.Println("\n按任意键退出...")
		fmt.Scanln()
		return
	}

	newExe := os.Args[1]
	targetExe := os.Args[2]
	var pid int
	fmt.Sscanf(os.Args[3], "%d", &pid)

	printInfo("任务信息")
	printKV("新版本文件", filepath.Base(newExe))
	printKV("当前程序", filepath.Base(targetExe))
	printKV("进程 PID", fmt.Sprintf("%d", pid))
	fmt.Println()

	// 检查管理员权限
	checkAdminRights()

	// 第一步：关闭主程序
	printStep(1, "正在关闭主程序")
	printDivider()
	if !forceKillProcess(pid, targetExe) {
		printWarning("无法确认主程序是否关闭，继续尝试...")
	} else {
		printSuccess("主程序已关闭")
	}

	// 等待进程退出
	printStep(2, "等待系统释放文件")
	printDivider()
	showSpinner("等待中", 2*time.Second)
	if !waitForFileRelease(targetExe) {
		printError("文件仍被占用，无法继续")
		fmt.Println("\n更新失败！按任意键退出...")
		fmt.Scanln()
		return
	}
	printSuccess("文件已释放")

	// 备份旧文件
	printStep(3, "备份当前版本")
	printDivider()
	backupPath := targetExe + ".backup"
	if err := copyFile(targetExe, backupPath); err != nil {
		printWarning(fmt.Sprintf("备份失败: %v", err))
	} else {
		printSuccess(fmt.Sprintf("备份完成: %s", filepath.Base(backupPath)))
	}

	// 删除旧文件
	printStep(4, "删除旧版本")
	printDivider()
	if err := os.Remove(targetExe); err != nil {
		printWarning(fmt.Sprintf("删除失败: %v", err))
		printInfo("尝试强制删除...")
		if err := forceDeleteFile(targetExe); err != nil {
			printError(fmt.Sprintf("无法删除: %v", err))
			fmt.Println("\n更新失败！按任意键退出...")
			fmt.Scanln()
			return
		}
	}
	printSuccess("旧版本已删除")
	time.Sleep(500 * time.Millisecond)

	// 安装新版本
	printStep(5, "安装新版本")
	printDivider()
	if err := copyFile(newExe, targetExe); err != nil {
		printError(fmt.Sprintf("安装失败: %v", err))
		printInfo("正在恢复旧版本...")
		copyFile(backupPath, targetExe)
		printSuccess("已恢复旧版本")
		fmt.Println("\n更新失败！按任意键退出...")
		fmt.Scanln()
		return
	}
	printSuccess("新版本安装完成")

	// 清理临时文件
	printStep(6, "清理临时文件")
	printDivider()
	os.Remove(newExe)
	os.Remove(backupPath)
	printSuccess("清理完成")

	// 启动新版本
	printStep(7, "启动新版本")
	printDivider()
	cmd := exec.Command(targetExe)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	if err := cmd.Start(); err != nil {
		printError(fmt.Sprintf("启动失败: %v", err))
		fmt.Println("\n更新完成但启动失败！请手动启动程序。")
		fmt.Println("按任意键退出...")
		fmt.Scanln()
		return
	}
	printSuccess(fmt.Sprintf("新版本已启动 (PID: %d)", cmd.Process.Pid))

	// 完成
	fmt.Println()
	printSeparator()
	fmt.Println("  悦速进销存 更新完成！")
	printSeparator()
	fmt.Println()
	fmt.Println("  窗口将在 3 秒后自动关闭...")
	time.Sleep(3 * time.Second)
}

// ========== 界面美化函数 ==========

func printHeader() {
	fmt.Println()
	printSeparator()
	fmt.Println("       悦速进销存 - 自动更新器 v1.0")
	printSeparator()
	fmt.Println()
}

func printSeparator() {
	fmt.Println("=================================================")
}

func printDivider() {
	fmt.Println("  ---------------------------------------------")
}

func printStep(step int, description string) {
	fmt.Printf("  [%d/7] %s\n", step, description)
}

func printKV(key, value string) {
	fmt.Printf("  %-20s %s\n", key+":", value)
}

func printSuccess(msg string) {
	fmt.Printf("  [成功] %s\n", msg)
	log.Printf("[成功] %s", msg)
}

func printError(msg string) {
	fmt.Printf("  [错误] %s\n", msg)
	log.Printf("[错误] %s", msg)
}

func printWarning(msg string) {
	fmt.Printf("  [警告] %s\n", msg)
	log.Printf("[警告] %s", msg)
}

func printInfo(msg string) {
	fmt.Printf("  [信息] %s\n", msg)
	log.Printf("[信息] %s", msg)
}

func showSpinner(prefix string, duration time.Duration) {
	spinner := []string{"|", "/", "-", "\\"}
	start := time.Now()
	for time.Since(start) < duration {
		for _, s := range spinner {
			fmt.Printf("\r  %s %s ...", s, prefix)
			time.Sleep(150 * time.Millisecond)
		}
	}
	fmt.Printf("\r  [完成] %s 完成        \n", prefix)
}

// ========== 功能函数 ==========

func setupLogging() {
	exePath, _ := os.Executable()
	logDir := filepath.Dir(exePath)
	logPath := filepath.Join(logDir, "updater.log")

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		multiWriter := io.MultiWriter(os.Stdout, f)
		log.SetOutput(multiWriter)
	} else {
		log.SetOutput(os.Stdout)
	}

	log.SetFlags(log.Ldate | log.Ltime)
}

func checkAdminRights() {
	cmd := exec.Command("net", "session")
	err := cmd.Run()
	if err != nil {
		printWarning("未以管理员权限运行，可能无法关闭主程序")
		fmt.Println()
	} else {
		printSuccess("已获得管理员权限")
		fmt.Println()
	}
}

func forceKillProcess(pid int, targetExe string) bool {
	processName := filepath.Base(targetExe)
	methods := []struct {
		name string
		fn   func() bool
	}{
		{"taskkill /F /PID", func() bool {
			cmd := exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", pid))
			output, err := cmd.CombinedOutput()
			if err == nil {
				printInfo(fmt.Sprintf("taskkill 输出: %s", strings.TrimSpace(string(output))))
				return true
			}
			return false
		}},
		{"taskkill /F /IM", func() bool {
			cmd := exec.Command("taskkill", "/F", "/IM", processName)
			output, err := cmd.CombinedOutput()
			if err == nil {
				printInfo(fmt.Sprintf("taskkill 输出: %s", strings.TrimSpace(string(output))))
				return true
			}
			return false
		}},
		{"Process.Kill", func() bool {
			proc, err := os.FindProcess(pid)
			if err != nil {
				return false
			}
			return proc.Kill() == nil
		}},
		{"SIGTERM", func() bool {
			proc, err := os.FindProcess(pid)
			if err != nil {
				return false
			}
			err = proc.Signal(syscall.SIGTERM)
			if err != nil {
				return false
			}
			time.Sleep(1 * time.Second)
			return proc.Signal(syscall.Signal(0)) != nil
		}},
		{"wmic delete", func() bool {
			cmd := exec.Command("wmic", "process", "where",
				fmt.Sprintf("name='%s'", processName), "delete")
			return cmd.Run() == nil
		}},
	}

	for _, method := range methods {
		printInfo(fmt.Sprintf("尝试: %s", method.name))
		if method.fn() {
			return true
		}
	}

	return false
}

func waitForFileRelease(filePath string) bool {
	for i := 0; i < 10; i++ {
		f, err := os.OpenFile(filePath, os.O_RDWR, 0)
		if err == nil {
			f.Close()
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}

func forceDeleteFile(filePath string) error {
	cmd := exec.Command("cmd", "/c", "del", "/f", "/q", filePath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		// 尝试重命名后删除
		tmpPath := filePath + ".todelete"
		if err := os.Rename(filePath, tmpPath); err == nil {
			return os.Remove(tmpPath)
		}
		return err
	}
	return nil
}

func copyFile(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("源文件不存在: %w", err)
	}

	s, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer d.Close()

	written, err := io.Copy(d, s)
	if err != nil {
		return fmt.Errorf("复制失败: %w", err)
	}

	if written != srcInfo.Size() {
		return fmt.Errorf("复制不完整: 期望 %d 字节, 实际 %d 字节", srcInfo.Size(), written)
	}

	return nil
}
