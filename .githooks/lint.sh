#!/bin/sh
lint() {
    make lint || exit 1
}

lint