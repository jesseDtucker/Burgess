﻿export var RAW_DB_URL = 'mongodb://localhost:27017/raw';
export var POS_DB_URL = 'mongodb://localhost:27017/retailers';
export var BLUETOOTHLOCATION_DB_URL = 'mongodb://localhost:27017/bluetoothLocations';
export var RAW_WIFI_COLLECTION = 'rawWifi';
export var RAW_BLUETOOTH_COLLECTION = 'rawBluetooth';
export var BLUETOOTHLOCATIONS_COLLECTION = 'locations';
export var POS_COLLECTION = 'position';
export var BT_ERROR = 0.5; // We are placing the regions 5 meters apart giving us an approx error of 2.5 meters on position