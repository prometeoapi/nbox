## develop - tareas en el proyecto

**configurar pre-commit**

```shell
./scripts/setup-precommit.sh
```



**instal deps and lint tools**

```shell
make install-all-deps install-tools gomod-tidy
```


## prod build docker
```bash
docker buildx build --platform=linux/amd64 --target production -t nbox:1  --progress=plain .
```