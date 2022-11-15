#!/usr/bin/env bash

basedir=$(dirname "$0")/..

export $(cat $basedir/.env | grep -ve ^\# | xargs -0)
