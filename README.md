#### User Entity Management Service

Run integrations tests:

    ./run-integration-test.sh

Start service:

    docker-compose up

Stop service:

    docker-compose down
    # or ^C in the docker-compose process

Stop service and drop the database:

    docker-compose down -v

Swagger Docs for the service: http://localhost:8080/swagger/index.html