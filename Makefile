.DEFAULT_GOAL = build

build:
	go build -o gofoto
.PHONY: build


clean:
	rm gofoto
.PHONY: clean

run:
	./gofoto ${HOME} 
.PHONY: clean
