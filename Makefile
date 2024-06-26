.Phony: build-fe
build-fe:
	@ echo "Build Frontend"
	@ cd web && npm install && npm run build

.Phony: run-client
run-client:
	@ echo "Run Client App"
	@ cd web && npm run dev

.Phony: run
run:
	FILE=db.sqlite3; \
		if [ ! -f "$$FILE" ]; then \
			cp ./db.sqlite3.example db.sqlite3; \
		fi
	@echo "Run Server  App"
	go mod tidy -compat=1.22
	go run main.go