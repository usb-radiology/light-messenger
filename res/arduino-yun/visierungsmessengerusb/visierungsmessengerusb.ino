
/* Dieser Source Code kann ueber ein USB - Mini USB Kabel auf einen beliebigen Arduino Yun uebertragen werden.
 * 
 * DIESER SOURCE CODE IST FUER DIE ARDUINO YUN MIT LED-STREIFEN DIE IN DER ABTEILUNG AOD STEHEN!
 */
#include "LPD8806.h"
#include <Bridge.h>
#include <HttpClient.h>
#include <FileIO.h>

// init variable to store what department the arduino fits
int arduinoDepartment = 0;
int AOD=1;
int CTD=2; 
int MSK=3;
int NR=4;

String serverRestPrefix = "http://10.5.72.154:9200/nce-rest/arduino-status/";

// init URL of serverside get LED status REST controller
String ledStatusRestControllerUrl = "";
String aodLedStatusRestControllerUrl = serverRestPrefix+"aod-open-notifications";
String ctdLedStatusRestControllerUrl = serverRestPrefix+"ctd-open-notifications";
String mskLedStatusRestControllerUrl = serverRestPrefix+"msk-open-notifications";
String nrLedStatusRestControllerUrl = serverRestPrefix+"nr-open-notifications";

// init URL of serverside post Arduino status REST controller
String arduinoStatusRestControllerUrl = "";
String aodArduinoStatusRestControllerUrl = serverRestPrefix+"aod-status";
String ctdArduinoStatusRestControllerUrl = serverRestPrefix+"ctd-status";
String mskArduinoStatusRestControllerUrl = serverRestPrefix+"msk-status";
String nrArduinoStatusRestControllerUrl = serverRestPrefix+"nr-status";

//LED INIT START ***********************************************
// Anzahl RGB LEDs
int nLEDs = 8;

// Die zwei Arduino-Pins, mit denen der LED Strip verbunden ist
int dataPin  = 8;
int clockPin = 9;

// Den Strip initialisieren
LPD8806 strip = LPD8806(nLEDs, dataPin, clockPin);
//LED INIT END ***********************************************


//initialize bools for lights
boolean isLowMode = false;
boolean isMediumMode = false;
boolean isHighMode = false;
boolean isLedsOff = false;


void setup() {
  //LED SETUP START ****************************
  // LED-Strip start
  strip.begin();
  // All LEDs off
  strip.show();
  //LED SETUP END ******************************
  
  // Bridge takes about two seconds to start up
  // it can be helpful to use the on-board LED
  // as an indicator for when it has initialized
  pinMode(13, OUTPUT);
  digitalWrite(13, LOW);
  Bridge.begin();
  delay(2000);
  digitalWrite(13, HIGH);

  //SerialUSB.begin(9600);
  //
  delay(1000);

  //start filesystem to read from file on sd card
  FileSystem.begin();
  delay(2000);
  
  //while (!SerialUSB); // wait for a serial connection

  //TODO read in arduinos department
  readDepartmentFromSDCard();
  //arduinoDepartment = AOD;

   //find the right urls
  findDepartmentSpecificUrls();
}


//MAIN LOOP
void loop() {
 
  //post Arduino Online status to Server
  postArduinoStatusAsYunGet();
  
  //get LED status from server
  httpGetLedStatus();

  //turn leds on or off
  manageLeds();
  
}

void readDepartmentFromSDCard(){
  //Serial.println("try connect to dep.txt: ");
  File departmentFile = FileSystem.open("/mnt/sda1/dep.txt");
  
  if (departmentFile.available()){
    //Serial.println("dep.txt: ");
    char department = departmentFile.read();
    // Serial.println(department);
    //this is a char to int cast... xD
    arduinoDepartment = department - '0';
    
  }
  else{
    //Serial.println("Error opening file dep.txt!");
  }
  departmentFile.close();
  //Serial.println();
  delay(1000);
}

void findDepartmentSpecificUrls(){
  String ledLogPrefix = "URL to get led status is: ";
  String arduinoStatusLogPrefix = "URL to post arduino status is: ";
  //SerialUSB.println("Try to find the department specific URL...");
  //SerialUSB.println("arduinoDepartment value is: " + String(arduinoDepartment));
  if(arduinoDepartment==AOD){
    arduinoStatusRestControllerUrl = aodArduinoStatusRestControllerUrl;
    //SerialUSB.println(arduinoStatusLogPrefix+arduinoStatusRestControllerUrl);
    ledStatusRestControllerUrl = aodLedStatusRestControllerUrl;
    //SerialUSB.println(ledLogPrefix+ledStatusRestControllerUrl);
  }
  else if (arduinoDepartment==CTD){
    arduinoStatusRestControllerUrl = ctdArduinoStatusRestControllerUrl;
    //SerialUSB.println(arduinoStatusLogPrefix+arduinoStatusRestControllerUrl);
    ledStatusRestControllerUrl = ctdLedStatusRestControllerUrl;
    //SerialUSB.println(ledLogPrefix+ledStatusRestControllerUrl);
  }
  else if (arduinoDepartment==MSK){
    arduinoStatusRestControllerUrl = mskArduinoStatusRestControllerUrl;
    //SerialUSB.println(arduinoStatusLogPrefix+arduinoStatusRestControllerUrl);
    ledStatusRestControllerUrl = mskLedStatusRestControllerUrl;
    //SerialUSB.println(ledLogPrefix+ledStatusRestControllerUrl);
  }
  else if (arduinoDepartment==NR){
    arduinoStatusRestControllerUrl = nrArduinoStatusRestControllerUrl;
    //SerialUSB.println(arduinoStatusLogPrefix+arduinoStatusRestControllerUrl);
    ledStatusRestControllerUrl = nrLedStatusRestControllerUrl;
    //SerialUSB.println(ledLogPrefix+ledStatusRestControllerUrl);
  }
  else{
    //SerialUSB.println("Invalid value for arduinoDepartment! Value is: "+String(arduinoDepartment));
  }
  //SerialUSB.println();
}

void postArduinoStatusAsYunGet(){
  // Initialize the client
  HttpClient clientOne;
  
  //SerialUSB.println("Making get request to URL "+arduinoStatusRestControllerUrl+"...");
   //Sends alive status to server
   clientOne.get(arduinoStatusRestControllerUrl);
     while (clientOne.available()) {
      char c = clientOne.read();
      //SerialUSB.print(c);         
      }
  //SerialUSB.println();
  //SerialUSB.println();
  clientOne.close();
  delay(1000);
}


void httpGetLedStatus(){
  // Initialize the client library
  HttpClient client;

   //Sends get that gets led status entry from server
   client.get(ledStatusRestControllerUrl);
     while (client.available()) {
      client.readStringUntil(';');
      String notificationStatus = client.readStringUntil(';');
      //SerialUSB.println("Notification status: " + notificationStatus);
      if(notificationStatus == "1"){
        
        //reading notification mode 
        String mode = client.readStringUntil(';');
        //SerialUSB.print("MODE: " + mode);
        //SerialUSB.println();
        delay(100);

        if(mode == "LOW"){
          isLowMode = true;
          //SerialUSB.println("bool is: " + mode);
        }
        else if(mode == "MEDIUM"){
          isMediumMode = true;
          //SerialUSB.println("bool is: " + mode);
        }
        else if (mode == "HIGH"){
          isHighMode = true;
          //SerialUSB.println("bool is: " + mode);
        }
        else{
          //error
          //SerialUSB.println("error: Mode Format Exception ");
        }
      }else if (notificationStatus == "0"){
        //No notificataion is open for department AOD
        isLedsOff = true;
        //SerialUSB.println("No notification is open for department AOD!");
      }else{
        //No notificataion is on
        //SerialUSB.println("error: status should be 0 or 1!");
      }
              
    }
  client.close();
  //SerialUSB.println();
  delay(1000);
}


void manageLeds(){
  if(isLowMode){
      //SerialUSB.println("Light turn on green ");
      //ledBlinkAndStayColor(0,127,0,50);
      makeLEDStripGreen();
    }
    else if (isMediumMode){
      //SerialUSB.println("Light turns on yellow");
      ledBlinkAndStayColor(127,127,0,50);
    }
    else if(isHighMode){
      //SerialUSB.println("Light turns on red");
      ledBlinkAndStayColor(127,0,0,50);
    }
    else if(isLedsOff){
      //SerialUSB.println("Light turns off");
      makeLEDStripOff();
    }
    else{
      //error
      //SerialUSB.println("NO LED LIGHT is on");
      ledBlinkAndStayColor(127,0,127,500);
    }
    //SerialUSB.println("At the end of manageLeds()!");
    //SerialUSB.println();
    delay(15000);

    isLowMode = false;
    isMediumMode = false;
    isHighMode = false;
    isLedsOff = false;
}

void ledBlinkAndStayColor(int red, int green, int blue, int delayTime){
  for(int i = 0; i<20; i++){
    //is even make led on in color defined
    if(i%2==0){
      showColor(red, green, blue);
      delay(delayTime);
    }
    //is uneven make led of 
    else{
      showColor(0, 0, 0);
      delay(delayTime);
    }
  }
  showColor(red, green, blue);
  delay(15000);
  
}

void makeLEDStripGreen(){
  showColor(0, 127, 0);
}

void makeLEDStripOff(){
  showColor(0, 0, 0);
}

void showColor(int r, int g, int b) {
  
  for (int i=0; i<nLEDs; i++) {
    strip.setPixelColor(i, strip.Color(r , g, b));    
    delay(20);
    strip.show();
  }
  strip.show();  
  
}
