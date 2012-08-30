#!/usr/bin/env bash

TEST_PKGS="
github.com/scotch/aego/v1/auth
github.com/scotch/aego/v1/auth/dev
github.com/scotch/aego/v1/auth/password
github.com/scotch/aego/v1/auth/profile
github.com/scotch/aego/v1/config
github.com/scotch/aego/v1/ds
github.com/scotch/aego/v1/ds/appengine/datastore
github.com/scotch/aego/v1/ds/appengine/memcache
github.com/scotch/aego/v1/ds/memory
github.com/scotch/aego/v1/user/email
github.com/scotch/aego/v1/user
"

echo '# Testing AEGo packages.'

for f in $TEST_PKGS
do
  go test -i $f
  go test $f
done
