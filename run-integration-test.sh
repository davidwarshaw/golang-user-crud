#! /bin/sh

docker-compose -f docker-compose.integration-tests.yaml down -v --remove-orphans
docker-compose -f docker-compose.integration-tests.yaml build
docker-compose -f docker-compose.integration-tests.yaml up --exit-code-from api-test
