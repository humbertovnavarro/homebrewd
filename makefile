client:
	docker-compose start database
	cd client && npm run dev

musicbot:
	docker-compose start database
	cd musicbot && go run main.go

database-schema:
	docker-compose start database
	cd client && npx prisma db push

.PHONY: database client musicbot database-schema