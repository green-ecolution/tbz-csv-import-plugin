module github.com/green-ecolution/tbz-csv-import-plugin

go 1.23.4

require (
	github.com/gofiber/fiber/v2 v2.52.5
	github.com/green-ecolution/green-ecolution-backend/client v0.0.0-00010101000000-000000000000
	github.com/green-ecolution/green-ecolution-backend/plugin v0.0.0-00010101000000-000000000000
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/omniscale/go-proj/v2 v2.0.0-20221006090944-6c8a5f5a510d
	github.com/pkg/errors v0.9.1
	golang.org/x/oauth2 v0.24.0
)

require (
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/green-ecolution/green-ecolution-backend v1.0.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pressly/goose v2.7.0+incompatible // indirect
	github.com/pressly/goose/v3 v3.23.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.58.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace github.com/green-ecolution/green-ecolution-backend => ../green-ecolution-management/green-ecolution-backend

replace github.com/green-ecolution/green-ecolution-backend/client => ../green-ecolution-management/green-ecolution-backend/pkg/client

replace github.com/green-ecolution/green-ecolution-backend/plugin => ../green-ecolution-management/green-ecolution-backend/pkg/plugin
