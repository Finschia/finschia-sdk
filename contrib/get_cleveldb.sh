set -e

PWD=$(pwd)
compressed=v1.20.tar.gz
leveldb_dir=leveldb-1.20

sudo="sudo"
if [ "$2" = "docker" ]; then
  sudo=""
fi

if [ "$1" = "yes" ]; then
  if [ ! -d ${PWD}/${leveldb_dir} ] || [ -e ${PWD}/${compressed} ] ; then
    rm -rf ${leveldb_dir}
    rm -f ${compressed}
    wget https://github.com/google/leveldb/archive/${compressed}
    tar -zxvf ${compressed}
    cd ${leveldb_dir}
    make
    cd ..
  fi
  $sudo cp -r ${leveldb_dir}/out-static/lib* ${leveldb_dir}/out-shared/lib* /usr/lib/
  $sudo cp -r ${leveldb_dir}/include/leveldb /usr/include/
  rm -f ${compressed}
fi
