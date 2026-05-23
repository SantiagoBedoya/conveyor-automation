#include <ESP32Servo.h>

#define SENSOR_GAS_PIN  34
#define SENSOR_HUM_PIN  39
#define BUZZER_PIN      27
#define RELAY_PIN       16
#define FAN_PIN         17
#define SERVO_PIN       18
#define TRIG_PIN        23
#define ECHO_PIN        19

int umbralGas        = 500;
int umbralHumedad    = 3000;
int distanciaUmbral  = 10;   // cm — objeto detectado si está a menos de esto
int contadorObjetos  = 0;
bool objetoPresente  = false;

Servo puerta;
bool puertaAbierta = false;
int pos_puerta_cerrada = 50;
int pos_puerta_abierta = 0;

long medirDistancia() {
  digitalWrite(TRIG_PIN, LOW);
  delayMicroseconds(2);
  digitalWrite(TRIG_PIN, HIGH);
  delayMicroseconds(10);
  digitalWrite(TRIG_PIN, LOW);
  long duracion = pulseIn(ECHO_PIN, HIGH, 30000); // timeout 30ms
  return duracion * 0.034 / 2;
}

void setup() {
  Serial.begin(115200);
  pinMode(BUZZER_PIN, OUTPUT);
  pinMode(RELAY_PIN, OUTPUT);
  pinMode(FAN_PIN, OUTPUT);
  pinMode(TRIG_PIN, OUTPUT);
  pinMode(ECHO_PIN, INPUT);
  digitalWrite(BUZZER_PIN, LOW);
  digitalWrite(RELAY_PIN, HIGH);
  digitalWrite(FAN_PIN, LOW);

  puerta.attach(SERVO_PIN);
  puerta.write(pos_puerta_cerrada);

  Serial.println("Calentando sensor MQ-02...");
  delay(20000);
}

void loop() {
  int valorGas     = analogRead(SENSOR_GAS_PIN);
  int valorHumedad = analogRead(SENSOR_HUM_PIN);
  long distancia   = medirDistancia();

  // --- Contador de objetos ---
  if (distancia > 0 && distancia < distanciaUmbral) {
    if (!objetoPresente) {
      contadorObjetos++;
      objetoPresente = true;
      Serial.print(">>> Objeto detectado! Total: ");
      Serial.println(contadorObjetos);
    }
  } else {
    objetoPresente = false; // Objeto salió, listo para el siguiente
  }

  Serial.print("Gas: ");
  Serial.print(valorGas);
  Serial.print(" | Humedad: ");
  Serial.print(valorHumedad);
  Serial.print(" | Distancia: ");
  Serial.print(distancia);
  Serial.print(" cm | Objetos: ");
  Serial.println(contadorObjetos);

  // --- Alarma por gas ---
  if (valorGas > umbralGas) {
    Serial.println(">>> GAS DETECTADO");
    digitalWrite(BUZZER_PIN, HIGH);
    digitalWrite(FAN_PIN, HIGH);
  }
  // --- Humedad alta ---
  else if (valorHumedad < umbralHumedad) {
    Serial.println(">>> HUMEDAD ALTA - Abriendo puerta");
    digitalWrite(BUZZER_PIN, HIGH);
    digitalWrite(RELAY_PIN, LOW); // apagando motor
    if (!puertaAbierta) {
      puerta.write(pos_puerta_abierta);
      puertaAbierta = true;
    }
  }
  // --- Todo normal ---
  else {
    digitalWrite(BUZZER_PIN, LOW);
    digitalWrite(RELAY_PIN, HIGH); // encendiendo motor
    digitalWrite(FAN_PIN, LOW);
    if (puertaAbierta) {
      puerta.write(pos_puerta_cerrada);
      puertaAbierta = false;
    }
  }

  delay(100); // Más rápido para no perder objetos en la banda
}