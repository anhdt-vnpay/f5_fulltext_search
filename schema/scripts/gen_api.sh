#!/bin/bash 
SCHEMA_PATH="./schema"
# SCHEMA_INCLUDE=$SCHEMA_PATH
SCHEMA_INCLUDE="./schema"


OUT_PATH="./types"
SWAGGER_PATH="./swagger"

export GOPATH="$(go env GOPATH)"
export PATH="$PATH:$(go env GOPATH)/bin"

mkdir -p $OUT_PATH
mkdir -p $SWAGGER_PATH 

source "./schema/scripts/gen_func.sh"

echo "1. Generate service: user.proto"
gen_service api/user.proto
gen_gateway api/user.proto
gen_swagger api/user.proto
gen_openapi api/user.proto
