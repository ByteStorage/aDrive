## /bin/bash
go build aDrive.go
nohup ./aDrive namenode --master --host 127.0.0.1 --port 9999 |tee namenode.log &
sleep 1
nohup ./aDrive datanode --namenode 127.0.0.1:9999 --host 127.0.0.1 --port 7000 --path data1/ |tee datanode1.log &
nohup ./aDrive datanode --namenode 127.0.0.1:9999 --host 127.0.0.1 --port 7001 --path data2/ |tee datanode2.log &
nohup ./aDrive datanode --namenode 127.0.0.1:9999 --host 127.0.0.1 --port 7002 --path data3/ |tee datanode3.log &
echo "Starting Server Successfully!"