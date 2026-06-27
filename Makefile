.PHONY: swagger run

swagger:
	swag init -g cmd/api/main.go -o docs
	sed -i '' '/LeftDelim\|RightDelim/d' docs/docs.go

run:
	go run cmd/api/main.go
