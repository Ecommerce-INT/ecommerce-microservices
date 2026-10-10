module com.ecommerce/inventory-service

go 1.27.0

require (
	com.ecommerce/shared v0.0.0
	github.com/go-chi/chi/v5 v5.3.2
	github.com/go-chi/cors v1.2.2
	github.com/jackc/pgx/v5 v5.11.0
	github.com/stephenafamo/bob v0.50.0
	github.com/stephenafamo/scan v0.9.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/aarondl/opt v0.0.0-20250607033636-982744e1bd65 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/qdm12/reprint v0.0.0-20200326205758-722754a53494 // indirect
	github.com/rogpeppe/go-internal v1.16.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace com.ecommerce/shared => ../shared
