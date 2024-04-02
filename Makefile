up:
	docker compose up --build -d

down:
	docker compose down

shell:
	docker compose exec -it wetholds-api bash

logs:
	docker compose logs -f wetholds-api
