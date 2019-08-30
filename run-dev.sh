#!/bin/sh
ag --go --json -l --ignore-dir=ui | entr -r -s "make run"
