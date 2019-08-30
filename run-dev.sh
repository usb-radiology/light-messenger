#!/bin/sh
ag --go --json --html -l --ignore-dir=ui | entr -r -s "make run"
