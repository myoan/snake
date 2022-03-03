for pid in `ps aux | grep bin/client | grep -v grep | awk '{print $2}'`
do
	echo kill $pid
	kill -9 $pid
done
