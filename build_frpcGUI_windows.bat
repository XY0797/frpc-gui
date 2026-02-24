:: build frpc GUI on windows

:: build web
call build_web.bat

:: amd64
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
@echo Build frpc_GUI windows amd64
go build -x -trimpath -ldflags "-s -w -H=windowsgui" -tags frpc -o ./release/windows/frpcGUI_%GOARCH%.exe ./cmd/frpc_GUI
:: Check the return value of the last command
if %ERRORLEVEL% neq 0 (
    echo Build failed with error level %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
:: i386
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=386
set CC=i686-w64-mingw32-gcc
set CXX=i686-w64-mingw32-g++
@echo Build frpc_GUI windows 386
go build -x -trimpath -ldflags "-s -w -H=windowsgui" -tags frpc -o ./release/windows/frpcGUI_%GOARCH%.exe ./cmd/frpc_GUI
:: Check the return value of the last command
if %ERRORLEVEL% neq 0 (
    echo Build failed with error level %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
