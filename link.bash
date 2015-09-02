#!/bin/bash -x

for f in *; do
    if [[ "$f" != _* ]] && [[ -d $f ]]; then
        yes | rm -r "$GOPATH/src/github.com/moo-mou/maingo/$f"
        ln -s "`pwd`/$f" "$GOPATH/src/github.com/moo-mou/maingo/$f"
    fi
done

for f in *; do
    if [[ "$f" != _* ]] && [[ -d $f ]]; then
        go build "github.com/moo-mou/maingo/$f"
    fi
done
