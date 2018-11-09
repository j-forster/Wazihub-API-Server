@echo off
go build && %~dp0\Wazihub-API-Server.exe -crt localhost.crt -key localhost.key
