#!/bin/sh
ag --go --json --html -l --ignore-dir=ui --ignore=src/server/rice-box.go | entr -r -s "make run"
