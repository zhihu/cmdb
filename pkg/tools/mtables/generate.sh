#!/bin/bash
set -e
echo 'package mtables' > ./tpl.gen.go
echo 'var bodyTemplate =`' >> ./tpl.gen.go
cat ./body.tpl >> ./tpl.gen.go
echo '`'>> ./tpl.gen.go


echo 'var headerTemplate =`' >> ./tpl.gen.go
cat ./header.tpl >> ./tpl.gen.go
echo '`'>> ./tpl.gen.go