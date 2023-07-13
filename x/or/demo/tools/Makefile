all: miniapp cannon preimage challenge

.PHONY: elf
miniapp: bin/mini.elf

.PHONY: cannon
cannon: bin/cannon

.PHONY: preimage
preimage: bin/preimage

.PHONY: challenge
preimage: bin/challenge

bin:
	mkdir bin

bin/mini.elf: bin
	cd miniapp && GOOS=linux GOARCH=mips GOMIPS=softfloat go build -o ../bin/mini.elf .

bin/cannon: bin
	cd cannon && go build -o ../bin/cannon .

bin/preimage: bin
	cd preimage && go build -o ../bin/preimage .

bin/challenge: bin
	cd challenge && go build -o ../bin/challenge .

clear:
	rm -rf bin
