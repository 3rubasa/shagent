# shagent
Smart Home Agent, running on Raspberry Pi

### Key Functions:
- Based on the current temperature indoors and outdoors and on the information about weather forcast, controls operation of a boiler to support specified temperatrue range in the building
- Integrates with OpenWeather API to obtain weather forecasts
- Collects data from sensors and sends them to a channel at thingspeak.com
- Provides REST API, GRPC API and CLI app for querying info about the current state and manual control of the devices
- Integrates with LTE SIM7600 module through serial interface (AT commands) to query basic info
- Turning room lights on/off by schedule or manually
- shwatchdog is an optional watchdog service that monitors the state of the smart home agent service and the presence of virtual network interface (openVPN) and restarts the Raspberry Pi in case of failure
- Access to the devices is available through SSH over VPN connection

### Supports a set of sensors and controllers:
- DHT20 temperature and humidity sensor
- ESP8266 WiFi temperature module
- SIM7600 LTE module controller
- 3-channel GPIO switch (relay)
- Sonoff R3RF WiFi switch (relay)

### CLI app commands:
```
$shagent_cli roomlight_on # turn on the light in the main room
$shagent_cli roomlight_off # turn off the light in the main room
$shagent_cli camlight_on # turn on the camera light
$shagent_cli camlight_off # turn off the camera light
$shagent_cli cell_balance # print sim balance
$shagent_cli cell_inet # print sim internet balance
$shagent_cli cell_tariff # print sim current tariff
$shagent_cli cell_phone # print sim phone number
$shagent_cli boiler_on # force turn the boiler on
$shagent_cli boiler_off # force turn the boiler off
```
