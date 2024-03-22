setup:
	ansible-playbook -i ansible/hosts ansible/setup.yml

build:
	templ generate
	GOOS=linux GOARCH=arm GOARM=5 go build -o rtmpScreenServer

deploy: build
	cp -v rtmpWebRTCRelay ansible/dist/
	ansible-playbook -i ansible/hosts ansible/deploy.yml
