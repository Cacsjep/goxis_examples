#!/bin/bash
go install github.com/Cacsjep/goxisbuilder@latest
goxisbuilder -appdir "./axevent/send"
goxisbuilder -appdir "./axevent/subscribe"
goxisbuilder -appdir "./axlarod/classify" -files converted_model.tflite
goxisbuilder -appdir "./axlarod/object_detection" -files ssd_mobilenet_v2_coco_quant_postprocess.tflite
goxisbuilder -appdir "./axlarod/yolov5" -files yolov5n.tflite
goxisbuilder -appdir "./axlicense" 
goxisbuilder -appdir "./axoverlay"
goxisbuilder -appdir "./axparameter"
goxisbuilder -appdir "./axstorage"
goxisbuilder -appdir "./vapix"
goxisbuilder -appdir "./vdostream"
goxisbuilder -appdir "./webserver"