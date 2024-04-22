## Build Examples 

### Prerequisites
- Docker for building the ACAP applications
- [goxisbuilder](https://github.com/Cacsjep/goxisbuilder)


### Build example axevent send
``` shell
go install github.com/Cacsjep/goxisbuilder@latest
git clone https://github.com/Cacsjep/goxis_examples
cd goxis_examples
goxisbuilder.exe --appdir="./axevent/send"
```

#### How each example should be builded
``` shell
go install github.com/Cacsjep/goxisbuilder@latest
goxisbuilder --appdir="./axevent/send"
goxisbuilder --appdir="./axevent/subscribe"
goxisbuilder --appdir="./axlarod/classify" -files converted_model.tflite
goxisbuilder --appdir="./axlarod/object_detection" -files ssd_mobilenet_v2_coco_quant_postprocess.tflite
goxisbuilder --appdir="./axlarod/yolov5" -files yolov5n.tflite
goxisbuilder --appdir="./axlicense" 
goxisbuilder --appdir="./axoverlay"
goxisbuilder --appdir="./axparameter"
goxisbuilder --appdir="./axstorage"
goxisbuilder --appdir="./vapix"
goxisbuilder --appdir="./vdostream"
goxisbuilder --appdir="./webserver"
```

Examples are really close to existing C examples of the [AXIS Native SDK repo](https://github.com/AxisCommunications/acap-native-sdk-examples).

> [!NOTE]  
> The examples use mainly acapapp package, all examples could also just written using the diffrent
go packages directly without using acapapp package.

| Example         | Description |
|-----------------|--------------|
| `axevent/send`	            | Demonstrate how to declare and send an event using acapapp package     |
| `axevent/subscribe`	        | Demonstrate how to subscribe to an Virutal Input state change          |
| `axoverlay`	                | Render rects via axolveray api                                         |
| `axlarod/classify`	        | Classification example with larod and vdo api  (artpec-8)              |
| `axlarod/object_detection`	| Object detection example with larod/vdo and overlay api api  (artpec-8)|
| `axlarod/yolov5           `	| Yolov5 example with larod/vdo and overlay api api  (artpec-8)          |
| `axparameter`                 | Demonstrate how to get an parameter and listen to changes              |
| `axstorage`                   | Interact with axstorage api                                            |
| `license` 	                | Show how to obtain the license state                                   |
| `vdostream` 	                | High level wrapper demonstration to get video frames (stream)          |
| `webserver`                   | Reverse proxy webserver with fiber                                     |