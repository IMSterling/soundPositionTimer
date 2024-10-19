@echo off
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo Golang not found, installing...
    set DOWNLOAD_URL=https://golang.org/dl/go1.21.1.windows-amd64.msi
    curl -LO %DOWNLOAD_URL%
    msiexec /i go1.21.1.windows-amd64.msi /quiet /norestart
    set PATH=%PATH%;C:\Go\bin
)
go build main.go
