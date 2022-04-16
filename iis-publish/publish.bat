@echo off
echo Script Started!
if not exist publish mkdir publish
cd publish
if not exist public mkdir public
echo Building app...
go build -o app.exe ..\..\src
xcopy /s/y ..\..\src\public .\public
xcopy /s/y ..\..\src\db.db
xcopy /s/y ..\..\src\openapi.json
echo ^<?xml version="1.0" encoding="utf-8"?^>> web.config
echo ^<configuration^>>> web.config
echo     ^<system.webServer^>>> web.config
echo         ^<handlers^>>> web.config
echo             ^<add name="aspNetCore" path="*" verb="*" modules="AspNetCoreModuleV2" resourceType="Unspecified" /^>>> web.config
echo         ^</handlers^>>> web.config
echo         ^<aspNetCore processPath=".\app.exe" stdoutLogEnabled="false" stdoutLogFile=".\logs\stdout"/^>>> web.config
echo     ^</system.webServer^>>> web.config
echo ^</configuration^>>> web.config
echo Build Completed!
pause