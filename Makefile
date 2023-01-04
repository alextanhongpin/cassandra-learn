up:
	@docker-compose up -d


down:
	@docker-compose down


cqlsh:
	docker exec -it $$(docker ps --filter='name=cassandra' -q) bash cqlsh --username cassandra --password cassandra


sh:
	docker exec -it $$(docker ps --filter='name=cassandra' -q) bash
