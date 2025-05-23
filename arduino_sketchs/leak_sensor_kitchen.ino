#include <ESP8266WiFi.h>
#include <WiFiClient.h>

const char* ssid = "MY_WIFI";
const char* password = "MY_WIFI_PASSWORD";
const char* server = "MY_SMART_HOME_SERVER_ENDPOINT";
const char* port = "MY_SMART_HOME_SERVER_PORT";

const int leakSensorPin = D1;
int leakStatus = 0;

void setup() {
    Serial.begin(115200);
    pinMode(leakSensorPin, INPUT);

    // Подключение к Wi-Fi
    WiFi.begin(ssid, password);
    Serial.print("Connecting WiFi");
    while (WiFi.status() != WL_CONNECTED) {
        delay(1000);
        Serial.print(".");
    }
    Serial.println("Connected WiFi");
}

void loop() {
    // Чтение статуса датчика протечки
    leakStatus = digitalRead(leakSensorPin);

    // Отправка статуса на сервер
    if (leakStatus == HIGH) {
        Serial.println("ALARM!!! Leak detected!");
        sendStatus("Leak Detected");
    } else {
        Serial.println("Leak not detected");
        sendStatus("No Leak");
    }

    // Проверка каждые 30 секунд
    delay(30000);
}

void sendStatus(const char* state) {
    String jsonData = "{"room": "kitchen", "sensor_state": " + String(state) + "}";

    WiFiClient client;
    if (client.connect(server, port)) {
        client.println("POST /api/arduino/senors/leak HTTP/1.1");
        client.println("Host: " + server + ":" + port);
        client.println("Content-Type: application/json");
        client.print("Content-Length: ");
        client.println(jsonData.length());
        client.println("Connection: close");
        client.println();
        client.println(jsonData);

        while (client.available()) {
            String resp = client.readStringUntil('\n');
            Serial.println("Response: " + resp);
        }
    } else {
        Serial.println("Connection failed");
    }
    
    client.stop();
}
