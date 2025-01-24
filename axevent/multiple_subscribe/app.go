package main

import (
	"github.com/Cacsjep/goxis/pkg/acapapp"
	"github.com/Cacsjep/goxis/pkg/axevent"
)

// This example demonstrates how to subscribe to a bunch of events, without any filters on them set.
// axevent holds a lot of predefinied events like DeviceIoVirtualInputEventKvs, just create your own Events
// when u need a specific event. Look how DeviceIoVirtualInputEventKvs is build in axevent package.
// using UnmarshalEvent you can convert the event like json.Unmarshal to a struct, here we just print the event name.
//
//
// Tipp: Use Axis Metadata Monitor to see live which events are produced by camera
// https://www.axis.com/developer-community/axis-metadata-monitor

var events = []struct {
	Name  string
	Event *axevent.AXEventKeyValueSet
}{
	{"DeviceIoVirtualInputEventKvs", axevent.DeviceIoVirtualInputEventKvs(nil, nil)},
	{"DeviceIoSupervisedPortEventKvs", axevent.DeviceIoSupervisedPortEventKvs(nil, nil, nil)},
	{"DeviceIoOutputPortEventKvs", axevent.DeviceIoOutputPortEventKvs(nil, nil)},
	{"DeviceIoPortEventKvs", axevent.DeviceIoPortEventKvs(nil, nil)},
	{"DeviceSensorPIREventKvs", axevent.DeviceSensorPIREventKvs(nil, nil)},
	{"DeviceLightStatusEventKvs", axevent.DeviceLightStatusEventKvs(nil, nil)},
	{"DeviceStatusSystemReadyEventKvs", axevent.DeviceStatusSystemReadyEventKvs(nil)},
	{"DeviceStatusTemperatureInsideEventKvs", axevent.DeviceStatusTemperatureInsideEventKvs(nil)},
	{"DeviceStatusTemperatureAboveEventKvs", axevent.DeviceStatusTemperatureAboveEventKvs(nil)},
	{"DeviceStatusTemperatureAboveOrBelowEventKvs", axevent.DeviceStatusTemperatureAboveOrBelowEventKvs(nil)},
	{"DeviceStatusTemperatureBelowEventKvs", axevent.DeviceStatusTemperatureBelowEventKvs(nil)},
	{"DeviceHardwareFailurePowerSupplyFailurePTZPowerFailureEventKvs", axevent.DeviceHardwareFailurePowerSupplyFailurePTZPowerFailureEventKvs(nil, nil)},
	{"DeviceTriggerDigitalInputEventKvs", axevent.DeviceTriggerDigitalInputEventKvs(nil, nil)},
	{"DeviceTriggerRelayEventKvs", axevent.DeviceTriggerRelayEventKvs(nil, nil)},
	{"DeviceRingPowerLimitExceededEventKvs", axevent.DeviceRingPowerLimitExceededEventKvs(nil, nil)},
	{"LightControlLightStatusChangedEventKvs", axevent.LightControlLightStatusChangedEventKvs(nil)},
	{"VideoSourceLiveStreamAccessedEventKvs", axevent.VideoSourceLiveStreamAccessedEventKvs(nil)},
	{"VideoSourceDayNightVisionEventKvs", axevent.VideoSourceDayNightVisionEventKvs(nil, nil)},
	{"VideoSourceTamperingEventKvs", axevent.VideoSourceTamperingEventKvs(nil, nil)},
	{"VideoSourceABREventKvs", axevent.VideoSourceABREventKvs(nil, nil)},
	{"VideoSourceGlobalSceneChangeEventKvs", axevent.VideoSourceGlobalSceneChangeEventKvs(nil, nil)},
	{"VideoSourceMotionAlarmEventKvs", axevent.VideoSourceMotionAlarmEventKvs(nil, nil)},
	{"PTZControllerPTZErrorEventKvs", axevent.PTZControllerPTZErrorEventKvs(nil, nil)},
	{"PTZControllerPTZReadyEventKvs", axevent.PTZControllerPTZReadyEventKvs(nil, nil)},
	{"MediaConfigurationChangedEventKvs", axevent.MediaConfigurationChangedEventKvs(nil, nil)},
	{"MediaProfileChangedEventKvs", axevent.MediaProfileChangedEventKvs(nil)},
	{"CameraApplicationPlatformDevice1Scenario1EventKvs", axevent.CameraApplicationPlatformDevice1Scenario1EventKvs(nil)},
	{"CameraApplicationPlatformDevice1ScenarioANYEventKvs", axevent.CameraApplicationPlatformDevice1ScenarioANYEventKvs(nil)},
	{"CameraApplicationPlatformXInternalDataEventKvs", axevent.CameraApplicationPlatformXInternalDataEventKvs(nil)},
	{"StorageAlertEventKvs", axevent.StorageAlertEventKvs(nil, nil, nil, nil, nil)},
	{"StorageDisruptionEventKvs", axevent.StorageDisruptionEventKvs(nil, nil)},
	{"StorageRecordingEventKvs", axevent.StorageRecordingEventKvs(nil)},
}

func main() {

	app := acapapp.NewAcapApplication()

	for _, event := range events {
		// Subscribe to each event and log details when triggered
		subscription_id, err := app.OnEvent(event.Event, func(e *axevent.Event) {

			switch event.Name {
			case "DeviceIoVirtualInputEventKvs":
				var vi axevent.DeviceIoVirtualInputEvent
				if err := acapapp.UnmarshalEvent(e, &vi); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("VirtualInput Port: %d, Active: %t", vi.Port, vi.Active)

			case "DeviceIoSupervisedPortEventKvs":
				var sp axevent.DeviceIoSupervisedPortEvent
				if err := acapapp.UnmarshalEvent(e, &sp); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("SupervisedPort Port: %d, Tampered: %t, State: %s", sp.Port, sp.Tampered, sp.State)

			case "LightControlLightStatusChangedEventKvs":
				var lightStatus axevent.LightControlLightStatusChangedEvent
				if err := acapapp.UnmarshalEvent(e, &lightStatus); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("LightControl Status: %s", lightStatus.State)

			case "VideoSourceDayNightVisionEventKvs":
				var dayNight axevent.VideoSourceDayNightVisionEvent
				if err := acapapp.UnmarshalEvent(e, &dayNight); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("DayNight Vision Token: %d, Day: %t", dayNight.VideoSourceConfigurationToken, dayNight.Day)

			case "VideoSourceABREventKvs":
				var abr axevent.VideoSourceABREvent
				if err := acapapp.UnmarshalEvent(e, &abr); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("ABR Token: %d, ABR Error: %t", abr.VideoSourceConfigurationToken, abr.AbrError)

			case "VideoSourceGlobalSceneChangeEventKvs":
				var globalScene axevent.VideoSourceGlobalSceneChangeEvent
				if err := acapapp.UnmarshalEvent(e, &globalScene); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Global Scene Source: %d, State: %t", globalScene.Source, globalScene.State)

			case "VideoSourceMotionAlarmEventKvs":
				var motionAlarm axevent.VideoSourceMotionAlarmEvent
				if err := acapapp.UnmarshalEvent(e, &motionAlarm); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Motion Alarm Source: %d, State: %t", motionAlarm.Source, motionAlarm.State)

			case "MediaConfigurationChangedEventKvs":
				var configChanged axevent.MediaConfigurationChangedEvent
				if err := acapapp.UnmarshalEvent(e, &configChanged); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Media Config Changed Type: %s, Token: %s", configChanged.Type, configChanged.Token)

			case "MediaProfileChangedEventKvs":
				var profileChanged axevent.MediaProfileChangedEvent
				if err := acapapp.UnmarshalEvent(e, &profileChanged); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Media Profile Changed Token: %s", profileChanged.Token)

			case "CameraApplicationPlatformDevice1Scenario1EventKvs":
				var scenario1 axevent.CameraApplicationPlatformDevice1Scenario1Event
				if err := acapapp.UnmarshalEvent(e, &scenario1); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Device1 Scenario1 Active: %t", scenario1.Active)

			case "CameraApplicationPlatformDevice1ScenarioANYEventKvs":
				var scenarioAny axevent.CameraApplicationPlatformDevice1ScenarioANYEvent
				if err := acapapp.UnmarshalEvent(e, &scenarioAny); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Device1 ScenarioANY Active: %t", scenarioAny.Active)

			case "CameraApplicationPlatformXInternalDataEventKvs":
				var xInternalData axevent.CameraApplicationPlatformXInternalDataEvent
				if err := acapapp.UnmarshalEvent(e, &xInternalData); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("X Internal Data SVGFrame: %s", xInternalData.SvgFrame)

			case "VideoSourceTamperingEventKvs":
				var tampering axevent.VideoSourceTamperingEvent
				if err := acapapp.UnmarshalEvent(e, &tampering); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Tampering Channel: %d, Tampering: %d", tampering.Channel, tampering.Tampering)

			case "VideoSourceLiveStreamAccessedEventKvs":
				var liveStream axevent.VideoSourceLiveStreamAccessedEvent
				if err := acapapp.UnmarshalEvent(e, &liveStream); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("LiveStream Accessed: %t", liveStream.Accessed)

			case "DeviceIoOutputPortEventKvs":
				var op axevent.DeviceIoOutputPortEvent
				if err := acapapp.UnmarshalEvent(e, &op); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("OutputPort Port: %d, State: %t", op.Port, op.State)

			case "DeviceIoPortEventKvs":
				var port axevent.DeviceIoPortEvent
				if err := acapapp.UnmarshalEvent(e, &port); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Port Port: %d, State: %t", port.Port, port.State)

			case "DeviceSensorPIREventKvs":
				var pir axevent.DeviceSensorPIREvent
				if err := acapapp.UnmarshalEvent(e, &pir); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("PIR Sensor: %d, State: %t", pir.Sensor, pir.State)

			case "DeviceLightStatusEventKvs":
				var light axevent.DeviceLightStatusEvent
				if err := acapapp.UnmarshalEvent(e, &light); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Light ID: %d, State: %s", light.Id, light.State)

			case "DeviceStatusSystemReadyEventKvs":
				var systemReady axevent.DeviceStatusSystemReadyEvent
				if err := acapapp.UnmarshalEvent(e, &systemReady); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("System Ready: %t", systemReady.Ready)

			case "PTZControllerPTZReadyEventKvs":
				var ptzReady axevent.PTZControllerPTZReadyEvent
				if err := acapapp.UnmarshalEvent(e, &ptzReady); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("PTZ Ready Channel: %d, Ready: %t", ptzReady.Channel, ptzReady.Ready)

			case "DeviceStatusTemperatureInsideEventKvs", "DeviceStatusTemperatureAboveEventKvs", "DeviceStatusTemperatureAboveOrBelowEventKvs", "DeviceStatusTemperatureBelowEventKvs":
				var temperature axevent.DeviceStatusTemperatureInsideEvent
				if err := acapapp.UnmarshalEvent(e, &temperature); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Temperature Sensor Level: %t", temperature.SensorLevel)

			case "DeviceHardwareFailurePowerSupplyFailurePTZPowerFailureEventKvs":
				var hardwareFailure axevent.DeviceHardwareFailurePowerSupplyFailurePTZPowerFailureEvent
				if err := acapapp.UnmarshalEvent(e, &hardwareFailure); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("PTZ Power Failure Token: %d, Failed: %t", hardwareFailure.Token, hardwareFailure.Failed)

			case "DeviceTriggerDigitalInputEventKvs":
				var digitalInput axevent.DeviceTriggerDigitalInputEvent
				if err := acapapp.UnmarshalEvent(e, &digitalInput); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Digital Input Token: %d, Logical State: %t", digitalInput.InputToken, digitalInput.LogicalState)

			case "DeviceTriggerRelayEventKvs":
				var relay axevent.DeviceTriggerRelayEvent
				if err := acapapp.UnmarshalEvent(e, &relay); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Relay Token: %d, Logical State: %t", relay.RelayToken, relay.LogicalState)

			case "DeviceRingPowerLimitExceededEventKvs":
				var ring axevent.RingPowerLimitExceededEvent
				if err := acapapp.UnmarshalEvent(e, &ring); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Ring Input: %d, Limit Exceeded: %t", ring.Input, ring.LimitExceeded)

			case "StorageAlertEventKvs":
				var alert axevent.StorageAlertEvent
				if err := acapapp.UnmarshalEvent(e, &alert); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Storage Alert Disk ID: %s, Alert: %t, Overall Health: %d, Temperature: %d, Wear: %d", alert.DiskID, alert.Alert, alert.OverallHealth, alert.Temperature, alert.Wear)

			case "StorageDisruptionEventKvs":
				var disruption axevent.StorageDisruptionEvent
				if err := acapapp.UnmarshalEvent(e, &disruption); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Storage Disruption Disk ID: %s, Disruption: %t", disruption.DiskID, disruption.Disruption)

			case "StorageRecordingEventKvs":
				var recording axevent.StorageRecordingEvent
				if err := acapapp.UnmarshalEvent(e, &recording); err != nil {
					app.Syslog.Errorf("Error unmarshalling %s: %s", event.Name, err)
					return
				}
				app.Syslog.Infof("Storage Recording: %t", recording.Recording)

			// Add more cases for remaining events here

			default:
				app.Syslog.Warnf("No handler defined for event %s", event.Name)
			}
		})

		if err != nil {
			app.Syslog.Critf("Failed to subscribe to event %s: %s", event.Name, err)
		} else {
			app.Syslog.Infof("Subscription created for event %s with subscription ID: %d", event.Name, subscription_id)
		}
	}

	// Run gmain loop with signal handler attached.
	// This will block the main thread until the application is stopped.
	// The application can be stopped by sending a signal to the process (e.g. SIGINT).
	// Axevent needs a running event loop to handle the events callbacks corretly
	app.Run()
}
