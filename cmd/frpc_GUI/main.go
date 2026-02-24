package main

// Copyright 2016 fatedier, fatedier@gmail.com
// Copyright 2025 XY0797, xy0797official@qq.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/policy/featuregate"
	"github.com/fatedier/frp/pkg/policy/security"
	"github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/pkg/util/system"
	_ "github.com/fatedier/frp/web/frpc"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Title            string   `toml:"title"`
	WebsiteURL       string   `toml:"websiteURL"`
	ConfigFilePath   string   `toml:"configFilePath"`
	StrictConfigMode bool     `toml:"strictConfigMode"`
	AllowUnsafe      []string `toml:"allowUnsafe"`
}

var GUIConfig Config

func loadGUIConfig() error {
	// 读取配置文件
	data, err := os.ReadFile("GUIConfig.toml")
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	// 设置默认值
	GUIConfig = Config{
		Title:            "frpc GUI",
		WebsiteURL:       "",
		ConfigFilePath:   "frpc.toml",
		StrictConfigMode: true,
		AllowUnsafe:      []string{},
	}
	// 解析 TOML 文件
	if err := toml.Unmarshal(data, &GUIConfig); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return nil
}

type AppState int

const (
	StateStopped AppState = iota
	StateStarting
	StateRunning
)

type FRPCApp struct {
	app         fyne.App
	window      fyne.Window
	state       AppState
	statusLabel *widget.Label
	logBox      *widget.Entry
	startBtn    *widget.Button
	stopBtn     *widget.Button
	openBtn     *widget.Button
	// 日志锁
	logLock sync.RWMutex
}

func NewFRPCApp() *FRPCApp {
	a := app.New()
	w := a.NewWindow(GUIConfig.Title)
	w.Resize(fyne.NewSize(400, 300))

	app := &FRPCApp{
		app:    a,
		window: w,
		state:  StateStopped,
	}

	app.createUI()
	app.updateUI()

	go app.startFRPC()

	return app
}

func (a *FRPCApp) createUI() {
	a.statusLabel = widget.NewLabel("运行状态：已停止")
	a.logBox = widget.NewMultiLineEntry()
	a.logBox.Wrapping = fyne.TextWrapWord

	a.startBtn = widget.NewButton("启动", a.startFRPC)
	a.stopBtn = widget.NewButton("停止", a.stopFRPC)
	a.openBtn = widget.NewButton("打开网站", a.openWebsite)

	// 使用Border布局，将按钮放在底部，logBox填充剩余空间
	content := container.NewBorder(
		a.statusLabel,
		container.NewHBox(a.startBtn, a.stopBtn, a.openBtn),
		nil,
		nil,
		a.logBox,
	)

	a.window.SetContent(content)
}

func (a *FRPCApp) updateUI() {
	switch a.state {
	case StateStopped:
		a.statusLabel.SetText("运行状态：已停止")
		a.startBtn.Show()
		a.stopBtn.Hide()
		a.openBtn.Hide()
	case StateStarting:
		a.statusLabel.SetText("运行状态：正在启动")
		a.startBtn.Hide()
		a.stopBtn.Hide()
		a.openBtn.Hide()
	case StateRunning:
		a.statusLabel.SetText("运行状态：运行中")
		a.startBtn.Hide()
		a.stopBtn.Show()
		if GUIConfig.WebsiteURL != "" {
			a.openBtn.Show()
		}
	}
}

func (a *FRPCApp) startFRPC() {
	a.state = StateStarting
	a.updateUI()

	// 删除旧日志
	a.ClearLog()

	a.state = StateRunning
	a.updateUI()
	go func() {
		unsafeFeatures := security.NewUnsafeFeatures(GUIConfig.AllowUnsafe)
		err := runFRPC(GUIConfig.ConfigFilePath, unsafeFeatures)
		if err != nil {
			FrpcGUIApp.AddLog(fmt.Sprintf("运行时遇到异常：%v\n", err))
		}
		a.state = StateStopped
		a.updateUI()
	}()
}

func (a *FRPCApp) stopFRPC() {
	if ServiceObj != nil {
		ServiceObj.Close()
	}
	a.state = StateStopped
	a.updateUI()
}

func (a *FRPCApp) openWebsite() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", GUIConfig.WebsiteURL)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", GUIConfig.WebsiteURL)
	default:
		cmd = exec.Command("xdg-open", GUIConfig.WebsiteURL)
	}
	err := cmd.Start()
	if err != nil {
		dialog.ShowError(err, a.window)
	}
}

func (a *FRPCApp) AddLog(log string) {
	a.logLock.Lock()
	a.logBox.Append(log)
	a.logBox.ScrollToBottom()
	a.logLock.Unlock()
}

func (a *FRPCApp) ClearLog() {
	a.logLock.Lock()
	a.logBox.SetText("")
	a.logLock.Unlock()
}

var FrpcGUIApp *FRPCApp
var ServiceObj *client.Service

type GUIWriter struct{}

func (w GUIWriter) Write(p []byte) (n int, err error) {
	if FrpcGUIApp != nil {
		FrpcGUIApp.AddLog(string(p))
	}
	return len(p), nil
}

// 参考cmd/frpc/sub/root.go中的runClient()
// 将`fmt.Printf`改成`FrpcGUIApp.AddLog(fmt.Sprintf())`
// 将`strictConfigMode`改成`GUIConfig.StrictConfigMode`
func runFRPC(cfgFilePath string, unsafeFeatures *security.UnsafeFeatures) error {
	cfg, proxyCfgs, visitorCfgs, isLegacyFormat, err := config.LoadClientConfig(cfgFilePath, GUIConfig.StrictConfigMode)
	if err != nil {
		return err
	}
	if isLegacyFormat {
		FrpcGUIApp.AddLog(fmt.Sprintf("WARNING: ini format is deprecated and the support will be removed in the future, " +
			"please use yaml/json/toml format instead!\n"))
	}

	if len(cfg.FeatureGates) > 0 {
		if err := featuregate.SetFromMap(cfg.FeatureGates); err != nil {
			return err
		}
	}

	warning, err := validation.ValidateAllClientConfig(cfg, proxyCfgs, visitorCfgs, unsafeFeatures)
	if warning != nil {
		FrpcGUIApp.AddLog(fmt.Sprintf("WARNING: %v\n", warning))
	}
	if err != nil {
		return err
	}

	return startService(cfg, proxyCfgs, visitorCfgs, unsafeFeatures, cfgFilePath)
}

// 参考cmd/frpc/sub/root.go中的startService()
// 将`log.InitLogger`改成`log.InitLoggerCustom`，末尾添加参数`GUIWriter{}`即可
func startService(
	cfg *v1.ClientCommonConfig,
	proxyCfgs []v1.ProxyConfigurer,
	visitorCfgs []v1.VisitorConfigurer,
	unsafeFeatures *security.UnsafeFeatures,
	cfgFile string,
) error {
	log.InitLoggerCustom(cfg.Log.To, cfg.Log.Level, int(cfg.Log.MaxDays), cfg.Log.DisablePrintColor, GUIWriter{})

	if cfgFile != "" {
		log.Infof("start frpc service for config file [%s]", cfgFile)
		defer log.Infof("frpc service for config file [%s] stopped", cfgFile)
	}
	svr, err := client.NewService(client.ServiceOptions{
		Common:         cfg,
		ProxyCfgs:      proxyCfgs,
		VisitorCfgs:    visitorCfgs,
		UnsafeFeatures: unsafeFeatures,
		ConfigFilePath: cfgFile,
	})
	if err != nil {
		return err
	}

	shouldGracefulClose := cfg.Transport.Protocol == "kcp" || cfg.Transport.Protocol == "quic"
	// Capture the exit signal if we use kcp or quic.
	if shouldGracefulClose {
		go handleTermSignal(svr)
	}
	return svr.Run(context.Background())
}

// 参考cmd/frpc/sub/root.go中的handleTermSignal()
func handleTermSignal(svr *client.Service) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
}

func createInfoWindow(title string, content string) {
	window := app.New().NewWindow(title)
	window.Resize(fyne.NewSize(400, 300))
	outEntry := widget.NewEntry()
	outEntry.MultiLine = true
	outEntry.Wrapping = fyne.TextWrapWord
	outEntry.SetText(content)
	window.SetContent(outEntry)
	window.ShowAndRun()
}

var GUIConfigHelp = fmt.Sprintf(`在运行路径下创建一个名为 GUIConfig.toml 的配置文件，内容如下：

title: 窗口标题，默认为“frpc GUI”
websiteURL: 网站地址，留空则不显示“打开网站”按钮
configFilePath: frpc配置文件路径，默认为“frpc.toml”
strictConfigMode: 是否严格校验frpc的配置，默认为true
allowUnsafe: 启用不安全的特性列表，默认为空列表。可用的选项有：%s

示例：
title = "xxxxxx站 - 访问器"
websiteURL = "http://localhost:1234/"
configFilePath = "frpc.toml"
strictConfigMode = true
allowUnsafe = []`, strings.Join(security.ClientUnsafeFeatures, ", "))

func main() {
	system.EnableCompatibilityMode()
	args := os.Args[1:]
	// 遍历参数，检查是否需要显示帮助
	for _, arg := range args {
		switch arg {
		case "-h", "--help", "/?", "?", "help":
			createInfoWindow("帮助", GUIConfigHelp)
			return
		}
	}
	err := loadGUIConfig()
	if err != nil {
		createInfoWindow("错误", fmt.Sprintf("配置文件加载失败：%v", err))
		return
	}
	FrpcGUIApp = NewFRPCApp()
	FrpcGUIApp.window.ShowAndRun()
}
