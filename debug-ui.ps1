$ErrorActionPreference = "Stop"

# Port dla Go Remote w GoLand
$port = 2345

# Budowanie binarki do debugowania
go build -gcflags="all=-N -l" -o wp2go-debug.exe .

# Uruchomienie pod Delve (headless) + przekazanie argumentów aplikacji po `--`
dlv --listen=":$port" --headless=true --api-version=2 exec .\wp2go-debug.exe -- restore --ui
