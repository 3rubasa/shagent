# shagent
Smart Home Agent, running on Raspberry Pi

Supports a set of sensors and controllers:
- DHT20 temperature and humidity sensor
- ESP8266 WiFi temperature module
- SIM7600 LTE module controller
- 3-channel GPIO switch (relay)
- Sonoff R3RF WiFi switch (relay)

Collects data from sensors and sends it to a channel at thingspeak.com
Exposes REST API for information querying and control
Provides CLI app for information querying and control

shwatchdog is an optional watchdog service that monitors the state of the smart home agent service and the presence of virtual network interface (openVPN) and restarts the Raspberry Pi in case of a failure
