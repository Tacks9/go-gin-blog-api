
# 伪目标 运行标签对应的命令 而非对应文件
.PHONY: build run-api run-cron clean help

# 编译整个项目
all: build

build:run-api run-cron

run-api:
	go build -o bin/api cmd/api/main.go
	./bin/api
	
run-cron:
	go build -o bin/cron cmd/cron/main.go
	./bin/cron

clean:
	rm -rf bin/*

help:
	@echo "make all: compile packages and dependencies"
	@echo "make run-api:  启动API服务"
	@echo "make run-cron: 启动定时任务"