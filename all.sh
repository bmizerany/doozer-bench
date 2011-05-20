for d in doozer-bench doozer-report
do
    pushd cmd/$d
    make install
    popd
done
