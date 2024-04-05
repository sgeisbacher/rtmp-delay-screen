setup:
	ansible-playbook -i ansible/hosts ansible/setup.yml

build:
	templ generate
	GOOS=linux GOARCH=arm GOARM=5 go build -o rtmpScreenServer

run:
	npx nodemon -e go --signal SIGTERM --exec 'go' run .

run-web-gen:
	npx nodemon -e templ --signal SIGTERM --exec 'templ' generate

deploy: build
	cp -v rtmpWebRTCRelay ansible/dist/
	ansible-playbook -i ansible/hosts ansible/deploy.yml
