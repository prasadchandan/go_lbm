all: pcmac run

init:
	go get golang.org/x/mobile/cmd/gomobile
	gomobile init # This could take a few minutes

ios:
	gomobile build -target=ios .

android:
	gomobile build -target=android .

pcmac:
	go build .

run:
	./lbm
