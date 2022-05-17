export tag=v0.1
export docker_registry=tylertang06

release: 
	echo "building myserver container"
	docker build -t ${docker_registry}/myserver:${tag} .

push: release
	echo "pushing tylertang06/myserver"
	docker push ${docker_registry}/myserver:${tag}