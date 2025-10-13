@echo off
wt ^
    --title "Templ Watch" cmd.exe /k "templ generate --watch --proxy=\"http://localhost:8080\" --open-browser=false" ^
; ^
    --title "Go Server (Air)" cmd.exe /k "go.exe run github.com/air-verse/air@v1.63.0 --build.cmd \"go.exe build -o tmp/bin/main.exe\" --build.bin \"tmp/bin/main.exe\" --build.delay \"100\" --build.exclude_dir \"node_modules\" --build.include_ext \"go\" --build.stop_on_error \"false\" --misc.clean_on_exit true" ^
; ^
    --title "Tailwind Watch" cmd.exe /k "pnpx @tailwindcss/cli -i ./input.css -o ./assets/css/styles.css --minify --watch" ^
; ^
    --title "Asset Reloader (Air)" cmd.exe /k "go.exe run github.com/air-verse/air@v1.63.0 --build.cmd \"templ.exe generate --notify-proxy\" --build.bin \"" --build.delay \"100\" --build.exclude_dir \"\" --build.include_dir \"assets\" --build.include_ext \"js,css\""
    REM --title "Asset Reloader (Air)" cmd.exe /k "go.exe run github.com/air-verse/air@v1.63.0 --build.cmd \"templ.exe generate --notify-proxy\" --build.bin \"cmd.exe /c echo 1\" --build.delay \"100\" --build.exclude_dir \"\" --build.include_dir \"assets\" --build.include_ext \"js,css\""


REM
REM wt ^
REM     --title "Templ watch" cmd.exe /k "templ generate --watch --proxy="http://localhost:8080" --open-browser=false -v" ^
REM ;   --tab ^
REM     --title "air"
