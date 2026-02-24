:: build frp on windows

:: clean and recreate release directory
rd /s /q release
md release
md release\windows
md release\darwin 
md release\freebsd
md release\openbsd
md release\linux_x86
md release\linux_arm
md release\linux_mips
md release\linux_riscv64
md release\linux_loong64
md release\android

:: build web
call build_web.bat

:: ==windows==
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
:: When building the first target, dependencies may need to be installed, so it is not concurrent.
@echo Build frpc windows amd64 ...
go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/windows/frpc_%GOARCH%.exe ./cmd/frpc
:: Check the return value of the last command
if %ERRORLEVEL% neq 0 (
    echo Build failed with error level %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)

:: concurrent build
start "Build frps windows amd64" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/windows/frps_%GOARCH%.exe ./cmd/frps
set GOARCH=386
start "Build frpc windows 386" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/windows/frpc_%GOARCH%.exe ./cmd/frpc
start "Build frps windows 386" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/windows/frps_%GOARCH%.exe ./cmd/frps

:: ==darwin==
set GOOS=darwin
set GOARCH=amd64
start "Build frpc darwin amd64" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/darwin/frpc_%GOOS%_%GOARCH% ./cmd/frpc
start "Build frps darwin amd64" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/darwin/frps_%GOOS%_%GOARCH% ./cmd/frps
set GOARCH=arm64
start "Build frpc darwin arm64" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/darwin/frpc_%GOOS%_%GOARCH% ./cmd/frpc
:: execute commands directly to avoid excessive concurrency
@echo Build frps darwin arm64 ...
go build -trimpath -ldflags "-s -w" -tags frps -o ./release/darwin/frps_%GOOS%_%GOARCH% ./cmd/frps
:: Check the return value of the last command
if %ERRORLEVEL% neq 0 (
    echo Build frps darwin arm64 with error level %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)

:: ==freebsd==
set GOOS=freebsd
set GOARCH=amd64
start "Build frpc freebsd amd64" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/freebsd/frpc_%GOOS%_%GOARCH% ./cmd/frpc
start "Build frps freebsd amd64" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/freebsd/frps_%GOOS%_%GOARCH% ./cmd/frps

:: ==openbsd==
set GOOS=openbsd
set GOARCH=amd64
start "Build frpc openbsd amd64" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/openbsd/frpc_%GOOS%_%GOARCH% ./cmd/frpc
start "Build frps openbsd amd64" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/openbsd/frps_%GOOS%_%GOARCH% ./cmd/frps

:: ==linux==
set GOOS=linux
set GOARCH=amd64
start "Build frpc linux amd64" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux_x86/frpc_%GOOS%_%GOARCH% ./cmd/frpc
start "Build frps linux amd64" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/linux_x86/frps_%GOOS%_%GOARCH% ./cmd/frps
set GOARCH=386
start "Build frpc linux 386" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux_x86/frpc_%GOOS%_%GOARCH% ./cmd/frpc
:: execute commands directly to avoid excessive concurrency
@echo Build frps linux 386 ...
go build -trimpath -ldflags "-s -w" -tags frps -o ./release/linux_x86/frps_%GOOS%_%GOARCH% ./cmd/frps
:: Check the return value of the last command
if %ERRORLEVEL% neq 0 (
    echo Build frps linux 386 with error level %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
set GOARCH=arm
set GOARM=7
start "Build frpc linux arm hf" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux_arm/frpc_%GOOS%_arm_hf ./cmd/frpc
start "Build frps linux arm hf" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/linux_arm/frps_%GOOS%_arm_hf ./cmd/frps
set GOARM=5
start "Build frpc linux arm" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux_arm/frpc_%GOOS%_arm ./cmd/frpc
start "Build frps linux arm" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/linux_arm/frps_%GOOS%_arm ./cmd/frps
set GOARCH=arm64
start "Build frpc linux arm64" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux_arm/frpc_%GOOS%_%GOARCH% ./cmd/frpc
start "Build frps linux arm64" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/linux_arm/frps_%GOOS%_%GOARCH% ./cmd/frps
set GOARCH=mips64
start "Build frpc linux mips64" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux_mips/frpc_%GOOS%_%GOARCH% ./cmd/frpc
:: execute commands directly to avoid excessive concurrency
@echo Build frps linux mips64 ...
go build -trimpath -ldflags "-s -w" -tags frps -o ./release/linux_mips/frps_%GOOS%_%GOARCH% ./cmd/frps
:: Check the return value of the last command
if %ERRORLEVEL% neq 0 (
    echo Build frps linux mips64 with error level %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
set GOARCH=mips64le
start "Build frpc linux mips64le" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux_mips/frpc_%GOOS%_%GOARCH% ./cmd/frpc
start "Build frps linux mips64le" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/linux_mips/frps_%GOOS%_%GOARCH% ./cmd/frps
set GOARCH=mips
set GOMIPS=softfloat
start "Build frpc linux mips" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux_mips/frpc_%GOOS%_%GOARCH% ./cmd/frpc
start "Build frps linux mips" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/linux_mips/frps_%GOOS%_%GOARCH% ./cmd/frps
set GOARCH=mipsle
set GOMIPS=softfloat
start "Build frpc linux mipsle" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux_mips/frpc_%GOOS%_%GOARCH% ./cmd/frpc
start "Build frps linux mipsle" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/linux_mips/frps_%GOOS%_%GOARCH% ./cmd/frps
set GOARCH=riscv64
start "Build frpc linux riscv64" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux_riscv64/frpc_%GOOS%_%GOARCH% ./cmd/frpc
:: execute commands directly to avoid excessive concurrency
@echo Build frps linux riscv64 ...
go build -trimpath -ldflags "-s -w" -tags frps -o ./release/linux_riscv64/frps_%GOOS%_%GOARCH% ./cmd/frps
:: Check the return value of the last command
if %ERRORLEVEL% neq 0 (
    echo Build frps linux riscv64 with error level %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
set GOARCH=loong64
start "Build frpc linux loong64" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/linux_loong64/frpc_%GOOS%_%GOARCH% ./cmd/frpc
start "Build frps linux loong64" go build -trimpath -ldflags "-s -w" -tags frps -o ./release/linux_loong64/frps_%GOOS%_%GOARCH% ./cmd/frps

:: ==android==
set GOOS=android
set GOARCH=arm64
start "Build frpc android arm64" go build -trimpath -ldflags "-s -w" -tags frpc -o ./release/android/frpc_%GOOS%_%GOARCH% ./cmd/frpc
:: execute commands directly to avoid excessive concurrency
@echo Build frps android arm64 ...
go build -trimpath -ldflags "-s -w" -tags frps -o ./release/android/frps_%GOOS%_%GOARCH% ./cmd/frps
:: Check the return value of the last command
if %ERRORLEVEL% neq 0 (
    echo Build frps android arm64 with error level %ERRORLEVEL%
    exit /b %ERRORLEVEL%
)
