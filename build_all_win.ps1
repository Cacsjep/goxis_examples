go install github.com/Cacsjep/goxisbuilder@v1.2.6
goxisbuilder.exe -appdir "./axevent/send"
goxisbuilder.exe -appdir "./axevent/subscribe"
goxisbuilder.exe -appdir "./axevent/multiple_subscribe"
goxisbuilder.exe -appdir "./axlarod/classify" -files converted_model.tflite
goxisbuilder.exe -appdir "./axlarod/object_detection" -files ssd_mobilenet_v2_coco_quant_postprocess.tflite
goxisbuilder.exe -appdir "./axlarod/face_tracking" -files ssd_mobilenet_v2_face_quant_postprocess.tflite
goxisbuilder.exe -appdir "./axlarod/yolov5" -files yolov5n.tflite
goxisbuilder.exe -appdir "./axlicense" 
goxisbuilder.exe -appdir "./axoverlay/pixel_array"
goxisbuilder.exe -appdir "./axoverlay/rects_text"
goxisbuilder.exe -appdir "./axoverlay/png_sequence" -files zinta
goxisbuilder.exe -appdir "./axparameter"
goxisbuilder.exe -appdir "./axstorage"
goxisbuilder.exe -appdir "./vapix/list_params"
goxisbuilder.exe -appdir "./vapix/websocket_metadatastream"
goxisbuilder.exe -appdir "./vdostream"
goxisbuilder.exe -appdir "./webserver"
goxisbuilder.exe -appdir "./axmdb/consume-scene-metadata"
goxisbuilder.exe -appdir "./axmdb/scene-metadata-overlay"