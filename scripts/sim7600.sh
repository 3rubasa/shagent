#! /usr/bin/bash

echo "Setting operating mode to online..."
qmicli -d /dev/cdc-wdm0 --dms-set-operating-mode='online'
if [ $? -ne 0 ]; then
	echo "Failed. Exiting..."
	exit 1
fi

echo "Verifying the operating mode..."
qmicli -d /dev/cdc-wdm0 --dms-get-operating-mode

echo "Bringing wwan0 down..."
ip link set wwan0 down
if [ $? -ne 0 ]; then
	echo "Failed. Exiting..."
	exit 1
fi

echo "Setting raw_ip=Y for wwan0..."
echo 'Y' | tee  /sys/class/net/wwan0/qmi/raw_ip
if [ $? -ne 0 ]; then
	echo "Failed. Exiting..."
	exit 1
fi

echo "Bringing wwan0 up..."
ip link set wwan0 up
if [ $? -ne 0 ]; then
	echo "Failed. Exiting..."
	exit 1
fi

echo "Starting network..."
qmicli --device=/dev/cdc-wdm0 --device-open-proxy --wds-start-network="ip-type=4,apn=internet" --client-no-release-cid
if [ $? -ne 0 ]; then
	echo "Failed. Exiting..."
	exit 1
fi

echo "DHCP..."
udhcpc -i wwan0
if [ $? -ne 0 ]; then
	echo "Failed. Exiting..."
	exit 1
fi


echo "Last step..."
ip a s wwan0
if [ $? -ne 0 ]; then
	echo "Failed. Exiting..."
	exit 1
fi

echo "SUCCESS"