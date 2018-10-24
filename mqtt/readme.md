## ENV variables

| ENV            | Default value |
|:---------------|:-------------:|
| MQTT_HOST      | localhost     |
| MQTT_PORT      | 1883          |
| MQTT_USERNAME  | guest         |
| MQTT_PASSWORD  | guest         |
| MQTT_CLIENT_ID | util.UUID()   |

## Usage

* **Create mqtt config**
```
mqttConfig := mqtt.NewConfig()
```

* **Optionally, customize config values (*these are defaults*)**
```
mqttConfig.KeepAlive = time.Second * 15
mqttConfig.CleanSession = true
mqttConfig.AutoReconnect = true
mqttConfig.MsgChanDept = 100
mqttConfig.PubQoS = 0
mqttConfig.SubQoS = 0
```

* **Create mqtt connection**
```
_mqtt, err := mqtt.NewConnection(mqttConfig)
if err != nil {
    return err
}
```

* **Publish payload to mqtt topic**
```
if token := _mqtt.Publish("my/topic", []byte("message")); token.Wait() && token.Error() != nil {
    log.Print("mqtt publish error: ", token.Error())
}
```

* **Subscribe to mqtt topic**
```
if token := _mqtt.Subscribe("my/topic", func(c mqtt.Client, m mqtt.Message) {
    log.Print(m)
}); token.Wait() && token.Error() != nil {
    log.Print("mqtt subscribe error: ", token.Error())
}
```