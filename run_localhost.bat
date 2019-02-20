@echo off
%~dp0\Wazihub-API-Server.exe -crt localhost.crt -key localhost.key -no-db=1 -www www %*
