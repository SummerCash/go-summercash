read -p "relative package directory (example: '../example'): "  directory # Get dir
read -p "proto file name (example: 'handler'): " name # Get file name

mkdir $name # Init support folder

cd $name # Cd into new folder

protoc --proto_path=$GOPATH/src:../$directory --twirp_out=. --go_out=. ../$directory/$name.proto # Compile proto file