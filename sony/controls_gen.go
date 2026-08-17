// Sony Controls.

package sony

// Sony control codes, issued with SendControl (0x9207 SetControlDeviceB).
//
// These are momentary actions rather than stored settings, and are how a Sony
// body is actually driven: there is no standard PTP InitiateCapture involved —
// a shot is an S2 button press and release.
//
// They are 32-bit. Most sit in the 0xD2xx range that looks like a property
// code, but Release and CancelContentsTransfer are above 0xFFFF, which is why
// ControlCode is its own type rather than sharing Prop.
const (
	CtrlS1AndS2Button                                       ControlCode = 0x0000D2E6
	CtrlRelease                                             ControlCode = 0x00010001
	CtrlMovieRecButton                                      ControlCode = 0x0000D2C8
	CtrlMovieRecButtonToggle                                ControlCode = 0x0000F001
	CtrlMovieRecButtonToggle2                               ControlCode = 0x0000D2FE
	CtrlSelectedMediaFormat                                 ControlCode = 0x0000D2E2
	CtrlCancelMediaFormat                                   ControlCode = 0x0000D2E7
	CtrlRECSettingsReset                                    ControlCode = 0x0000D2F3
	CtrlAPSCOrFullSwitching                                 ControlCode = 0x0000D2EB
	CtrlCancelRemoteTouchOperation                          ControlCode = 0x0000D2E5
	CtrlPixelMapping                                        ControlCode = 0x0000D300
	CtrlTimeCodePresetReset                                 ControlCode = 0x0000D302
	CtrlUserBitPresetReset                                  ControlCode = 0x0000D303
	CtrlSensorCleaning                                      ControlCode = 0x0000D304
	CtrlResetPictureProfile                                 ControlCode = 0x0000D305
	CtrlResetCreativeLook                                   ControlCode = 0x0000D306
	CtrlStreamButton                                        ControlCode = 0x0000D307
	CtrlFlickerScan                                         ControlCode = 0x0000D2F1
	CtrlContinuousShootingSpotBoostButton                   ControlCode = 0x0000D2F6
	CtrlTrackingOnAndAFOnButton                             ControlCode = 0x0000D30D
	CtrlForcedFileNumberReset                               ControlCode = 0x0000D30E
	CtrlCameraStandBy                                       ControlCode = 0x0000D315
	CtrlPowerOff                                            ControlCode = 0x0000D301
	CtrlPowerOn                                             ControlCode = 0x0000D316
	CtrlRemoteKeyUp                                         ControlCode = 0x0000D2CD
	CtrlRemoteKeyDown                                       ControlCode = 0x0000D2CE
	CtrlRemoteKeyLeft                                       ControlCode = 0x0000D2CF
	CtrlRemoteKeyRight                                      ControlCode = 0x0000D2D0
	CtrlRemoteKeyCancelBackButton                           ControlCode = 0x0000D2F7
	CtrlRemoteKeyDisplayButton                              ControlCode = 0x0000D2F8
	CtrlRemoteKeySet                                        ControlCode = 0x0000D2F9
	CtrlRemoteKeyRightUp                                    ControlCode = 0x0000D2FA
	CtrlRemoteKeyRightDown                                  ControlCode = 0x0000D2FB
	CtrlRemoteKeyLeftUp                                     ControlCode = 0x0000D2FC
	CtrlRemoteKeyLeftDown                                   ControlCode = 0x0000D2FD
	CtrlRemoteKeyMenuButton                                 ControlCode = 0x0000D2FF
	CtrlResetMultiMatrix                                    ControlCode = 0x0000D2EE
	CtrlCancelFocusPosition                                 ControlCode = 0x0000F002
	CtrlCancelZoomPosition                                  ControlCode = 0x0000F00C
	CtrlCancelContentsTransfer                              ControlCode = 0x00020002
	CtrlS1Button                                            ControlCode = 0x0000D2C1
	CtrlS2Button                                            ControlCode = 0x0000D2C2
	CtrlAELButton                                           ControlCode = 0x0000D2C3
	CtrlFELButton                                           ControlCode = 0x0000D2C9
	CtrlAWBLButton                                          ControlCode = 0x0000D2D9
	CtrlNearFar                                             ControlCode = 0x0000D2D1
	CtrlAFAreaPosition                                      ControlCode = 0x0000D2DC
	CtrlZoomOperation                                       ControlCode = 0x0000D2DD
	CtrlCustomWBCaptureStandby                              ControlCode = 0x0000D2DF
	CtrlCustomWBCaptureStandbyCancel                        ControlCode = 0x0000D2E0
	CtrlCustomWBCapture                                     ControlCode = 0x0000D2E1
	CtrlHighResolutionShutterSpeedAdjust                    ControlCode = 0x0000D2E3
	CtrlHighResolutionShutterSpeedAdjustInIntegralMultiples ControlCode = 0x0000D2F0
	CtrlFocusOperation                                      ControlCode = 0x0000D2EF
	CtrlRemoteTouchOperation                                ControlCode = 0x0000D2E4
	CtrlSaveZoomAndFocusPosition                            ControlCode = 0x0000D2E9
	CtrlLoadZoomAndFocusPosition                            ControlCode = 0x0000D2EA
	CtrlColorTemperatureStep                                ControlCode = 0x0000D2EC
	CtrlWhiteBalanceTintStep                                ControlCode = 0x0000D2ED
	CtrlSetPresetInfoZoomOnlyValue                          ControlCode = 0x0000D2F2
	CtrlCameraButtonFunction                                ControlCode = 0x0000D309
	CtrlCameraButtonFunctionMulti                           ControlCode = 0x0000D30A
	CtrlCameraDialFunction                                  ControlCode = 0x0000D30B
	CtrlCameraLeverFunction                                 ControlCode = 0x0000D30C
	CtrlCreateNewFolder                                     ControlCode = 0x0000D314
	CtrlShutterECSNumberStep                                ControlCode = 0x0000F000
	CtrlZoomOperationWithInt16                              ControlCode = 0x0000F003
	CtrlFocusOperationWithInt16                             ControlCode = 0x0000F004
	CtrlPresetPTZFRecall                                    ControlCode = 0x0000F015
	CtrlUSBConnectionModeRequest                            ControlCode = 0x0000F012
)

var controlNames = map[ControlCode]string{
	CtrlS1AndS2Button:                     "S1AndS2Button",
	CtrlRelease:                           "Release",
	CtrlMovieRecButton:                    "MovieRecButton",
	CtrlMovieRecButtonToggle:              "MovieRecButtonToggle",
	CtrlMovieRecButtonToggle2:             "MovieRecButtonToggle2",
	CtrlSelectedMediaFormat:               "SelectedMediaFormat",
	CtrlCancelMediaFormat:                 "CancelMediaFormat",
	CtrlRECSettingsReset:                  "RECSettingsReset",
	CtrlAPSCOrFullSwitching:               "APS_C_or_Full_Switching",
	CtrlCancelRemoteTouchOperation:        "CancelRemoteTouchOperation",
	CtrlPixelMapping:                      "PixelMapping",
	CtrlTimeCodePresetReset:               "TimeCodePresetReset",
	CtrlUserBitPresetReset:                "UserBitPresetReset",
	CtrlSensorCleaning:                    "SensorCleaning",
	CtrlResetPictureProfile:               "ResetPictureProfile",
	CtrlResetCreativeLook:                 "ResetCreativeLook",
	CtrlStreamButton:                      "StreamButton",
	CtrlFlickerScan:                       "FlickerScan",
	CtrlContinuousShootingSpotBoostButton: "ContinuousShootingSpotBoostButton",
	CtrlTrackingOnAndAFOnButton:           "TrackingOnAndAFOnButton",
	CtrlForcedFileNumberReset:             "ForcedFileNumberReset",
	CtrlCameraStandBy:                     "CameraStandBy",
	CtrlPowerOff:                          "PowerOff",
	CtrlPowerOn:                           "PowerOn",
	CtrlRemoteKeyUp:                       "RemoteKeyUp",
	CtrlRemoteKeyDown:                     "RemoteKeyDown",
	CtrlRemoteKeyLeft:                     "RemoteKeyLeft",
	CtrlRemoteKeyRight:                    "RemoteKeyRight",
	CtrlRemoteKeyCancelBackButton:         "RemoteKeyCancelBackButton",
	CtrlRemoteKeyDisplayButton:            "RemoteKeyDisplayButton",
	CtrlRemoteKeySet:                      "RemoteKeySet",
	CtrlRemoteKeyRightUp:                  "RemoteKeyRightUp",
	CtrlRemoteKeyRightDown:                "RemoteKeyRightDown",
	CtrlRemoteKeyLeftUp:                   "RemoteKeyLeftUp",
	CtrlRemoteKeyLeftDown:                 "RemoteKeyLeftDown",
	CtrlRemoteKeyMenuButton:               "RemoteKeyMenuButton",
	CtrlResetMultiMatrix:                  "ResetMultiMatrix",
	CtrlCancelFocusPosition:               "CancelFocusPosition",
	CtrlCancelZoomPosition:                "CancelZoomPosition",
	CtrlCancelContentsTransfer:            "CancelContentsTransfer",
	CtrlS1Button:                          "S1Button",
	CtrlS2Button:                          "S2Button",
	CtrlAELButton:                         "AELButton",
	CtrlFELButton:                         "FELButton",
	CtrlAWBLButton:                        "AWBLButton",
	CtrlNearFar:                           "NearFar",
	CtrlAFAreaPosition:                    "AFAreaPosition",
	CtrlZoomOperation:                     "ZoomOperation",
	CtrlCustomWBCaptureStandby:            "CustomWBCaptureStandby",
	CtrlCustomWBCaptureStandbyCancel:      "CustomWBCaptureStandbyCancel",
	CtrlCustomWBCapture:                   "CustomWBCapture",
	CtrlHighResolutionShutterSpeedAdjust:  "HighResolutionShutterSpeedAdjust",
	CtrlHighResolutionShutterSpeedAdjustInIntegralMultiples: "HighResolutionShutterSpeedAdjustInIntegralMultiples",
	CtrlFocusOperation:             "FocusOperation",
	CtrlRemoteTouchOperation:       "RemoteTouchOperation",
	CtrlSaveZoomAndFocusPosition:   "SaveZoomAndFocusPosition",
	CtrlLoadZoomAndFocusPosition:   "LoadZoomAndFocusPosition",
	CtrlColorTemperatureStep:       "ColorTemperatureStep",
	CtrlWhiteBalanceTintStep:       "WhiteBalanceTintStep",
	CtrlSetPresetInfoZoomOnlyValue: "SetPresetInfoZoomOnlyValue",
	CtrlCameraButtonFunction:       "CameraButtonFunction",
	CtrlCameraButtonFunctionMulti:  "CameraButtonFunctionMulti",
	CtrlCameraDialFunction:         "CameraDialFunction",
	CtrlCameraLeverFunction:        "CameraLeverFunction",
	CtrlCreateNewFolder:            "CreateNewFolder",
	CtrlShutterECSNumberStep:       "ShutterECSNumberStep",
	CtrlZoomOperationWithInt16:     "ZoomOperationWithInt16",
	CtrlFocusOperationWithInt16:    "FocusOperationWithInt16",
	CtrlPresetPTZFRecall:           "PresetPTZFRecall",
	CtrlUSBConnectionModeRequest:   "USBConnectionModeRequest",
}
