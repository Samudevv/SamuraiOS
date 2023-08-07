#! /bin/sh

bluetoothctl show | grep 'Powered: yes'
bluetooth_is_on=$?

if [ $bluetooth_is_on -eq 0 ]
then
  bluetoothctl power off
else
  bluetoothctl power on
fi
