
::Build For Linux

echo "Clean for build..."
go clean

echo "Build For Linux..."

::set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64

go build -o forward-server

::go build -ldflags "-s -w" -o forward-server
::-s ȥ�����ű�
::-w ȥ��DWARF������Ϣ���õ��ĳ�������gdb����

echo "--------- Build For Linux Success!"


pause

