#!/usr/bin/env python
import paho.mqtt.client as mqtt
import json

# Imports the Google Cloud client library
from google.cloud import bigquery

# Instantiates a client
client = bigquery.Client()
ds = client.dataset('iothackdata')
table = ds.table('roadreflectorevent')
table.reload()


def on_connect(client, userdata, flags, rc):
    client.subscribe("reflectors")
    print("Successfully connected to MQTT server!")


def on_message(client, userdata, msg):
    with open("data.json", "a") as f:
        f.write(str(msg.payload) + ", ")
    data = json.loads(msg.payload)
    raw_record = [data['Temperature'], data['Source'], data['Time'], data['Name']]
    record = tuple(raw_record)
    row = [record]
    table.insert_data(row)


client = mqtt.Client()
client.on_connect = on_connect
client.on_message = on_message

client.connect("104.154.46.156", 1883, 60)

client.loop_forever()
