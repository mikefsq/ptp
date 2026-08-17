// Fuji Properties

package fuji

import "github.com/mikefsq/ptp"

// Fujifilm vendor device properties.
const (
	PropFilmSimulation             ptp.Prop = 0xD001
	PropFilmSimulationTune         ptp.Prop = 0xD002
	PropDRangeMode                 ptp.Prop = 0xD007
	PropColorMode                  ptp.Prop = 0xD008
	PropColorSpace                 ptp.Prop = 0xD00A
	PropWhitebalanceTune1          ptp.Prop = 0xD00B
	PropWhitebalanceTune2          ptp.Prop = 0xD00C
	PropColorTemperature           ptp.Prop = 0xD017
	PropQuality                    ptp.Prop = 0xD018
	PropRecMode                    ptp.Prop = 0xD019
	PropLiveViewBrightness         ptp.Prop = 0xD01A
	PropThroughImageZoom           ptp.Prop = 0xD01B
	PropNoiseReduction             ptp.Prop = 0xD01C
	PropMacroMode                  ptp.Prop = 0xD01D
	PropLiveViewStyle              ptp.Prop = 0xD01E
	PropFaceDetectionMode          ptp.Prop = 0xD020
	PropRedEyeCorrectionMode       ptp.Prop = 0xD021
	PropRawCompression             ptp.Prop = 0xD022
	PropGrainEffect                ptp.Prop = 0xD023
	PropSetEyeAFMode               ptp.Prop = 0xD024
	PropFocusPoints                ptp.Prop = 0xD025
	PropMFAssistMode               ptp.Prop = 0xD026
	PropInterlockAEAFArea          ptp.Prop = 0xD027
	PropCommandDialMode            ptp.Prop = 0xD028
	PropShadowing                  ptp.Prop = 0xD029
	PropExposureIndex              ptp.Prop = 0xD02A
	PropMovieISO                   ptp.Prop = 0xD02B
	PropWideDynamicRange           ptp.Prop = 0xD02E
	PropTNumber                    ptp.Prop = 0xD02F
	PropComment                    ptp.Prop = 0xD100
	PropSerialMode                 ptp.Prop = 0xD101
	PropExposureDelay              ptp.Prop = 0xD102
	PropPreviewTime                ptp.Prop = 0xD103
	PropBlackImageTone             ptp.Prop = 0xD104
	PropIllumination               ptp.Prop = 0xD105
	PropFrameGuideMode             ptp.Prop = 0xD106
	PropViewfinderWarning          ptp.Prop = 0xD107
	PropAutoImageRotation          ptp.Prop = 0xD108
	PropDetectImageRotation        ptp.Prop = 0xD109
	PropShutterPriorityMode1       ptp.Prop = 0xD10A
	PropShutterPriorityMode2       ptp.Prop = 0xD10B
	PropAFIlluminator              ptp.Prop = 0xD112
	PropBeep                       ptp.Prop = 0xD113
	PropAELock                     ptp.Prop = 0xD114
	PropISOAutoSetting1            ptp.Prop = 0xD115
	PropISOAutoSetting2            ptp.Prop = 0xD116
	PropISOAutoSetting3            ptp.Prop = 0xD117
	PropExposureStep               ptp.Prop = 0xD118
	PropCompensationStep           ptp.Prop = 0xD119
	PropExposureSimpleSet          ptp.Prop = 0xD11A
	PropCenterPhotometryRange      ptp.Prop = 0xD11B
	PropPhotometryLevel1           ptp.Prop = 0xD11C
	PropPhotometryLevel2           ptp.Prop = 0xD11D
	PropPhotometryLevel3           ptp.Prop = 0xD11E
	PropFlashTuneSpeed             ptp.Prop = 0xD11F
	PropFlashShutterLimit          ptp.Prop = 0xD120
	PropBuiltinFlashMode           ptp.Prop = 0xD121
	PropFlashManualMode            ptp.Prop = 0xD122
	PropFlashRepeatingMode1        ptp.Prop = 0xD123
	PropFlashRepeatingMode2        ptp.Prop = 0xD124
	PropFlashRepeatingMode3        ptp.Prop = 0xD125
	PropFlashCommanderMode1        ptp.Prop = 0xD126
	PropFlashCommanderMode2        ptp.Prop = 0xD127
	PropFlashCommanderMode3        ptp.Prop = 0xD128
	PropFlashCommanderMode4        ptp.Prop = 0xD129
	PropFlashCommanderMode5        ptp.Prop = 0xD12A
	PropFlashCommanderMode6        ptp.Prop = 0xD12B
	PropFlashCommanderMode7        ptp.Prop = 0xD12C
	PropModelingFlash              ptp.Prop = 0xD12D
	PropBKT                        ptp.Prop = 0xD12E
	PropBKTChange                  ptp.Prop = 0xD12F
	PropBKTOrder                   ptp.Prop = 0xD130
	PropBKTSelection               ptp.Prop = 0xD131
	PropAEAFLockButton             ptp.Prop = 0xD132
	PropCenterButton               ptp.Prop = 0xD133
	PropMultiSelectorButton        ptp.Prop = 0xD134
	PropFunctionLock               ptp.Prop = 0xD136
	PropPassword                   ptp.Prop = 0xD145
	PropChangePassword             ptp.Prop = 0xD146
	PropCommandDialSetting1        ptp.Prop = 0xD147
	PropCommandDialSetting2        ptp.Prop = 0xD148
	PropCommandDialSetting3        ptp.Prop = 0xD149
	PropCommandDialSetting4        ptp.Prop = 0xD14A
	PropButtonsAndDials            ptp.Prop = 0xD14B
	PropNonCPULensData             ptp.Prop = 0xD14C
	PropMBD200Batteries            ptp.Prop = 0xD14E
	PropAFOnForMBD200Batteries     ptp.Prop = 0xD14F
	PropFirmwareVersion            ptp.Prop = 0xD153
	PropShotCount                  ptp.Prop = 0xD154
	PropShutterExchangeCount       ptp.Prop = 0xD155
	PropWorldClock                 ptp.Prop = 0xD157
	PropTimeDifference1            ptp.Prop = 0xD158
	PropTimeDifference2            ptp.Prop = 0xD159
	PropLanguage                   ptp.Prop = 0xD15A
	PropFrameNumberSequence        ptp.Prop = 0xD15B
	PropVideoMode                  ptp.Prop = 0xD15C
	PropSetUSBMode                 ptp.Prop = 0xD15D
	PropCommentWriteSetting        ptp.Prop = 0xD161
	PropBCRAppendDelimiter         ptp.Prop = 0xD162
	PropCommentEx                  ptp.Prop = 0xD167
	PropVideoOutOnOff              ptp.Prop = 0xD168
	PropCropMode                   ptp.Prop = 0xD16F
	PropLensZoomPos                ptp.Prop = 0xD170
	PropFocusPosition              ptp.Prop = 0xD171
	PropLiveViewImageQuality       ptp.Prop = 0xD173
	PropLiveViewImageSize          ptp.Prop = 0xD174
	PropLiveViewCondition          ptp.Prop = 0xD175
	PropStandbyMode                ptp.Prop = 0xD176
	PropLiveViewExposure           ptp.Prop = 0xD177
	PropLiveViewWhiteBalance       ptp.Prop = 0xD178
	PropLiveViewWhiteBalanceGain   ptp.Prop = 0xD179
	PropLiveViewTuning             ptp.Prop = 0xD17A
	PropFocusMeteringMode          ptp.Prop = 0xD17C
	PropFocusLength                ptp.Prop = 0xD17D
	PropCropAreaFrameInfo          ptp.Prop = 0xD17E
	PropResetSetting               ptp.Prop = 0xD17F
	PropIOPCode                    ptp.Prop = 0xD184
	PropTetherRawConditionCode     ptp.Prop = 0xD186
	PropTetherRawCompatibilityCode ptp.Prop = 0xD187
	PropLightTune                  ptp.Prop = 0xD200
	PropReleaseMode                ptp.Prop = 0xD201
	PropBKTFrame1                  ptp.Prop = 0xD202
	PropBKTFrame2                  ptp.Prop = 0xD203
	PropBKTStep                    ptp.Prop = 0xD204
	PropProgramShift               ptp.Prop = 0xD205
	PropFocusAreas                 ptp.Prop = 0xD206
	PropPriorityMode               ptp.Prop = 0xD207
	PropAFStatus                   ptp.Prop = 0xD209
	PropDeviceName                 ptp.Prop = 0xD20B
	PropMediaRecord                ptp.Prop = 0xD20C
	PropMediaCapacity              ptp.Prop = 0xD20D
	PropFreeSDRAMImages            ptp.Prop = 0xD20E
	PropMediaStatus                ptp.Prop = 0xD211
	PropCurrentState               ptp.Prop = 0xD212
	PropAELock2                    ptp.Prop = 0xD213
	PropCopyright                  ptp.Prop = 0xD215
	PropCopyright2                 ptp.Prop = 0xD216
	PropAperture                   ptp.Prop = 0xD218
	PropShutterSpeed               ptp.Prop = 0xD219
	PropDeviceError                ptp.Prop = 0xD21B
	PropSensitivityFineTune1       ptp.Prop = 0xD222
	PropSensitivityFineTune2       ptp.Prop = 0xD223
	PropCaptureRemaining           ptp.Prop = 0xD229
	PropMovieRemainingTime         ptp.Prop = 0xD22A
	PropForceMode                  ptp.Prop = 0xD230
	PropShutterSpeed2              ptp.Prop = 0xD240
	PropImageAspectRatio           ptp.Prop = 0xD241
	PropBatteryLevel               ptp.Prop = 0xD242
	PropTotalShotCount             ptp.Prop = 0xD310
	PropHighLightTone              ptp.Prop = 0xD320
	PropShadowTone                 ptp.Prop = 0xD321
	PropLongExposureNR             ptp.Prop = 0xD322
	PropFullTimeManualFocus        ptp.Prop = 0xD323
	PropISODialHn1                 ptp.Prop = 0xD332
	PropISODialHn2                 ptp.Prop = 0xD333
	PropViewMode1                  ptp.Prop = 0xD33F
	PropViewMode2                  ptp.Prop = 0xD340
	PropDispInfoMode               ptp.Prop = 0xD343
	PropLensISSwitch               ptp.Prop = 0xD346
	PropFocusPoint                 ptp.Prop = 0xD347
	PropInstantAFMode              ptp.Prop = 0xD34A
	PropPreAFMode                  ptp.Prop = 0xD34B
	PropCustomSetting              ptp.Prop = 0xD34C
	PropLMOMode                    ptp.Prop = 0xD34D
	PropLockButtonMode             ptp.Prop = 0xD34E
	PropAFLockMode                 ptp.Prop = 0xD34F
	PropMicJackMode                ptp.Prop = 0xD350
	PropISMode                     ptp.Prop = 0xD351
	PropDateTimeDispFormat         ptp.Prop = 0xD352
	PropAeAfLockKeyAssign          ptp.Prop = 0xD353
	PropCrossKeyAssign             ptp.Prop = 0xD354
	PropSilentMode                 ptp.Prop = 0xD355
	PropPBSound                    ptp.Prop = 0xD356
	PropEVFDispAutoRotate          ptp.Prop = 0xD358
	PropExposurePreview            ptp.Prop = 0xD359
	PropDispBrightness1            ptp.Prop = 0xD35A
	PropDispBrightness2            ptp.Prop = 0xD35B
	PropDispChroma1                ptp.Prop = 0xD35C
	PropDispChroma2                ptp.Prop = 0xD35D
	PropFocusCheckMode             ptp.Prop = 0xD35E
	PropFocusScaleUnit             ptp.Prop = 0xD35F
	PropSetFunctionButton          ptp.Prop = 0xD361
	PropSensorCleanTiming          ptp.Prop = 0xD363
	PropCustomAutoPowerOff         ptp.Prop = 0xD364
	PropFileNamePrefix1            ptp.Prop = 0xD365
	PropFileNamePrefix2            ptp.Prop = 0xD366
	PropBatteryInfo1               ptp.Prop = 0xD36A
	PropBatteryInfo2               ptp.Prop = 0xD36B
	PropLensNameAndSerial          ptp.Prop = 0xD36D
	PropCustomDispInfo             ptp.Prop = 0xD36E
	PropFunctionLockCategory1      ptp.Prop = 0xD36F
	PropFunctionLockCategory2      ptp.Prop = 0xD370
	PropCustomPreviewTime          ptp.Prop = 0xD371
	PropFocusArea1                 ptp.Prop = 0xD372
	PropFocusArea2                 ptp.Prop = 0xD373
	PropFocusArea3                 ptp.Prop = 0xD374
	PropFrameGuideGridInfo1        ptp.Prop = 0xD375
	PropFrameGuideGridInfo2        ptp.Prop = 0xD376
	PropFrameGuideGridInfo3        ptp.Prop = 0xD377
	PropFrameGuideGridInfo4        ptp.Prop = 0xD378
	PropLensUnknownData            ptp.Prop = 0xD38A
	PropLensZoomPosCaps            ptp.Prop = 0xD38C
	PropLensFNumberList            ptp.Prop = 0xD38D
	PropLensFocalLengthList        ptp.Prop = 0xD38E
	PropFocusLimiter               ptp.Prop = 0xD390
	PropFocusArea4                 ptp.Prop = 0xD395
	PropInitSequence               ptp.Prop = 0xDF01
	PropAppVersion                 ptp.Prop = 0xDF24
)

// propNames gives a readable name per vendor property code.
var propNames = map[ptp.Prop]string{
	PropFilmSimulation:             "FilmSimulation",
	PropFilmSimulationTune:         "FilmSimulationTune",
	PropDRangeMode:                 "DRangeMode",
	PropColorMode:                  "ColorMode",
	PropColorSpace:                 "ColorSpace",
	PropWhitebalanceTune1:          "WhitebalanceTune1",
	PropWhitebalanceTune2:          "WhitebalanceTune2",
	PropColorTemperature:           "ColorTemperature",
	PropQuality:                    "Quality",
	PropRecMode:                    "RecMode",
	PropLiveViewBrightness:         "LiveViewBrightness",
	PropThroughImageZoom:           "ThroughImageZoom",
	PropNoiseReduction:             "NoiseReduction",
	PropMacroMode:                  "MacroMode",
	PropLiveViewStyle:              "LiveViewStyle",
	PropFaceDetectionMode:          "FaceDetectionMode",
	PropRedEyeCorrectionMode:       "RedEyeCorrectionMode",
	PropRawCompression:             "RawCompression",
	PropGrainEffect:                "GrainEffect",
	PropSetEyeAFMode:               "SetEyeAFMode",
	PropFocusPoints:                "FocusPoints",
	PropMFAssistMode:               "MFAssistMode",
	PropInterlockAEAFArea:          "InterlockAEAFArea",
	PropCommandDialMode:            "CommandDialMode",
	PropShadowing:                  "Shadowing",
	PropExposureIndex:              "ExposureIndex",
	PropMovieISO:                   "MovieISO",
	PropWideDynamicRange:           "WideDynamicRange",
	PropTNumber:                    "TNumber",
	PropComment:                    "Comment",
	PropSerialMode:                 "SerialMode",
	PropExposureDelay:              "ExposureDelay",
	PropPreviewTime:                "PreviewTime",
	PropBlackImageTone:             "BlackImageTone",
	PropIllumination:               "Illumination",
	PropFrameGuideMode:             "FrameGuideMode",
	PropViewfinderWarning:          "ViewfinderWarning",
	PropAutoImageRotation:          "AutoImageRotation",
	PropDetectImageRotation:        "DetectImageRotation",
	PropShutterPriorityMode1:       "ShutterPriorityMode1",
	PropShutterPriorityMode2:       "ShutterPriorityMode2",
	PropAFIlluminator:              "AFIlluminator",
	PropBeep:                       "Beep",
	PropAELock:                     "AELock",
	PropISOAutoSetting1:            "ISOAutoSetting1",
	PropISOAutoSetting2:            "ISOAutoSetting2",
	PropISOAutoSetting3:            "ISOAutoSetting3",
	PropExposureStep:               "ExposureStep",
	PropCompensationStep:           "CompensationStep",
	PropExposureSimpleSet:          "ExposureSimpleSet",
	PropCenterPhotometryRange:      "CenterPhotometryRange",
	PropPhotometryLevel1:           "PhotometryLevel1",
	PropPhotometryLevel2:           "PhotometryLevel2",
	PropPhotometryLevel3:           "PhotometryLevel3",
	PropFlashTuneSpeed:             "FlashTuneSpeed",
	PropFlashShutterLimit:          "FlashShutterLimit",
	PropBuiltinFlashMode:           "BuiltinFlashMode",
	PropFlashManualMode:            "FlashManualMode",
	PropFlashRepeatingMode1:        "FlashRepeatingMode1",
	PropFlashRepeatingMode2:        "FlashRepeatingMode2",
	PropFlashRepeatingMode3:        "FlashRepeatingMode3",
	PropFlashCommanderMode1:        "FlashCommanderMode1",
	PropFlashCommanderMode2:        "FlashCommanderMode2",
	PropFlashCommanderMode3:        "FlashCommanderMode3",
	PropFlashCommanderMode4:        "FlashCommanderMode4",
	PropFlashCommanderMode5:        "FlashCommanderMode5",
	PropFlashCommanderMode6:        "FlashCommanderMode6",
	PropFlashCommanderMode7:        "FlashCommanderMode7",
	PropModelingFlash:              "ModelingFlash",
	PropBKT:                        "BKT",
	PropBKTChange:                  "BKTChange",
	PropBKTOrder:                   "BKTOrder",
	PropBKTSelection:               "BKTSelection",
	PropAEAFLockButton:             "AEAFLockButton",
	PropCenterButton:               "CenterButton",
	PropMultiSelectorButton:        "MultiSelectorButton",
	PropFunctionLock:               "FunctionLock",
	PropPassword:                   "Password",
	PropChangePassword:             "ChangePassword",
	PropCommandDialSetting1:        "CommandDialSetting1",
	PropCommandDialSetting2:        "CommandDialSetting2",
	PropCommandDialSetting3:        "CommandDialSetting3",
	PropCommandDialSetting4:        "CommandDialSetting4",
	PropButtonsAndDials:            "ButtonsAndDials",
	PropNonCPULensData:             "NonCPULensData",
	PropMBD200Batteries:            "MBD200Batteries",
	PropAFOnForMBD200Batteries:     "AFOnForMBD200Batteries",
	PropFirmwareVersion:            "FirmwareVersion",
	PropShotCount:                  "ShotCount",
	PropShutterExchangeCount:       "ShutterExchangeCount",
	PropWorldClock:                 "WorldClock",
	PropTimeDifference1:            "TimeDifference1",
	PropTimeDifference2:            "TimeDifference2",
	PropLanguage:                   "Language",
	PropFrameNumberSequence:        "FrameNumberSequence",
	PropVideoMode:                  "VideoMode",
	PropSetUSBMode:                 "SetUSBMode",
	PropCommentWriteSetting:        "CommentWriteSetting",
	PropBCRAppendDelimiter:         "BCRAppendDelimiter",
	PropCommentEx:                  "CommentEx",
	PropVideoOutOnOff:              "VideoOutOnOff",
	PropCropMode:                   "CropMode",
	PropLensZoomPos:                "LensZoomPos",
	PropFocusPosition:              "FocusPosition",
	PropLiveViewImageQuality:       "LiveViewImageQuality",
	PropLiveViewImageSize:          "LiveViewImageSize",
	PropLiveViewCondition:          "LiveViewCondition",
	PropStandbyMode:                "StandbyMode",
	PropLiveViewExposure:           "LiveViewExposure",
	PropLiveViewWhiteBalance:       "LiveViewWhiteBalance",
	PropLiveViewWhiteBalanceGain:   "LiveViewWhiteBalanceGain",
	PropLiveViewTuning:             "LiveViewTuning",
	PropFocusMeteringMode:          "FocusMeteringMode",
	PropFocusLength:                "FocusLength",
	PropCropAreaFrameInfo:          "CropAreaFrameInfo",
	PropResetSetting:               "ResetSetting",
	PropIOPCode:                    "IOPCode",
	PropTetherRawConditionCode:     "TetherRawConditionCode",
	PropTetherRawCompatibilityCode: "TetherRawCompatibilityCode",
	PropLightTune:                  "LightTune",
	PropReleaseMode:                "ReleaseMode",
	PropBKTFrame1:                  "BKTFrame1",
	PropBKTFrame2:                  "BKTFrame2",
	PropBKTStep:                    "BKTStep",
	PropProgramShift:               "ProgramShift",
	PropFocusAreas:                 "FocusAreas",
	PropPriorityMode:               "PriorityMode",
	PropAFStatus:                   "AFStatus",
	PropDeviceName:                 "DeviceName",
	PropMediaRecord:                "MediaRecord",
	PropMediaCapacity:              "MediaCapacity",
	PropFreeSDRAMImages:            "FreeSDRAMImages",
	PropMediaStatus:                "MediaStatus",
	PropCurrentState:               "CurrentState",
	PropAELock2:                    "AELock2",
	PropCopyright:                  "Copyright",
	PropCopyright2:                 "Copyright2",
	PropAperture:                   "Aperture",
	PropShutterSpeed:               "ShutterSpeed",
	PropDeviceError:                "DeviceError",
	PropSensitivityFineTune1:       "SensitivityFineTune1",
	PropSensitivityFineTune2:       "SensitivityFineTune2",
	PropCaptureRemaining:           "CaptureRemaining",
	PropMovieRemainingTime:         "MovieRemainingTime",
	PropForceMode:                  "ForceMode",
	PropShutterSpeed2:              "ShutterSpeed2",
	PropImageAspectRatio:           "ImageAspectRatio",
	PropBatteryLevel:               "BatteryLevel",
	PropTotalShotCount:             "TotalShotCount",
	PropHighLightTone:              "HighLightTone",
	PropShadowTone:                 "ShadowTone",
	PropLongExposureNR:             "LongExposureNR",
	PropFullTimeManualFocus:        "FullTimeManualFocus",
	PropISODialHn1:                 "ISODialHn1",
	PropISODialHn2:                 "ISODialHn2",
	PropViewMode1:                  "ViewMode1",
	PropViewMode2:                  "ViewMode2",
	PropDispInfoMode:               "DispInfoMode",
	PropLensISSwitch:               "LensISSwitch",
	PropFocusPoint:                 "FocusPoint",
	PropInstantAFMode:              "InstantAFMode",
	PropPreAFMode:                  "PreAFMode",
	PropCustomSetting:              "CustomSetting",
	PropLMOMode:                    "LMOMode",
	PropLockButtonMode:             "LockButtonMode",
	PropAFLockMode:                 "AFLockMode",
	PropMicJackMode:                "MicJackMode",
	PropISMode:                     "ISMode",
	PropDateTimeDispFormat:         "DateTimeDispFormat",
	PropAeAfLockKeyAssign:          "AeAfLockKeyAssign",
	PropCrossKeyAssign:             "CrossKeyAssign",
	PropSilentMode:                 "SilentMode",
	PropPBSound:                    "PBSound",
	PropEVFDispAutoRotate:          "EVFDispAutoRotate",
	PropExposurePreview:            "ExposurePreview",
	PropDispBrightness1:            "DispBrightness1",
	PropDispBrightness2:            "DispBrightness2",
	PropDispChroma1:                "DispChroma1",
	PropDispChroma2:                "DispChroma2",
	PropFocusCheckMode:             "FocusCheckMode",
	PropFocusScaleUnit:             "FocusScaleUnit",
	PropSetFunctionButton:          "SetFunctionButton",
	PropSensorCleanTiming:          "SensorCleanTiming",
	PropCustomAutoPowerOff:         "CustomAutoPowerOff",
	PropFileNamePrefix1:            "FileNamePrefix1",
	PropFileNamePrefix2:            "FileNamePrefix2",
	PropBatteryInfo1:               "BatteryInfo1",
	PropBatteryInfo2:               "BatteryInfo2",
	PropLensNameAndSerial:          "LensNameAndSerial",
	PropCustomDispInfo:             "CustomDispInfo",
	PropFunctionLockCategory1:      "FunctionLockCategory1",
	PropFunctionLockCategory2:      "FunctionLockCategory2",
	PropCustomPreviewTime:          "CustomPreviewTime",
	PropFocusArea1:                 "FocusArea1",
	PropFocusArea2:                 "FocusArea2",
	PropFocusArea3:                 "FocusArea3",
	PropFrameGuideGridInfo1:        "FrameGuideGridInfo1",
	PropFrameGuideGridInfo2:        "FrameGuideGridInfo2",
	PropFrameGuideGridInfo3:        "FrameGuideGridInfo3",
	PropFrameGuideGridInfo4:        "FrameGuideGridInfo4",
	PropLensUnknownData:            "LensUnknownData",
	PropLensZoomPosCaps:            "LensZoomPosCaps",
	PropLensFNumberList:            "LensFNumberList",
	PropLensFocalLengthList:        "LensFocalLengthList",
	PropFocusLimiter:               "FocusLimiter",
	PropFocusArea4:                 "FocusArea4",
	PropInitSequence:               "InitSequence",
	PropAppVersion:                 "AppVersion",
}

// xt5Props lists the properties the X-T5's SDK plugin names.
var xt5Props = map[ptp.Prop]bool{
	PropFilmSimulationTune:       true,
	PropColorMode:                true,
	PropColorSpace:               true,
	PropLiveViewBrightness:       true,
	PropThroughImageZoom:         true,
	PropNoiseReduction:           true,
	PropMacroMode:                true,
	PropLiveViewStyle:            true,
	PropFaceDetectionMode:        true,
	PropGrainEffect:              true,
	PropFocusPoints:              true,
	PropMFAssistMode:             true,
	PropInterlockAEAFArea:        true,
	PropShadowing:                true,
	PropWideDynamicRange:         true,
	PropTNumber:                  true,
	PropComment:                  true,
	PropSerialMode:               true,
	PropExposureDelay:            true,
	PropPreviewTime:              true,
	PropBlackImageTone:           true,
	PropIllumination:             true,
	PropFrameGuideMode:           true,
	PropViewfinderWarning:        true,
	PropAutoImageRotation:        true,
	PropAFIlluminator:            true,
	PropBeep:                     true,
	PropAELock:                   true,
	PropExposureStep:             true,
	PropCompensationStep:         true,
	PropExposureSimpleSet:        true,
	PropCenterPhotometryRange:    true,
	PropFlashTuneSpeed:           true,
	PropFlashShutterLimit:        true,
	PropBuiltinFlashMode:         true,
	PropFlashManualMode:          true,
	PropModelingFlash:            true,
	PropBKT:                      true,
	PropBKTChange:                true,
	PropBKTOrder:                 true,
	PropBKTSelection:             true,
	PropAEAFLockButton:           true,
	PropCenterButton:             true,
	PropMultiSelectorButton:      true,
	PropFunctionLock:             true,
	PropPassword:                 true,
	PropButtonsAndDials:          true,
	PropNonCPULensData:           true,
	PropMBD200Batteries:          true,
	PropFirmwareVersion:          true,
	PropWorldClock:               true,
	PropLanguage:                 true,
	PropFrameNumberSequence:      true,
	PropVideoMode:                true,
	PropCommentWriteSetting:      true,
	PropBCRAppendDelimiter:       true,
	PropCommentEx:                true,
	PropVideoOutOnOff:            true,
	PropCropMode:                 true,
	PropLiveViewImageQuality:     true,
	PropLiveViewImageSize:        true,
	PropLiveViewCondition:        true,
	PropStandbyMode:              true,
	PropLiveViewExposure:         true,
	PropLiveViewWhiteBalanceGain: true,
	PropLiveViewTuning:           true,
	PropFocusLength:              true,
	PropCropAreaFrameInfo:        true,
	PropResetSetting:             true,
	PropLightTune:                true,
	PropBKTStep:                  true,
	PropProgramShift:             true,
	PropPriorityMode:             true,
	PropAFStatus:                 true,
	PropMediaRecord:              true,
	PropMediaCapacity:            true,
	PropMediaStatus:              true,
	PropCopyright:                true,
	PropAperture:                 true,
	PropShutterSpeed:             true,
	PropMovieRemainingTime:       true,
	PropTotalShotCount:           true,
	PropHighLightTone:            true,
	PropShadowTone:               true,
	PropLongExposureNR:           true,
	PropFullTimeManualFocus:      true,
	PropDispInfoMode:             true,
	PropLensISSwitch:             true,
	PropInstantAFMode:            true,
	PropPreAFMode:                true,
	PropCustomSetting:            true,
	PropLMOMode:                  true,
	PropLockButtonMode:           true,
	PropAFLockMode:               true,
	PropMicJackMode:              true,
	PropISMode:                   true,
	PropDateTimeDispFormat:       true,
	PropAeAfLockKeyAssign:        true,
	PropCrossKeyAssign:           true,
	PropSilentMode:               true,
	PropPBSound:                  true,
	PropEVFDispAutoRotate:        true,
	PropExposurePreview:          true,
	PropFocusCheckMode:           true,
	PropFocusScaleUnit:           true,
	PropSensorCleanTiming:        true,
	PropCustomAutoPowerOff:       true,
	PropCustomDispInfo:           true,
	PropCustomPreviewTime:        true,
	PropLensZoomPosCaps:          true,
	PropLensFNumberList:          true,
	PropFocusLimiter:             true,
}

// SupportedOnXT5 reports whether the X-T5 plugin names this property. It is a
// function rather than a method because ptp.Prop is not this package's type.
func SupportedOnXT5(p ptp.Prop) bool { return xt5Props[p] }

// Fujifilm vendor operations.
const (
	OpFujiInitiateMovieCapture    ptp.OpCode = 0x9020
	OpFujiTerminateMovieCapture   ptp.OpCode = 0x9021
	OpFujiGetCapturePreview       ptp.OpCode = 0x9022
	OpFujiSetFocusPoint           ptp.OpCode = 0x9026
	OpFujiResetFocusPoint         ptp.OpCode = 0x9027
	OpFujiGetDeviceInfo           ptp.OpCode = 0x902B
	OpFujiSetShutterSpeed         ptp.OpCode = 0x902C
	OpFujiSetAperture             ptp.OpCode = 0x902D
	OpFujiSetExposureCompensation ptp.OpCode = 0x902E
	OpFujiCancelInitiateCapture   ptp.OpCode = 0x9030
	OpFujiFmSendObjectInfo        ptp.OpCode = 0x9040
	OpFujiFmSendObject            ptp.OpCode = 0x9041
	OpFujiFmSendPartialObject     ptp.OpCode = 0x9042
)
