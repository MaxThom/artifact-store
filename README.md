# Artifact Store

## Builer pattern

- with_store
  - with_local_disk
  - with_s3



## Proxy of a store



go run main.go upload ../eg/manifests/tmtc_v1.0.0.tgz
./package.sh test v0.0.1 metadata.yaml -- hello.txt world.txt aurevoir.txt

./package.sh text v0.0.1 metadata_txt.yaml -- hello.txt world.txt aurevoir.txt
go run main.go upload ../eg/manifests/text_v0.0.1.tgz