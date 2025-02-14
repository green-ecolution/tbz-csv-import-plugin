module github.com/green-ecolution/tbz-csv-import-plugin

go 1.23.6

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/gofiber/fiber/v2 v2.52.5
	github.com/green-ecolution/green-ecolution-backend/pkg/client v0.0.0-20250212180816-c7a5f88209c8
	github.com/green-ecolution/green-ecolution-backend/pkg/plugin v0.0.0-20250212180816-c7a5f88209c8
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/twpayne/go-proj/v10 v10.5.0
	golang.org/x/oauth2 v0.24.0
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/green-ecolution/green-ecolution-backend v1.1.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.58.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace github.com/green-ecolution/green-ecolution-backend => ../green-ecolution-management/green-ecolution-backend

replace github.com/green-ecolution/green-ecolution-backend/pkg/client => ../green-ecolution-management/green-ecolution-backend/pkg/client

replace github.com/green-ecolution/green-ecolution-backend/pkg/plugin => ../green-ecolution-management/green-ecolution-backend/pkg/plugin
