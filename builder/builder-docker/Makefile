default:
	docker build -t ginS-builder:1.0 .
	@echo "build ok"

pull:
	docker pull hohice/ginS-builder:1.0
	@echo "pull ok"

push:
	docker push hohice/ginS-builder:1.0
	@echo "push ok"

run:
	docker run --rm -it -v `pwd`:/root hohice/ginS-builder:1.0

clean:
	@echo "clean ok"
