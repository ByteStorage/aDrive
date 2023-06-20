##/bin/bash
# killall aDrive
kill -9 $(ps aux | grep '[a]Drive' | awk '{print $2}') &
rm -rf data1
rm -rf data2
rm -rf data3
rm -rf Cluster
rm -rf logs
rm -rf datanode1.log
rm -rf datanode2.log
rm -rf datanode3.log
rm -rf namenode.log
rm -rf aDrive
echo "Stop Server Successfully!"
