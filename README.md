# frp 修改版

## 原项目

[fatedier](https://github.com/fatedier)的[frp](https://github.com/fatedier/frp)项目

## 修改说明

1. 为frpc添加GUI，详见：`cmd/frpc_GUI`
2. 添加了GUI相关的依赖：

```bash
go get fyne.io/fyne/v2@latest
go mod tidy
```

3. 添加适用于windows的编译脚本：

- `build.bat`用于交叉编译各个平台下的`frps`、`frpc`
- `build_frpcGUI_windows.bat`用于编译windows版的`frpcGUI`

编译`frpcGUI`需要安装gcc编译器，否则无法正常编译

特别注意：需要安装`i686-w64-mingw32-gcc`用于编译32位版本

> 可下载`小熊猫C++`64位和32位的绿色版，解压出`mingw32`和`mingw64`文件夹，移动到合适的地方，最后把`mingw64/bin`和`mingw32/bin`添加到环境变量。
>
> 这种方法并非标准方法，因此在配环境变量时，`mingw64/bin`必须在`mingw32/bin`的上面。
>
> 而且编译完成后最好把`mingw32/bin`的环境变量删除，避免日后出现奇奇怪怪的问题。

4. GUI中的entry需要实现自动滚动，因此修改了库：

手动编辑`~/go/pkg/mod/fyne.io/fyne/v2@v2.7.3/widget/entry.go`:

在第531行(`Append`方法底下)添加下面代码:

```go
// ScrollToBottom scrolls the entry's content to the bottom.
//
// Since: 2.7
func (e *Entry) ScrollToBottom() {
	e.scroll.ScrollToBottom()
}
```

### Linux下编译GUI版本说明

Debian(x-face):

```
sudo apt install libx11-dev libxcursor-dev libxrandr-dev libxinerama-dev libgl-dev libxi-dev pkg-config xorg-dev
go build -x -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux/frpcGUI_amd64 ./cmd/frpc_GUI
```
