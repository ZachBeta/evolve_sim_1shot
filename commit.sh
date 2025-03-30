#! /bin/bash

git add .

# take all args and pass them to git commit
git commit -m "$*"

