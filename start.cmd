@echo off
rem Installs what grabb needs, builds it and opens the dashboard.
rem Usage: start.cmd [flags]   e.g. start.cmd -out D:\Videos
setlocal
cd /d "%~dp0"

where winget >nul 2>nul || (
  echo grabb needs winget ^(App Installer^) from the Microsoft Store: https://aka.ms/getwinget
  exit /b 1
)

rem Make tools installed during this run usable without reopening the terminal.
set "PATH=%LOCALAPPDATA%\Microsoft\WinGet\Links;%ProgramFiles%\Go\bin;%USERPROFILE%\go\bin;%PATH%"

call :need go      GoLang.Go
call :need yt-dlp  yt-dlp.yt-dlp
call :need ffmpeg  Gyan.FFmpeg
call :need aria2c  aria2.aria2
call :need deno    DenoLand.Deno
if errorlevel 1 exit /b 1

echo Building grabb...
go build -o bin\grabb.exe .\cmd\grabb || exit /b 1
bin\grabb.exe %*
exit /b %errorlevel%

:need
where %1 >nul 2>nul && exit /b 0
echo Installing %1...
winget install --id %2 -e --silent --accept-source-agreements --accept-package-agreements >nul
where %1 >nul 2>nul && exit /b 0
echo Could not install %1. Install it with: winget install %2
exit /b 1
