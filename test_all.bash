#!/usr/bin/env bash

TEST_PKGS="
github.com/scotch/hal/auth
github.com/scotch/hal/auth/dev
github.com/scotch/hal/auth/password
github.com/scotch/hal/config
github.com/scotch/hal/ds
github.com/scotch/hal/ds/appengine/datastore
github.com/scotch/hal/ds/appengine/memcache
github.com/scotch/hal/ds/memory
github.com/scotch/hal/user
github.com/scotch/hal/user_profile
"

echo '# Testing HAL packages.'

for f in $TEST_PKGS
do
  go test -i $f
  go test $f
done
