# nbox

**Es una solición para la administración de variables de entornos, secretos, templates y reemplazo de variables en los templates**

### Características

- Almacena las variables de entorno en **AWS Dynamodb**
- Historial de cambios en las variables
- Almcena los secretos en **AWS Parameter Store** con un Key de encriptación propia (**AWS KMS**)
- Almacen centralizado de templates (***ECS task definition*** )
- Reemplazo de nombre de la variable por su valor


### Estructura de nombre de la variable

El nombre de las variables debe estar estructurado por **`stage/service-name/var-name`**

Es necesario que al momento de crear la variable presenta algunos de los siguientes prefijos, en caso de no presentar se usará el valor por defecto `global/`

**stage configurados en el servicio:**

- `development/`

- `qa/`
- `beta/`
- `sandbox/`
- `production/`
- `global/`


## Endpoints para variables

**Para el uso del servicio se requieren**

- Credenciales para autenticación. En esta versión la autenticación es del tipo *Http Basic*


### Endpoint upsert

Este endpoint permite un batch de creación / actualización de variables. En el caso de variables marcadas como secretos se guarda el valor en *AWS Parameter Store* y se almacena en **AWS Dynamodb**  la referencia a *AWS Paramenter Store*

```sh
PAYLOAD=$(<<EOF 
[
   {
      "key": "global/example/email_password",
      "value": "xxxxxxxxxx",
      "secure": true
   },
   {
      "key": "global/example/email_user",
      "value": "test@gmail.com"
   }
]
EOF
)

NBOX_CREDENTIALS='user:pass'

curl -X POST --location -v "https://nbox.example.com/api/entry" \
    -H "Content-Type: application/json" \
    -d "${PAYLOAD}" \
    --basic --user "$NBOX_CREDENTIALS" -sSf
```



### Endpoint entries

Este endpoint permite listar las variables que pertenezcan a determinado namespace

`stage/service`

```shell
NBOX_CREDENTIALS='user:pass'
curl -X GET --location -s "https://nbox.example.com/api/entry/prefix?v=global/example" \
    -H "Content-Type: application/json" \
    --basic --user "$NBOX_CREDENTIALS" -sSf |jq
```

**Response**

```json
[
  {
    "path": "global/example",
    "key": "email_password",
    "value": "/global/example/email_password",
    "secure": true
  },
  {
    "path": "global/example",
    "key": "email_user",
    "value": "test@gmail.com",
    "secure": false
  }
]
```

### Endpoint /entry

Permite obtener el valor de una variable

```shell
NBOX_CREDENTIALS='user:pass'
curl -X GET --location "https://nbox.example.com/api/entry/key?v=global/example/email_user" \
    -H "Content-Type: application/json" \
    --basic --user "$NBOX_CREDENTIALS" -sSf |jq
```

**response**

```json
{
  "path": "",
  "key": "global/example/email_user",
  "value": "imap.gmail.com",
  "secure": false
}
```


## Endpoints para templates

Los templates son almacenados en **AWS S3** donde están versionados, también se mantiene guardado en una tabla de dynamodb la metadata de los templates almacenados

### Endpoint upsert

Permite la creación / modificación de templates.

Para la creción de un template se necesita codificar en base64 el template. Para este ejemplo se usa un JSON para AWS ECS task definition

**task-definition.json**
```json
{
  "requiresCompatibilities": [
    "EC2"
  ],
  "containerDefinitions": [
    {
      "name": "nginx",
      "image": "nginx:latest",
      "memory": 256,
      "cpu": 256,
      "essential": true,
      "portMappings": [
        {
          "containerPort": 80,
          "protocol": "tcp"
        }
      ],
      "secrets": [
        {
          "name": "EMAIL_PASSWORD",
          "valueFrom": "{{global/example/email_password}}"
        }
      ],
      "environment": [
        {
          "name": "ENVIRONMENT_NAME",
          "value": ":stage"
        },
        {
          "name": "DEBUG",
          "value": "{{ global/example/email_user }}"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/nginx_:stage",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "nginx"
        }
      },
      "healthCheck": {
        "command": [
          "CMD-SHELL",
          "wget --no-verbose --tries=1 -O /dev/null --quiet http://localhost || exit 1"
        ],
        "interval": 30,
        "timeout": 10,
        "retries": 3,
        "startPeriod": 10
      }
    }
  ],
  "volumes": [],
  "placementConstraints": [],
  "family": "nginx"
}
```

```shell
TASK_DEFINITION=$(cat task-definition.json| base64)

PAYLOAD=$(<<EOF 
{
  "payload": {
    "service": "example",
    "stage": {
      "development": {
        "template": {
          "name": "task_definition.json",
          "value": "${TASK_DEFINITION}"
        }
      },
      "production": {
        "template": {
          "name": "task_definition.json",
          "value": "${TASK_DEFINITION}"
        }
      }
    }
  }
}
EOF
)

NBOX_CREDENTIALS='user:pass'

curl -X POST --location -v "https://nbox.example.com/api/box" \
    -H "Content-Type: application/json" \
    -d "${PAYLOAD}" \
    --basic --user "$NBOX_CREDENTIALS" -sSf | jq
```

**response**

```json
[
  "example/development/task_definition.json",
  "example/production/task_definition.json"
]
```

### Endpoint obtener template

```shell
curl -X GET --location "https://nbox.example.com/api/box/example/development/task_definition.json" \
    -H "Content-Type: application/json" \
    --basic --user "$NBOX_CREDENTIALS"  -sSf | jq
```

### Endpoint existe template

Permite validar la existencia de un template

```shell
curl -I --head --location "https://nbox.example.com/api/box/token-api/production/task_definition.json" \
    -H "Content-Type: application/json" \
    --basic --user "$NBOX_CREDENTIALS"
```

**response**

```shell
HTTP/2 200 
date: Tue, 27 Aug 2024 14:14:21 GMT
content-type: application/json
content-length: 14
vary: Origin
```

### Endpoint build

Toma un template y reemplaza las varaibles en el template

En la contrucción de los templates se cuentan con dos tipos de variables

- variables de templates ej: `global/development/email_user`
- variables de stages del endpoint: `:service`
  - stage
  - service
  - template
  - cualquier variable tipo querystring

```shell
curl -X GET --location "https://nbox.prometeoapi.com/api/box/token-api/development/task_definition.json/build?image-name=nginx:latest" \
	-H "Content-Type: application/json" \
	--basic --user "$NBOX_CREDENTIALS" -sSf | jq
```

**response**
```json
{
  "requiresCompatibilities": [
    "EC2"
  ],
  "containerDefinitions": [
    {
      "name": "nginx",
      "image": "nginx:latest",
      "memory": 256,
      "cpu": 256,
      "essential": true,
      "portMappings": [
        {
          "containerPort": 80,
          "protocol": "tcp"
        }
      ],
      "secrets": [
        {
          "name": "EMAIL_PASSWORD",
          "valueFrom": "/global/example/email_password"
        }
      ],
      "environment": [
        {
          "name": "ENVIRONMENT_NAME",
          "value": "development"
        },
        {
          "name": "EMAIL_USER",
          "value": "test"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/nginx_development",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "nginx"
        }
      },
      "healthCheck": {
        "command": [
          "CMD-SHELL",
          "wget --no-verbose --tries=1 -O /dev/null --quiet http://localhost || exit 1"
        ],
        "interval": 30,
        "timeout": 10,
        "retries": 3,
        "startPeriod": 10
      }
    }
  ],
  "volumes": [],
  "placementConstraints": [],
  "family": "nginx"
}
```


## Configuración del servicio

```ini
# stages permitidos
NBOX_ALLOWED_PREFIXES = development/,qa/,beta/,sandbox/,production/

# stage por defecto
NBOX_DEFAULT_PREFIX = global

# secret manager para credenciales del tipo http basic
NBOX_BASIC_AUTH_CREDENTIALS =

# tabla de dynamodb para inventario de templates
NBOX_BOX_TABLE_NAME = 

# bucket para almacenar los templates
NBOX_BUCKET_NAME = 

# tabla de dynamodb para almecenar las variable
NBOX_ENTRIES_TABLE_NAME = 

# tabla en dynamodb para almacenar historial de cambios en las variables
NBOX_TRACKING_ENTRIES_TABLE_NAME = 

# tipos de parameter store Standard | Advanced
# Standard: gratuitos, tamaño 4kb
# Advanced: requieren pagos, tamaño 8kb
NBOX_PARAMETER_STORE_DEFAULT_TIER = Standard

# key de KMS para encriptar los secretos
NBOX_PARAMETER_STORE_KEY_ID = 

# determinar el formato de la referencia del secreto en parameter store guardada la tabla de dynamodb
# true: almacena el nombre del parameter store
# false: almancena el ARN del recurso
NBOX_PARAMETER_STORE_SHORT_ARN = true

```






## Deployment

**configurar pre-commit**

```shell
./scripts/setup-precommit.sh
```

**install deps and lint tools**

```shell
make install-all-deps install-tools gomod-tidy
```


### prod build docker
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
export NBOX_ALLOWED_PREFIXES=development/,qa/,beta/,sandbox/,production/
export NBOX_DEFAULT_PREFIX=global
export NBOX_PARAMETER_STORE_DEFAULT_TIER=Standard
export NBOX_PARAMETER_STORE_SHORT_ARN=true

go run cmd/nbox/main.go

curl -X GET --location "http://localhost:7337/health" -H "Content-Type: application/json" 
```


## TODO
- [ ] Enables HTTP Basic authentication. 

- [ ] It accepts a comma-separated list of username:password pairs. 
  Each pair represents a valid username and password combination for authentication.
