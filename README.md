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

## environments variables
```ini
export AWS_REGION=us-east-1
export NBOX_ENTRIES_TABLE_NAME=nbox-entries-production
export NBOX_BOX_TABLE_NAME=nbox-box-production
export NBOX_BUCKET_NAME=xx-nbox-box-production
export NBOX_BASIC_AUTH_CREDENTIALS='{"user":"pass"}'

go run cmd/nbox/main.go

curl -X GET --location "http://localhost:7337/api/entry/prefix?v=widget-x%2Fdevelopment" \
    -H "Content-Type: application/json" \
    --basic --user user:pass
```

### example
```shell
base64 <<EOF 
{
   "ENV_1": "{widget-x.development.key}",
   "ENV_2": "{widget-x.development.debug}",
   "GLOBAL_SERVICE": "{widget-x.sentry}",
   "domain": "{private-domain}",
   "version": "1"
}
EOF
  
echo $PAYLOAD | base64
```