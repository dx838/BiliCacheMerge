@echo off
chcp 936 >nul

rem set GO111MODULE=on
rem set GOPROXY=https://goproxy.cn,direct

rem set GOSUMDB=off

rem echo 检查 winres 工具...
rem where winres >nul 2>&1
rem if %errorlevel% neq 0 (
rem     echo 未找到 winres，正在安装...
rem     go get github.com/tc-hib/winres@latest
rem     if %errorlevel% neq 0 (
rem         echo 安装 winres 失败，请检查 Go 环境
rem         pause
rem         exit /b 1
rem     )
rem )




echo 生成新的 syso 文件...

cd /d %~dp0

go-winres.exe make
if %errorlevel% neq 0 (
    echo 生成 syso 文件失败
    pause
    exit /b 1
)


pause
