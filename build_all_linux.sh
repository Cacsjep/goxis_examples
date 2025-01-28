#!/bin/bash
go install github.com/Cacsjep/goxisbuilder@latest
goxisbuilder -appdir "./axevent/send"
goxisbuilder -appdir "./axevent/subscribe"
goxisbuilder -appdir "./axevent/multiple_subscribe"
goxisbuilder -appdir "./axlarod/classify" -files converted_model.tflite
goxisbuilder -appdir "./axlarod/object_detection" -files ssd_mobilenet_v2_coco_quant_postprocess.tflite
goxisbuilder -appdir "./axlarod/face_tracking" -files ssd_mobilenet_v2_face_quant_postprocess.tflite
goxisbuilder -appdir "./axlarod/yolov5" -files yolov5n.tflite
goxisbuilder -appdir "./axlicense" 
goxisbuilder -appdir "./axoverlay/rects_text"
goxisbuilder -appdir "./axoverlay/pixel_array"
goxisbuilder -appdir "./axoverlay/png_sequence" -files zinta
goxisbuilder -appdir "./axparameter"
goxisbuilder -appdir "./axstorage"
goxisbuilder -appdir "./vapix/list_params"
goxisbuilder -appdir "./vapix/websocket_metadatastream"
goxisbuilder -appdir "./vdostream"
goxisbuilder -appdir "./webserver"
goxisbuilder -appdir "./axmdb/consume-scene-metadata"
goxisbuilder -appdir "./axmdb/scene-metadata-overlay"