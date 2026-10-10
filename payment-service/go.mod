module com.ecommerce/payment-service

go 1.23

require (
	com.ecommerce/pkg/common v0.0.0
	github.com/gofiber/fiber/v2 v2.52.6
	github.com/jackc/pgx/v5 v5.7.2
	github.com/segmentio/kafka-go v0.4.47
	github.com/stretchr/testify v1.8.1
)

replace com.ecommerce/pkg/common => ../pkg/common
