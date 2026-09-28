@echo off

REM Go path config
REM set GOROOT=C:\Go
REM set GOPATH=%USERPROFILE%\go
REM set PATH=%PATH%;%GOROOT%\bin;%GOPATH%\bin

REM cd work dir
cd /d "%~dp0"

echo start build...
set GOOS=windows
set GOARCH=arm64
set CGO_ENABLED=0


@echo off
set "content="
set n=1
set count=0
for /f "delims=" %%i in (env) do (
    set /a count+=1
    echo %%i
    rem set "content=!content!%%i\n"
    set "content=%%i"
    :: 
    if !count! equ %n% goto end
)
:end

go build -o %content%_arm.exe

if %errorlevel% equ 0 (
    echo build ok . %content%.exe
    echo press any key to out...
    pause >nul
) else (
    echo build error
    echo press any key to out...
    pause >nul
)