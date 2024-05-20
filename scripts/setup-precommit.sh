#!/bin/sh

: "${GIT_HOOKS_PATH:=.git/hooks}"

make install-tools

echo "creating pre-commit ${GIT_HOOKS_PATH}/pre-commit"

cat <<EOF > ${GIT_HOOKS_PATH}/pre-commit
#!/bin/sh
set -euo pipefail

GOIMPORTS=./bin/goimports
STATIC_CHECK_OUT=./bin/staticcheck


echo "  _____ _____ _____ _____ _____ _____ _____ _____ _____ _____ "
echo "|  running                                                    |"
echo "|  -> lint-static-check                                       |"
echo "|  -> imports-check                                           |"
echo "| _____ _____ _____ _____ _____ _____ _____ _____ _____ _____ |"


if [ ! -f "\$GOIMPORTS" ]; then
  echo "Failed to find goimports binary at \${GOIMPORTS}"
  echo "Use command below to install the binary and try again."
  echo "\n\tmake install-tools\n"
  exit 1
fi

if [ ! -f "\$STATIC_CHECK_OUT" ]; then
  echo "Failed to find goimports binary at \${STATIC_CHECK_OUT}"
  echo "Use command below to install the binary and try again."
  echo "\n\tmake install-tools\n"
  exit 1
fi

make lint-static-check imports-check
EOF

chmod +x ${GIT_HOOKS_PATH}/pre-commit