# RTMP Delay Server

## start browser

open chrome-browser with

```
--autoplay-policy=no-user-gesture-required
```

otherwise `autoplay` of the video wont work

## setup raspberry-pi

```
make setup
```

## build and deploy server-app

```
go mod tidy
make build
make deploy
```

## develop

run build in watchmode

```
make run-web-gen
make run
```
