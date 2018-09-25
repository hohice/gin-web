default:
	docker build -t gins-builder:1.1 .
	@echo "build ok"

pull:
	docker pull hohice/gins-builder:1.1
	@echo "pull ok"

push:
	docker push hohice/gins-builder:1.1
	@echo "push ok"

run:
	docker run --rm -it -v `pwd`:/root hohice/ginS-builder:1.1 

clean:
	@echo "clean ok"
