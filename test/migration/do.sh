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

function apply () {
    case "$1" in
        sqlite)
            ./dbmigrate -v --config ./sqlite/dbmigrate.toml apply
            ;;
        *)
            echo "Usage: $0 apply {sqlite}"
            ;;
    esac
}

setup

case "$1" in
    init)
        init $2
        ;;
    apply)
        apply $2
        ;;
    *)
        echo "Usage: $0 {init|apply}"
        ;;
esac