echo off

echo Checking required build folders...
if not exist ".tmp" mkdir .tmp
if not exist ".tmp\win64" mkdir .tmp\win64

echo Setting OS and Architecture for windows / amd64...
set GOOS=windows
set GOARCH=amd64


echo Building: GO-POST...
go build -o .tmp\win64\go-post.exe .\cmd\post

echo Building: GO-GET...
go build -o .tmp\win64\go-get.exe .\cmd\get


echo Creating archive for version %*...
7z a -tzip .\.dist\go-request--windows-amd64--%*.zip .\.tmp\win64\*.exe readme.md

echo Copying application to current directory...
copy /b/v/y .\.tmp\win64\*.exe .\

echo All done!
