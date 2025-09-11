#!/bin/bash

# Default old module name
OLD_MODULE="github.com/alimosavifard/zyros-backend"

# Use MODULE_NAME from environment variable
NEW_MODULE=${MODULE_NAME:-github.com/alimosavifard/zyros-backend}

# Edit go.mod to change module name
go mod edit -module "$NEW_MODULE"

# Replace all import paths in .go files
find . -type f -name '*.go' -exec sed -i -e "s,$OLD_MODULE,$NEW_MODULE,g" {} \;

# Tidy dependencies
go mod tidy