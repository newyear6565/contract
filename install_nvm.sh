# Copyright (C) 2017-2019 go-nebulas authors
#
# This file is part of the go-nebulas library.
#
# the go-nebulas library is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.

# the go-nebulas library is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.

#!/bin/bash

# usage: source native-libs.sh

CUR_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}"  )" >/dev/null && pwd  )"
#CUR_DIR="$( pwd )"
OS="$(uname -s)"
#default region is China
REGION="China"

if [ "$OS" = "Darwin" ]; then
  LOGICAL_CPU=$(sysctl -n hw.ncpu)
  DYLIB="dylib"
else
  LOGICAL_CPU=$(cat /proc/cpuinfo |grep "processor"|wc -l)
  DYLIB="so"
fi
PARALLEL=$LOGICAL_CPU

if [ "$REGION" = "China" ]; then
  SOURCE_URL="http://develop-center.oss-cn-zhangjiakou.aliyuncs.com"
else
  SOURCE_URL="https://s3-us-west-1.amazonaws.com/develop-center"
fi

install_nvm() {
  nvm_lib=$CUR_DIR/native-lib
  rm -rf $nvm_lib
  mkdir -p $nvm_lib
  pushd $nvm_lib
  wget $SOURCE_URL/setup/nvm/lib_nvm_$OS.tar.gz -O lib_nvm_$OS.tar.gz
  tar -zxvf lib_nvm_$OS.tar.gz
  sudo cp -Rf lib_nvm_$OS/* /usr/local/lib
  
  libs=`ls $nvm_lib/lib_nvm_$OS|grep .$DYLIB`
  for lib in $libs; do
    cp $nvm_lib/lib_nvm_$OS/$lib  $nvm_lib/$lib
  done
  
  rm -rf lib_nvm_$OS
  rm -rf lib_nvm_$OS.tar.gz
  popd
}

install_nvm
