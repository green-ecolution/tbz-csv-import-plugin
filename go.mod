module github.com/green-ecolution/tbz-csv-import-plugin

go 1.23.6

require (
	github.com/caarlos0/env/v11 v11.3.1
	github.com/go-chi/chi/v5 v5.2.1
	github.com/green-ecolution/green-ecolution-backend/pkg/client v0.0.0-20250212180816-c7a5f88209c8
	github.com/green-ecolution/green-ecolution-backend/pkg/plugin v0.0.0-20250212180816-c7a5f88209c8
	github.com/joho/godotenv v1.5.1
	github.com/twpayne/go-proj/v10 v10.5.0
	golang.org/x/oauth2 v0.24.0
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/green-ecolution/green-ecolution-backend v1.1.0 // indirect
)

replace github.com/green-ecolution/green-ecolution-backend => ../green-ecolution-management/green-ecolution-backend

replace github.com/green-ecolution/green-ecolution-backend/pkg/client => ../green-ecolution-management/green-ecolution-backend/pkg/client

replace github.com/green-ecolution/green-ecolution-backend/pkg/plugin => ../green-ecolution-management/green-ecolution-backend/pkg/plugin
