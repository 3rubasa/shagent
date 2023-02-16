#!/usr/bin/bash

TRIES=1
while true
do
    echo "Setting operating mode to online, try #$TRIES"
    qmicli -d /dev/cdc-wdm0 --dms-set-operating-mode='online'
    if [ $? -ne 0 ]; then
        echo "Failed. Retrying after delay..."
        TRIES=$TRIES+1
        sleep 5
    else
        break
    fi
done

echo "Verifying the operating mode..."
qmicli -d /dev/cdc-wdm0 --dms-get-operating-mode

TRIES=1
while true
do
    echo "Bringing wwan0 down, try #$TRIES"
    
    ip link set wwan0 down
    
    if [ $? -ne 0 ]; then
        echo "Failed. Retrying after delay..."
        TRIES=$TRIES+1
        sleep 5
    else
        break
    fi
done


TRIES=1
while true
do
    echo "Setting raw_ip=Y for wwan0, try #$TRIES"
    
    echo 'Y' | tee  /sys/class/net/wwan0/qmi/raw_ip
    
    if [ $? -ne 0 ]; then
        echo "Failed. Retrying after delay..."
        TRIES=$TRIES+1
        sleep 5
    else
        break
    fi
done

TRIES=1
while true
do
    echo "Bringing wwan0 up, try #$TRIES"
    
    ip link set wwan0 up
    
    if [ $? -ne 0 ]; then
        echo "Failed. Retrying after delay..."
        TRIES=$TRIES+1
        sleep 5
    else
        break
    fi
done

TRIES=1
while true
do
    echo "Starting network, try #$TRIES"
    
    qmicli --device=/dev/cdc-wdm0 --device-open-proxy --wds-start-network="ip-type=4,apn=internet" --client-no-release-cid
    
    if [ $? -ne 0 ]; then
        echo "Failed. Retrying after delay..."
        TRIES=$TRIES+1
        sleep 5
    else
        break
    fi
done


TRIES=1
while true
do
    echo "DHCP, try #$TRIES"
    
    udhcpc -i wwan0
    
    if [ $? -ne 0 ]; then
        echo "Failed. Retrying after delay..."
        TRIES=$TRIES+1
        sleep 5
    else
        break
    fi
done

TRIES=1
while true
do
    echo "Last step, try #$TRIES"
    
    ip a s wwan0
    
    if [ $? -ne 0 ]; then
        echo "Failed. Retrying after delay..."
        TRIES=$TRIES+1
        sleep 5
    else
        break
    fi
done

echo "SUCCESS"