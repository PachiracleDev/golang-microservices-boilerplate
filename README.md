### INSTALAR PROTOC EN TU TERMINAL
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc
source ~/.bashrc


## GENERAR LOS PROTOBUFFERS
protoc --go_out=. --go_opt=paths=source_relative protobuffers/*.proto