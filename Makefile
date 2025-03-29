up:
	docker compose up -d
sh:
	docker compose exec app sh
docs:
	swag init -g main.go --output docs
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto