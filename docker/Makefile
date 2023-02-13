all: kill clean build run 
		
build: build_base build_gitpod

build_gitpod:
	docker tag jcabrero/hefactory:latest jcabrero/hefactory:gitpod

build_base:
	docker build --network=host \
	-f Dockerfile -t jcabrero/hefactory:latest ..

hefactory:
	docker run  \
		--net host \
		--hostname=$@ --name=$@ \
		-it --rm \
		jcabrero/hefactory:latest bash

run: hefactory 

clean:
	docker rmi $(shell docker images --filter dangling=true -q --no-trunc) || true

kill:
	docker kill hefactory || true

