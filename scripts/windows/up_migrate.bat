@echo off

call "%~dp0setvars.bat"

cd "%~dp0..\..\db\schema"
goose postgres %MIGRATION_URL% up
