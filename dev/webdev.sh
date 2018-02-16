#!/bin/bash

# Must be run from the repo root.

# Kill the background processes (watching JS and SASS) on exit.
trap 'kill $(jobs -p)' EXIT

mkdir -p build/web
ln -s -f ../../web/index.html build/web/index.html
ln -s -f ../../web/style.css build/web/style.css

echo 'watching'

# Watch the JS and SASS for changes and automatically rebuild.
./node_modules/.bin/watchify -v -d -e web/index.js -o build/web/build.js

# block forever
cat
