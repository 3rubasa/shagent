# shagent
Smart Home Agent, running on Raspberry Pi

### Supports a set of sensors and controllers:
- DHT20 temperature and humidity sensor
- ESP8266 WiFi temperature module
- SIM7600 LTE module controller
- 3-channel GPIO switch (relay)
- Sonoff R3RF WiFi switch (relay)

### Key Functions:
- Based on the current temperature indoors and outdoors and on the information about weather forcast, controls operation of a boiler to support specified temperatrue range in the building
- Integrates with OpenWeather API to obtain weather forecasts
- Collects data from sensors and sends them to a channel at thingspeak.com
- Provides REST API, GRPC API and CLI app for information querying and control
- shwatchdog is an optional watchdog service that monitors the state of the smart home agent service and the presence of virtual network interface (openVPN) and restarts the Raspberry Pi in case of failure

### CLI app commands:
```
$shagent_cli roomlight_on *# turn on the light in the main room*
$shagent_cli roomlight_off *# turn off the light in the main room*
$shagent_cli camlight_on *# turn on the camera light*
$camlight_off *# turn off the camera light*
$cell_balance *# print sim balance*
$cell_inet *# print sim internet balance*
$cell_tariff *# print sim current tariff*
$cell_phone *# print sim phone number*
$boiler_on *# force turn the boiler on*
$boiler_off *# force turn the boiler off*
```
