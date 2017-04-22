import paho.mqtt.client as mqtt


def on_connect(client, userdata, flags, rc):
    client.subscribe("reflectors")
    print("Successfully connected to MQTT server!")


def on_message(client, userdata, msg):
    with open("data.json", "a") as f:
        f.write(str(msg.payload))


client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect("104.154.46.156", 1883, 60)

client.loop_forever()
