setup:
	mkdir pgdata
	chmod 777 pgdata
	chown -R 1000:1000 pgdata

run:
	make setup
	docker-compose up

clean:
	docker-compose down
	rm -rf pgdata
