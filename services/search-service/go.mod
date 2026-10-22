module com.ecommerce/search-service

go 1.27.0

require (
	com.ecommerce/shared v0.0.0
	github.com/go-chi/chi/v5 v5.3.2
	github.com/go-chi/cors v1.2.2
	github.com/opensearch-project/opensearch-go/v5 v5.0.3
	github.com/segmentio/kafka-go v0.4.47
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/rs/dnscache v0.0.0-20230804202142-fc85eb664529 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
)

replace com.ecommerce/shared => ../shared
