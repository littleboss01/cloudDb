::设置交叉编译的环境linux
set GOOS=linux
set GOARCH=amd64
::编译
go build -o cloud -ldflags "-s -w"
