### currently

``` bash
./package.sh <name> <version> <metadata.yaml> -- <file1> <file2> <file3>
go run main.go upload <bucket> ../eg/manifests/tmtc_v1.0.0.tgz


./package.sh text v0.0.1 metadata_txt.yaml -- hello.txt world.txt aurevoir.txt
go run main.go upload ../eg/manifests/text_v0.0.1.tgz