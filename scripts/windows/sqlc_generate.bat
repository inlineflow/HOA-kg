@echo off

call "%~dp0setvars.bat"
cd "%~dp0..\.."

sqlc.exe generate
