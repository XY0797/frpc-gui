pushd %~dp0web\frpc
call npm install
call npm run build
popd
pushd %~dp0web\frps
call npm install
call npm run build
popd
