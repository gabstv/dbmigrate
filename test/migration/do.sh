#!/bin/sh

set -e

function setup () {
    (cd ../../cmd/dbmigrate && go build && mv ./dbmigrate ../../test/migration/dbmigrate)
}

function init () {
    case "$1" in
        sqlite)
            ./dbmigrate init -p ./sqlite
            ;;
        *)
            echo "Usage: $0 init {sqlite}"
            ;;
    esac
}

setup

case "$1" in
    init)
        init $2
        ;;
    *)
        echo "Usage: $0 {init}"
        ;;
esac