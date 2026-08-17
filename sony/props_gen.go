// Sony Properties

package sony

// SDKProp is the Camera Remote SDK's internal property code.
type SDKProp uint32

// PropInfo is what the SDK's static table knows about a property before any
// camera is attached.
type PropInfo struct {
	SDK      SDKProp // Camera Remote SDK internal code
	Name     string  // CrDeviceProperty_* name, minus the prefix
	Min, Max uint64  // valid value range; both zero when the SDK records none
	A7RV     bool    // supported on ILCE-7RM5
	A7RVI    bool    // supported on ILCE-7RM6
}

// Wire property codes.
const (
	PropFNumber                                             Prop = 0x5007 // FNumber
	PropExposureBiasCompensation                            Prop = 0x5010 // ExposureBiasCompensation
	PropFlashCompensation                                   Prop = 0xD200 // FlashCompensation
	PropShutterSpeed                                        Prop = 0xD20D // ShutterSpeed
	PropIsoSensitivity                                      Prop = 0xD21E // IsoSensitivity
	PropExposureProgramMode                                 Prop = 0x500E // ExposureProgramMode
	PropFileType                                            Prop = 0xD253 // FileType
	PropMediaSLOT1FileType                                  Prop = 0xD28B // MediaSLOT1_FileType
	PropMediaSLOT2FileType                                  Prop = 0xD28C // MediaSLOT2_FileType
	PropStillImageQuality                                   Prop = 0xD252 // StillImageQuality
	PropMediaSLOT1ImageQuality                              Prop = 0xD28D // MediaSLOT1_ImageQuality
	PropMediaSLOT2ImageQuality                              Prop = 0xD28E // MediaSLOT2_ImageQuality
	PropWhiteBalance                                        Prop = 0x5005 // WhiteBalance
	PropFocusMode                                           Prop = 0x500A // FocusMode
	PropMeteringMode                                        Prop = 0x500B // MeteringMode
	PropFlashMode                                           Prop = 0x500C // FlashMode
	PropWirelessFlash                                       Prop = 0xD262 // WirelessFlash
	PropRedEyeReduction                                     Prop = 0xD263 // RedEyeReduction
	PropDriveMode                                           Prop = 0x5013 // DriveMode
	PropDRO                                                 Prop = 0xD201 // DRO
	PropImageSize                                           Prop = 0xD203 // ImageSize
	PropMediaSLOT1ImageSize                                 Prop = 0xD28F // MediaSLOT1_ImageSize
	PropMediaSLOT2ImageSize                                 Prop = 0xD290 // MediaSLOT2_ImageSize
	PropAspectRatio                                         Prop = 0xD211 // AspectRatio
	PropPictureEffect                                       Prop = 0xD21B // PictureEffect
	PropFocusArea                                           Prop = 0xD22C // FocusArea
	PropReserved4                                           Prop = 0xD219 // reserved4
	PropColortemp                                           Prop = 0xD20F // Colortemp
	PropColorTuningAB                                       Prop = 0xD21C // ColorTuningAB
	PropColorTuningGM                                       Prop = 0xD210 // ColorTuningGM
	PropLiveViewDisplayEffect                               Prop = 0xD231 // LiveViewDisplayEffect
	PropStillImageStoreDestination                          Prop = 0xD222 // StillImageStoreDestination
	PropPriorityKeySettings                                 Prop = 0xD25A // PriorityKeySettings
	PropAFTrackingSensitivity                               Prop = 0xD255 // AFTrackingSensitivity
	PropFocusMagnifierSetting                               Prop = 0xD254 // Focus_Magnifier_Setting
	PropDateTimeSettings                                    Prop = 0xD223 // DateTime_Settings
	PropZoomScale                                           Prop = 0xD25C // Zoom_Scale
	PropZoomSetting                                         Prop = 0xD25F // Zoom_Setting
	PropMovieFileFormat                                     Prop = 0xD241 // Movie_File_Format
	PropMovieRecordingSetting                               Prop = 0xD242 // Movie_Recording_Setting
	PropMovieRecordingFrameRateSetting                      Prop = 0xD286 // Movie_Recording_FrameRateSetting
	PropCompressionFileFormatStill                          Prop = 0xD287 // CompressionFileFormatStill
	PropRAWFileCompressionType                              Prop = 0xD288 // RAW_FileCompressionType
	PropMediaSLOT1RAWFileCompressionType                    Prop = 0xD289 // MediaSLOT1_RAW_FileCompressionType
	PropMediaSLOT2RAWFileCompressionType                    Prop = 0xD28A // MediaSLOT2_RAW_FileCompressionType
	PropIrisModeSetting                                     Prop = 0xD001 // IrisModeSetting
	PropShutterModeSetting                                  Prop = 0xD013 // ShutterModeSetting
	PropGainControlSetting                                  Prop = 0xD01C // GainControlSetting
	PropGainBaseIsoSensitivity                              Prop = 0xD020 // GainBaseIsoSensitivity
	PropGainBaseSensitivity                                 Prop = 0xD021 // GainBaseSensitivity
	PropExposureIndex                                       Prop = 0xD022 // ExposureIndex
	PropBaseLookValue                                       Prop = 0xD03C // BaseLookValue
	PropPlaybackMedia                                       Prop = 0xD042 // PlaybackMedia
	PropDispModeSetting                                     Prop = 0xD045 // DispModeSetting
	PropDispMode                                            Prop = 0xD046 // DispMode
	PropTouchOperation                                      Prop = 0xD047 // TouchOperation
	PropSelectFinderMonitor                                 Prop = 0xD048 // SelectFinderMonitor
	PropAutoPowerOffTemperature                             Prop = 0xD049 // AutoPowerOffTemperature
	PropBodyKeyLock                                         Prop = 0xD04A // BodyKeyLock
	PropImageIDNumSetting                                   Prop = 0xD092 // ImageID_Num_Setting
	PropImageIDNum                                          Prop = 0xD04B // ImageID_Num
	PropImageIDString                                       Prop = 0xD04C // ImageID_String
	PropExposureCtrlType                                    Prop = 0xD099 // ExposureCtrlType
	PropMonitorLUTSetting                                   Prop = 0xD04D // MonitorLUTSetting
	PropFocalDistanceInMeter                                Prop = 0xD004 // FocalDistanceInMeter
	PropFocalDistanceInFeet                                 Prop = 0xD005 // FocalDistanceInFeet
	PropFocalDistanceUnitSetting                            Prop = 0xD006 // FocalDistanceUnitSetting
	PropDigitalZoomScale                                    Prop = 0xD00A // DigitalZoomScale
	PropZoomDistance                                        Prop = 0xD00B // ZoomDistance
	PropZoomDistanceUnitSetting                             Prop = 0xD029 // ZoomDistanceUnitSetting
	PropShutterModeStatus                                   Prop = 0xD011 // ShutterModeStatus
	PropShutterSlow                                         Prop = 0xD014 // ShutterSlow
	PropShutterSlowFrames                                   Prop = 0xD015 // ShutterSlowFrames
	PropMovieRecordingResolutionForMain                     Prop = 0xD024 // Movie_Recording_ResolutionForMain
	PropMovieRecordingResolutionForProxy                    Prop = 0xD025 // Movie_Recording_ResolutionForProxy
	PropMovieRecordingFrameRateProxySetting                 Prop = 0xD028 // Movie_Recording_FrameRateProxySetting
	PropBatteryRemainDisplayUnit                            Prop = 0xD037 // BatteryRemainDisplayUnit
	PropPowerSource                                         Prop = 0xD03A // PowerSource
	PropMovieShootingMode                                   Prop = 0xE000 // MovieShootingMode
	PropMovieShootingModeColorGamut                         Prop = 0xE001 // MovieShootingModeColorGamut
	PropMovieShootingModeTargetDisplay                      Prop = 0xE002 // MovieShootingModeTargetDisplay
	PropDepthOfFieldAdjustmentMode                          Prop = 0xE009 // DepthOfFieldAdjustmentMode
	PropWhiteBalanceModeSetting                             Prop = 0xD00C // WhiteBalanceModeSetting
	PropWhiteBalanceTint                                    Prop = 0xD00D // WhiteBalanceTint
	PropShutterECSSetting                                   Prop = 0xE006 // ShutterECSSetting
	PropShutterECSNumber                                    Prop = 0xE007 // ShutterECSNumber
	PropShutterECSFrequency                                 Prop = 0xE008 // ShutterECSFrequency
	PropRecorderControlProxySetting                         Prop = 0xE00D // RecorderControlProxySetting
	PropButtonAssignmentAssignable1                         Prop = 0xE014 // ButtonAssignmentAssignable1
	PropButtonAssignmentAssignable2                         Prop = 0xE015 // ButtonAssignmentAssignable2
	PropButtonAssignmentAssignable3                         Prop = 0xE016 // ButtonAssignmentAssignable3
	PropButtonAssignmentAssignable4                         Prop = 0xE017 // ButtonAssignmentAssignable4
	PropButtonAssignmentAssignable5                         Prop = 0xE018 // ButtonAssignmentAssignable5
	PropButtonAssignmentAssignable6                         Prop = 0xE019 // ButtonAssignmentAssignable6
	PropButtonAssignmentAssignable7                         Prop = 0xE01A // ButtonAssignmentAssignable7
	PropButtonAssignmentAssignable8                         Prop = 0xE01B // ButtonAssignmentAssignable8
	PropButtonAssignmentAssignable9                         Prop = 0xE01C // ButtonAssignmentAssignable9
	PropButtonAssignmentAssignable10                        Prop = 0xE01D // ButtonAssignmentAssignable10
	PropButtonAssignmentAssignable11                        Prop = 0xE08D // ButtonAssignmentAssignable11
	PropButtonAssignmentLensAssignable1                     Prop = 0xE01E // ButtonAssignmentLensAssignable1
	PropAssignableButton1                                   Prop = 0xE029 // AssignableButton1
	PropAssignableButton2                                   Prop = 0xE02A // AssignableButton2
	PropAssignableButton3                                   Prop = 0xE02B // AssignableButton3
	PropAssignableButton4                                   Prop = 0xE02C // AssignableButton4
	PropAssignableButton5                                   Prop = 0xE02D // AssignableButton5
	PropAssignableButton6                                   Prop = 0xE02E // AssignableButton6
	PropAssignableButton7                                   Prop = 0xE02F // AssignableButton7
	PropAssignableButton8                                   Prop = 0xE030 // AssignableButton8
	PropAssignableButton9                                   Prop = 0xE031 // AssignableButton9
	PropAssignableButton10                                  Prop = 0xE032 // AssignableButton10
	PropAssignableButton11                                  Prop = 0xE08E // AssignableButton11
	PropLensAssignableButton1                               Prop = 0xE033 // LensAssignableButton1
	PropFocusModeSetting                                    Prop = 0xD007 // FocusModeSetting
	PropShutterAngle                                        Prop = 0xD00E // ShutterAngle
	PropShutterSetting                                      Prop = 0xD00F // ShutterSetting
	PropShutterMode                                         Prop = 0xD010 // ShutterMode
	PropShutterSpeedValue                                   Prop = 0xD016 // ShutterSpeedValue
	PropNDFilter                                            Prop = 0xD018 // NDFilter
	PropNDFilterModeSetting                                 Prop = 0xD01A // NDFilterModeSetting
	PropNDFilterValue                                       Prop = 0xD01B // NDFilterValue
	PropGainUnitSetting                                     Prop = 0xD01D // GainUnitSetting
	PropGaindBValue                                         Prop = 0xD01E // GaindBValue
	PropAWB                                                 Prop = 0xD03B // AWB
	PropSceneFileIndex                                      Prop = 0xE01F // SceneFileIndex
	PropMoviePlayButton                                     Prop = 0xE021 // MoviePlayButton
	PropMoviePlayPauseButton                                Prop = 0xE022 // MoviePlayPauseButton
	PropMoviePlayStopButton                                 Prop = 0xE023 // MoviePlayStopButton
	PropMovieForwardButton                                  Prop = 0xE024 // MovieForwardButton
	PropMovieRewindButton                                   Prop = 0xE025 // MovieRewindButton
	PropMovieNextButton                                     Prop = 0xE026 // MovieNextButton
	PropMoviePrevButton                                     Prop = 0xE027 // MoviePrevButton
	PropMovieRecReviewButton                                Prop = 0xE028 // MovieRecReviewButton
	PropSubjectRecognitionAF                                Prop = 0xD060 // SubjectRecognitionAF
	PropAFTransitionSpeed                                   Prop = 0xD061 // AFTransitionSpeed
	PropAFSubjShiftSens                                     Prop = 0xD062 // AFSubjShiftSens
	PropAFAssist                                            Prop = 0xE084 // AFAssist
	PropNDFilterSwitchingSetting                            Prop = 0xD073 // NDFilterSwitchingSetting
	PropFunctionOfRemoteTouchOperation                      Prop = 0xE083 // FunctionOfRemoteTouchOperation
	PropFollowFocusPositionSetting                          Prop = 0xE088 // FollowFocusPositionSetting
	PropFocusBracketShotNumber                              Prop = 0xD2A1 // FocusBracketShotNumber
	PropFocusBracketFocusRange                              Prop = 0xD2A2 // FocusBracketFocusRange
	PropExtendedInterfaceMode                               Prop = 0xD0C5 // ExtendedInterfaceMode
	PropSQRecordingFrameRateSetting                         Prop = 0xD0D0 // SQRecordingFrameRateSetting
	PropSQFrameRate                                         Prop = 0xD052 // SQFrameRate
	PropSQRecordingSetting                                  Prop = 0xD0D1 // SQRecordingSetting
	PropAudioRecording                                      Prop = 0xD0D2 // AudioRecording
	PropAudioInputMasterLevel                               Prop = 0xE050 // AudioInputMasterLevel
	PropTimeCodePreset                                      Prop = 0xD0D3 // TimeCodePreset
	PropTimeCodeFormat                                      Prop = 0xD0D5 // TimeCodeFormat
	PropTimeCodeRun                                         Prop = 0xD0D6 // TimeCodeRun
	PropTimeCodeMake                                        Prop = 0xD0D7 // TimeCodeMake
	PropUserBitPreset                                       Prop = 0xD0D4 // UserBitPreset
	PropUserBitTimeRec                                      Prop = 0xD0D8 // UserBitTimeRec
	PropImageStabilizationSteadyShot                        Prop = 0xD0D9 // ImageStabilizationSteadyShot
	PropMovieImageStabilizationSteadyShot                   Prop = 0xD0DA // Movie_ImageStabilizationSteadyShot
	PropSilentMode                                          Prop = 0xD0DB // SilentMode
	PropSilentModeApertureDriveInAF                         Prop = 0xD0DC // SilentModeApertureDriveInAF
	PropSilentModeShutterWhenPowerOff                       Prop = 0xD0DD // SilentModeShutterWhenPowerOff
	PropSilentModeAutoPixelMapping                          Prop = 0xD0DE // SilentModeAutoPixelMapping
	PropShutterType                                         Prop = 0xD0DF // ShutterType
	PropPictureProfile                                      Prop = 0xD23F // PictureProfile
	PropPictureProfileBlackLevel                            Prop = 0xD0E0 // PictureProfile_BlackLevel
	PropPictureProfileGamma                                 Prop = 0xD0E1 // PictureProfile_Gamma
	PropPictureProfileBlackGammaRange                       Prop = 0xD0E2 // PictureProfile_BlackGammaRange
	PropPictureProfileBlackGammaLevel                       Prop = 0xD0E3 // PictureProfile_BlackGammaLevel
	PropPictureProfileKneeMode                              Prop = 0xD0E4 // PictureProfile_KneeMode
	PropPictureProfileKneeAutoSetMaxPoint                   Prop = 0xD0E5 // PictureProfile_KneeAutoSet_MaxPoint
	PropPictureProfileKneeAutoSetSensitivity                Prop = 0xD0E6 // PictureProfile_KneeAutoSet_Sensitivity
	PropPictureProfileKneeManualSetPoint                    Prop = 0xD0E7 // PictureProfile_KneeManualSet_Point
	PropPictureProfileKneeManualSetSlope                    Prop = 0xD0E8 // PictureProfile_KneeManualSet_Slope
	PropPictureProfileColorMode                             Prop = 0xD0E9 // PictureProfile_ColorMode
	PropPictureProfileSaturation                            Prop = 0xD0EA // PictureProfile_Saturation
	PropPictureProfileColorPhase                            Prop = 0xD0EB // PictureProfile_ColorPhase
	PropPictureProfileColorDepthRed                         Prop = 0xD0EC // PictureProfile_ColorDepthRed
	PropPictureProfileColorDepthGreen                       Prop = 0xD0ED // PictureProfile_ColorDepthGreen
	PropPictureProfileColorDepthBlue                        Prop = 0xD0EE // PictureProfile_ColorDepthBlue
	PropPictureProfileColorDepthCyan                        Prop = 0xD0EF // PictureProfile_ColorDepthCyan
	PropPictureProfileColorDepthMagenta                     Prop = 0xD0F0 // PictureProfile_ColorDepthMagenta
	PropPictureProfileColorDepthYellow                      Prop = 0xD0F1 // PictureProfile_ColorDepthYellow
	PropPictureProfileDetailLevel                           Prop = 0xD0F2 // PictureProfile_DetailLevel
	PropPictureProfileDetailAdjustMode                      Prop = 0xD0F3 // PictureProfile_DetailAdjustMode
	PropPictureProfileDetailAdjustVHBalance                 Prop = 0xD0F4 // PictureProfile_DetailAdjustVHBalance
	PropPictureProfileDetailAdjustBWBalance                 Prop = 0xD0F5 // PictureProfile_DetailAdjustBWBalance
	PropPictureProfileDetailAdjustLimit                     Prop = 0xD0F6 // PictureProfile_DetailAdjustLimit
	PropPictureProfileDetailAdjustCrispening                Prop = 0xD0F7 // PictureProfile_DetailAdjustCrispening
	PropPictureProfileDetailAdjustHiLightDetail             Prop = 0xD0F8 // PictureProfile_DetailAdjustHiLightDetail
	PropPictureProfileCopy                                  Prop = 0xD0F9 // PictureProfile_Copy
	PropCreativeLook                                        Prop = 0xD0FA // CreativeLook
	PropCreativeLookContrast                                Prop = 0xD0FB // CreativeLook_Contrast
	PropCreativeLookHighlights                              Prop = 0xD0FC // CreativeLook_Highlights
	PropCreativeLookShadows                                 Prop = 0xD0FD // CreativeLook_Shadows
	PropCreativeLookFade                                    Prop = 0xD0FE // CreativeLook_Fade
	PropCreativeLookSaturation                              Prop = 0xD0FF // CreativeLook_Saturation
	PropCreativeLookSharpness                               Prop = 0xD100 // CreativeLook_Sharpness
	PropCreativeLookSharpnessRange                          Prop = 0xD101 // CreativeLook_SharpnessRange
	PropCreativeLookClarity                                 Prop = 0xD102 // CreativeLook_Clarity
	PropCreativeLookCustomLook                              Prop = 0xD103 // CreativeLook_CustomLook
	PropMovieProxyFileFormat                                Prop = 0xD027 // Movie_ProxyFileFormat
	PropProxyRecordingSetting                               Prop = 0xD109 // ProxyRecordingSetting
	PropFunctionOfTouchOperation                            Prop = 0xD283 // FunctionOfTouchOperation
	PropHighResolutionShutterSpeedSetting                   Prop = 0xD281 // HighResolutionShutterSpeedSetting
	PropDeleteUserBaseLook                                  Prop = 0xD0C7 // DeleteUserBaseLook
	PropSelectUserBaseLookToEdit                            Prop = 0xD0C8 // SelectUserBaseLookToEdit
	PropSelectUserBaseLookToSetInPPLUT                      Prop = 0xD0CC // SelectUserBaseLookToSetInPPLUT
	PropUserBaseLookInput                                   Prop = 0xD0C9 // UserBaseLookInput
	PropUserBaseLookAELevelOffset                           Prop = 0xD0CA // UserBaseLookAELevelOffset
	PropBaseISOSwitchEI                                     Prop = 0xD0CB // BaseISOSwitchEI
	PropFlickerLessShooting                                 Prop = 0xD133 // FlickerLessShooting
	PropAudioLevelDisplay                                   Prop = 0xD172 // AudioLevelDisplay
	PropPlaybackVolumeSettings                              Prop = 0xD17C // PlaybackVolumeSettings
	PropAutoReview                                          Prop = 0xD17D // AutoReview
	PropAudioSignals                                        Prop = 0xD17E // AudioSignals
	PropHDMIResolutionStillPlay                             Prop = 0xD17F // HDMIResolutionStillPlay
	PropMovieHDMIOutputRecMedia                             Prop = 0xD180 // Movie_HDMIOutputRecMedia
	PropMovieHDMIOutputResolution                           Prop = 0xD181 // Movie_HDMIOutputResolution
	PropMovieHDMIOutput4KSetting                            Prop = 0xD182 // Movie_HDMIOutput4KSetting
	PropMovieHDMIOutputRAW                                  Prop = 0xD183 // Movie_HDMIOutputRAW
	PropMovieHDMIOutputRawSetting                           Prop = 0xD184 // Movie_HDMIOutputRawSetting
	PropMovieHDMIOutputColorGamutForRAWOut                  Prop = 0xD185 // Movie_HDMIOutputColorGamutForRAWOut
	PropMovieHDMIOutputTimeCode                             Prop = 0xD186 // Movie_HDMIOutputTimeCode
	PropMovieHDMIOutputRecControl                           Prop = 0xD187 // Movie_HDMIOutputRecControl
	PropMonitoringOutputDisplayHDMI                         Prop = 0xD079 // MonitoringOutputDisplayHDMI
	PropMovieHDMIOutputAudioCH                              Prop = 0xE059 // Movie_HDMIOutputAudioCH
	PropMovieIntervalRecIntervalTime                        Prop = 0xD055 // Movie_IntervalRec_IntervalTime
	PropMovieIntervalRecFrameRateSetting                    Prop = 0xD151 // Movie_IntervalRec_FrameRateSetting
	PropMovieIntervalRecRecordingSetting                    Prop = 0xD152 // Movie_IntervalRec_RecordingSetting
	PropEframingScaleAuto                                   Prop = 0xD0CD // EframingScaleAuto
	PropEframingSpeedAuto                                   Prop = 0xD0CE // EframingSpeedAuto
	PropEframingModeAuto                                    Prop = 0xD124 // EframingModeAuto
	PropEframingRecordingImageCrop                          Prop = 0xD153 // EframingRecordingImageCrop
	PropEframingHDMICrop                                    Prop = 0xD154 // EframingHDMICrop
	PropCameraEframing                                      Prop = 0xD0CF // CameraEframing
	PropUSBPowerSupply                                      Prop = 0xD150 // USBPowerSupply
	PropLongExposureNR                                      Prop = 0xD15B // LongExposureNR
	PropHighIsoNR                                           Prop = 0xD15C // HighIsoNR
	PropHLGStillImage                                       Prop = 0xD15D // HLGStillImage
	PropColorSpace                                          Prop = 0xD15E // ColorSpace
	PropBracketOrder                                        Prop = 0xD166 // BracketOrder
	PropFocusBracketOrder                                   Prop = 0xD167 // FocusBracketOrder
	PropFocusBracketExposureLock1stImg                      Prop = 0xD168 // FocusBracketExposureLock1stImg
	PropFocusBracketIntervalUntilNextShot                   Prop = 0xD169 // FocusBracketIntervalUntilNextShot
	PropIntervalRecShootingStartTime                        Prop = 0xD16A // IntervalRec_ShootingStartTime
	PropIntervalRecShootingInterval                         Prop = 0xD16B // IntervalRec_ShootingInterval
	PropIntervalRecShootIntervalPriority                    Prop = 0xD16F // IntervalRec_ShootIntervalPriority
	PropIntervalRecNumberOfShots                            Prop = 0xD16C // IntervalRec_NumberOfShots
	PropIntervalRecAETrackingSensitivity                    Prop = 0xD16D // IntervalRec_AETrackingSensitivity
	PropIntervalRecShutterType                              Prop = 0xD16E // IntervalRec_ShutterType
	PropElectricFrontCurtainShutter                         Prop = 0xD170 // ElectricFrontCurtainShutter
	PropWindNoiseReduct                                     Prop = 0xD171 // WindNoiseReduct
	PropRecordingSelfTimer                                  Prop = 0xD29C // RecordingSelfTimer
	PropRecordingSelfTimerCountTime                         Prop = 0xD29D // RecordingSelfTimerCountTime
	PropRecordingSelfTimerContinuous                        Prop = 0xD29F // RecordingSelfTimerContinuous
	PropRecordingSelfTimerStatus                            Prop = 0xD2A0 // RecordingSelfTimerStatus
	PropBulbTimerSetting                                    Prop = 0xD2A4 // BulbTimerSetting
	PropBulbExposureTimeSetting                             Prop = 0xD2A5 // BulbExposureTimeSetting
	PropAutoSlowShutter                                     Prop = 0xD173 // AutoSlowShutter
	PropIsoAutoMinShutterSpeedMode                          Prop = 0xD14D // IsoAutoMinShutterSpeedMode
	PropIsoAutoMinShutterSpeedManual                        Prop = 0xD176 // IsoAutoMinShutterSpeedManual
	PropIsoAutoMinShutterSpeedPreset                        Prop = 0xD177 // IsoAutoMinShutterSpeedPreset
	PropFocusPositionSetting                                Prop = 0xE042 // FocusPositionSetting
	PropSoftSkinEffect                                      Prop = 0xD178 // SoftSkinEffect
	PropPrioritySetInAFS                                    Prop = 0xD179 // PrioritySetInAF_S
	PropPrioritySetInAFC                                    Prop = 0xD17A // PrioritySetInAF_C
	PropFocusMagnificationTime                              Prop = 0xD17B // FocusMagnificationTime
	PropSubjectRecognitionInAF                              Prop = 0xD157 // SubjectRecognitionInAF
	PropRecognitionTarget                                   Prop = 0xD158 // RecognitionTarget
	PropRightLeftEyeSelect                                  Prop = 0xD159 // RightLeftEyeSelect
	PropSelectFTPServer                                     Prop = 0xD27C // SelectFTPServer
	PropSelectFTPServerID                                   Prop = 0xD02E // SelectFTPServerID
	PropFTPFunction                                         Prop = 0xD041 // FTP_Function
	PropFTPAutoTransfer                                     Prop = 0xD04E // FTP_AutoTransfer
	PropFTPAutoTransferTarget                               Prop = 0xD04F // FTP_AutoTransferTarget
	PropMovieFTPAutoTransferTarget                          Prop = 0xD199 // Movie_FTP_AutoTransferTarget
	PropFTPTransferTarget                                   Prop = 0xD19A // FTP_TransferTarget
	PropMovieFTPTransferTarget                              Prop = 0xD14B // Movie_FTP_TransferTarget
	PropFTPPowerSave                                        Prop = 0xD14C // FTP_PowerSave
	PropNDFilterUnitSetting                                 Prop = 0xD14E // NDFilterUnitSetting
	PropNDFilterOpticalDensityValue                         Prop = 0xD14F // NDFilterOpticalDensityValue
	PropTNumber                                             Prop = 0xD000 // TNumber
	PropIrisDisplayUnit                                     Prop = 0xD003 // IrisDisplayUnit
	PropMovieImageStabilizationLevel                        Prop = 0xE080 // Movie_ImageStabilizationLevel
	PropImageStabilizationSteadyShotAdjust                  Prop = 0xD192 // ImageStabilizationSteadyShotAdjust
	PropImageStabilizationSteadyShotFocalLength             Prop = 0xD193 // ImageStabilizationSteadyShotFocalLength
	PropExtendedShutterSpeed                                Prop = 0xD19F // ExtendedShutterSpeed
	PropCameraButtonFunction                                Prop = 0xD208 // CameraButtonFunction
	PropCameraButtonFunctionMulti                           Prop = 0xD209 // CameraButtonFunctionMulti
	PropCameraDialFunction                                  Prop = 0xD20A // CameraDialFunction
	PropSynchroterminalForcedOutput                         Prop = 0xD155 // SynchroterminalForcedOutput
	PropShutterReleaseTimeLagControl                        Prop = 0xD156 // ShutterReleaseTimeLagControl
	PropContinuousShootingSpotBoostFrameSpeed               Prop = 0xD12F // ContinuousShootingSpotBoostFrameSpeed
	PropTimeShiftShooting                                   Prop = 0xD189 // TimeShiftShooting
	PropTimeShiftTriggerSetting                             Prop = 0xD18A // TimeShiftTriggerSetting
	PropTimeShiftPreShootingTimeSetting                     Prop = 0xD18B // TimeShiftPreShootingTimeSetting
	PropEmbedLUTFile                                        Prop = 0xD196 // EmbedLUTFile
	PropAPSCS35                                             Prop = 0xD1C7 // APS_C_S35
	PropLensCompensationShading                             Prop = 0xD1A2 // LensCompensationShading
	PropLensCompensationChromaticAberration                 Prop = 0xD1A3 // LensCompensationChromaticAberration
	PropLensCompensationDistortion                          Prop = 0xD1A4 // LensCompensationDistortion
	PropLensCompensationBreathing                           Prop = 0xD1A5 // LensCompensationBreathing
	PropRecordingMedia                                      Prop = 0xD15F // RecordingMedia
	PropMovieRecordingMedia                                 Prop = 0xD160 // Movie_RecordingMedia
	PropAutoSwitchMedia                                     Prop = 0xD161 // AutoSwitchMedia
	PropRecordingFileNumber                                 Prop = 0xD1C8 // RecordingFileNumber
	PropMovieRecordingFileNumber                            Prop = 0xD2B5 // Movie_RecordingFileNumber
	PropRecordingSettingFileName                            Prop = 0xD1CA // RecordingSettingFileName
	PropRecordingFolderFormat                               Prop = 0xD1CB // RecordingFolderFormat
	PropSelectIPTCMetadata                                  Prop = 0xD132 // SelectIPTCMetadata
	PropWriteCopyrightInfo                                  Prop = 0xD1CD // WriteCopyrightInfo
	PropSetPhotographer                                     Prop = 0xD1CE // SetPhotographer
	PropSetCopyright                                        Prop = 0xD1CF // SetCopyright
	PropFileSettingsTitleNameSettings                       Prop = 0xD1DC // FileSettingsTitleNameSettings
	PropFocusBracketRecordingFolder                         Prop = 0xD1A6 // FocusBracketRecordingFolder
	PropReleaseWithoutLens                                  Prop = 0xD1A7 // ReleaseWithoutLens
	PropReleaseWithoutCard                                  Prop = 0xD1A8 // ReleaseWithoutCard
	PropGridLineDisplay                                     Prop = 0xD1DE // GridLineDisplay
	PropContinuousShootingSpeedInElectricShutterHiPlus      Prop = 0xD162 // ContinuousShootingSpeedInElectricShutterHiPlus
	PropContinuousShootingSpeedInElectricShutterHi          Prop = 0xD163 // ContinuousShootingSpeedInElectricShutterHi
	PropContinuousShootingSpeedInElectricShutterMid         Prop = 0xD164 // ContinuousShootingSpeedInElectricShutterMid
	PropContinuousShootingSpeedInElectricShutterLo          Prop = 0xD165 // ContinuousShootingSpeedInElectricShutterLo
	PropIsoAutoRangeLimitMin                                Prop = 0xD1B6 // IsoAutoRangeLimitMin
	PropIsoAutoRangeLimitMax                                Prop = 0xD1B7 // IsoAutoRangeLimitMax
	PropFacePriorityInMultiMetering                         Prop = 0xD1E5 // FacePriorityInMultiMetering
	PropPrioritySetInAWB                                    Prop = 0xD1AA // PrioritySetInAWB
	PropCustomWBSizeSetting                                 Prop = 0xD135 // CustomWB_Size_Setting
	PropAFIlluminator                                       Prop = 0xD1AB // AFIlluminator
	PropApertureDriveInAF                                   Prop = 0xD1AC // ApertureDriveInAF
	PropAFWithShutter                                       Prop = 0xD1AD // AFWithShutter
	PropFullTimeDMF                                         Prop = 0xD1AE // FullTimeDMF
	PropPreAF                                               Prop = 0xD1AF // PreAF
	PropSubjectRecognitionPersonTrackingSubjectShiftRange   Prop = 0xD1E6 // SubjectRecognitionPersonTrackingSubjectShiftRange
	PropSubjectRecognitionAnimalBirdPriority                Prop = 0xD1E7 // SubjectRecognitionAnimalBirdPriority
	PropSubjectRecognitionAnimalBirdDetectionParts          Prop = 0xD1E8 // SubjectRecognitionAnimalBirdDetectionParts
	PropSubjectRecognitionAnimalTrackingSubjectShiftRange   Prop = 0xD1E9 // SubjectRecognitionAnimalTrackingSubjectShiftRange
	PropSubjectRecognitionAnimalTrackingSensitivity         Prop = 0xD1EA // SubjectRecognitionAnimalTrackingSensitivity
	PropSubjectRecognitionAnimalDetectionSensitivity        Prop = 0xD1EB // SubjectRecognitionAnimalDetectionSensitivity
	PropSubjectRecognitionAnimalDetectionParts              Prop = 0xD1EC // SubjectRecognitionAnimalDetectionParts
	PropSubjectRecognitionBirdTrackingSubjectShiftRange     Prop = 0xD1ED // SubjectRecognitionBirdTrackingSubjectShiftRange
	PropSubjectRecognitionBirdTrackingSensitivity           Prop = 0xD1EE // SubjectRecognitionBirdTrackingSensitivity
	PropSubjectRecognitionBirdDetectionSensitivity          Prop = 0xD1EF // SubjectRecognitionBirdDetectionSensitivity
	PropSubjectRecognitionBirdDetectionParts                Prop = 0xD1F0 // SubjectRecognitionBirdDetectionParts
	PropSubjectRecognitionInsectTrackingSubjectShiftRange   Prop = 0xD1F1 // SubjectRecognitionInsectTrackingSubjectShiftRange
	PropSubjectRecognitionInsectTrackingSensitivity         Prop = 0xD1F2 // SubjectRecognitionInsectTrackingSensitivity
	PropSubjectRecognitionInsectDetectionSensitivity        Prop = 0xD1F3 // SubjectRecognitionInsectDetectionSensitivity
	PropSubjectRecognitionCarTrainTrackingSubjectShiftRange Prop = 0xD1F4 // SubjectRecognitionCarTrainTrackingSubjectShiftRange
	PropSubjectRecognitionCarTrainTrackingSensitivity       Prop = 0xD1F5 // SubjectRecognitionCarTrainTrackingSensitivity
	PropSubjectRecognitionCarTrainDetectionSensitivity      Prop = 0xD1F6 // SubjectRecognitionCarTrainDetectionSensitivity
	PropSubjectRecognitionPlaneTrackingSubjectShiftRange    Prop = 0xD1F7 // SubjectRecognitionPlaneTrackingSubjectShiftRange
	PropSubjectRecognitionPlaneTrackingSensitivity          Prop = 0xD1F8 // SubjectRecognitionPlaneTrackingSensitivity
	PropSubjectRecognitionPlaneDetectionSensitivity         Prop = 0xD1F9 // SubjectRecognitionPlaneDetectionSensitivity
	PropSubjectRecognitionPriorityOnRegisteredFace          Prop = 0xD1FA // SubjectRecognitionPriorityOnRegisteredFace
	PropFaceEyeFrameDisplay                                 Prop = 0xD1B8 // FaceEyeFrameDisplay
	PropFocusMap                                            Prop = 0xD1B9 // FocusMap
	PropInitialFocusMagnifier                               Prop = 0xD1BD // InitialFocusMagnifier
	PropAFInFocusMagnifier                                  Prop = 0xD1BA // AFInFocusMagnifier
	PropAFTrackForSpeedChange                               Prop = 0xD2AD // AFTrackForSpeedChange
	PropAFFreeSizeAndPositionSetting                        Prop = 0xD138 // AFFreeSizeAndPositionSetting
	PropPlaySetOfMultiMedia                                 Prop = 0xD2B2 // PlaySetOfMultiMedia
	PropRemoteSaveImageSize                                 Prop = 0xD1BE // RemoteSaveImageSize
	PropFTPTransferStillImageQualitySize                    Prop = 0xD14A // FTP_TransferStillImageQualitySize
	PropFTPAutoTransferTargetStillImage                     Prop = 0xD216 // FTP_AutoTransferTarget_StillImage
	PropProtectImageInFTPTransfer                           Prop = 0xD225 // ProtectImageInFTPTransfer
	PropMonitorBrightnessType                               Prop = 0xD1FB // MonitorBrightnessType
	PropMonitorBrightnessManual                             Prop = 0xD1FC // MonitorBrightnessManual
	PropDisplayQualityFinderMonitor                         Prop = 0xD1B0 // DisplayQualityFinderMonitor
	PropTCUBDisplaySetting                                  Prop = 0xD1FD // TCUBDisplaySetting
	PropGammaDisplayAssist                                  Prop = 0xD1FE // GammaDisplayAssist
	PropGammaDisplayAssistType                              Prop = 0xD1FF // GammaDisplayAssistType
	PropAudioSignalsStartEnd                                Prop = 0xD220 // AudioSignalsStartEnd
	PropAudioSignalsVolume                                  Prop = 0xD1B1 // AudioSignalsVolume
	PropControlForHDMI                                      Prop = 0xD1B2 // ControlForHDMI
	PropAntidustShutterWhenPowerOff                         Prop = 0xD1B3 // AntidustShutterWhenPowerOff
	PropWakeOnLAN                                           Prop = 0xD1C1 // WakeOnLAN
	PropReserved10                                          Prop = 0xD239 // reserved10
	PropReserved11                                          Prop = 0xD23A // reserved11
	PropReserved12                                          Prop = 0xD23B // reserved12
	PropIntervalRecMode                                     Prop = 0xD24F // Interval_Rec_Mode
	PropStillImageTransSize                                 Prop = 0xD268 // Still_Image_Trans_Size
	PropRAWJPCSaveImage                                     Prop = 0xD269 // RAW_J_PC_Save_Image
	PropLiveViewImageQuality                                Prop = 0xD26A // LiveView_Image_Quality
	PropRemoconZoomSpeedType                                Prop = 0xD299 // Remocon_Zoom_Speed_Type
	PropSnapshotInfo                                        Prop = 0xD215 // SnapshotInfo
	PropBatteryRemain                                       Prop = 0xD218 // BatteryRemain
	PropBatteryLevel                                        Prop = 0xD20E // BatteryLevel
	PropEstimatePictureSize                                 Prop = 0xD214 // EstimatePictureSize
	PropRecordingState                                      Prop = 0xD21D // RecordingState
	PropLiveViewStatus                                      Prop = 0xD221 // LiveViewStatus
	PropFocusIndication                                     Prop = 0xD213 // FocusIndication
	PropMediaSLOT1Status                                    Prop = 0xD248 // MediaSLOT1_Status
	PropMediaSLOT1RemainingNumber                           Prop = 0xD249 // MediaSLOT1_RemainingNumber
	PropMediaSLOT1RemainingTime                             Prop = 0xD24A // MediaSLOT1_RemainingTime
	PropMediaSLOT1FormatEnableStatus                        Prop = 0xD279 // MediaSLOT1_FormatEnableStatus
	PropReserved20                                          Prop = 0xD27D // reserved20
	PropMediaSLOT2Status                                    Prop = 0xD256 // MediaSLOT2_Status
	PropMediaSLOT2FormatEnableStatus                        Prop = 0xD27A // MediaSLOT2_FormatEnableStatus
	PropMediaSLOT2RemainingNumber                           Prop = 0xD257 // MediaSLOT2_RemainingNumber
	PropMediaSLOT2RemainingTime                             Prop = 0xD258 // MediaSLOT2_RemainingTime
	PropReserved22                                          Prop = 0xD27E // reserved22
	PropMediaFormatProgressRate                             Prop = 0xD27B // Media_FormatProgressRate
	PropFTPConnectionStatus                                 Prop = 0xD27F // FTP_ConnectionStatus
	PropFTPConnectionErrorInfo                              Prop = 0xD280 // FTP_ConnectionErrorInfo
	PropLiveViewArea                                        Prop = 0xD267 // LiveView_Area
	PropReserved26                                          Prop = 0xD23C // reserved26
	PropReserved27                                          Prop = 0xD23D // reserved27
	PropIntervalRecStatus                                   Prop = 0xD250 // Interval_Rec_Status
	PropCustomWBExecutionState                              Prop = 0xD270 // CustomWB_Execution_State
	PropCustomWBCapturableArea                              Prop = 0xD26B // CustomWB_Capturable_Area
	PropCustomWBCaptureFrameSize                            Prop = 0xD26C // CustomWB_Capture_Frame_Size
	PropCustomWBCaptureOperation                            Prop = 0xD26F // CustomWB_Capture_Operation
	PropZoomOperationStatus                                 Prop = 0xD25B // Zoom_Operation_Status
	PropZoomBarInformation                                  Prop = 0xD25D // Zoom_Bar_Information
	PropZoomTypeStatus                                      Prop = 0xD260 // Zoom_Type_Status
	PropMediaSLOT1QuickFormatEnableStatus                   Prop = 0xD292 // MediaSLOT1_QuickFormatEnableStatus
	PropMediaSLOT2QuickFormatEnableStatus                   Prop = 0xD293 // MediaSLOT2_QuickFormatEnableStatus
	PropCancelMediaFormatEnableStatus                       Prop = 0xD294 // Cancel_Media_FormatEnableStatus
	PropZoomSpeedRange                                      Prop = 0xD25E // Zoom_Speed_Range
	PropIsoCurrentSensitivity                               Prop = 0xD023 // IsoCurrentSensitivity
	PropCameraSettingSaveOperationEnableStatus              Prop = 0xD271 // CameraSetting_SaveOperationEnableStatus
	PropCameraSettingReadOperationEnableStatus              Prop = 0xD272 // CameraSetting_ReadOperationEnableStatus
	PropCameraSettingSaveReadState                          Prop = 0xD273 // CameraSetting_SaveRead_State
	PropCameraSettingsResetEnableStatus                     Prop = 0xD043 // CameraSettingsResetEnableStatus
	PropAPSCOrFullSwitchingSetting                          Prop = 0xD29A // APS_C_or_Full_SwitchingSetting
	PropAPSCOrFullSwitchingEnableStatus                     Prop = 0xD29B // APS_C_or_Full_SwitchingEnableStatus
	PropDispModeCandidate                                   Prop = 0xD044 // DispModeCandidate
	PropShutterSpeedCurrentValue                            Prop = 0xD017 // ShutterSpeedCurrentValue
	PropFocusSpeedRange                                     Prop = 0xD008 // Focus_Speed_Range
	PropNDFilterMode                                        Prop = 0xD019 // NDFilterMode
	PropMoviePlayingSpeed                                   Prop = 0xD030 // MoviePlayingSpeed
	PropMediaSLOT1Player                                    Prop = 0xD035 // MediaSLOT1Player
	PropMediaSLOT2Player                                    Prop = 0xD036 // MediaSLOT2Player
	PropBatteryRemainingInMinutes                           Prop = 0xD038 // BatteryRemainingInMinutes
	PropBatteryRemainingInVoltage                           Prop = 0xD039 // BatteryRemainingInVoltage
	PropDCVoltage                                           Prop = 0xD03E // DCVoltage
	PropMoviePlayingState                                   Prop = 0xD02F // MoviePlayingState
	PropFocusTouchSpotStatus                                Prop = 0xE004 // FocusTouchSpotStatus
	PropFocusTrackingStatus                                 Prop = 0xE005 // FocusTrackingStatus
	PropDepthOfFieldAdjustmentInterlockingMode              Prop = 0xE00A // DepthOfFieldAdjustmentInterlockingMode
	PropRecorderClipName                                    Prop = 0xE00B // RecorderClipName
	PropRecorderControlMainSetting                          Prop = 0xE00C // RecorderControlMainSetting
	PropRecorderStartMain                                   Prop = 0xE00E // RecorderStartMain
	PropRecorderStartProxy                                  Prop = 0xE00F // RecorderStartProxy
	PropRecorderMainStatus                                  Prop = 0xE010 // RecorderMainStatus
	PropRecorderProxyStatus                                 Prop = 0xE011 // RecorderProxyStatus
	PropRecorderExtRawStatus                                Prop = 0xE012 // RecorderExtRawStatus
	PropRecorderSaveDestination                             Prop = 0xE013 // RecorderSaveDestination
	PropAssignableButtonIndicator1                          Prop = 0xE035 // AssignableButtonIndicator1
	PropAssignableButtonIndicator2                          Prop = 0xE036 // AssignableButtonIndicator2
	PropAssignableButtonIndicator3                          Prop = 0xE037 // AssignableButtonIndicator3
	PropAssignableButtonIndicator4                          Prop = 0xE038 // AssignableButtonIndicator4
	PropAssignableButtonIndicator5                          Prop = 0xE039 // AssignableButtonIndicator5
	PropAssignableButtonIndicator6                          Prop = 0xE03A // AssignableButtonIndicator6
	PropAssignableButtonIndicator7                          Prop = 0xE03B // AssignableButtonIndicator7
	PropAssignableButtonIndicator8                          Prop = 0xE03C // AssignableButtonIndicator8
	PropAssignableButtonIndicator9                          Prop = 0xE03D // AssignableButtonIndicator9
	PropAssignableButtonIndicator10                         Prop = 0xE03E // AssignableButtonIndicator10
	PropAssignableButtonIndicator11                         Prop = 0xE08F // AssignableButtonIndicator11
	PropLensAssignableButtonIndicator1                      Prop = 0xE03F // LensAssignableButtonIndicator1
	PropGaindBCurrentValue                                  Prop = 0xD01F // GaindBCurrentValue
	PropSoftwareVersion                                     Prop = 0xD040 // SoftwareVersion
	PropCurrentSceneFileEdited                              Prop = 0xE020 // CurrentSceneFileEdited
	PropMovieRecButtonToggleEnableStatus                    Prop = 0xE061 // MovieRecButtonToggleEnableStatus
	PropRemoteTouchOperationEnableStatus                    Prop = 0xD284 // RemoteTouchOperationEnableStatus
	PropCancelRemoteTouchOperationEnableStatus              Prop = 0xD285 // CancelRemoteTouchOperationEnableStatus
	PropLensInformationEnableStatus                         Prop = 0xE086 // LensInformationEnableStatus
	PropFollowFocusPositionCurrentValue                     Prop = 0xE089 // FollowFocusPositionCurrentValue
	PropFocusBracketShootingStatus                          Prop = 0xD0AB // FocusBracketShootingStatus
	PropPixelMappingEnableStatus                            Prop = 0xD0C6 // PixelMappingEnableStatus
	PropTimeCodePresetResetEnableStatus                     Prop = 0xD104 // TimeCodePresetResetEnableStatus
	PropUserBitPresetResetEnableStatus                      Prop = 0xD105 // UserBitPresetResetEnableStatus
	PropSensorCleaningEnableStatus                          Prop = 0xD106 // SensorCleaningEnableStatus
	PropPictureProfileResetEnableStatus                     Prop = 0xD107 // PictureProfileResetEnableStatus
	PropCreativeLookResetEnableStatus                       Prop = 0xD108 // CreativeLookResetEnableStatus
	PropLensVersionNumber                                   Prop = 0xD07D // LensVersionNumber
	PropDeviceOverheatingState                              Prop = 0xD251 // DeviceOverheatingState
	PropMovieIntervalRecCountDownIntervalTime               Prop = 0xD11F // Movie_IntervalRec_CountDownIntervalTime
	PropMovieIntervalRecRecordingDuration                   Prop = 0xD120 // Movie_IntervalRec_RecordingDuration
	PropHighResolutionShutterSpeed                          Prop = 0xD282 // HighResolutionShutterSpeed
	PropBaseLookImportOperationEnableStatus                 Prop = 0xD08B // BaseLookImportOperationEnableStatus
	PropLensModelName                                       Prop = 0xD07B // LensModelName
	PropFocusPositionCurrentValue                           Prop = 0xE043 // FocusPositionCurrentValue
	PropFocusDrivingStatus                                  Prop = 0xD19C // FocusDrivingStatus
	PropFlickerScanStatus                                   Prop = 0xD2BA // FlickerScanStatus
	PropFlickerScanEnableStatus                             Prop = 0xD2BB // FlickerScanEnableStatus
	PropFTPServerSettingOperationEnableStatus               Prop = 0xD09A // FTPServerSettingOperationEnableStatus
	PropFTPTransferSettingSaveOperationEnableStatus         Prop = 0xD274 // FTPTransferSetting_SaveOperationEnableStatus
	PropFTPTransferSettingReadOperationEnableStatus         Prop = 0xD275 // FTPTransferSetting_ReadOperationEnableStatus
	PropFTPTransferSettingSaveReadState                     Prop = 0xD276 // FTPTransferSetting_SaveRead_State
	PropCameraShakeStatus                                   Prop = 0xD194 // CameraShakeStatus
	PropUpdateBodyStatus                                    Prop = 0xD195 // UpdateBodyStatus
	PropMediaSLOT1WritingState                              Prop = 0xD197 // MediaSLOT1_WritingState
	PropMediaSLOT2WritingState                              Prop = 0xD198 // MediaSLOT2_WritingState
	PropMediaSLOT1RecordingAvailableType                    Prop = 0xD02B // MediaSLOT1_RecordingAvailableType
	PropMediaSLOT2RecordingAvailableType                    Prop = 0xD02C // MediaSLOT2_RecordingAvailableType
	PropMediaSLOT3RecordingAvailableType                    Prop = 0xD190 // MediaSLOT3_RecordingAvailableType
	PropCameraOperatingMode                                 Prop = 0xD0BC // CameraOperatingMode
	PropPlaybackViewMode                                    Prop = 0xD0BD // PlaybackViewMode
	PropMediaSLOT3Status                                    Prop = 0xD18E // MediaSLOT3_Status
	PropMediaSLOT3RemainingTime                             Prop = 0xD18F // MediaSLOT3_RemainingTime
	PropMonitoringDeliveringStatus                          Prop = 0xE098 // MonitoringDeliveringStatus
	PropMonitoringIsDelivering                              Prop = 0xE099 // MonitoringIsDelivering
	PropMonitoringSettingVersion                            Prop = 0xE09D // MonitoringSettingVersion
	PropMonitoringDeliveryTypeSupportInfo                   Prop = 0xE09F // MonitoringDeliveryTypeSupportInfo
	PropCameraErrorCautionStatus                            Prop = 0xD1BB // CameraErrorCautionStatus
	PropSystemErrorCautionStatus                            Prop = 0xD1BC // SystemErrorCautionStatus
	PropCameraButtonFunctionStatus                          Prop = 0xD20C // CameraButtonFunctionStatus
	PropFlickerLessShootingStatus                           Prop = 0xD134 // FlickerLessShootingStatus
	PropContinuousShootingSpotBoostStatus                   Prop = 0xD130 // ContinuousShootingSpotBoostStatus
	PropContinuousShootingSpotBoostEnableStatus             Prop = 0xD131 // ContinuousShootingSpotBoostEnableStatus
	PropTimeShiftShootingStatus                             Prop = 0xD19B // TimeShiftShootingStatus
	PropZoomDrivingStatus                                   Prop = 0xD19D // ZoomDrivingStatus
	PropShootingSelfTimerStatus                             Prop = 0xD1B4 // ShootingSelfTimerStatus
	PropCreateNewFolderEnableStatus                         Prop = 0xD1DB // CreateNewFolderEnableStatus
	PropForcedFileNumberResetEnableStatus                   Prop = 0xD1C9 // ForcedFileNumberResetEnableStatus
	PropDefaultAFFreeSizeAndPositionSetting                 Prop = 0xD19E // DefaultAFFreeSizeAndPositionSetting
	PropTrackingOnAndAFOnEnableStatus                       Prop = 0xD1C6 // TrackingOnAndAFOnEnableStatus
	PropProgramShiftStatus                                  Prop = 0xD1BF // ProgramShiftStatus
	PropMeteredManualLevel                                  Prop = 0xD1B5 // MeteredManualLevel
	PropSecondBatteryRemain                                 Prop = 0xD12D // SecondBatteryRemain
	PropSecondBatteryLevel                                  Prop = 0xD12E // SecondBatteryLevel
	PropTotalBatteryRemain                                  Prop = 0xD204 // TotalBatteryRemain
	PropTotalBatteryLevel                                   Prop = 0xD205 // TotalBatteryLevel
	PropCameraLeverFunction                                 Prop = 0xD20B // CameraLeverFunction
	PropShootingTimingPreNotificationMode                   Prop = 0xD202 // ShootingTimingPreNotificationMode
	PropMicrophoneDirectivity                               Prop = 0xD1DD // MicrophoneDirectivity
	PropProductShowcaseSet                                  Prop = 0xD1DF // ProductShowcaseSet
	PropAmountOfDefocusSetting                              Prop = 0xD1E0 // AmountOfDefocusSetting
	PropCinematicVlogSetting                                Prop = 0xD1E1 // CinematicVlogSetting
	PropCinematicVlogLook                                   Prop = 0xD1E2 // CinematicVlogLook
	PropCinematicVlogMood                                   Prop = 0xD1E3 // CinematicVlogMood
	PropCinematicVlogAFTransitionSpeed                      Prop = 0xD1E4 // CinematicVlogAFTransitionSpeed
	PropMonitoringTransportProtocol                         Prop = 0xE0A1 // MonitoringTransportProtocol
	PropMonitoringAvailableFormat                           Prop = 0xE0A2 // MonitoringAvailableFormat
	PropMonitoringFormatSupportInformation                  Prop = 0xE0A3 // MonitoringFormatSupportInformation
	PropDeSqueezeDisplayRatio                               Prop = 0xE0A4 // DeSqueezeDisplayRatio
	PropZoomPositionSetting                                 Prop = 0xE040 // ZoomPositionSetting
	PropZoomPositionCurrentValue                            Prop = 0xE041 // ZoomPositionCurrentValue
	PropPriv0F07                                            Prop = 0xE0CF // private
	PropPriv0F08                                            Prop = 0xE0D0 // private
	PropPriv0F09                                            Prop = 0xE0D1 // private
	PropPriv0F0A                                            Prop = 0xD2B4 // private
	PropPriv0F0C                                            Prop = 0xE0B3 // private
	PropMonitoringOutputDisplaySDI                          Prop = 0xD078 // MonitoringOutputDisplaySDI
	PropCameraSystemErrorInfo                               Prop = 0xD07A // CameraSystemErrorInfo
	PropAFAreaPositionAFC                                   Prop = 0xE079 // AFAreaPositionAF_C
	PropAFAreaPositionAFS                                   Prop = 0xE078 // AFAreaPositionAF_S
	PropFaceEyeDetectionAFStatus                            Prop = 0xD09B // FaceEyeDetectionAFStatus
	PropAutoFocusHold                                       Prop = 0xD091 // AutoFocusHold
	PropPushAFModeSetting                                   Prop = 0xD06E // PushAFModeSetting
	PropTouchFunctionInMF                                   Prop = 0xD064 // TouchFunctionInMF
	PropPushAutoFocus                                       Prop = 0xD05E // PushAutoFocus
	PropPushAGC                                             Prop = 0xD05C // PushAGC
	PropPushAutoIris                                        Prop = 0xD05B // PushAutoIris
	PropNDFilterPreset3Value                                Prop = 0xD077 // NDFilterPreset3Value
	PropNDFilterPreset2Value                                Prop = 0xD076 // NDFilterPreset2Value
	PropNDFilterPreset1Value                                Prop = 0xD075 // NDFilterPreset1Value
	PropNDFilterPresetSelect                                Prop = 0xD074 // NDFilterPresetSelect
	PropPushAutoNDFilter                                    Prop = 0xD05D // PushAutoNDFilter
	PropWhiteBalanceOffsetColorTemp                         Prop = 0xD0AA // WhiteBalanceOffsetColorTemp
	PropWhiteBalanceOffsetSetting                           Prop = 0xD0A9 // WhiteBalanceOffsetSetting
	PropWhiteBalanceOffsetTintATW                           Prop = 0xD08A // WhiteBalanceOffsetTintATW
	PropWhiteBalanceOffsetColorTempATW                      Prop = 0xD089 // WhiteBalanceOffsetColorTempATW
	PropWhiteBalanceBGain                                   Prop = 0xD088 // WhiteBalanceBGain
	PropWhiteBalanceRGain                                   Prop = 0xD087 // WhiteBalanceRGain
	PropWhiteBalancePresetColorTemperature                  Prop = 0xD086 // WhiteBalancePresetColorTemperature
	PropWhiteBalanceSwitch                                  Prop = 0xD085 // WhiteBalanceSwitch
	PropPaintLookDetailLevel                                Prop = 0xD06D // PaintLookDetailLevel
	PropPaintLookDetailSetting                              Prop = 0xD06C // PaintLookDetailSetting
	PropPaintLookKneeSlope                                  Prop = 0xD06B // PaintLookKneeSlope
	PropPaintLookKneePoint                                  Prop = 0xD06A // PaintLookKneePoint
	PropPaintLookAutoKnee                                   Prop = 0xD069 // PaintLookAutoKnee
	PropPaintLookKneeSetting                                Prop = 0xD068 // PaintLookKneeSetting
	PropPaintLookBBlack                                     Prop = 0xD067 // PaintLookBBlack
	PropPaintLookRBlack                                     Prop = 0xD066 // PaintLookRBlack
	PropPaintLookMasterBlack                                Prop = 0xD065 // PaintLookMasterBlack
	PropUploadDatasetVersion                                Prop = 0xD057 // UploadDatasetVersion
	PropUserBaseLookOutput                                  Prop = 0xD2B1 // UserBaseLookOutput
	PropMonitorLUTSettingOutputDestAssign                   Prop = 0xD2A9 // MonitorLUTSettingOutputDestAssign
	PropMonitorLUTSetting1                                  Prop = 0xD2AA // MonitorLUTSetting1
	PropMonitorLUTSetting2                                  Prop = 0xD2AB // MonitorLUTSetting2
	PropMonitorLUTSetting3                                  Prop = 0xD2AC // MonitorLUTSetting3
	PropMaximumNumberOfBytes                                Prop = 0xD0C1 // MaximumNumberOfBytes
	PropSQModeSetting                                       Prop = 0xD051 // SQModeSetting
	PropMovieQualityFullAutoMode                            Prop = 0xE0A0 // MovieQualityFullAutoMode
	PropFileSettingsCameraId                                Prop = 0xD1D0 // FileSettingsCameraId
	PropFileSettingsReelNumber                              Prop = 0xD1D1 // FileSettingsReelNumber
	PropFileSettingsCameraPosition                          Prop = 0xD1D2 // FileSettingsCameraPosition
	PropImageStabilizationFramingStabilizer                 Prop = 0xD1A9 // ImageStabilizationFramingStabilizer
	PropExposureStep                                        Prop = 0xD237 // ExposureStep
	PropTeleWideLeverValueCapability                        Prop = 0xD2BD // TeleWideLeverValueCapability
	PropEnlargeScreenSetting                                Prop = 0xE0CD // EnlargeScreenSetting
	PropMediaSLOT1ContentsInfoListEnableStatus              Prop = 0xD1D4 // MediaSLOT1_ContentsInfoListEnableStatus
	PropMediaSLOT2ContentsInfoListEnableStatus              Prop = 0xD1D5 // MediaSLOT2_ContentsInfoListEnableStatus
	PropMediaSLOT1ContentsInfoListRegenerateUpdateTime      Prop = 0xD1D6 // MediaSLOT1_ContentsInfoListRegenerateUpdateTime
	PropMediaSLOT2ContentsInfoListRegenerateUpdateTime      Prop = 0xD1D7 // MediaSLOT2_ContentsInfoListRegenerateUpdateTime
	PropMediaSLOT1ContentsInfoListUpdateTime                Prop = 0xE0D9 // MediaSLOT1_ContentsInfoListUpdateTime
	PropMediaSLOT2ContentsInfoListUpdateTime                Prop = 0xE0DA // MediaSLOT2_ContentsInfoListUpdateTime
	PropPostViewTransferResourceStatus                      Prop = 0xD1DA // PostViewTransferResourceStatus
	PropSimulRecSetting                                     Prop = 0xD050 // SimulRecSetting
	PropSimulRecSettingMovieRecButton                       Prop = 0xD07F // SimulRecSettingMovieRecButton
	PropShutterSelectMode                                   Prop = 0xD012 // ShutterSelectMode
	PropOSDImageMode                                        Prop = 0xD207 // OSDImageMode
	PropFirmwareUpdateCommandVersion                        Prop = 0xD0BF // FirmwareUpdateCommandVersion
	PropDebugMode                                           Prop = 0xE0D5 // DebugMode
	PropPriv0F0B                                            Prop = 0xD0C0 // private
	PropReserved18                                          Prop = 0xE0DF // reserved18
	PropReserved19                                          Prop = 0xE0E0 // reserved19
	PropPriv0F0D                                            Prop = 0xE0C0 // private
	PropSetPresetPTZFBinaryVersion                          Prop = 0xE0CC // SetPresetPTZFBinaryVersion
	PropPanPositionStatus                                   Prop = 0xE0B7 // PanPositionStatus
	PropTiltPositionStatus                                  Prop = 0xE0B9 // TiltPositionStatus
	PropPanPositionCurrentValue                             Prop = 0xE0B6 // PanPositionCurrentValue
	PropTiltPositionCurrentValue                            Prop = 0xE0B8 // TiltPositionCurrentValue
	PropPanTiltAccelerationRampCurve                        Prop = 0xE05C // PanTiltAccelerationRampCurve
	PropPanLimitMode                                        Prop = 0xE0BE // PanLimitMode
	PropPanLimitRangeMinimum                                Prop = 0xE0BA // PanLimitRangeMinimum
	PropPanLimitRangeMaximum                                Prop = 0xE0BB // PanLimitRangeMaximum
	PropTiltLimitMode                                       Prop = 0xE0BF // TiltLimitMode
	PropTiltLimitRangeMinimum                               Prop = 0xE0BC // TiltLimitRangeMinimum
	PropTiltLimitRangeMaximum                               Prop = 0xE0BD // TiltLimitRangeMaximum
	PropPresetPTZFSlotNumber                                Prop = 0xE0CB // PresetPTZFSlotNumber
	PropCameraPowerStatus                                   Prop = 0xE0AF // CameraPowerStatus
	PropTargetStreamingDestinationSelect                    Prop = 0xD2BE // TargetStreamingDestinationSelect
	PropStreamStatus                                        Prop = 0xD119 // StreamStatus
	PropIRRemoteSetting                                     Prop = 0xE0A5 // IRRemoteSetting
	PropIPSetupProtocolSetting                              Prop = 0xE0A6 // IPSetupProtocolSetting
	PropRecordablePowerSources                              Prop = 0xE0B0 // RecordablePowerSources
	PropStreamSettingListOperationStatus                    Prop = 0xD1CC // StreamSettingListOperationStatus
	PropPaintLookMultiMatrixAreaIndication                  Prop = 0xD29E // PaintLookMultiMatrixAreaIndication
	PropIrisCloseSetting                                    Prop = 0xD002 // IrisCloseSetting
	PropDisplayedMenuStatus                                 Prop = 0xD080 // DisplayedMenuStatus
	PropLanguageSetting                                     Prop = 0xD08F // LanguageSetting
	PropPlaybackContentsRecordingDateTime                   Prop = 0xD09C // PlaybackContentsRecordingDateTime
	PropPlaybackContentsName                                Prop = 0xD09D // PlaybackContentsName
	PropPlaybackContentsNumber                              Prop = 0xD09E // PlaybackContentsNumber
	PropPlaybackContentsTotalNumber                         Prop = 0xD09F // PlaybackContentsTotalNumber
	PropPlaybackContentsRecordingResolution                 Prop = 0xD0A0 // PlaybackContentsRecordingResolution
	PropPlaybackContentsRecordingFrameRate                  Prop = 0xD0A1 // PlaybackContentsRecordingFrameRate
	PropPlaybackContentsRecordingFileFormat                 Prop = 0xD0A2 // PlaybackContentsRecordingFileFormat
	PropPlaybackContentsGammaType                           Prop = 0xD0A4 // PlaybackContentsGammaType
	PropBaseLookNameofPlayback                              Prop = 0xD0C3 // BaseLookNameofPlayback
	PropBaseLookAppliedofPlayback                           Prop = 0xD0C4 // BaseLookAppliedofPlayback
	PropPaintLookUserMatrixSetting                          Prop = 0xD139 // PaintLookUserMatrixSetting
	PropPaintLookUserMatrixLevel                            Prop = 0xD13A // PaintLookUserMatrixLevel
	PropPaintLookUserMatrixPhase                            Prop = 0xD13B // PaintLookUserMatrixPhase
	PropPaintLookUserMatrixRG                               Prop = 0xD13C // PaintLookUserMatrixRG
	PropPaintLookUserMatrixRB                               Prop = 0xD13D // PaintLookUserMatrixRB
	PropPaintLookUserMatrixGR                               Prop = 0xD13E // PaintLookUserMatrixGR
	PropPaintLookUserMatrixGB                               Prop = 0xD13F // PaintLookUserMatrixGB
	PropPaintLookUserMatrixBR                               Prop = 0xD140 // PaintLookUserMatrixBR
	PropPaintLookUserMatrixBG                               Prop = 0xD141 // PaintLookUserMatrixBG
	PropPaintLookMultiMatrixSetting                         Prop = 0xD142 // PaintLookMultiMatrixSetting
	PropPaintLookMultiMatrixAxis                            Prop = 0xD143 // PaintLookMultiMatrixAxis
	PropPaintLookMultiMatrixHue                             Prop = 0xD144 // PaintLookMultiMatrixHue
	PropPaintLookMultiMatrixSaturation                      Prop = 0xD145 // PaintLookMultiMatrixSaturation
	PropMonitoringOutputDisplaySettingDestAssign            Prop = 0xD2AE // MonitoringOutputDisplaySettingDestAssign
	PropMonitoringOutputDisplaySetting1                     Prop = 0xD2AF // MonitoringOutputDisplaySetting1
	PropMonitoringOutputDisplaySetting2                     Prop = 0xD2B0 // MonitoringOutputDisplaySetting2
	PropFocusModeStatus                                     Prop = 0xE044 // FocusModeStatus
	PropFocusOperationWithInt16EnableStatus                 Prop = 0xE045 // FocusOperationWithInt16EnableStatus
	PropAudioInputCH1LevelControl                           Prop = 0xE048 // AudioInputCH1LevelControl
	PropAudioInputCH2LevelControl                           Prop = 0xE049 // AudioInputCH2LevelControl
	PropAudioInputCH3LevelControl                           Prop = 0xE04A // AudioInputCH3LevelControl
	PropAudioInputCH4LevelControl                           Prop = 0xE04B // AudioInputCH4LevelControl
	PropAudioInputCH1Level                                  Prop = 0xE04C // AudioInputCH1Level
	PropAudioInputCH2Level                                  Prop = 0xE04D // AudioInputCH2Level
	PropAudioInputCH3Level                                  Prop = 0xE04E // AudioInputCH3Level
	PropAudioInputCH4Level                                  Prop = 0xE04F // AudioInputCH4Level
	PropAudioInputCH1InputSelect                            Prop = 0xE051 // AudioInputCH1InputSelect
	PropAudioInputCH2InputSelect                            Prop = 0xE052 // AudioInputCH2InputSelect
	PropAudioInputCH3InputSelect                            Prop = 0xE053 // AudioInputCH3InputSelect
	PropAudioInputCH4InputSelect                            Prop = 0xE054 // AudioInputCH4InputSelect
	PropAudioInputCH1WindFilter                             Prop = 0xE055 // AudioInputCH1WindFilter
	PropAudioInputCH2WindFilter                             Prop = 0xE056 // AudioInputCH2WindFilter
	PropAudioInputCH3WindFilter                             Prop = 0xE057 // AudioInputCH3WindFilter
	PropAudioInputCH4WindFilter                             Prop = 0xE058 // AudioInputCH4WindFilter
	PropRemoteKeyThumbnailButton                            Prop = 0xE05F // RemoteKeyThumbnailButton
	PropRemoteKeySLOTSelectButton                           Prop = 0xE060 // RemoteKeySLOTSelectButton
	PropVideoRecordingFormatBitrateSetting                  Prop = 0xE067 // VideoRecordingFormatBitrateSetting
	PropValidRecordingVideoFormat                           Prop = 0xE068 // ValidRecordingVideoFormat
	PropMonitoringOutputFormat                              Prop = 0xE06C // MonitoringOutputFormat
	PropFocusSpeedDirectSync                                Prop = 0xE092 // FocusSpeedDirectSync
	PropAudioInput1TypeSelect                               Prop = 0xE0AA // AudioInput1TypeSelect
	PropAudioInput2TypeSelect                               Prop = 0xE0AB // AudioInput2TypeSelect
	PropVideoRecordingFormatQuality                         Prop = 0xE0B5 // VideoRecordingFormatQuality
	PropLiveViewImageQualityByNumericalValue                Prop = 0xE0CE // LiveViewImageQualityByNumericalValue
	PropTallyLampControlRed                                 Prop = 0xE0D2 // TallyLampControlRed
	PropTallyLampControlGreen                               Prop = 0xE0D3 // TallyLampControlGreen
	PropTallyLampControlYellow                              Prop = 0xE0D4 // TallyLampControlYellow
	PropMovieRecordingResolutionForRTSP                     Prop = 0xD026 // Movie_RecordingResolutionForRTSP
	PropMovieRecordingFrameRateRTSPSetting                  Prop = 0xD02D // Movie_RecordingFrameRateRTSPSetting
	PropPictureCacheRecSetting                              Prop = 0xD053 // PictureCacheRecSetting
	PropPictureCacheRecSizeAndTime                          Prop = 0xD054 // PictureCacheRecSizeAndTime
	PropMovieIntervalRecFrames                              Prop = 0xD056 // Movie_IntervalRecFrames
	PropImagerScanMode                                      Prop = 0xD05F // ImagerScanMode
	PropMovieRecordingResolutionForRAW                      Prop = 0xD063 // Movie_RecordingResolutionForRAW
	PropLensSerialNumber                                    Prop = 0xD07C // LensSerialNumber
	PropShootingEnableSettingLicense                        Prop = 0xE0FA // ShootingEnableSettingLicense
	PropGridLineDisplayPlayback                             Prop = 0xE10E // GridLineDisplayPlayback
	PropGridLineType                                        Prop = 0xE10F // GridLineType
	PropCustomGridLineFileCommandVersion                    Prop = 0xE110 // CustomGridLineFileCommandVersion
	PropMaximumSizeOfImageIDString                          Prop = 0xE0FE // MaximumSizeOfImageIDString
	PropStreamButtonEnableStatus                            Prop = 0xE112 // StreamButtonEnableStatus
	PropAutoRecognitionTargetCandidates                     Prop = 0xD229 // AutoRecognitionTargetCandidates
	PropAutoRecognitionTargetSetting                        Prop = 0xD234 // AutoRecognitionTargetSetting
	PropDeleteContentOperationEnableStatusSLOT1             Prop = 0xE0F3 // DeleteContentOperationEnableStatusSLOT1
	PropDeleteContentOperationEnableStatusSLOT2             Prop = 0xE0F4 // DeleteContentOperationEnableStatusSLOT2
	PropDifferentSetForSQMovie                              Prop = 0xE0F5 // DifferentSetForSQMovie
	PropManualInputForNDFilterValue                         Prop = 0xE0E2 // ManualInputForNDFilterValue
	PropLogShootingMode                                     Prop = 0xE0E3 // LogShootingMode
	PropLogShootingModeColorGamut                           Prop = 0xE0E4 // LogShootingModeColorGamut
	PropVideoStreamSelect                                   Prop = 0xD0AC // VideoStreamSelect
	PropStreamDisplayName                                   Prop = 0xD2BF // StreamDisplayName
	PropVideoStreamResolution                               Prop = 0xD0AE // VideoStreamResolution
	PropVideoStreamMaxBitRate                               Prop = 0xD2C0 // VideoStreamMaxBitRate
	PropVideoStreamAdaptiveRateControl                      Prop = 0xD10E // VideoStreamAdaptiveRateControl
	PropVideoStreamCodec                                    Prop = 0xD0AD // VideoStreamCodec
	PropStreamLatency                                       Prop = 0xD116 // StreamLatency
	PropStreamTTL                                           Prop = 0xD117 // StreamTTL
	PropStreamCipherType                                    Prop = 0xD113 // StreamCipherType
	PropStreamModeSetting                                   Prop = 0xD118 // StreamModeSetting
	PropVideoStreamResolutionMethod                         Prop = 0xD10F // VideoStreamResolutionMethod
	PropVideoStreamMovieRecPermission                       Prop = 0xD11A // VideoStreamMovieRecPermission
	PropVideoStreamBitRateCompressionMode                   Prop = 0xD0B5 // VideoStreamBitRateCompressionMode
	PropVideoStreamBitRateVBRMode                           Prop = 0xD2B3 // VideoStreamBitRateVBRMode
	PropAudioStreamCodecType                                Prop = 0xD11B // AudioStreamCodecType
	PropAudioStreamSamplingFrequency                        Prop = 0xD11C // AudioStreamSamplingFrequency
	PropAudioStreamBitDepth                                 Prop = 0xD11D // AudioStreamBitDepth
	PropAudioStreamChannel                                  Prop = 0xD11E // AudioStreamChannel
	PropHomeMenuSetting                                     Prop = 0xE111 // HomeMenuSetting
	PropCallSetting                                         Prop = 0xE0E1 // CallSetting
	PropNDFilterPositionSetting                             Prop = 0xE0E5 // NDFilterPositionSetting
	PropSceneFileCommandVersion                             Prop = 0xE0E6 // SceneFileCommandVersion
	PropSceneFileUploadOperationEnableStatus                Prop = 0xE101 // SceneFileUploadOperationEnableStatus
	PropSceneFileDownloadOperationEnableStatus              Prop = 0xE102 // SceneFileDownloadOperationEnableStatus
	PropSceneFileIndexesAvailableForDownload                Prop = 0xE0F9 // SceneFileIndexesAvailableForDownload
	PropEframingType                                        Prop = 0xD122 // EframingType
	PropEframingCommandVersion                              Prop = 0xD123 // EframingCommandVersion
	PropEframingAutoFraming                                 Prop = 0xE0F8 // EframingAutoFraming
	PropEframingTrackingStartMode                           Prop = 0xE0F0 // EframingTrackingStartMode
	PropEframingProductionEffect                            Prop = 0xE0F1 // EframingProductionEffect
	PropEframingSpeedPTZ                                    Prop = 0xD127 // EframingSpeedPTZ
	PropPriv0601                                            Prop = 0xD009 // private
	PropTopOfTheGroupShootingMarkSetting                    Prop = 0xD18C // TopOfTheGroupShootingMarkSetting
	PropPriv0603                                            Prop = 0xD18D // private
	PropCompRAWShootingNR                                   Prop = 0xD148 // CompRAWShootingNR
	PropCompRAWShootingNRFileCompressionType                Prop = 0xD149 // CompRAWShootingNRFileCompressionType
	PropCompRAWShootingNRNumberOfSheets                     Prop = 0xD15A // CompRAWShootingNRNumberOfSheets
	PropElapsedBulbExposureTime                             Prop = 0xD2A6 // ElapsedBulbExposureTime
	PropRemainingBulbExposureTime                           Prop = 0xD2A7 // RemainingBulbExposureTime
	PropRemainingNoiseReductionTime                         Prop = 0xD2A8 // RemainingNoiseReductionTime
	PropPriv0F01                                            Prop = 0xD261 // private
	PropPriv0F02                                            Prop = 0xD22A // private
	PropDigitalExtenderMagnificationSetting                 Prop = 0xE100 // DigitalExtenderMagnificationSetting
	PropMovieRecReviewPlayingState                          Prop = 0xE07A // MovieRecReviewPlayingState
	PropNearFar                                             Prop = 0xD2D1 // NearFar
	PropReserved7                                           Prop = 0xD2D2 // reserved7
	PropAFAreaPosition                                      Prop = 0xD2DC // AF_Area_Position
	PropZoomOperation                                       Prop = 0xD2DD // Zoom_Operation
	PropZoomAndFocusPositionSave                            Prop = 0xD2E9 // ZoomAndFocusPosition_Save
	PropZoomAndFocusPositionLoad                            Prop = 0xD2EA // ZoomAndFocusPosition_Load
	PropColortempStep                                       Prop = 0xD2EC // ColortempStep
	PropWhiteBalanceTintStep                                Prop = 0xD2ED // WhiteBalanceTintStep
	PropFocusOperation                                      Prop = 0xD2EF // Focus_Operation
	PropShutterECSNumberStep                                Prop = 0xF000 // ShutterECSNumberStep
	PropRemoteTouchOperation                                Prop = 0xD2E4 // RemoteTouchOperation
	PropZoomAndFocusPresetZoomOnlySet                       Prop = 0xD2F2 // ZoomAndFocusPresetZoomOnly_Set
	PropCustomWBCaptureStandby                              Prop = 0xD2DF // CustomWB_Capture_Standby
	PropCustomWBCaptureStandbyCancel                        Prop = 0xD2E0 // CustomWB_Capture_Standby_Cancel
	PropCustomWBCapture                                     Prop = 0xD2E1 // CustomWB_Capture
	PropZoomOperationWithInt16                              Prop = 0xF003 // ZoomOperationWithInt16
	PropFocusOperationWithInt16                             Prop = 0xF004 // FocusOperationWithInt16
	PropHighResolutionShutterSpeedAdjust                    Prop = 0xD2E3 // HighResolutionShutterSpeedAdjust
	PropHighResolutionShutterSpeedAdjustInIntegralMultiples Prop = 0xD2F0 // HighResolutionShutterSpeedAdjustInIntegralMultiples
	PropMovieAngleOfViewPriority                            Prop = 0xE0F6 // Movie_AngleOfViewPriority
	PropWindNoiseReductForExternalMic                       Prop = 0xE0D6 // WindNoiseReductForExternalMic
	PropNoiseCutFilter                                      Prop = 0xE0D7 // NoiseCutFilter
	PropNoiseCutFilterForExternalMic                        Prop = 0xE0D8 // NoiseCutFilterForExternalMic
	PropDispModeCandidateStill                              Prop = 0xE107 // DispModeCandidateStill
	PropDispModeSettingStill                                Prop = 0xE108 // DispModeSettingStill
	PropDispModeStill                                       Prop = 0xE109 // DispModeStill
	PropDispModeCandidateMovie                              Prop = 0xE10A // DispModeCandidateMovie
	PropDispModeSettingMovie                                Prop = 0xE10B // DispModeSettingMovie
	PropDispModeMovie                                       Prop = 0xE10C // DispModeMovie
	PropCompRAWShootingHDR                                  Prop = 0xE0DB // CompRAWShootingHDR
	PropCompRAWShootingHDRDRSetting                         Prop = 0xE0DC // CompRAWShootingHDRDRSetting
	PropCompRAWShootingHDRFileCompressionType               Prop = 0xE0DD // CompRAWShootingHDRFileCompressionType
	PropCompRAWShootingHDRNumberOfSheets                    Prop = 0xE0DE // CompRAWShootingHDRNumberOfSheets
	PropControlGeneralSettingFileEnableStatus               Prop = 0xE081 // ControlGeneralSettingFileEnableStatus
	PropPeakingDisplay                                      Prop = 0xE11A // PeakingDisplay
	PropPeakingLevel                                        Prop = 0xE11B // PeakingLevel
	PropPeakingColor                                        Prop = 0xE11C // PeakingColor
	PropZebraDisplay                                        Prop = 0xE11D // ZebraDisplay
	PropZebraLevel                                          Prop = 0xE11E // ZebraLevel
	PropZebraLevelTypeCustom                                Prop = 0xE11F // ZebraLevelTypeCustom
	PropZebraLevelStandardCustom                            Prop = 0xE120 // ZebraLevelStandardCustom
	PropZebraLevelRangeCustom                               Prop = 0xE121 // ZebraLevelRangeCustom
	PropZebraLevelLowerLimitCustom                          Prop = 0xE122 // ZebraLevelLowerLimitCustom
	PropMarkerDisplay                                       Prop = 0xE123 // MarkerDisplay
	PropCenterMarkerDisplay                                 Prop = 0xE124 // CenterMarkerDisplay
	PropAspectMarkerRatioMovie                              Prop = 0xE125 // AspectMarkerRatioMovie
	PropSafetyZoneDisplay                                   Prop = 0xE126 // SafetyZoneDisplay
	PropGuideframeDisplay                                   Prop = 0xE127 // GuideframeDisplay
	PropDualGain                                            Prop = 0xE12A // DualGain
	PropImagerMode                                          Prop = 0xE131 // ImagerMode
	PropDisplayQualityForFinderOnly                         Prop = 0xE133 // DisplayQualityForFinderOnly
	PropDisplayQualityForMonitorOnly                        Prop = 0xE134 // DisplayQualityForMonitorOnly
	PropPriv0F06                                            Prop = 0xD264 // private
	PropPriv0602                                            Prop = 0xD1C2 // private
	PropPriv0F03                                            Prop = 0xD1C3 // private
	PropPriv0F04                                            Prop = 0xD1C4 // private
	PropPriv0F05                                            Prop = 0xD1C5 // private
)

// PropTable is everything the SDK's static table knows.
var PropTable = map[Prop]PropInfo{
	PropFNumber:                                             {SDK: 0x0100, Name: "FNumber", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropExposureBiasCompensation:                            {SDK: 0x0101, Name: "ExposureBiasCompensation", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropFlashCompensation:                                   {SDK: 0x0102, Name: "FlashCompensation", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropShutterSpeed:                                        {SDK: 0x0103, Name: "ShutterSpeed", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropIsoSensitivity:                                      {SDK: 0x0104, Name: "IsoSensitivity", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropExposureProgramMode:                                 {SDK: 0x0105, Name: "ExposureProgramMode", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropFileType:                                            {SDK: 0x0106, Name: "FileType", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMediaSLOT1FileType:                                  {SDK: 0x012B, Name: "MediaSLOT1_FileType", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMediaSLOT2FileType:                                  {SDK: 0x012C, Name: "MediaSLOT2_FileType", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropStillImageQuality:                                   {SDK: 0x0107, Name: "StillImageQuality", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMediaSLOT1ImageQuality:                              {SDK: 0x012D, Name: "MediaSLOT1_ImageQuality", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMediaSLOT2ImageQuality:                              {SDK: 0x012E, Name: "MediaSLOT2_ImageQuality", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropWhiteBalance:                                        {SDK: 0x0108, Name: "WhiteBalance", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropFocusMode:                                           {SDK: 0x0109, Name: "FocusMode", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMeteringMode:                                        {SDK: 0x010A, Name: "MeteringMode", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropFlashMode:                                           {SDK: 0x010B, Name: "FlashMode", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropWirelessFlash:                                       {SDK: 0x010C, Name: "WirelessFlash", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropRedEyeReduction:                                     {SDK: 0x010D, Name: "RedEyeReduction", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropDriveMode:                                           {SDK: 0x010E, Name: "DriveMode", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropDRO:                                                 {SDK: 0x010F, Name: "DRO", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropImageSize:                                           {SDK: 0x0110, Name: "ImageSize", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMediaSLOT1ImageSize:                                 {SDK: 0x012F, Name: "MediaSLOT1_ImageSize", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMediaSLOT2ImageSize:                                 {SDK: 0x0130, Name: "MediaSLOT2_ImageSize", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropAspectRatio:                                         {SDK: 0x0111, Name: "AspectRatio", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropPictureEffect:                                       {SDK: 0x0112, Name: "PictureEffect", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropFocusArea:                                           {SDK: 0x0113, Name: "FocusArea", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropReserved4:                                           {SDK: 0x0114, Name: "reserved4", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropColortemp:                                           {SDK: 0x0115, Name: "Colortemp", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropColorTuningAB:                                       {SDK: 0x0116, Name: "ColorTuningAB", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropColorTuningGM:                                       {SDK: 0x0117, Name: "ColorTuningGM", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropLiveViewDisplayEffect:                               {SDK: 0x0118, Name: "LiveViewDisplayEffect", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropStillImageStoreDestination:                          {SDK: 0x0119, Name: "StillImageStoreDestination", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropPriorityKeySettings:                                 {SDK: 0x011A, Name: "PriorityKeySettings", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropAFTrackingSensitivity:                               {SDK: 0x011B, Name: "AFTrackingSensitivity", Min: 1, Max: 5, A7RV: false, A7RVI: true},
	PropFocusMagnifierSetting:                               {SDK: 0x011D, Name: "Focus_Magnifier_Setting", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropDateTimeSettings:                                    {SDK: 0x011E, Name: "DateTime_Settings", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropZoomScale:                                           {SDK: 0x0124, Name: "Zoom_Scale", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropZoomSetting:                                         {SDK: 0x0125, Name: "Zoom_Setting", Min: 1, Max: 4, A7RV: true, A7RVI: true},
	PropMovieFileFormat:                                     {SDK: 0x0127, Name: "Movie_File_Format", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMovieRecordingSetting:                               {SDK: 0x0128, Name: "Movie_Recording_Setting", Min: 1, Max: 69, A7RV: true, A7RVI: true},
	PropMovieRecordingFrameRateSetting:                      {SDK: 0x0129, Name: "Movie_Recording_FrameRateSetting", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropCompressionFileFormatStill:                          {SDK: 0x012A, Name: "CompressionFileFormatStill", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropRAWFileCompressionType:                              {SDK: 0x0131, Name: "RAW_FileCompressionType", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMediaSLOT1RAWFileCompressionType:                    {SDK: 0x0132, Name: "MediaSLOT1_RAW_FileCompressionType", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMediaSLOT2RAWFileCompressionType:                    {SDK: 0x0133, Name: "MediaSLOT2_RAW_FileCompressionType", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropIrisModeSetting:                                     {SDK: 0x0136, Name: "IrisModeSetting", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropShutterModeSetting:                                  {SDK: 0x0137, Name: "ShutterModeSetting", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropGainControlSetting:                                  {SDK: 0x0138, Name: "GainControlSetting", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropGainBaseIsoSensitivity:                              {SDK: 0x0139, Name: "GainBaseIsoSensitivity", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropGainBaseSensitivity:                                 {SDK: 0x013A, Name: "GainBaseSensitivity", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropExposureIndex:                                       {SDK: 0x013B, Name: "ExposureIndex", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropBaseLookValue:                                       {SDK: 0x013C, Name: "BaseLookValue", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropPlaybackMedia:                                       {SDK: 0x013D, Name: "PlaybackMedia", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropDispModeSetting:                                     {SDK: 0x013E, Name: "DispModeSetting", Min: 0, Max: 0, A7RV: true, A7RVI: false},
	PropDispMode:                                            {SDK: 0x013F, Name: "DispMode", Min: 1, Max: 7, A7RV: true, A7RVI: false},
	PropTouchOperation:                                      {SDK: 0x0140, Name: "TouchOperation", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropSelectFinderMonitor:                                 {SDK: 0x0141, Name: "SelectFinderMonitor", Min: 1, Max: 4, A7RV: true, A7RVI: true},
	PropAutoPowerOffTemperature:                             {SDK: 0x0142, Name: "AutoPowerOffTemperature", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropBodyKeyLock:                                         {SDK: 0x0143, Name: "BodyKeyLock", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropImageIDNumSetting:                                   {SDK: 0x0144, Name: "ImageID_Num_Setting", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropImageIDNum:                                          {SDK: 0x0145, Name: "ImageID_Num", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropImageIDString:                                       {SDK: 0x0146, Name: "ImageID_String", Min: 1, Max: 1, A7RV: true, A7RVI: true},
	PropExposureCtrlType:                                    {SDK: 0x0147, Name: "ExposureCtrlType", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMonitorLUTSetting:                                   {SDK: 0x0148, Name: "MonitorLUTSetting", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropFocalDistanceInMeter:                                {SDK: 0x0149, Name: "FocalDistanceInMeter", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropFocalDistanceInFeet:                                 {SDK: 0x014A, Name: "FocalDistanceInFeet", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropFocalDistanceUnitSetting:                            {SDK: 0x014B, Name: "FocalDistanceUnitSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropDigitalZoomScale:                                    {SDK: 0x014C, Name: "DigitalZoomScale", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropZoomDistance:                                        {SDK: 0x014D, Name: "ZoomDistance", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropZoomDistanceUnitSetting:                             {SDK: 0x014E, Name: "ZoomDistanceUnitSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropShutterModeStatus:                                   {SDK: 0x014F, Name: "ShutterModeStatus", Min: 1, Max: 5, A7RV: false, A7RVI: false},
	PropShutterSlow:                                         {SDK: 0x0150, Name: "ShutterSlow", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropShutterSlowFrames:                                   {SDK: 0x0151, Name: "ShutterSlowFrames", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropMovieRecordingResolutionForMain:                     {SDK: 0x0152, Name: "Movie_Recording_ResolutionForMain", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropMovieRecordingResolutionForProxy:                    {SDK: 0x0153, Name: "Movie_Recording_ResolutionForProxy", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropMovieRecordingFrameRateProxySetting:                 {SDK: 0x0154, Name: "Movie_Recording_FrameRateProxySetting", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropBatteryRemainDisplayUnit:                            {SDK: 0x0155, Name: "BatteryRemainDisplayUnit", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropPowerSource:                                         {SDK: 0x0156, Name: "PowerSource", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropMovieShootingMode:                                   {SDK: 0x0157, Name: "MovieShootingMode", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropMovieShootingModeColorGamut:                         {SDK: 0x0158, Name: "MovieShootingModeColorGamut", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropMovieShootingModeTargetDisplay:                      {SDK: 0x0159, Name: "MovieShootingModeTargetDisplay", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropDepthOfFieldAdjustmentMode:                          {SDK: 0x015A, Name: "DepthOfFieldAdjustmentMode", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropWhiteBalanceModeSetting:                             {SDK: 0x015C, Name: "WhiteBalanceModeSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropWhiteBalanceTint:                                    {SDK: 0x015D, Name: "WhiteBalanceTint", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropShutterECSSetting:                                   {SDK: 0x0160, Name: "ShutterECSSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropShutterECSNumber:                                    {SDK: 0x0161, Name: "ShutterECSNumber", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropShutterECSFrequency:                                 {SDK: 0x0163, Name: "ShutterECSFrequency", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropRecorderControlProxySetting:                         {SDK: 0x0164, Name: "RecorderControlProxySetting", Min: 0, Max: 1, A7RV: true, A7RVI: true},
	PropButtonAssignmentAssignable1:                         {SDK: 0x0165, Name: "ButtonAssignmentAssignable1", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropButtonAssignmentAssignable2:                         {SDK: 0x0166, Name: "ButtonAssignmentAssignable2", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropButtonAssignmentAssignable3:                         {SDK: 0x0167, Name: "ButtonAssignmentAssignable3", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropButtonAssignmentAssignable4:                         {SDK: 0x0168, Name: "ButtonAssignmentAssignable4", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropButtonAssignmentAssignable5:                         {SDK: 0x0169, Name: "ButtonAssignmentAssignable5", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropButtonAssignmentAssignable6:                         {SDK: 0x016A, Name: "ButtonAssignmentAssignable6", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropButtonAssignmentAssignable7:                         {SDK: 0x016B, Name: "ButtonAssignmentAssignable7", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropButtonAssignmentAssignable8:                         {SDK: 0x016C, Name: "ButtonAssignmentAssignable8", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropButtonAssignmentAssignable9:                         {SDK: 0x016D, Name: "ButtonAssignmentAssignable9", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropButtonAssignmentAssignable10:                        {SDK: 0x021F, Name: "ButtonAssignmentAssignable10", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropButtonAssignmentAssignable11:                        {SDK: 0x0220, Name: "ButtonAssignmentAssignable11", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropButtonAssignmentLensAssignable1:                     {SDK: 0x016E, Name: "ButtonAssignmentLensAssignable1", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropAssignableButton1:                                   {SDK: 0x016F, Name: "AssignableButton1", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAssignableButton2:                                   {SDK: 0x0170, Name: "AssignableButton2", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAssignableButton3:                                   {SDK: 0x0171, Name: "AssignableButton3", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAssignableButton4:                                   {SDK: 0x0172, Name: "AssignableButton4", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAssignableButton5:                                   {SDK: 0x0173, Name: "AssignableButton5", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAssignableButton6:                                   {SDK: 0x0174, Name: "AssignableButton6", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAssignableButton7:                                   {SDK: 0x0175, Name: "AssignableButton7", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAssignableButton8:                                   {SDK: 0x0176, Name: "AssignableButton8", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAssignableButton9:                                   {SDK: 0x0177, Name: "AssignableButton9", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAssignableButton10:                                  {SDK: 0x0225, Name: "AssignableButton10", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAssignableButton11:                                  {SDK: 0x0226, Name: "AssignableButton11", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropLensAssignableButton1:                               {SDK: 0x0178, Name: "LensAssignableButton1", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropFocusModeSetting:                                    {SDK: 0x0179, Name: "FocusModeSetting", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropShutterAngle:                                        {SDK: 0x017A, Name: "ShutterAngle", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropShutterSetting:                                      {SDK: 0x017B, Name: "ShutterSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropShutterMode:                                         {SDK: 0x017C, Name: "ShutterMode", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropShutterSpeedValue:                                   {SDK: 0x017D, Name: "ShutterSpeedValue", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropNDFilter:                                            {SDK: 0x017E, Name: "NDFilter", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropNDFilterModeSetting:                                 {SDK: 0x017F, Name: "NDFilterModeSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropNDFilterValue:                                       {SDK: 0x0180, Name: "NDFilterValue", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropGainUnitSetting:                                     {SDK: 0x0181, Name: "GainUnitSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropGaindBValue:                                         {SDK: 0x0182, Name: "GaindBValue", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropAWB:                                                 {SDK: 0x0183, Name: "AWB", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropSceneFileIndex:                                      {SDK: 0x0184, Name: "SceneFileIndex", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropMoviePlayButton:                                     {SDK: 0x0185, Name: "MoviePlayButton", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMoviePlayPauseButton:                                {SDK: 0x0186, Name: "MoviePlayPauseButton", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMoviePlayStopButton:                                 {SDK: 0x0187, Name: "MoviePlayStopButton", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMovieForwardButton:                                  {SDK: 0x0188, Name: "MovieForwardButton", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMovieRewindButton:                                   {SDK: 0x0189, Name: "MovieRewindButton", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMovieNextButton:                                     {SDK: 0x018A, Name: "MovieNextButton", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMoviePrevButton:                                     {SDK: 0x018B, Name: "MoviePrevButton", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMovieRecReviewButton:                                {SDK: 0x018C, Name: "MovieRecReviewButton", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropSubjectRecognitionAF:                                {SDK: 0x018D, Name: "SubjectRecognitionAF", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropAFTransitionSpeed:                                   {SDK: 0x018E, Name: "AFTransitionSpeed", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropAFSubjShiftSens:                                     {SDK: 0x018F, Name: "AFSubjShiftSens", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropAFAssist:                                            {SDK: 0x0190, Name: "AFAssist", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropNDFilterSwitchingSetting:                            {SDK: 0x0191, Name: "NDFilterSwitchingSetting", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropFunctionOfRemoteTouchOperation:                      {SDK: 0x0192, Name: "FunctionOfRemoteTouchOperation", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropFollowFocusPositionSetting:                          {SDK: 0x0194, Name: "FollowFocusPositionSetting", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropFocusBracketShotNumber:                              {SDK: 0x0195, Name: "FocusBracketShotNumber", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropFocusBracketFocusRange:                              {SDK: 0x0196, Name: "FocusBracketFocusRange", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropExtendedInterfaceMode:                               {SDK: 0x0197, Name: "ExtendedInterfaceMode", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropSQRecordingFrameRateSetting:                         {SDK: 0x0198, Name: "SQRecordingFrameRateSetting", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropSQFrameRate:                                         {SDK: 0x0199, Name: "SQFrameRate", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropSQRecordingSetting:                                  {SDK: 0x019A, Name: "SQRecordingSetting", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropAudioRecording:                                      {SDK: 0x019B, Name: "AudioRecording", Min: 0, Max: 1, A7RV: true, A7RVI: true},
	PropAudioInputMasterLevel:                               {SDK: 0x019C, Name: "AudioInputMasterLevel", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropTimeCodePreset:                                      {SDK: 0x019D, Name: "TimeCodePreset", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropTimeCodeFormat:                                      {SDK: 0x019E, Name: "TimeCodeFormat", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropTimeCodeRun:                                         {SDK: 0x019F, Name: "TimeCodeRun", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropTimeCodeMake:                                        {SDK: 0x01A0, Name: "TimeCodeMake", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropUserBitPreset:                                       {SDK: 0x01A1, Name: "UserBitPreset", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropUserBitTimeRec:                                      {SDK: 0x01A2, Name: "UserBitTimeRec", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropImageStabilizationSteadyShot:                        {SDK: 0x01A3, Name: "ImageStabilizationSteadyShot", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMovieImageStabilizationSteadyShot:                   {SDK: 0x01A4, Name: "Movie_ImageStabilizationSteadyShot", Min: 1, Max: 4, A7RV: true, A7RVI: true},
	PropSilentMode:                                          {SDK: 0x01A5, Name: "SilentMode", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropSilentModeApertureDriveInAF:                         {SDK: 0x01A6, Name: "SilentModeApertureDriveInAF", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropSilentModeShutterWhenPowerOff:                       {SDK: 0x01A7, Name: "SilentModeShutterWhenPowerOff", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropSilentModeAutoPixelMapping:                          {SDK: 0x01A8, Name: "SilentModeAutoPixelMapping", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropShutterType:                                         {SDK: 0x01A9, Name: "ShutterType", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropPictureProfile:                                      {SDK: 0x01AA, Name: "PictureProfile", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropPictureProfileBlackLevel:                            {SDK: 0x01AB, Name: "PictureProfile_BlackLevel", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropPictureProfileGamma:                                 {SDK: 0x01AC, Name: "PictureProfile_Gamma", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropPictureProfileBlackGammaRange:                       {SDK: 0x01AD, Name: "PictureProfile_BlackGammaRange", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropPictureProfileBlackGammaLevel:                       {SDK: 0x01AE, Name: "PictureProfile_BlackGammaLevel", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropPictureProfileKneeMode:                              {SDK: 0x01AF, Name: "PictureProfile_KneeMode", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropPictureProfileKneeAutoSetMaxPoint:                   {SDK: 0x01B0, Name: "PictureProfile_KneeAutoSet_MaxPoint", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropPictureProfileKneeAutoSetSensitivity:                {SDK: 0x01B1, Name: "PictureProfile_KneeAutoSet_Sensitivity", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropPictureProfileKneeManualSetPoint:                    {SDK: 0x01B2, Name: "PictureProfile_KneeManualSet_Point", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropPictureProfileKneeManualSetSlope:                    {SDK: 0x01B3, Name: "PictureProfile_KneeManualSet_Slope", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropPictureProfileColorMode:                             {SDK: 0x01B4, Name: "PictureProfile_ColorMode", Min: 1, Max: 13, A7RV: true, A7RVI: true},
	PropPictureProfileSaturation:                            {SDK: 0x01B5, Name: "PictureProfile_Saturation", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropPictureProfileColorPhase:                            {SDK: 0x01B6, Name: "PictureProfile_ColorPhase", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropPictureProfileColorDepthRed:                         {SDK: 0x01B7, Name: "PictureProfile_ColorDepthRed", Min: 0, Max: 0, A7RV: true, A7RVI: false},
	PropPictureProfileColorDepthGreen:                       {SDK: 0x01B8, Name: "PictureProfile_ColorDepthGreen", Min: 0, Max: 0, A7RV: true, A7RVI: false},
	PropPictureProfileColorDepthBlue:                        {SDK: 0x01B9, Name: "PictureProfile_ColorDepthBlue", Min: 0, Max: 0, A7RV: true, A7RVI: false},
	PropPictureProfileColorDepthCyan:                        {SDK: 0x01BA, Name: "PictureProfile_ColorDepthCyan", Min: 0, Max: 0, A7RV: true, A7RVI: false},
	PropPictureProfileColorDepthMagenta:                     {SDK: 0x01BB, Name: "PictureProfile_ColorDepthMagenta", Min: 0, Max: 0, A7RV: true, A7RVI: false},
	PropPictureProfileColorDepthYellow:                      {SDK: 0x01BC, Name: "PictureProfile_ColorDepthYellow", Min: 0, Max: 0, A7RV: true, A7RVI: false},
	PropPictureProfileDetailLevel:                           {SDK: 0x01BD, Name: "PictureProfile_DetailLevel", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropPictureProfileDetailAdjustMode:                      {SDK: 0x01BE, Name: "PictureProfile_DetailAdjustMode", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropPictureProfileDetailAdjustVHBalance:                 {SDK: 0x01BF, Name: "PictureProfile_DetailAdjustVHBalance", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropPictureProfileDetailAdjustBWBalance:                 {SDK: 0x01C0, Name: "PictureProfile_DetailAdjustBWBalance", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropPictureProfileDetailAdjustLimit:                     {SDK: 0x01C1, Name: "PictureProfile_DetailAdjustLimit", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropPictureProfileDetailAdjustCrispening:                {SDK: 0x01C2, Name: "PictureProfile_DetailAdjustCrispening", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropPictureProfileDetailAdjustHiLightDetail:             {SDK: 0x01C3, Name: "PictureProfile_DetailAdjustHiLightDetail", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropPictureProfileCopy:                                  {SDK: 0x01C4, Name: "PictureProfile_Copy", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropCreativeLook:                                        {SDK: 0x01C5, Name: "CreativeLook", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropCreativeLookContrast:                                {SDK: 0x01C6, Name: "CreativeLook_Contrast", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropCreativeLookHighlights:                              {SDK: 0x01C7, Name: "CreativeLook_Highlights", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropCreativeLookShadows:                                 {SDK: 0x01C8, Name: "CreativeLook_Shadows", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropCreativeLookFade:                                    {SDK: 0x01C9, Name: "CreativeLook_Fade", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropCreativeLookSaturation:                              {SDK: 0x01CA, Name: "CreativeLook_Saturation", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropCreativeLookSharpness:                               {SDK: 0x01CB, Name: "CreativeLook_Sharpness", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropCreativeLookSharpnessRange:                          {SDK: 0x01CC, Name: "CreativeLook_SharpnessRange", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropCreativeLookClarity:                                 {SDK: 0x01CD, Name: "CreativeLook_Clarity", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropCreativeLookCustomLook:                              {SDK: 0x01CE, Name: "CreativeLook_CustomLook", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropMovieProxyFileFormat:                                {SDK: 0x01CF, Name: "Movie_ProxyFileFormat", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropProxyRecordingSetting:                               {SDK: 0x01D0, Name: "ProxyRecordingSetting", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropFunctionOfTouchOperation:                            {SDK: 0x01D1, Name: "FunctionOfTouchOperation", Min: 1, Max: 11, A7RV: true, A7RVI: true},
	PropHighResolutionShutterSpeedSetting:                   {SDK: 0x01D2, Name: "HighResolutionShutterSpeedSetting", Min: 0, Max: 1, A7RV: false, A7RVI: true},
	PropDeleteUserBaseLook:                                  {SDK: 0x01D3, Name: "DeleteUserBaseLook", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropSelectUserBaseLookToEdit:                            {SDK: 0x01D4, Name: "SelectUserBaseLookToEdit", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropSelectUserBaseLookToSetInPPLUT:                      {SDK: 0x01D5, Name: "SelectUserBaseLookToSetInPPLUT", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropUserBaseLookInput:                                   {SDK: 0x01D6, Name: "UserBaseLookInput", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropUserBaseLookAELevelOffset:                           {SDK: 0x01D7, Name: "UserBaseLookAELevelOffset", Min: 2, Max: 2, A7RV: false, A7RVI: true},
	PropBaseISOSwitchEI:                                     {SDK: 0x01D8, Name: "BaseISOSwitchEI", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropFlickerLessShooting:                                 {SDK: 0x01D9, Name: "FlickerLessShooting", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropAudioLevelDisplay:                                   {SDK: 0x01DA, Name: "AudioLevelDisplay", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPlaybackVolumeSettings:                              {SDK: 0x01DB, Name: "PlaybackVolumeSettings", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropAutoReview:                                          {SDK: 0x01DC, Name: "AutoReview", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropAudioSignals:                                        {SDK: 0x01DD, Name: "AudioSignals", Min: 1, Max: 4, A7RV: true, A7RVI: true},
	PropHDMIResolutionStillPlay:                             {SDK: 0x01DE, Name: "HDMIResolutionStillPlay", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMovieHDMIOutputRecMedia:                             {SDK: 0x01DF, Name: "Movie_HDMIOutputRecMedia", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMovieHDMIOutputResolution:                           {SDK: 0x01E0, Name: "Movie_HDMIOutputResolution", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMovieHDMIOutput4KSetting:                            {SDK: 0x01E1, Name: "Movie_HDMIOutput4KSetting", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMovieHDMIOutputRAW:                                  {SDK: 0x01E2, Name: "Movie_HDMIOutputRAW", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMovieHDMIOutputRawSetting:                           {SDK: 0x01E3, Name: "Movie_HDMIOutputRawSetting", Min: 1, Max: 6, A7RV: false, A7RVI: true},
	PropMovieHDMIOutputColorGamutForRAWOut:                  {SDK: 0x01E4, Name: "Movie_HDMIOutputColorGamutForRAWOut", Min: 1, Max: 2, A7RV: true, A7RVI: false},
	PropMovieHDMIOutputTimeCode:                             {SDK: 0x01E5, Name: "Movie_HDMIOutputTimeCode", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMovieHDMIOutputRecControl:                           {SDK: 0x01E6, Name: "Movie_HDMIOutputRecControl", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMonitoringOutputDisplayHDMI:                         {SDK: 0x01E8, Name: "MonitoringOutputDisplayHDMI", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMovieHDMIOutputAudioCH:                              {SDK: 0x01E9, Name: "Movie_HDMIOutputAudioCH", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMovieIntervalRecIntervalTime:                        {SDK: 0x01EA, Name: "Movie_IntervalRec_IntervalTime", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropMovieIntervalRecFrameRateSetting:                    {SDK: 0x01EB, Name: "Movie_IntervalRec_FrameRateSetting", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropMovieIntervalRecRecordingSetting:                    {SDK: 0x01EC, Name: "Movie_IntervalRec_RecordingSetting", Min: 38, Max: 69, A7RV: false, A7RVI: true},
	PropEframingScaleAuto:                                   {SDK: 0x01ED, Name: "EframingScaleAuto", Min: 1, Max: 3, A7RV: false, A7RVI: true},
	PropEframingSpeedAuto:                                   {SDK: 0x01EE, Name: "EframingSpeedAuto", Min: 0, Max: 0, A7RV: false, A7RVI: true},
	PropEframingModeAuto:                                    {SDK: 0x01EF, Name: "EframingModeAuto", Min: 1, Max: 4, A7RV: false, A7RVI: true},
	PropEframingRecordingImageCrop:                          {SDK: 0x01F0, Name: "EframingRecordingImageCrop", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropEframingHDMICrop:                                    {SDK: 0x01F1, Name: "EframingHDMICrop", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropCameraEframing:                                      {SDK: 0x01F2, Name: "CameraEframing", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropUSBPowerSupply:                                      {SDK: 0x01F3, Name: "USBPowerSupply", Min: 1, Max: 5, A7RV: true, A7RVI: true},
	PropLongExposureNR:                                      {SDK: 0x01F4, Name: "LongExposureNR", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropHighIsoNR:                                           {SDK: 0x01F5, Name: "HighIsoNR", Min: 1, Max: 4, A7RV: true, A7RVI: true},
	PropHLGStillImage:                                       {SDK: 0x01F6, Name: "HLGStillImage", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropColorSpace:                                          {SDK: 0x01F7, Name: "ColorSpace", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropBracketOrder:                                        {SDK: 0x01F8, Name: "BracketOrder", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropFocusBracketOrder:                                   {SDK: 0x01F9, Name: "FocusBracketOrder", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropFocusBracketExposureLock1stImg:                      {SDK: 0x01FA, Name: "FocusBracketExposureLock1stImg", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropFocusBracketIntervalUntilNextShot:                   {SDK: 0x01FB, Name: "FocusBracketIntervalUntilNextShot", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropIntervalRecShootingStartTime:                        {SDK: 0x01FC, Name: "IntervalRec_ShootingStartTime", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropIntervalRecShootingInterval:                         {SDK: 0x01FD, Name: "IntervalRec_ShootingInterval", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropIntervalRecShootIntervalPriority:                    {SDK: 0x01FE, Name: "IntervalRec_ShootIntervalPriority", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropIntervalRecNumberOfShots:                            {SDK: 0x01FF, Name: "IntervalRec_NumberOfShots", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropIntervalRecAETrackingSensitivity:                    {SDK: 0x0200, Name: "IntervalRec_AETrackingSensitivity", Min: 1, Max: 4, A7RV: true, A7RVI: true},
	PropIntervalRecShutterType:                              {SDK: 0x0201, Name: "IntervalRec_ShutterType", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropElectricFrontCurtainShutter:                         {SDK: 0x0202, Name: "ElectricFrontCurtainShutter", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropWindNoiseReduct:                                     {SDK: 0x0203, Name: "WindNoiseReduct", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropRecordingSelfTimer:                                  {SDK: 0x0204, Name: "RecordingSelfTimer", Min: 0, Max: 1, A7RV: true, A7RVI: true},
	PropRecordingSelfTimerCountTime:                         {SDK: 0x0205, Name: "RecordingSelfTimerCountTime", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropRecordingSelfTimerContinuous:                        {SDK: 0x0206, Name: "RecordingSelfTimerContinuous", Min: 0, Max: 1, A7RV: true, A7RVI: true},
	PropRecordingSelfTimerStatus:                            {SDK: 0x0207, Name: "RecordingSelfTimerStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropBulbTimerSetting:                                    {SDK: 0x0208, Name: "BulbTimerSetting", Min: 0, Max: 1, A7RV: true, A7RVI: true},
	PropBulbExposureTimeSetting:                             {SDK: 0x0209, Name: "BulbExposureTimeSetting", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropAutoSlowShutter:                                     {SDK: 0x020A, Name: "AutoSlowShutter", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropIsoAutoMinShutterSpeedMode:                          {SDK: 0x020B, Name: "IsoAutoMinShutterSpeedMode", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropIsoAutoMinShutterSpeedManual:                        {SDK: 0x020C, Name: "IsoAutoMinShutterSpeedManual", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropIsoAutoMinShutterSpeedPreset:                        {SDK: 0x020D, Name: "IsoAutoMinShutterSpeedPreset", Min: 1, Max: 5, A7RV: true, A7RVI: true},
	PropFocusPositionSetting:                                {SDK: 0x020E, Name: "FocusPositionSetting", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSoftSkinEffect:                                      {SDK: 0x020F, Name: "SoftSkinEffect", Min: 1, Max: 4, A7RV: true, A7RVI: true},
	PropPrioritySetInAFS:                                    {SDK: 0x0210, Name: "PrioritySetInAF_S", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropPrioritySetInAFC:                                    {SDK: 0x0211, Name: "PrioritySetInAF_C", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropFocusMagnificationTime:                              {SDK: 0x0212, Name: "FocusMagnificationTime", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropSubjectRecognitionInAF:                              {SDK: 0x0213, Name: "SubjectRecognitionInAF", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropRecognitionTarget:                                   {SDK: 0x0214, Name: "RecognitionTarget", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropRightLeftEyeSelect:                                  {SDK: 0x0215, Name: "RightLeftEyeSelect", Min: 1, Max: 3, A7RV: false, A7RVI: true},
	PropSelectFTPServer:                                     {SDK: 0x0216, Name: "SelectFTPServer", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropSelectFTPServerID:                                   {SDK: 0x0217, Name: "SelectFTPServerID", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropFTPFunction:                                         {SDK: 0x0218, Name: "FTP_Function", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropFTPAutoTransfer:                                     {SDK: 0x0219, Name: "FTP_AutoTransfer", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropFTPAutoTransferTarget:                               {SDK: 0x021A, Name: "FTP_AutoTransferTarget", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropMovieFTPAutoTransferTarget:                          {SDK: 0x021B, Name: "Movie_FTP_AutoTransferTarget", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropFTPTransferTarget:                                   {SDK: 0x021C, Name: "FTP_TransferTarget", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropMovieFTPTransferTarget:                              {SDK: 0x021D, Name: "Movie_FTP_TransferTarget", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropFTPPowerSave:                                        {SDK: 0x021E, Name: "FTP_PowerSave", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropNDFilterUnitSetting:                                 {SDK: 0x022B, Name: "NDFilterUnitSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropNDFilterOpticalDensityValue:                         {SDK: 0x022C, Name: "NDFilterOpticalDensityValue", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropTNumber:                                             {SDK: 0x022D, Name: "TNumber", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropIrisDisplayUnit:                                     {SDK: 0x022E, Name: "IrisDisplayUnit", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropMovieImageStabilizationLevel:                        {SDK: 0x022F, Name: "Movie_ImageStabilizationLevel", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropImageStabilizationSteadyShotAdjust:                  {SDK: 0x0230, Name: "ImageStabilizationSteadyShotAdjust", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropImageStabilizationSteadyShotFocalLength:             {SDK: 0x0231, Name: "ImageStabilizationSteadyShotFocalLength", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropExtendedShutterSpeed:                                {SDK: 0x0232, Name: "ExtendedShutterSpeed", Min: 2, Max: 2, A7RV: false, A7RVI: true},
	PropCameraButtonFunction:                                {SDK: 0x0233, Name: "CameraButtonFunction", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropCameraButtonFunctionMulti:                           {SDK: 0x0234, Name: "CameraButtonFunctionMulti", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropCameraDialFunction:                                  {SDK: 0x0235, Name: "CameraDialFunction", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropSynchroterminalForcedOutput:                         {SDK: 0x0236, Name: "SynchroterminalForcedOutput", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropShutterReleaseTimeLagControl:                        {SDK: 0x0237, Name: "ShutterReleaseTimeLagControl", Min: 1, Max: 3, A7RV: false, A7RVI: true},
	PropContinuousShootingSpotBoostFrameSpeed:               {SDK: 0x0238, Name: "ContinuousShootingSpotBoostFrameSpeed", Min: 2, Max: 2, A7RV: false, A7RVI: true},
	PropTimeShiftShooting:                                   {SDK: 0x0239, Name: "TimeShiftShooting", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropTimeShiftTriggerSetting:                             {SDK: 0x023A, Name: "TimeShiftTriggerSetting", Min: 1, Max: 3, A7RV: false, A7RVI: true},
	PropTimeShiftPreShootingTimeSetting:                     {SDK: 0x023B, Name: "TimeShiftPreShootingTimeSetting", Min: 2, Max: 2, A7RV: false, A7RVI: true},
	PropEmbedLUTFile:                                        {SDK: 0x023C, Name: "EmbedLUTFile", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropAPSCS35:                                             {SDK: 0x023D, Name: "APS_C_S35", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropLensCompensationShading:                             {SDK: 0x023E, Name: "LensCompensationShading", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropLensCompensationChromaticAberration:                 {SDK: 0x023F, Name: "LensCompensationChromaticAberration", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropLensCompensationDistortion:                          {SDK: 0x0240, Name: "LensCompensationDistortion", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropLensCompensationBreathing:                           {SDK: 0x0241, Name: "LensCompensationBreathing", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropRecordingMedia:                                      {SDK: 0x0242, Name: "RecordingMedia", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMovieRecordingMedia:                                 {SDK: 0x0243, Name: "Movie_RecordingMedia", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropAutoSwitchMedia:                                     {SDK: 0x0244, Name: "AutoSwitchMedia", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropRecordingFileNumber:                                 {SDK: 0x0245, Name: "RecordingFileNumber", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMovieRecordingFileNumber:                            {SDK: 0x0246, Name: "Movie_RecordingFileNumber", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropRecordingSettingFileName:                            {SDK: 0x0247, Name: "RecordingSettingFileName", Min: 1, Max: 1, A7RV: true, A7RVI: true},
	PropRecordingFolderFormat:                               {SDK: 0x0248, Name: "RecordingFolderFormat", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropSelectIPTCMetadata:                                  {SDK: 0x024A, Name: "SelectIPTCMetadata", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropWriteCopyrightInfo:                                  {SDK: 0x024B, Name: "WriteCopyrightInfo", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropSetPhotographer:                                     {SDK: 0x024C, Name: "SetPhotographer", Min: 1, Max: 1, A7RV: true, A7RVI: true},
	PropSetCopyright:                                        {SDK: 0x024D, Name: "SetCopyright", Min: 1, Max: 1, A7RV: true, A7RVI: true},
	PropFileSettingsTitleNameSettings:                       {SDK: 0x024E, Name: "FileSettingsTitleNameSettings", Min: 1, Max: 1, A7RV: true, A7RVI: true},
	PropFocusBracketRecordingFolder:                         {SDK: 0x024F, Name: "FocusBracketRecordingFolder", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropReleaseWithoutLens:                                  {SDK: 0x0250, Name: "ReleaseWithoutLens", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropReleaseWithoutCard:                                  {SDK: 0x0251, Name: "ReleaseWithoutCard", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropGridLineDisplay:                                     {SDK: 0x0252, Name: "GridLineDisplay", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropContinuousShootingSpeedInElectricShutterHiPlus:      {SDK: 0x0253, Name: "ContinuousShootingSpeedInElectricShutterHiPlus", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropContinuousShootingSpeedInElectricShutterHi:          {SDK: 0x0254, Name: "ContinuousShootingSpeedInElectricShutterHi", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropContinuousShootingSpeedInElectricShutterMid:         {SDK: 0x0255, Name: "ContinuousShootingSpeedInElectricShutterMid", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropContinuousShootingSpeedInElectricShutterLo:          {SDK: 0x0256, Name: "ContinuousShootingSpeedInElectricShutterLo", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropIsoAutoRangeLimitMin:                                {SDK: 0x0257, Name: "IsoAutoRangeLimitMin", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropIsoAutoRangeLimitMax:                                {SDK: 0x0258, Name: "IsoAutoRangeLimitMax", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropFacePriorityInMultiMetering:                         {SDK: 0x0259, Name: "FacePriorityInMultiMetering", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropPrioritySetInAWB:                                    {SDK: 0x025A, Name: "PrioritySetInAWB", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropCustomWBSizeSetting:                                 {SDK: 0x025B, Name: "CustomWB_Size_Setting", Min: 1, Max: 3, A7RV: false, A7RVI: true},
	PropAFIlluminator:                                       {SDK: 0x025C, Name: "AFIlluminator", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropApertureDriveInAF:                                   {SDK: 0x025D, Name: "ApertureDriveInAF", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropAFWithShutter:                                       {SDK: 0x025E, Name: "AFWithShutter", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropFullTimeDMF:                                         {SDK: 0x025F, Name: "FullTimeDMF", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropPreAF:                                               {SDK: 0x0260, Name: "PreAF", Min: 1, Max: 2, A7RV: true, A7RVI: false},
	PropSubjectRecognitionPersonTrackingSubjectShiftRange:   {SDK: 0x0261, Name: "SubjectRecognitionPersonTrackingSubjectShiftRange", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionAnimalBirdPriority:                {SDK: 0x0262, Name: "SubjectRecognitionAnimalBirdPriority", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropSubjectRecognitionAnimalBirdDetectionParts:          {SDK: 0x0263, Name: "SubjectRecognitionAnimalBirdDetectionParts", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropSubjectRecognitionAnimalTrackingSubjectShiftRange:   {SDK: 0x0264, Name: "SubjectRecognitionAnimalTrackingSubjectShiftRange", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionAnimalTrackingSensitivity:         {SDK: 0x0265, Name: "SubjectRecognitionAnimalTrackingSensitivity", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionAnimalDetectionSensitivity:        {SDK: 0x0266, Name: "SubjectRecognitionAnimalDetectionSensitivity", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionAnimalDetectionParts:              {SDK: 0x0267, Name: "SubjectRecognitionAnimalDetectionParts", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropSubjectRecognitionBirdTrackingSubjectShiftRange:     {SDK: 0x0268, Name: "SubjectRecognitionBirdTrackingSubjectShiftRange", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionBirdTrackingSensitivity:           {SDK: 0x0269, Name: "SubjectRecognitionBirdTrackingSensitivity", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionBirdDetectionSensitivity:          {SDK: 0x026A, Name: "SubjectRecognitionBirdDetectionSensitivity", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionBirdDetectionParts:                {SDK: 0x026B, Name: "SubjectRecognitionBirdDetectionParts", Min: 1, Max: 3, A7RV: true, A7RVI: true},
	PropSubjectRecognitionInsectTrackingSubjectShiftRange:   {SDK: 0x026C, Name: "SubjectRecognitionInsectTrackingSubjectShiftRange", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionInsectTrackingSensitivity:         {SDK: 0x026D, Name: "SubjectRecognitionInsectTrackingSensitivity", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionInsectDetectionSensitivity:        {SDK: 0x026E, Name: "SubjectRecognitionInsectDetectionSensitivity", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionCarTrainTrackingSubjectShiftRange: {SDK: 0x026F, Name: "SubjectRecognitionCarTrainTrackingSubjectShiftRange", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionCarTrainTrackingSensitivity:       {SDK: 0x0270, Name: "SubjectRecognitionCarTrainTrackingSensitivity", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionCarTrainDetectionSensitivity:      {SDK: 0x0271, Name: "SubjectRecognitionCarTrainDetectionSensitivity", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionPlaneTrackingSubjectShiftRange:    {SDK: 0x0272, Name: "SubjectRecognitionPlaneTrackingSubjectShiftRange", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionPlaneTrackingSensitivity:          {SDK: 0x0273, Name: "SubjectRecognitionPlaneTrackingSensitivity", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionPlaneDetectionSensitivity:         {SDK: 0x0274, Name: "SubjectRecognitionPlaneDetectionSensitivity", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropSubjectRecognitionPriorityOnRegisteredFace:          {SDK: 0x0275, Name: "SubjectRecognitionPriorityOnRegisteredFace", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropFaceEyeFrameDisplay:                                 {SDK: 0x0276, Name: "FaceEyeFrameDisplay", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropFocusMap:                                            {SDK: 0x0277, Name: "FocusMap", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropInitialFocusMagnifier:                               {SDK: 0x0278, Name: "InitialFocusMagnifier", Min: 2, Max: 2, A7RV: true, A7RVI: true},
	PropAFInFocusMagnifier:                                  {SDK: 0x0279, Name: "AFInFocusMagnifier", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropAFTrackForSpeedChange:                               {SDK: 0x027A, Name: "AFTrackForSpeedChange", Min: 1, Max: 3, A7RV: false, A7RVI: true},
	PropAFFreeSizeAndPositionSetting:                        {SDK: 0x027B, Name: "AFFreeSizeAndPositionSetting", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropPlaySetOfMultiMedia:                                 {SDK: 0x027D, Name: "PlaySetOfMultiMedia", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropRemoteSaveImageSize:                                 {SDK: 0x027E, Name: "RemoteSaveImageSize", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropFTPTransferStillImageQualitySize:                    {SDK: 0x027F, Name: "FTP_TransferStillImageQualitySize", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropFTPAutoTransferTargetStillImage:                     {SDK: 0x0280, Name: "FTP_AutoTransferTarget_StillImage", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropProtectImageInFTPTransfer:                           {SDK: 0x0281, Name: "ProtectImageInFTPTransfer", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMonitorBrightnessType:                               {SDK: 0x0282, Name: "MonitorBrightnessType", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMonitorBrightnessManual:                             {SDK: 0x0283, Name: "MonitorBrightnessManual", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropDisplayQualityFinderMonitor:                         {SDK: 0x0284, Name: "DisplayQualityFinderMonitor", Min: 1, Max: 2, A7RV: true, A7RVI: false},
	PropTCUBDisplaySetting:                                  {SDK: 0x0285, Name: "TCUBDisplaySetting", Min: 1, Max: 4, A7RV: true, A7RVI: true},
	PropGammaDisplayAssist:                                  {SDK: 0x0286, Name: "GammaDisplayAssist", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropGammaDisplayAssistType:                              {SDK: 0x0287, Name: "GammaDisplayAssistType", Min: 1, Max: 1026, A7RV: true, A7RVI: true},
	PropAudioSignalsStartEnd:                                {SDK: 0x0288, Name: "AudioSignalsStartEnd", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropAudioSignalsVolume:                                  {SDK: 0x0289, Name: "AudioSignalsVolume", Min: 0, Max: 0, A7RV: false, A7RVI: true},
	PropControlForHDMI:                                      {SDK: 0x028A, Name: "ControlForHDMI", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropAntidustShutterWhenPowerOff:                         {SDK: 0x028B, Name: "AntidustShutterWhenPowerOff", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropWakeOnLAN:                                           {SDK: 0x028C, Name: "WakeOnLAN", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropReserved10:                                          {SDK: 0x0501, Name: "reserved10", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropReserved11:                                          {SDK: 0x0502, Name: "reserved11", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropReserved12:                                          {SDK: 0x0503, Name: "reserved12", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropIntervalRecMode:                                     {SDK: 0x0505, Name: "Interval_Rec_Mode", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropStillImageTransSize:                                 {SDK: 0x0506, Name: "Still_Image_Trans_Size", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropRAWJPCSaveImage:                                     {SDK: 0x0507, Name: "RAW_J_PC_Save_Image", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropLiveViewImageQuality:                                {SDK: 0x0508, Name: "LiveView_Image_Quality", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropRemoconZoomSpeedType:                                {SDK: 0x050C, Name: "Remocon_Zoom_Speed_Type", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropSnapshotInfo:                                        {SDK: 0x0701, Name: "SnapshotInfo", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropBatteryRemain:                                       {SDK: 0x0702, Name: "BatteryRemain", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropBatteryLevel:                                        {SDK: 0x0703, Name: "BatteryLevel", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropEstimatePictureSize:                                 {SDK: 0x0704, Name: "EstimatePictureSize", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropRecordingState:                                      {SDK: 0x0705, Name: "RecordingState", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropLiveViewStatus:                                      {SDK: 0x0706, Name: "LiveViewStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFocusIndication:                                     {SDK: 0x0707, Name: "FocusIndication", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT1Status:                                    {SDK: 0x0708, Name: "MediaSLOT1_Status", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT1RemainingNumber:                           {SDK: 0x0709, Name: "MediaSLOT1_RemainingNumber", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT1RemainingTime:                             {SDK: 0x070A, Name: "MediaSLOT1_RemainingTime", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT1FormatEnableStatus:                        {SDK: 0x070B, Name: "MediaSLOT1_FormatEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropReserved20:                                          {SDK: 0x070C, Name: "reserved20", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMediaSLOT2Status:                                    {SDK: 0x070D, Name: "MediaSLOT2_Status", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT2FormatEnableStatus:                        {SDK: 0x070E, Name: "MediaSLOT2_FormatEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT2RemainingNumber:                           {SDK: 0x070F, Name: "MediaSLOT2_RemainingNumber", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT2RemainingTime:                             {SDK: 0x0710, Name: "MediaSLOT2_RemainingTime", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropReserved22:                                          {SDK: 0x0711, Name: "reserved22", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMediaFormatProgressRate:                             {SDK: 0x0712, Name: "Media_FormatProgressRate", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFTPConnectionStatus:                                 {SDK: 0x0713, Name: "FTP_ConnectionStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFTPConnectionErrorInfo:                              {SDK: 0x0714, Name: "FTP_ConnectionErrorInfo", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropLiveViewArea:                                        {SDK: 0x0715, Name: "LiveView_Area", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropReserved26:                                          {SDK: 0x0716, Name: "reserved26", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropReserved27:                                          {SDK: 0x0717, Name: "reserved27", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropIntervalRecStatus:                                   {SDK: 0x0718, Name: "Interval_Rec_Status", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCustomWBExecutionState:                              {SDK: 0x0719, Name: "CustomWB_Execution_State", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCustomWBCapturableArea:                              {SDK: 0x071A, Name: "CustomWB_Capturable_Area", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCustomWBCaptureFrameSize:                            {SDK: 0x071B, Name: "CustomWB_Capture_Frame_Size", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCustomWBCaptureOperation:                            {SDK: 0x071C, Name: "CustomWB_Capture_Operation", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropZoomOperationStatus:                                 {SDK: 0x071E, Name: "Zoom_Operation_Status", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropZoomBarInformation:                                  {SDK: 0x071F, Name: "Zoom_Bar_Information", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropZoomTypeStatus:                                      {SDK: 0x0720, Name: "Zoom_Type_Status", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT1QuickFormatEnableStatus:                   {SDK: 0x0721, Name: "MediaSLOT1_QuickFormatEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT2QuickFormatEnableStatus:                   {SDK: 0x0722, Name: "MediaSLOT2_QuickFormatEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCancelMediaFormatEnableStatus:                       {SDK: 0x0723, Name: "Cancel_Media_FormatEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropZoomSpeedRange:                                      {SDK: 0x0724, Name: "Zoom_Speed_Range", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropIsoCurrentSensitivity:                               {SDK: 0x0729, Name: "IsoCurrentSensitivity", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCameraSettingSaveOperationEnableStatus:              {SDK: 0x072A, Name: "CameraSetting_SaveOperationEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCameraSettingReadOperationEnableStatus:              {SDK: 0x072B, Name: "CameraSetting_ReadOperationEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCameraSettingSaveReadState:                          {SDK: 0x072C, Name: "CameraSetting_SaveRead_State", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCameraSettingsResetEnableStatus:                     {SDK: 0x072D, Name: "CameraSettingsResetEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropAPSCOrFullSwitchingSetting:                          {SDK: 0x072E, Name: "APS_C_or_Full_SwitchingSetting", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropAPSCOrFullSwitchingEnableStatus:                     {SDK: 0x072F, Name: "APS_C_or_Full_SwitchingEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropDispModeCandidate:                                   {SDK: 0x0730, Name: "DispModeCandidate", Min: 4, Max: 4, A7RV: true, A7RVI: false},
	PropShutterSpeedCurrentValue:                            {SDK: 0x0731, Name: "ShutterSpeedCurrentValue", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropFocusSpeedRange:                                     {SDK: 0x0732, Name: "Focus_Speed_Range", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropNDFilterMode:                                        {SDK: 0x0733, Name: "NDFilterMode", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMoviePlayingSpeed:                                   {SDK: 0x0734, Name: "MoviePlayingSpeed", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMediaSLOT1Player:                                    {SDK: 0x0735, Name: "MediaSLOT1Player", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMediaSLOT2Player:                                    {SDK: 0x0736, Name: "MediaSLOT2Player", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropBatteryRemainingInMinutes:                           {SDK: 0x0737, Name: "BatteryRemainingInMinutes", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropBatteryRemainingInVoltage:                           {SDK: 0x0738, Name: "BatteryRemainingInVoltage", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropDCVoltage:                                           {SDK: 0x0739, Name: "DCVoltage", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMoviePlayingState:                                   {SDK: 0x073A, Name: "MoviePlayingState", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropFocusTouchSpotStatus:                                {SDK: 0x073B, Name: "FocusTouchSpotStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFocusTrackingStatus:                                 {SDK: 0x073C, Name: "FocusTrackingStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropDepthOfFieldAdjustmentInterlockingMode:              {SDK: 0x073D, Name: "DepthOfFieldAdjustmentInterlockingMode", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropRecorderClipName:                                    {SDK: 0x073E, Name: "RecorderClipName", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropRecorderControlMainSetting:                          {SDK: 0x073F, Name: "RecorderControlMainSetting", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropRecorderStartMain:                                   {SDK: 0x0740, Name: "RecorderStartMain", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropRecorderStartProxy:                                  {SDK: 0x0741, Name: "RecorderStartProxy", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropRecorderMainStatus:                                  {SDK: 0x0742, Name: "RecorderMainStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropRecorderProxyStatus:                                 {SDK: 0x0743, Name: "RecorderProxyStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropRecorderExtRawStatus:                                {SDK: 0x0744, Name: "RecorderExtRawStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropRecorderSaveDestination:                             {SDK: 0x0745, Name: "RecorderSaveDestination", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAssignableButtonIndicator1:                          {SDK: 0x0746, Name: "AssignableButtonIndicator1", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAssignableButtonIndicator2:                          {SDK: 0x0747, Name: "AssignableButtonIndicator2", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAssignableButtonIndicator3:                          {SDK: 0x0748, Name: "AssignableButtonIndicator3", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAssignableButtonIndicator4:                          {SDK: 0x0749, Name: "AssignableButtonIndicator4", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAssignableButtonIndicator5:                          {SDK: 0x074A, Name: "AssignableButtonIndicator5", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAssignableButtonIndicator6:                          {SDK: 0x074B, Name: "AssignableButtonIndicator6", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAssignableButtonIndicator7:                          {SDK: 0x074C, Name: "AssignableButtonIndicator7", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAssignableButtonIndicator8:                          {SDK: 0x074D, Name: "AssignableButtonIndicator8", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAssignableButtonIndicator9:                          {SDK: 0x074E, Name: "AssignableButtonIndicator9", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAssignableButtonIndicator10:                         {SDK: 0x077B, Name: "AssignableButtonIndicator10", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAssignableButtonIndicator11:                         {SDK: 0x077C, Name: "AssignableButtonIndicator11", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropLensAssignableButtonIndicator1:                      {SDK: 0x074F, Name: "LensAssignableButtonIndicator1", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropGaindBCurrentValue:                                  {SDK: 0x0750, Name: "GaindBCurrentValue", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropSoftwareVersion:                                     {SDK: 0x0751, Name: "SoftwareVersion", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCurrentSceneFileEdited:                              {SDK: 0x0752, Name: "CurrentSceneFileEdited", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMovieRecButtonToggleEnableStatus:                    {SDK: 0x0753, Name: "MovieRecButtonToggleEnableStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropRemoteTouchOperationEnableStatus:                    {SDK: 0x0754, Name: "RemoteTouchOperationEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCancelRemoteTouchOperationEnableStatus:              {SDK: 0x0755, Name: "CancelRemoteTouchOperationEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropLensInformationEnableStatus:                         {SDK: 0x0756, Name: "LensInformationEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFollowFocusPositionCurrentValue:                     {SDK: 0x0757, Name: "FollowFocusPositionCurrentValue", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFocusBracketShootingStatus:                          {SDK: 0x0758, Name: "FocusBracketShootingStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropPixelMappingEnableStatus:                            {SDK: 0x0759, Name: "PixelMappingEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropTimeCodePresetResetEnableStatus:                     {SDK: 0x075A, Name: "TimeCodePresetResetEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropUserBitPresetResetEnableStatus:                      {SDK: 0x075B, Name: "UserBitPresetResetEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropSensorCleaningEnableStatus:                          {SDK: 0x075C, Name: "SensorCleaningEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropPictureProfileResetEnableStatus:                     {SDK: 0x075D, Name: "PictureProfileResetEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCreativeLookResetEnableStatus:                       {SDK: 0x075E, Name: "CreativeLookResetEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropLensVersionNumber:                                   {SDK: 0x075F, Name: "LensVersionNumber", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropDeviceOverheatingState:                              {SDK: 0x0760, Name: "DeviceOverheatingState", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMovieIntervalRecCountDownIntervalTime:               {SDK: 0x0761, Name: "Movie_IntervalRec_CountDownIntervalTime", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropMovieIntervalRecRecordingDuration:                   {SDK: 0x0762, Name: "Movie_IntervalRec_RecordingDuration", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropHighResolutionShutterSpeed:                          {SDK: 0x0763, Name: "HighResolutionShutterSpeed", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropBaseLookImportOperationEnableStatus:                 {SDK: 0x0764, Name: "BaseLookImportOperationEnableStatus", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropLensModelName:                                       {SDK: 0x0765, Name: "LensModelName", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFocusPositionCurrentValue:                           {SDK: 0x0766, Name: "FocusPositionCurrentValue", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFocusDrivingStatus:                                  {SDK: 0x0767, Name: "FocusDrivingStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFlickerScanStatus:                                   {SDK: 0x0768, Name: "FlickerScanStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFlickerScanEnableStatus:                             {SDK: 0x0769, Name: "FlickerScanEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFTPServerSettingOperationEnableStatus:               {SDK: 0x076B, Name: "FTPServerSettingOperationEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFTPTransferSettingSaveOperationEnableStatus:         {SDK: 0x076C, Name: "FTPTransferSetting_SaveOperationEnableStatus", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropFTPTransferSettingReadOperationEnableStatus:         {SDK: 0x076D, Name: "FTPTransferSetting_ReadOperationEnableStatus", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropFTPTransferSettingSaveReadState:                     {SDK: 0x076E, Name: "FTPTransferSetting_SaveRead_State", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropCameraShakeStatus:                                   {SDK: 0x0770, Name: "CameraShakeStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropUpdateBodyStatus:                                    {SDK: 0x0771, Name: "UpdateBodyStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT1WritingState:                              {SDK: 0x0773, Name: "MediaSLOT1_WritingState", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT2WritingState:                              {SDK: 0x0774, Name: "MediaSLOT2_WritingState", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT1RecordingAvailableType:                    {SDK: 0x0776, Name: "MediaSLOT1_RecordingAvailableType", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMediaSLOT2RecordingAvailableType:                    {SDK: 0x0777, Name: "MediaSLOT2_RecordingAvailableType", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMediaSLOT3RecordingAvailableType:                    {SDK: 0x0778, Name: "MediaSLOT3_RecordingAvailableType", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropCameraOperatingMode:                                 {SDK: 0x0779, Name: "CameraOperatingMode", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropPlaybackViewMode:                                    {SDK: 0x077A, Name: "PlaybackViewMode", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMediaSLOT3Status:                                    {SDK: 0x0781, Name: "MediaSLOT3_Status", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMediaSLOT3RemainingTime:                             {SDK: 0x0783, Name: "MediaSLOT3_RemainingTime", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMonitoringDeliveringStatus:                          {SDK: 0x0786, Name: "MonitoringDeliveringStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMonitoringIsDelivering:                              {SDK: 0x0787, Name: "MonitoringIsDelivering", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMonitoringSettingVersion:                            {SDK: 0x0788, Name: "MonitoringSettingVersion", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMonitoringDeliveryTypeSupportInfo:                   {SDK: 0x0789, Name: "MonitoringDeliveryTypeSupportInfo", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCameraErrorCautionStatus:                            {SDK: 0x078B, Name: "CameraErrorCautionStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropSystemErrorCautionStatus:                            {SDK: 0x078C, Name: "SystemErrorCautionStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCameraButtonFunctionStatus:                          {SDK: 0x078D, Name: "CameraButtonFunctionStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropFlickerLessShootingStatus:                           {SDK: 0x078E, Name: "FlickerLessShootingStatus", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropContinuousShootingSpotBoostStatus:                   {SDK: 0x078F, Name: "ContinuousShootingSpotBoostStatus", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropContinuousShootingSpotBoostEnableStatus:             {SDK: 0x0790, Name: "ContinuousShootingSpotBoostEnableStatus", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropTimeShiftShootingStatus:                             {SDK: 0x0791, Name: "TimeShiftShootingStatus", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropZoomDrivingStatus:                                   {SDK: 0x0792, Name: "ZoomDrivingStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropShootingSelfTimerStatus:                             {SDK: 0x0793, Name: "ShootingSelfTimerStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropCreateNewFolderEnableStatus:                         {SDK: 0x0794, Name: "CreateNewFolderEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropForcedFileNumberResetEnableStatus:                   {SDK: 0x0795, Name: "ForcedFileNumberResetEnableStatus", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropDefaultAFFreeSizeAndPositionSetting:                 {SDK: 0x0796, Name: "DefaultAFFreeSizeAndPositionSetting", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropTrackingOnAndAFOnEnableStatus:                       {SDK: 0x0797, Name: "TrackingOnAndAFOnEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropProgramShiftStatus:                                  {SDK: 0x0798, Name: "ProgramShiftStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMeteredManualLevel:                                  {SDK: 0x0799, Name: "MeteredManualLevel", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropSecondBatteryRemain:                                 {SDK: 0x079B, Name: "SecondBatteryRemain", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropSecondBatteryLevel:                                  {SDK: 0x079C, Name: "SecondBatteryLevel", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropTotalBatteryRemain:                                  {SDK: 0x079D, Name: "TotalBatteryRemain", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropTotalBatteryLevel:                                   {SDK: 0x079E, Name: "TotalBatteryLevel", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropCameraLeverFunction:                                 {SDK: 0x028D, Name: "CameraLeverFunction", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropShootingTimingPreNotificationMode:                   {SDK: 0x028E, Name: "ShootingTimingPreNotificationMode", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMicrophoneDirectivity:                               {SDK: 0x028F, Name: "MicrophoneDirectivity", Min: 1, Max: 4, A7RV: false, A7RVI: false},
	PropProductShowcaseSet:                                  {SDK: 0x0290, Name: "ProductShowcaseSet", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAmountOfDefocusSetting:                              {SDK: 0x0291, Name: "AmountOfDefocusSetting", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropCinematicVlogSetting:                                {SDK: 0x0292, Name: "CinematicVlogSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropCinematicVlogLook:                                   {SDK: 0x0293, Name: "CinematicVlogLook", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropCinematicVlogMood:                                   {SDK: 0x0294, Name: "CinematicVlogMood", Min: 1, Max: 4, A7RV: false, A7RVI: false},
	PropCinematicVlogAFTransitionSpeed:                      {SDK: 0x0295, Name: "CinematicVlogAFTransitionSpeed", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropMonitoringTransportProtocol:                         {SDK: 0x07A2, Name: "MonitoringTransportProtocol", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMonitoringAvailableFormat:                           {SDK: 0x07A3, Name: "MonitoringAvailableFormat", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMonitoringFormatSupportInformation:                  {SDK: 0x07A4, Name: "MonitoringFormatSupportInformation", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropDeSqueezeDisplayRatio:                               {SDK: 0x02BB, Name: "DeSqueezeDisplayRatio", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropZoomPositionSetting:                                 {SDK: 0x02BC, Name: "ZoomPositionSetting", Min: 0, Max: 0, A7RV: true, A7RVI: true},
	PropZoomPositionCurrentValue:                            {SDK: 0x07A1, Name: "ZoomPositionCurrentValue", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropPriv0F07:                                            {SDK: 0x0F07, Name: "Priv0F07", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPriv0F08:                                            {SDK: 0x0F08, Name: "Priv0F08", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPriv0F09:                                            {SDK: 0x0F09, Name: "Priv0F09", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPriv0F0A:                                            {SDK: 0x0F0A, Name: "Priv0F0A", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPriv0F0C:                                            {SDK: 0x0F0C, Name: "Priv0F0C", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMonitoringOutputDisplaySDI:                          {SDK: 0x0296, Name: "MonitoringOutputDisplaySDI", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropCameraSystemErrorInfo:                               {SDK: 0x07AC, Name: "CameraSystemErrorInfo", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAFAreaPositionAFC:                                   {SDK: 0x0297, Name: "AFAreaPositionAF_C", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropAFAreaPositionAFS:                                   {SDK: 0x0298, Name: "AFAreaPositionAF_S", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropFaceEyeDetectionAFStatus:                            {SDK: 0x07AD, Name: "FaceEyeDetectionAFStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAutoFocusHold:                                       {SDK: 0x0299, Name: "AutoFocusHold", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPushAFModeSetting:                                   {SDK: 0x029A, Name: "PushAFModeSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropTouchFunctionInMF:                                   {SDK: 0x029B, Name: "TouchFunctionInMF", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPushAutoFocus:                                       {SDK: 0x029C, Name: "PushAutoFocus", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPushAGC:                                             {SDK: 0x029D, Name: "PushAGC", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPushAutoIris:                                        {SDK: 0x029E, Name: "PushAutoIris", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropNDFilterPreset3Value:                                {SDK: 0x02A3, Name: "NDFilterPreset3Value", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropNDFilterPreset2Value:                                {SDK: 0x02A2, Name: "NDFilterPreset2Value", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropNDFilterPreset1Value:                                {SDK: 0x02A1, Name: "NDFilterPreset1Value", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropNDFilterPresetSelect:                                {SDK: 0x02A0, Name: "NDFilterPresetSelect", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropPushAutoNDFilter:                                    {SDK: 0x029F, Name: "PushAutoNDFilter", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropWhiteBalanceOffsetColorTemp:                         {SDK: 0x02AB, Name: "WhiteBalanceOffsetColorTemp", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropWhiteBalanceOffsetSetting:                           {SDK: 0x02AA, Name: "WhiteBalanceOffsetSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropWhiteBalanceOffsetTintATW:                           {SDK: 0x02A9, Name: "WhiteBalanceOffsetTintATW", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropWhiteBalanceOffsetColorTempATW:                      {SDK: 0x02A8, Name: "WhiteBalanceOffsetColorTempATW", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropWhiteBalanceBGain:                                   {SDK: 0x02A7, Name: "WhiteBalanceBGain", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropWhiteBalanceRGain:                                   {SDK: 0x02A6, Name: "WhiteBalanceRGain", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropWhiteBalancePresetColorTemperature:                  {SDK: 0x02A5, Name: "WhiteBalancePresetColorTemperature", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropWhiteBalanceSwitch:                                  {SDK: 0x02A4, Name: "WhiteBalanceSwitch", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropPaintLookDetailLevel:                                {SDK: 0x02B4, Name: "PaintLookDetailLevel", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookDetailSetting:                              {SDK: 0x02B3, Name: "PaintLookDetailSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPaintLookKneeSlope:                                  {SDK: 0x02B2, Name: "PaintLookKneeSlope", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookKneePoint:                                  {SDK: 0x02B1, Name: "PaintLookKneePoint", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookAutoKnee:                                   {SDK: 0x02B0, Name: "PaintLookAutoKnee", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPaintLookKneeSetting:                                {SDK: 0x02AF, Name: "PaintLookKneeSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPaintLookBBlack:                                     {SDK: 0x02AE, Name: "PaintLookBBlack", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookRBlack:                                     {SDK: 0x02AD, Name: "PaintLookRBlack", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookMasterBlack:                                {SDK: 0x02AC, Name: "PaintLookMasterBlack", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropUploadDatasetVersion:                                {SDK: 0x07AE, Name: "UploadDatasetVersion", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropUserBaseLookOutput:                                  {SDK: 0x02B5, Name: "UserBaseLookOutput", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropMonitorLUTSettingOutputDestAssign:                   {SDK: 0x07AF, Name: "MonitorLUTSettingOutputDestAssign", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMonitorLUTSetting1:                                  {SDK: 0x02B6, Name: "MonitorLUTSetting1", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMonitorLUTSetting2:                                  {SDK: 0x02B7, Name: "MonitorLUTSetting2", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMonitorLUTSetting3:                                  {SDK: 0x02B8, Name: "MonitorLUTSetting3", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMaximumNumberOfBytes:                                {SDK: 0x07B0, Name: "MaximumNumberOfBytes", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropSQModeSetting:                                       {SDK: 0x02B9, Name: "SQModeSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMovieQualityFullAutoMode:                            {SDK: 0x02BA, Name: "MovieQualityFullAutoMode", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropFileSettingsCameraId:                                {SDK: 0x02BD, Name: "FileSettingsCameraId", Min: 1, Max: 26, A7RV: false, A7RVI: false},
	PropFileSettingsReelNumber:                              {SDK: 0x02BE, Name: "FileSettingsReelNumber", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropFileSettingsCameraPosition:                          {SDK: 0x02BF, Name: "FileSettingsCameraPosition", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropImageStabilizationFramingStabilizer:                 {SDK: 0x02C0, Name: "ImageStabilizationFramingStabilizer", Min: 1, Max: 3, A7RV: false, A7RVI: true},
	PropExposureStep:                                        {SDK: 0x02C1, Name: "ExposureStep", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropTeleWideLeverValueCapability:                        {SDK: 0x07A0, Name: "TeleWideLeverValueCapability", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropEnlargeScreenSetting:                                {SDK: 0x02C2, Name: "EnlargeScreenSetting", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropMediaSLOT1ContentsInfoListEnableStatus:              {SDK: 0x07A5, Name: "MediaSLOT1_ContentsInfoListEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT2ContentsInfoListEnableStatus:              {SDK: 0x07A6, Name: "MediaSLOT2_ContentsInfoListEnableStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT1ContentsInfoListRegenerateUpdateTime:      {SDK: 0x07A7, Name: "MediaSLOT1_ContentsInfoListRegenerateUpdateTime", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT2ContentsInfoListRegenerateUpdateTime:      {SDK: 0x07A8, Name: "MediaSLOT2_ContentsInfoListRegenerateUpdateTime", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT1ContentsInfoListUpdateTime:                {SDK: 0x07A9, Name: "MediaSLOT1_ContentsInfoListUpdateTime", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMediaSLOT2ContentsInfoListUpdateTime:                {SDK: 0x07AA, Name: "MediaSLOT2_ContentsInfoListUpdateTime", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropPostViewTransferResourceStatus:                      {SDK: 0x07AB, Name: "PostViewTransferResourceStatus", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropSimulRecSetting:                                     {SDK: 0x02C3, Name: "SimulRecSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropSimulRecSettingMovieRecButton:                       {SDK: 0x02C4, Name: "SimulRecSettingMovieRecButton", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropShutterSelectMode:                                   {SDK: 0x02C9, Name: "ShutterSelectMode", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropOSDImageMode:                                        {SDK: 0x02C5, Name: "OSDImageMode", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropFirmwareUpdateCommandVersion:                        {SDK: 0x07B1, Name: "FirmwareUpdateCommandVersion", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropDebugMode:                                           {SDK: 0x02C6, Name: "DebugMode", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropPriv0F0B:                                            {SDK: 0x0F0B, Name: "Priv0F0B", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropReserved18:                                          {SDK: 0x02C7, Name: "reserved18", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropReserved19:                                          {SDK: 0x02C8, Name: "reserved19", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropPriv0F0D:                                            {SDK: 0x0F0D, Name: "Priv0F0D", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropSetPresetPTZFBinaryVersion:                          {SDK: 0x07B5, Name: "SetPresetPTZFBinaryVersion", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPanPositionStatus:                                   {SDK: 0x07B6, Name: "PanPositionStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropTiltPositionStatus:                                  {SDK: 0x07B7, Name: "TiltPositionStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPanPositionCurrentValue:                             {SDK: 0x07B8, Name: "PanPositionCurrentValue", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropTiltPositionCurrentValue:                            {SDK: 0x07B9, Name: "TiltPositionCurrentValue", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPanTiltAccelerationRampCurve:                        {SDK: 0x02CB, Name: "PanTiltAccelerationRampCurve", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropPanLimitMode:                                        {SDK: 0x02CC, Name: "PanLimitMode", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPanLimitRangeMinimum:                                {SDK: 0x02CD, Name: "PanLimitRangeMinimum", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPanLimitRangeMaximum:                                {SDK: 0x02CE, Name: "PanLimitRangeMaximum", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropTiltLimitMode:                                       {SDK: 0x02CF, Name: "TiltLimitMode", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropTiltLimitRangeMinimum:                               {SDK: 0x02D0, Name: "TiltLimitRangeMinimum", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropTiltLimitRangeMaximum:                               {SDK: 0x02D1, Name: "TiltLimitRangeMaximum", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPresetPTZFSlotNumber:                                {SDK: 0x02D2, Name: "PresetPTZFSlotNumber", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropCameraPowerStatus:                                   {SDK: 0x07BA, Name: "CameraPowerStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropTargetStreamingDestinationSelect:                    {SDK: 0x0313, Name: "TargetStreamingDestinationSelect", Min: 1, Max: 16, A7RV: false, A7RVI: false},
	PropStreamStatus:                                        {SDK: 0x07CD, Name: "StreamStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropIRRemoteSetting:                                     {SDK: 0x02D3, Name: "IRRemoteSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropIPSetupProtocolSetting:                              {SDK: 0x02D4, Name: "IPSetupProtocolSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropRecordablePowerSources:                              {SDK: 0x07BB, Name: "RecordablePowerSources", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropStreamSettingListOperationStatus:                    {SDK: 0x07CF, Name: "StreamSettingListOperationStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPaintLookMultiMatrixAreaIndication:                  {SDK: 0x02E4, Name: "PaintLookMultiMatrixAreaIndication", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropIrisCloseSetting:                                    {SDK: 0x02D5, Name: "IrisCloseSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropDisplayedMenuStatus:                                 {SDK: 0x07BC, Name: "DisplayedMenuStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropLanguageSetting:                                     {SDK: 0x02D6, Name: "LanguageSetting", Min: 1, Max: 24, A7RV: false, A7RVI: false},
	PropPlaybackContentsRecordingDateTime:                   {SDK: 0x07BD, Name: "PlaybackContentsRecordingDateTime", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPlaybackContentsName:                                {SDK: 0x07BE, Name: "PlaybackContentsName", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPlaybackContentsNumber:                              {SDK: 0x07BF, Name: "PlaybackContentsNumber", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPlaybackContentsTotalNumber:                         {SDK: 0x07C0, Name: "PlaybackContentsTotalNumber", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPlaybackContentsRecordingResolution:                 {SDK: 0x07C1, Name: "PlaybackContentsRecordingResolution", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPlaybackContentsRecordingFrameRate:                  {SDK: 0x07C2, Name: "PlaybackContentsRecordingFrameRate", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPlaybackContentsRecordingFileFormat:                 {SDK: 0x07C3, Name: "PlaybackContentsRecordingFileFormat", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPlaybackContentsGammaType:                           {SDK: 0x07C4, Name: "PlaybackContentsGammaType", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropBaseLookNameofPlayback:                              {SDK: 0x07C5, Name: "BaseLookNameofPlayback", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropBaseLookAppliedofPlayback:                           {SDK: 0x07C6, Name: "BaseLookAppliedofPlayback", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPaintLookUserMatrixSetting:                          {SDK: 0x02D7, Name: "PaintLookUserMatrixSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPaintLookUserMatrixLevel:                            {SDK: 0x02D8, Name: "PaintLookUserMatrixLevel", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookUserMatrixPhase:                            {SDK: 0x02D9, Name: "PaintLookUserMatrixPhase", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookUserMatrixRG:                               {SDK: 0x02DA, Name: "PaintLookUserMatrixRG", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookUserMatrixRB:                               {SDK: 0x02DB, Name: "PaintLookUserMatrixRB", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookUserMatrixGR:                               {SDK: 0x02DC, Name: "PaintLookUserMatrixGR", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookUserMatrixGB:                               {SDK: 0x02DD, Name: "PaintLookUserMatrixGB", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookUserMatrixBR:                               {SDK: 0x02DE, Name: "PaintLookUserMatrixBR", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookUserMatrixBG:                               {SDK: 0x02DF, Name: "PaintLookUserMatrixBG", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookMultiMatrixSetting:                         {SDK: 0x02E0, Name: "PaintLookMultiMatrixSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPaintLookMultiMatrixAxis:                            {SDK: 0x02E1, Name: "PaintLookMultiMatrixAxis", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropPaintLookMultiMatrixHue:                             {SDK: 0x02E2, Name: "PaintLookMultiMatrixHue", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPaintLookMultiMatrixSaturation:                      {SDK: 0x02E3, Name: "PaintLookMultiMatrixSaturation", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropMonitoringOutputDisplaySettingDestAssign:            {SDK: 0x07C7, Name: "MonitoringOutputDisplaySettingDestAssign", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMonitoringOutputDisplaySetting1:                     {SDK: 0x02FC, Name: "MonitoringOutputDisplaySetting1", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMonitoringOutputDisplaySetting2:                     {SDK: 0x02FD, Name: "MonitoringOutputDisplaySetting2", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropFocusModeStatus:                                     {SDK: 0x07C8, Name: "FocusModeStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropFocusOperationWithInt16EnableStatus:                 {SDK: 0x07C9, Name: "FocusOperationWithInt16EnableStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAudioInputCH1LevelControl:                           {SDK: 0x02E5, Name: "AudioInputCH1LevelControl", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAudioInputCH2LevelControl:                           {SDK: 0x02E6, Name: "AudioInputCH2LevelControl", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAudioInputCH3LevelControl:                           {SDK: 0x02E7, Name: "AudioInputCH3LevelControl", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAudioInputCH4LevelControl:                           {SDK: 0x02E8, Name: "AudioInputCH4LevelControl", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAudioInputCH1Level:                                  {SDK: 0x02E9, Name: "AudioInputCH1Level", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropAudioInputCH2Level:                                  {SDK: 0x02EA, Name: "AudioInputCH2Level", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropAudioInputCH3Level:                                  {SDK: 0x02EB, Name: "AudioInputCH3Level", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropAudioInputCH4Level:                                  {SDK: 0x02EC, Name: "AudioInputCH4Level", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropAudioInputCH1InputSelect:                            {SDK: 0x02ED, Name: "AudioInputCH1InputSelect", Min: 1, Max: 8, A7RV: false, A7RVI: false},
	PropAudioInputCH2InputSelect:                            {SDK: 0x02EE, Name: "AudioInputCH2InputSelect", Min: 1, Max: 8, A7RV: false, A7RVI: false},
	PropAudioInputCH3InputSelect:                            {SDK: 0x02EF, Name: "AudioInputCH3InputSelect", Min: 1, Max: 8, A7RV: false, A7RVI: false},
	PropAudioInputCH4InputSelect:                            {SDK: 0x02F0, Name: "AudioInputCH4InputSelect", Min: 1, Max: 8, A7RV: false, A7RVI: false},
	PropAudioInputCH1WindFilter:                             {SDK: 0x02F1, Name: "AudioInputCH1WindFilter", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAudioInputCH2WindFilter:                             {SDK: 0x02F2, Name: "AudioInputCH2WindFilter", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAudioInputCH3WindFilter:                             {SDK: 0x02F3, Name: "AudioInputCH3WindFilter", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAudioInputCH4WindFilter:                             {SDK: 0x02F4, Name: "AudioInputCH4WindFilter", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropRemoteKeyThumbnailButton:                            {SDK: 0x02F7, Name: "RemoteKeyThumbnailButton", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropRemoteKeySLOTSelectButton:                           {SDK: 0x02F8, Name: "RemoteKeySLOTSelectButton", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropVideoRecordingFormatBitrateSetting:                  {SDK: 0x02F9, Name: "VideoRecordingFormatBitrateSetting", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropValidRecordingVideoFormat:                           {SDK: 0x07CA, Name: "ValidRecordingVideoFormat", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropMonitoringOutputFormat:                              {SDK: 0x02FB, Name: "MonitoringOutputFormat", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropFocusSpeedDirectSync:                                {SDK: 0x02FE, Name: "FocusSpeedDirectSync", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropAudioInput1TypeSelect:                               {SDK: 0x02F5, Name: "AudioInput1TypeSelect", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAudioInput2TypeSelect:                               {SDK: 0x02F6, Name: "AudioInput2TypeSelect", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropVideoRecordingFormatQuality:                         {SDK: 0x02FA, Name: "VideoRecordingFormatQuality", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropLiveViewImageQualityByNumericalValue:                {SDK: 0x02FF, Name: "LiveViewImageQualityByNumericalValue", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropTallyLampControlRed:                                 {SDK: 0x0300, Name: "TallyLampControlRed", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropTallyLampControlGreen:                               {SDK: 0x0301, Name: "TallyLampControlGreen", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropTallyLampControlYellow:                              {SDK: 0x0302, Name: "TallyLampControlYellow", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropMovieRecordingResolutionForRTSP:                     {SDK: 0x0305, Name: "Movie_RecordingResolutionForRTSP", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropMovieRecordingFrameRateRTSPSetting:                  {SDK: 0x0307, Name: "Movie_RecordingFrameRateRTSPSetting", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropPictureCacheRecSetting:                              {SDK: 0x0308, Name: "PictureCacheRecSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPictureCacheRecSizeAndTime:                          {SDK: 0x0309, Name: "PictureCacheRecSizeAndTime", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropMovieIntervalRecFrames:                              {SDK: 0x030A, Name: "Movie_IntervalRecFrames", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropImagerScanMode:                                      {SDK: 0x030B, Name: "ImagerScanMode", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropMovieRecordingResolutionForRAW:                      {SDK: 0x0306, Name: "Movie_RecordingResolutionForRAW", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropLensSerialNumber:                                    {SDK: 0x07CB, Name: "LensSerialNumber", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropShootingEnableSettingLicense:                        {SDK: 0x030C, Name: "ShootingEnableSettingLicense", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropGridLineDisplayPlayback:                             {SDK: 0x030D, Name: "GridLineDisplayPlayback", Min: 1, Max: 2, A7RV: true, A7RVI: true},
	PropGridLineType:                                        {SDK: 0x030E, Name: "GridLineType", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropCustomGridLineFileCommandVersion:                    {SDK: 0x07D9, Name: "CustomGridLineFileCommandVersion", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropMaximumSizeOfImageIDString:                          {SDK: 0x07CC, Name: "MaximumSizeOfImageIDString", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropStreamButtonEnableStatus:                            {SDK: 0x07CE, Name: "StreamButtonEnableStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropAutoRecognitionTargetCandidates:                     {SDK: 0x07B3, Name: "AutoRecognitionTargetCandidates", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropAutoRecognitionTargetSetting:                        {SDK: 0x02CA, Name: "AutoRecognitionTargetSetting", Min: 0, Max: 0, A7RV: false, A7RVI: true},
	PropDeleteContentOperationEnableStatusSLOT1:             {SDK: 0x07D3, Name: "DeleteContentOperationEnableStatusSLOT1", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropDeleteContentOperationEnableStatusSLOT2:             {SDK: 0x07D4, Name: "DeleteContentOperationEnableStatusSLOT2", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropDifferentSetForSQMovie:                              {SDK: 0x030F, Name: "DifferentSetForSQMovie", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropManualInputForNDFilterValue:                         {SDK: 0x0310, Name: "ManualInputForNDFilterValue", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropLogShootingMode:                                     {SDK: 0x0311, Name: "LogShootingMode", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropLogShootingModeColorGamut:                           {SDK: 0x0312, Name: "LogShootingModeColorGamut", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropVideoStreamSelect:                                   {SDK: 0x0315, Name: "VideoStreamSelect", Min: 1, Max: 16, A7RV: false, A7RVI: false},
	PropStreamDisplayName:                                   {SDK: 0x0314, Name: "StreamDisplayName", Min: 1, Max: 1, A7RV: false, A7RVI: false},
	PropVideoStreamResolution:                               {SDK: 0x0316, Name: "VideoStreamResolution", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropVideoStreamMaxBitRate:                               {SDK: 0x0317, Name: "VideoStreamMaxBitRate", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropVideoStreamAdaptiveRateControl:                      {SDK: 0x0318, Name: "VideoStreamAdaptiveRateControl", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropVideoStreamCodec:                                    {SDK: 0x0319, Name: "VideoStreamCodec", Min: 1, Max: 4, A7RV: false, A7RVI: false},
	PropStreamLatency:                                       {SDK: 0x031A, Name: "StreamLatency", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropStreamTTL:                                           {SDK: 0x031B, Name: "StreamTTL", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropStreamCipherType:                                    {SDK: 0x031C, Name: "StreamCipherType", Min: 1, Max: 4, A7RV: false, A7RVI: false},
	PropStreamModeSetting:                                   {SDK: 0x031D, Name: "StreamModeSetting", Min: 1, Max: 4, A7RV: false, A7RVI: false},
	PropVideoStreamResolutionMethod:                         {SDK: 0x031E, Name: "VideoStreamResolutionMethod", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropVideoStreamMovieRecPermission:                       {SDK: 0x031F, Name: "VideoStreamMovieRecPermission", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropVideoStreamBitRateCompressionMode:                   {SDK: 0x0320, Name: "VideoStreamBitRateCompressionMode", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropVideoStreamBitRateVBRMode:                           {SDK: 0x0321, Name: "VideoStreamBitRateVBRMode", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAudioStreamCodecType:                                {SDK: 0x0322, Name: "AudioStreamCodecType", Min: 32, Max: 32, A7RV: false, A7RVI: false},
	PropAudioStreamSamplingFrequency:                        {SDK: 0x0323, Name: "AudioStreamSamplingFrequency", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropAudioStreamBitDepth:                                 {SDK: 0x0324, Name: "AudioStreamBitDepth", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropAudioStreamChannel:                                  {SDK: 0x0325, Name: "AudioStreamChannel", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropHomeMenuSetting:                                     {SDK: 0x0326, Name: "HomeMenuSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropCallSetting:                                         {SDK: 0x0327, Name: "CallSetting", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropNDFilterPositionSetting:                             {SDK: 0x0328, Name: "NDFilterPositionSetting", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropSceneFileCommandVersion:                             {SDK: 0x07D5, Name: "SceneFileCommandVersion", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropSceneFileUploadOperationEnableStatus:                {SDK: 0x07D6, Name: "SceneFileUploadOperationEnableStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropSceneFileDownloadOperationEnableStatus:              {SDK: 0x07D7, Name: "SceneFileDownloadOperationEnableStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropSceneFileIndexesAvailableForDownload:                {SDK: 0x07D8, Name: "SceneFileIndexesAvailableForDownload", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropEframingType:                                        {SDK: 0x07D1, Name: "EframingType", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropEframingCommandVersion:                              {SDK: 0x07D2, Name: "EframingCommandVersion", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropEframingAutoFraming:                                 {SDK: 0x0329, Name: "EframingAutoFraming", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropEframingTrackingStartMode:                           {SDK: 0x032A, Name: "EframingTrackingStartMode", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropEframingProductionEffect:                            {SDK: 0x032B, Name: "EframingProductionEffect", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropEframingSpeedPTZ:                                    {SDK: 0x032C, Name: "EframingSpeedPTZ", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropPriv0601:                                            {SDK: 0x0601, Name: "Priv0601", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropTopOfTheGroupShootingMarkSetting:                    {SDK: 0x0331, Name: "TopOfTheGroupShootingMarkSetting", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropPriv0603:                                            {SDK: 0x0603, Name: "Priv0603", Min: 0, Max: 5, A7RV: false, A7RVI: false},
	PropCompRAWShootingNR:                                   {SDK: 0x0332, Name: "CompRAWShootingNR", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropCompRAWShootingNRFileCompressionType:                {SDK: 0x0333, Name: "CompRAWShootingNRFileCompressionType", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropCompRAWShootingNRNumberOfSheets:                     {SDK: 0x0334, Name: "CompRAWShootingNRNumberOfSheets", Min: 2, Max: 2, A7RV: false, A7RVI: true},
	PropElapsedBulbExposureTime:                             {SDK: 0x07DA, Name: "ElapsedBulbExposureTime", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropRemainingBulbExposureTime:                           {SDK: 0x07DB, Name: "RemainingBulbExposureTime", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropRemainingNoiseReductionTime:                         {SDK: 0x07DC, Name: "RemainingNoiseReductionTime", Min: 4, Max: 4, A7RV: true, A7RVI: true},
	PropPriv0F01:                                            {SDK: 0x0F01, Name: "Priv0F01", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPriv0F02:                                            {SDK: 0x0F02, Name: "Priv0F02", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropDigitalExtenderMagnificationSetting:                 {SDK: 0x032D, Name: "DigitalExtenderMagnificationSetting", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropMovieRecReviewPlayingState:                          {SDK: 0x032E, Name: "MovieRecReviewPlayingState", Min: 0, Max: 1, A7RV: false, A7RVI: false},
	PropNearFar:                                             {SDK: 0x011F, Name: "NearFar", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropReserved7:                                           {SDK: 0x0120, Name: "reserved7", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropAFAreaPosition:                                      {SDK: 0x0121, Name: "AF_Area_Position", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropZoomOperation:                                       {SDK: 0x0126, Name: "Zoom_Operation", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropZoomAndFocusPositionSave:                            {SDK: 0x0134, Name: "ZoomAndFocusPosition_Save", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropZoomAndFocusPositionLoad:                            {SDK: 0x0135, Name: "ZoomAndFocusPosition_Load", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropColortempStep:                                       {SDK: 0x015B, Name: "ColortempStep", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropWhiteBalanceTintStep:                                {SDK: 0x015E, Name: "WhiteBalanceTintStep", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropFocusOperation:                                      {SDK: 0x015F, Name: "Focus_Operation", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropShutterECSNumberStep:                                {SDK: 0x0162, Name: "ShutterECSNumberStep", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropRemoteTouchOperation:                                {SDK: 0x0193, Name: "RemoteTouchOperation", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropZoomAndFocusPresetZoomOnlySet:                       {SDK: 0x027C, Name: "ZoomAndFocusPresetZoomOnly_Set", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropCustomWBCaptureStandby:                              {SDK: 0x0509, Name: "CustomWB_Capture_Standby", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropCustomWBCaptureStandbyCancel:                        {SDK: 0x050A, Name: "CustomWB_Capture_Standby_Cancel", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropCustomWBCapture:                                     {SDK: 0x050B, Name: "CustomWB_Capture", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropZoomOperationWithInt16:                              {SDK: 0x0303, Name: "ZoomOperationWithInt16", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropFocusOperationWithInt16:                             {SDK: 0x0304, Name: "FocusOperationWithInt16", Min: 3, Max: 3, A7RV: false, A7RVI: false},
	PropHighResolutionShutterSpeedAdjust:                    {SDK: 0x032F, Name: "HighResolutionShutterSpeedAdjust", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropHighResolutionShutterSpeedAdjustInIntegralMultiples: {SDK: 0x0330, Name: "HighResolutionShutterSpeedAdjustInIntegralMultiples", Min: 3, Max: 3, A7RV: true, A7RVI: true},
	PropMovieAngleOfViewPriority:                            {SDK: 0x034B, Name: "Movie_AngleOfViewPriority", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropWindNoiseReductForExternalMic:                       {SDK: 0x034C, Name: "WindNoiseReductForExternalMic", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropNoiseCutFilter:                                      {SDK: 0x034D, Name: "NoiseCutFilter", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropNoiseCutFilterForExternalMic:                        {SDK: 0x034E, Name: "NoiseCutFilterForExternalMic", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropDispModeCandidateStill:                              {SDK: 0x07DD, Name: "DispModeCandidateStill", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropDispModeSettingStill:                                {SDK: 0x0339, Name: "DispModeSettingStill", Min: 0, Max: 0, A7RV: false, A7RVI: true},
	PropDispModeStill:                                       {SDK: 0x033A, Name: "DispModeStill", Min: 1, Max: 8, A7RV: false, A7RVI: true},
	PropDispModeCandidateMovie:                              {SDK: 0x07DE, Name: "DispModeCandidateMovie", Min: 4, Max: 4, A7RV: false, A7RVI: true},
	PropDispModeSettingMovie:                                {SDK: 0x033B, Name: "DispModeSettingMovie", Min: 0, Max: 0, A7RV: false, A7RVI: true},
	PropDispModeMovie:                                       {SDK: 0x033C, Name: "DispModeMovie", Min: 1, Max: 8, A7RV: false, A7RVI: true},
	PropCompRAWShootingHDR:                                  {SDK: 0x0335, Name: "CompRAWShootingHDR", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropCompRAWShootingHDRDRSetting:                         {SDK: 0x0338, Name: "CompRAWShootingHDRDRSetting", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropCompRAWShootingHDRFileCompressionType:               {SDK: 0x0336, Name: "CompRAWShootingHDRFileCompressionType", Min: 3, Max: 3, A7RV: false, A7RVI: true},
	PropCompRAWShootingHDRNumberOfSheets:                    {SDK: 0x0337, Name: "CompRAWShootingHDRNumberOfSheets", Min: 2, Max: 2, A7RV: false, A7RVI: true},
	PropControlGeneralSettingFileEnableStatus:               {SDK: 0x07DF, Name: "ControlGeneralSettingFileEnableStatus", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPeakingDisplay:                                      {SDK: 0x033D, Name: "PeakingDisplay", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropPeakingLevel:                                        {SDK: 0x033E, Name: "PeakingLevel", Min: 1, Max: 3, A7RV: false, A7RVI: false},
	PropPeakingColor:                                        {SDK: 0x033F, Name: "PeakingColor", Min: 1, Max: 4, A7RV: false, A7RVI: false},
	PropZebraDisplay:                                        {SDK: 0x0340, Name: "ZebraDisplay", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropZebraLevel:                                          {SDK: 0x0341, Name: "ZebraLevel", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropZebraLevelTypeCustom:                                {SDK: 0x0342, Name: "ZebraLevelTypeCustom", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropZebraLevelStandardCustom:                            {SDK: 0x0343, Name: "ZebraLevelStandardCustom", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropZebraLevelRangeCustom:                               {SDK: 0x0344, Name: "ZebraLevelRangeCustom", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropZebraLevelLowerLimitCustom:                          {SDK: 0x0345, Name: "ZebraLevelLowerLimitCustom", Min: 0, Max: 0, A7RV: false, A7RVI: false},
	PropMarkerDisplay:                                       {SDK: 0x0346, Name: "MarkerDisplay", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropCenterMarkerDisplay:                                 {SDK: 0x0347, Name: "CenterMarkerDisplay", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropAspectMarkerRatioMovie:                              {SDK: 0x0348, Name: "AspectMarkerRatioMovie", Min: 1, Max: 16, A7RV: false, A7RVI: false},
	PropSafetyZoneDisplay:                                   {SDK: 0x0349, Name: "SafetyZoneDisplay", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropGuideframeDisplay:                                   {SDK: 0x034A, Name: "GuideframeDisplay", Min: 1, Max: 2, A7RV: false, A7RVI: false},
	PropDualGain:                                            {SDK: 0x034F, Name: "DualGain", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropImagerMode:                                          {SDK: 0x0350, Name: "ImagerMode", Min: 2, Max: 2, A7RV: false, A7RVI: false},
	PropDisplayQualityForFinderOnly:                         {SDK: 0x0351, Name: "DisplayQualityForFinderOnly", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropDisplayQualityForMonitorOnly:                        {SDK: 0x0352, Name: "DisplayQualityForMonitorOnly", Min: 1, Max: 2, A7RV: false, A7RVI: true},
	PropPriv0F06:                                            {SDK: 0x0F06, Name: "Priv0F06", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPriv0602:                                            {SDK: 0x0602, Name: "Priv0602", Min: 0, Max: 2, A7RV: false, A7RVI: false},
	PropPriv0F03:                                            {SDK: 0x0F03, Name: "Priv0F03", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPriv0F04:                                            {SDK: 0x0F04, Name: "Priv0F04", Min: 4, Max: 4, A7RV: false, A7RVI: false},
	PropPriv0F05:                                            {SDK: 0x0F05, Name: "Priv0F05", Min: 4, Max: 4, A7RV: false, A7RVI: false},
}

// sdkToWire maps the SDK's internal codes to wire codes.
var sdkToWire = map[SDKProp]Prop{
	0x0100: PropFNumber,
	0x0101: PropExposureBiasCompensation,
	0x0102: PropFlashCompensation,
	0x0103: PropShutterSpeed,
	0x0104: PropIsoSensitivity,
	0x0105: PropExposureProgramMode,
	0x0106: PropFileType,
	0x012B: PropMediaSLOT1FileType,
	0x012C: PropMediaSLOT2FileType,
	0x0107: PropStillImageQuality,
	0x012D: PropMediaSLOT1ImageQuality,
	0x012E: PropMediaSLOT2ImageQuality,
	0x0108: PropWhiteBalance,
	0x0109: PropFocusMode,
	0x010A: PropMeteringMode,
	0x010B: PropFlashMode,
	0x010C: PropWirelessFlash,
	0x010D: PropRedEyeReduction,
	0x010E: PropDriveMode,
	0x010F: PropDRO,
	0x0110: PropImageSize,
	0x012F: PropMediaSLOT1ImageSize,
	0x0130: PropMediaSLOT2ImageSize,
	0x0111: PropAspectRatio,
	0x0112: PropPictureEffect,
	0x0113: PropFocusArea,
	0x0114: PropReserved4,
	0x0115: PropColortemp,
	0x0116: PropColorTuningAB,
	0x0117: PropColorTuningGM,
	0x0118: PropLiveViewDisplayEffect,
	0x0119: PropStillImageStoreDestination,
	0x011A: PropPriorityKeySettings,
	0x011B: PropAFTrackingSensitivity,
	0x011D: PropFocusMagnifierSetting,
	0x011E: PropDateTimeSettings,
	0x0124: PropZoomScale,
	0x0125: PropZoomSetting,
	0x0127: PropMovieFileFormat,
	0x0128: PropMovieRecordingSetting,
	0x0129: PropMovieRecordingFrameRateSetting,
	0x012A: PropCompressionFileFormatStill,
	0x0131: PropRAWFileCompressionType,
	0x0132: PropMediaSLOT1RAWFileCompressionType,
	0x0133: PropMediaSLOT2RAWFileCompressionType,
	0x0136: PropIrisModeSetting,
	0x0137: PropShutterModeSetting,
	0x0138: PropGainControlSetting,
	0x0139: PropGainBaseIsoSensitivity,
	0x013A: PropGainBaseSensitivity,
	0x013B: PropExposureIndex,
	0x013C: PropBaseLookValue,
	0x013D: PropPlaybackMedia,
	0x013E: PropDispModeSetting,
	0x013F: PropDispMode,
	0x0140: PropTouchOperation,
	0x0141: PropSelectFinderMonitor,
	0x0142: PropAutoPowerOffTemperature,
	0x0143: PropBodyKeyLock,
	0x0144: PropImageIDNumSetting,
	0x0145: PropImageIDNum,
	0x0146: PropImageIDString,
	0x0147: PropExposureCtrlType,
	0x0148: PropMonitorLUTSetting,
	0x0149: PropFocalDistanceInMeter,
	0x014A: PropFocalDistanceInFeet,
	0x014B: PropFocalDistanceUnitSetting,
	0x014C: PropDigitalZoomScale,
	0x014D: PropZoomDistance,
	0x014E: PropZoomDistanceUnitSetting,
	0x014F: PropShutterModeStatus,
	0x0150: PropShutterSlow,
	0x0151: PropShutterSlowFrames,
	0x0152: PropMovieRecordingResolutionForMain,
	0x0153: PropMovieRecordingResolutionForProxy,
	0x0154: PropMovieRecordingFrameRateProxySetting,
	0x0155: PropBatteryRemainDisplayUnit,
	0x0156: PropPowerSource,
	0x0157: PropMovieShootingMode,
	0x0158: PropMovieShootingModeColorGamut,
	0x0159: PropMovieShootingModeTargetDisplay,
	0x015A: PropDepthOfFieldAdjustmentMode,
	0x015C: PropWhiteBalanceModeSetting,
	0x015D: PropWhiteBalanceTint,
	0x0160: PropShutterECSSetting,
	0x0161: PropShutterECSNumber,
	0x0163: PropShutterECSFrequency,
	0x0164: PropRecorderControlProxySetting,
	0x0165: PropButtonAssignmentAssignable1,
	0x0166: PropButtonAssignmentAssignable2,
	0x0167: PropButtonAssignmentAssignable3,
	0x0168: PropButtonAssignmentAssignable4,
	0x0169: PropButtonAssignmentAssignable5,
	0x016A: PropButtonAssignmentAssignable6,
	0x016B: PropButtonAssignmentAssignable7,
	0x016C: PropButtonAssignmentAssignable8,
	0x016D: PropButtonAssignmentAssignable9,
	0x021F: PropButtonAssignmentAssignable10,
	0x0220: PropButtonAssignmentAssignable11,
	0x016E: PropButtonAssignmentLensAssignable1,
	0x016F: PropAssignableButton1,
	0x0170: PropAssignableButton2,
	0x0171: PropAssignableButton3,
	0x0172: PropAssignableButton4,
	0x0173: PropAssignableButton5,
	0x0174: PropAssignableButton6,
	0x0175: PropAssignableButton7,
	0x0176: PropAssignableButton8,
	0x0177: PropAssignableButton9,
	0x0225: PropAssignableButton10,
	0x0226: PropAssignableButton11,
	0x0178: PropLensAssignableButton1,
	0x0179: PropFocusModeSetting,
	0x017A: PropShutterAngle,
	0x017B: PropShutterSetting,
	0x017C: PropShutterMode,
	0x017D: PropShutterSpeedValue,
	0x017E: PropNDFilter,
	0x017F: PropNDFilterModeSetting,
	0x0180: PropNDFilterValue,
	0x0181: PropGainUnitSetting,
	0x0182: PropGaindBValue,
	0x0183: PropAWB,
	0x0184: PropSceneFileIndex,
	0x0185: PropMoviePlayButton,
	0x0186: PropMoviePlayPauseButton,
	0x0187: PropMoviePlayStopButton,
	0x0188: PropMovieForwardButton,
	0x0189: PropMovieRewindButton,
	0x018A: PropMovieNextButton,
	0x018B: PropMoviePrevButton,
	0x018C: PropMovieRecReviewButton,
	0x018D: PropSubjectRecognitionAF,
	0x018E: PropAFTransitionSpeed,
	0x018F: PropAFSubjShiftSens,
	0x0190: PropAFAssist,
	0x0191: PropNDFilterSwitchingSetting,
	0x0192: PropFunctionOfRemoteTouchOperation,
	0x0194: PropFollowFocusPositionSetting,
	0x0195: PropFocusBracketShotNumber,
	0x0196: PropFocusBracketFocusRange,
	0x0197: PropExtendedInterfaceMode,
	0x0198: PropSQRecordingFrameRateSetting,
	0x0199: PropSQFrameRate,
	0x019A: PropSQRecordingSetting,
	0x019B: PropAudioRecording,
	0x019C: PropAudioInputMasterLevel,
	0x019D: PropTimeCodePreset,
	0x019E: PropTimeCodeFormat,
	0x019F: PropTimeCodeRun,
	0x01A0: PropTimeCodeMake,
	0x01A1: PropUserBitPreset,
	0x01A2: PropUserBitTimeRec,
	0x01A3: PropImageStabilizationSteadyShot,
	0x01A4: PropMovieImageStabilizationSteadyShot,
	0x01A5: PropSilentMode,
	0x01A6: PropSilentModeApertureDriveInAF,
	0x01A7: PropSilentModeShutterWhenPowerOff,
	0x01A8: PropSilentModeAutoPixelMapping,
	0x01A9: PropShutterType,
	0x01AA: PropPictureProfile,
	0x01AB: PropPictureProfileBlackLevel,
	0x01AC: PropPictureProfileGamma,
	0x01AD: PropPictureProfileBlackGammaRange,
	0x01AE: PropPictureProfileBlackGammaLevel,
	0x01AF: PropPictureProfileKneeMode,
	0x01B0: PropPictureProfileKneeAutoSetMaxPoint,
	0x01B1: PropPictureProfileKneeAutoSetSensitivity,
	0x01B2: PropPictureProfileKneeManualSetPoint,
	0x01B3: PropPictureProfileKneeManualSetSlope,
	0x01B4: PropPictureProfileColorMode,
	0x01B5: PropPictureProfileSaturation,
	0x01B6: PropPictureProfileColorPhase,
	0x01B7: PropPictureProfileColorDepthRed,
	0x01B8: PropPictureProfileColorDepthGreen,
	0x01B9: PropPictureProfileColorDepthBlue,
	0x01BA: PropPictureProfileColorDepthCyan,
	0x01BB: PropPictureProfileColorDepthMagenta,
	0x01BC: PropPictureProfileColorDepthYellow,
	0x01BD: PropPictureProfileDetailLevel,
	0x01BE: PropPictureProfileDetailAdjustMode,
	0x01BF: PropPictureProfileDetailAdjustVHBalance,
	0x01C0: PropPictureProfileDetailAdjustBWBalance,
	0x01C1: PropPictureProfileDetailAdjustLimit,
	0x01C2: PropPictureProfileDetailAdjustCrispening,
	0x01C3: PropPictureProfileDetailAdjustHiLightDetail,
	0x01C4: PropPictureProfileCopy,
	0x01C5: PropCreativeLook,
	0x01C6: PropCreativeLookContrast,
	0x01C7: PropCreativeLookHighlights,
	0x01C8: PropCreativeLookShadows,
	0x01C9: PropCreativeLookFade,
	0x01CA: PropCreativeLookSaturation,
	0x01CB: PropCreativeLookSharpness,
	0x01CC: PropCreativeLookSharpnessRange,
	0x01CD: PropCreativeLookClarity,
	0x01CE: PropCreativeLookCustomLook,
	0x01CF: PropMovieProxyFileFormat,
	0x01D0: PropProxyRecordingSetting,
	0x01D1: PropFunctionOfTouchOperation,
	0x01D2: PropHighResolutionShutterSpeedSetting,
	0x01D3: PropDeleteUserBaseLook,
	0x01D4: PropSelectUserBaseLookToEdit,
	0x01D5: PropSelectUserBaseLookToSetInPPLUT,
	0x01D6: PropUserBaseLookInput,
	0x01D7: PropUserBaseLookAELevelOffset,
	0x01D8: PropBaseISOSwitchEI,
	0x01D9: PropFlickerLessShooting,
	0x01DA: PropAudioLevelDisplay,
	0x01DB: PropPlaybackVolumeSettings,
	0x01DC: PropAutoReview,
	0x01DD: PropAudioSignals,
	0x01DE: PropHDMIResolutionStillPlay,
	0x01DF: PropMovieHDMIOutputRecMedia,
	0x01E0: PropMovieHDMIOutputResolution,
	0x01E1: PropMovieHDMIOutput4KSetting,
	0x01E2: PropMovieHDMIOutputRAW,
	0x01E3: PropMovieHDMIOutputRawSetting,
	0x01E4: PropMovieHDMIOutputColorGamutForRAWOut,
	0x01E5: PropMovieHDMIOutputTimeCode,
	0x01E6: PropMovieHDMIOutputRecControl,
	0x01E8: PropMonitoringOutputDisplayHDMI,
	0x01E9: PropMovieHDMIOutputAudioCH,
	0x01EA: PropMovieIntervalRecIntervalTime,
	0x01EB: PropMovieIntervalRecFrameRateSetting,
	0x01EC: PropMovieIntervalRecRecordingSetting,
	0x01ED: PropEframingScaleAuto,
	0x01EE: PropEframingSpeedAuto,
	0x01EF: PropEframingModeAuto,
	0x01F0: PropEframingRecordingImageCrop,
	0x01F1: PropEframingHDMICrop,
	0x01F2: PropCameraEframing,
	0x01F3: PropUSBPowerSupply,
	0x01F4: PropLongExposureNR,
	0x01F5: PropHighIsoNR,
	0x01F6: PropHLGStillImage,
	0x01F7: PropColorSpace,
	0x01F8: PropBracketOrder,
	0x01F9: PropFocusBracketOrder,
	0x01FA: PropFocusBracketExposureLock1stImg,
	0x01FB: PropFocusBracketIntervalUntilNextShot,
	0x01FC: PropIntervalRecShootingStartTime,
	0x01FD: PropIntervalRecShootingInterval,
	0x01FE: PropIntervalRecShootIntervalPriority,
	0x01FF: PropIntervalRecNumberOfShots,
	0x0200: PropIntervalRecAETrackingSensitivity,
	0x0201: PropIntervalRecShutterType,
	0x0202: PropElectricFrontCurtainShutter,
	0x0203: PropWindNoiseReduct,
	0x0204: PropRecordingSelfTimer,
	0x0205: PropRecordingSelfTimerCountTime,
	0x0206: PropRecordingSelfTimerContinuous,
	0x0207: PropRecordingSelfTimerStatus,
	0x0208: PropBulbTimerSetting,
	0x0209: PropBulbExposureTimeSetting,
	0x020A: PropAutoSlowShutter,
	0x020B: PropIsoAutoMinShutterSpeedMode,
	0x020C: PropIsoAutoMinShutterSpeedManual,
	0x020D: PropIsoAutoMinShutterSpeedPreset,
	0x020E: PropFocusPositionSetting,
	0x020F: PropSoftSkinEffect,
	0x0210: PropPrioritySetInAFS,
	0x0211: PropPrioritySetInAFC,
	0x0212: PropFocusMagnificationTime,
	0x0213: PropSubjectRecognitionInAF,
	0x0214: PropRecognitionTarget,
	0x0215: PropRightLeftEyeSelect,
	0x0216: PropSelectFTPServer,
	0x0217: PropSelectFTPServerID,
	0x0218: PropFTPFunction,
	0x0219: PropFTPAutoTransfer,
	0x021A: PropFTPAutoTransferTarget,
	0x021B: PropMovieFTPAutoTransferTarget,
	0x021C: PropFTPTransferTarget,
	0x021D: PropMovieFTPTransferTarget,
	0x021E: PropFTPPowerSave,
	0x022B: PropNDFilterUnitSetting,
	0x022C: PropNDFilterOpticalDensityValue,
	0x022D: PropTNumber,
	0x022E: PropIrisDisplayUnit,
	0x022F: PropMovieImageStabilizationLevel,
	0x0230: PropImageStabilizationSteadyShotAdjust,
	0x0231: PropImageStabilizationSteadyShotFocalLength,
	0x0232: PropExtendedShutterSpeed,
	0x0233: PropCameraButtonFunction,
	0x0234: PropCameraButtonFunctionMulti,
	0x0235: PropCameraDialFunction,
	0x0236: PropSynchroterminalForcedOutput,
	0x0237: PropShutterReleaseTimeLagControl,
	0x0238: PropContinuousShootingSpotBoostFrameSpeed,
	0x0239: PropTimeShiftShooting,
	0x023A: PropTimeShiftTriggerSetting,
	0x023B: PropTimeShiftPreShootingTimeSetting,
	0x023C: PropEmbedLUTFile,
	0x023D: PropAPSCS35,
	0x023E: PropLensCompensationShading,
	0x023F: PropLensCompensationChromaticAberration,
	0x0240: PropLensCompensationDistortion,
	0x0241: PropLensCompensationBreathing,
	0x0242: PropRecordingMedia,
	0x0243: PropMovieRecordingMedia,
	0x0244: PropAutoSwitchMedia,
	0x0245: PropRecordingFileNumber,
	0x0246: PropMovieRecordingFileNumber,
	0x0247: PropRecordingSettingFileName,
	0x0248: PropRecordingFolderFormat,
	0x024A: PropSelectIPTCMetadata,
	0x024B: PropWriteCopyrightInfo,
	0x024C: PropSetPhotographer,
	0x024D: PropSetCopyright,
	0x024E: PropFileSettingsTitleNameSettings,
	0x024F: PropFocusBracketRecordingFolder,
	0x0250: PropReleaseWithoutLens,
	0x0251: PropReleaseWithoutCard,
	0x0252: PropGridLineDisplay,
	0x0253: PropContinuousShootingSpeedInElectricShutterHiPlus,
	0x0254: PropContinuousShootingSpeedInElectricShutterHi,
	0x0255: PropContinuousShootingSpeedInElectricShutterMid,
	0x0256: PropContinuousShootingSpeedInElectricShutterLo,
	0x0257: PropIsoAutoRangeLimitMin,
	0x0258: PropIsoAutoRangeLimitMax,
	0x0259: PropFacePriorityInMultiMetering,
	0x025A: PropPrioritySetInAWB,
	0x025B: PropCustomWBSizeSetting,
	0x025C: PropAFIlluminator,
	0x025D: PropApertureDriveInAF,
	0x025E: PropAFWithShutter,
	0x025F: PropFullTimeDMF,
	0x0260: PropPreAF,
	0x0261: PropSubjectRecognitionPersonTrackingSubjectShiftRange,
	0x0262: PropSubjectRecognitionAnimalBirdPriority,
	0x0263: PropSubjectRecognitionAnimalBirdDetectionParts,
	0x0264: PropSubjectRecognitionAnimalTrackingSubjectShiftRange,
	0x0265: PropSubjectRecognitionAnimalTrackingSensitivity,
	0x0266: PropSubjectRecognitionAnimalDetectionSensitivity,
	0x0267: PropSubjectRecognitionAnimalDetectionParts,
	0x0268: PropSubjectRecognitionBirdTrackingSubjectShiftRange,
	0x0269: PropSubjectRecognitionBirdTrackingSensitivity,
	0x026A: PropSubjectRecognitionBirdDetectionSensitivity,
	0x026B: PropSubjectRecognitionBirdDetectionParts,
	0x026C: PropSubjectRecognitionInsectTrackingSubjectShiftRange,
	0x026D: PropSubjectRecognitionInsectTrackingSensitivity,
	0x026E: PropSubjectRecognitionInsectDetectionSensitivity,
	0x026F: PropSubjectRecognitionCarTrainTrackingSubjectShiftRange,
	0x0270: PropSubjectRecognitionCarTrainTrackingSensitivity,
	0x0271: PropSubjectRecognitionCarTrainDetectionSensitivity,
	0x0272: PropSubjectRecognitionPlaneTrackingSubjectShiftRange,
	0x0273: PropSubjectRecognitionPlaneTrackingSensitivity,
	0x0274: PropSubjectRecognitionPlaneDetectionSensitivity,
	0x0275: PropSubjectRecognitionPriorityOnRegisteredFace,
	0x0276: PropFaceEyeFrameDisplay,
	0x0277: PropFocusMap,
	0x0278: PropInitialFocusMagnifier,
	0x0279: PropAFInFocusMagnifier,
	0x027A: PropAFTrackForSpeedChange,
	0x027B: PropAFFreeSizeAndPositionSetting,
	0x027D: PropPlaySetOfMultiMedia,
	0x027E: PropRemoteSaveImageSize,
	0x027F: PropFTPTransferStillImageQualitySize,
	0x0280: PropFTPAutoTransferTargetStillImage,
	0x0281: PropProtectImageInFTPTransfer,
	0x0282: PropMonitorBrightnessType,
	0x0283: PropMonitorBrightnessManual,
	0x0284: PropDisplayQualityFinderMonitor,
	0x0285: PropTCUBDisplaySetting,
	0x0286: PropGammaDisplayAssist,
	0x0287: PropGammaDisplayAssistType,
	0x0288: PropAudioSignalsStartEnd,
	0x0289: PropAudioSignalsVolume,
	0x028A: PropControlForHDMI,
	0x028B: PropAntidustShutterWhenPowerOff,
	0x028C: PropWakeOnLAN,
	0x0501: PropReserved10,
	0x0502: PropReserved11,
	0x0503: PropReserved12,
	0x0505: PropIntervalRecMode,
	0x0506: PropStillImageTransSize,
	0x0507: PropRAWJPCSaveImage,
	0x0508: PropLiveViewImageQuality,
	0x050C: PropRemoconZoomSpeedType,
	0x0701: PropSnapshotInfo,
	0x0702: PropBatteryRemain,
	0x0703: PropBatteryLevel,
	0x0704: PropEstimatePictureSize,
	0x0705: PropRecordingState,
	0x0706: PropLiveViewStatus,
	0x0707: PropFocusIndication,
	0x0708: PropMediaSLOT1Status,
	0x0709: PropMediaSLOT1RemainingNumber,
	0x070A: PropMediaSLOT1RemainingTime,
	0x070B: PropMediaSLOT1FormatEnableStatus,
	0x070C: PropReserved20,
	0x070D: PropMediaSLOT2Status,
	0x070E: PropMediaSLOT2FormatEnableStatus,
	0x070F: PropMediaSLOT2RemainingNumber,
	0x0710: PropMediaSLOT2RemainingTime,
	0x0711: PropReserved22,
	0x0712: PropMediaFormatProgressRate,
	0x0713: PropFTPConnectionStatus,
	0x0714: PropFTPConnectionErrorInfo,
	0x0715: PropLiveViewArea,
	0x0716: PropReserved26,
	0x0717: PropReserved27,
	0x0718: PropIntervalRecStatus,
	0x0719: PropCustomWBExecutionState,
	0x071A: PropCustomWBCapturableArea,
	0x071B: PropCustomWBCaptureFrameSize,
	0x071C: PropCustomWBCaptureOperation,
	0x071E: PropZoomOperationStatus,
	0x071F: PropZoomBarInformation,
	0x0720: PropZoomTypeStatus,
	0x0721: PropMediaSLOT1QuickFormatEnableStatus,
	0x0722: PropMediaSLOT2QuickFormatEnableStatus,
	0x0723: PropCancelMediaFormatEnableStatus,
	0x0724: PropZoomSpeedRange,
	0x0729: PropIsoCurrentSensitivity,
	0x072A: PropCameraSettingSaveOperationEnableStatus,
	0x072B: PropCameraSettingReadOperationEnableStatus,
	0x072C: PropCameraSettingSaveReadState,
	0x072D: PropCameraSettingsResetEnableStatus,
	0x072E: PropAPSCOrFullSwitchingSetting,
	0x072F: PropAPSCOrFullSwitchingEnableStatus,
	0x0730: PropDispModeCandidate,
	0x0731: PropShutterSpeedCurrentValue,
	0x0732: PropFocusSpeedRange,
	0x0733: PropNDFilterMode,
	0x0734: PropMoviePlayingSpeed,
	0x0735: PropMediaSLOT1Player,
	0x0736: PropMediaSLOT2Player,
	0x0737: PropBatteryRemainingInMinutes,
	0x0738: PropBatteryRemainingInVoltage,
	0x0739: PropDCVoltage,
	0x073A: PropMoviePlayingState,
	0x073B: PropFocusTouchSpotStatus,
	0x073C: PropFocusTrackingStatus,
	0x073D: PropDepthOfFieldAdjustmentInterlockingMode,
	0x073E: PropRecorderClipName,
	0x073F: PropRecorderControlMainSetting,
	0x0740: PropRecorderStartMain,
	0x0741: PropRecorderStartProxy,
	0x0742: PropRecorderMainStatus,
	0x0743: PropRecorderProxyStatus,
	0x0744: PropRecorderExtRawStatus,
	0x0745: PropRecorderSaveDestination,
	0x0746: PropAssignableButtonIndicator1,
	0x0747: PropAssignableButtonIndicator2,
	0x0748: PropAssignableButtonIndicator3,
	0x0749: PropAssignableButtonIndicator4,
	0x074A: PropAssignableButtonIndicator5,
	0x074B: PropAssignableButtonIndicator6,
	0x074C: PropAssignableButtonIndicator7,
	0x074D: PropAssignableButtonIndicator8,
	0x074E: PropAssignableButtonIndicator9,
	0x077B: PropAssignableButtonIndicator10,
	0x077C: PropAssignableButtonIndicator11,
	0x074F: PropLensAssignableButtonIndicator1,
	0x0750: PropGaindBCurrentValue,
	0x0751: PropSoftwareVersion,
	0x0752: PropCurrentSceneFileEdited,
	0x0753: PropMovieRecButtonToggleEnableStatus,
	0x0754: PropRemoteTouchOperationEnableStatus,
	0x0755: PropCancelRemoteTouchOperationEnableStatus,
	0x0756: PropLensInformationEnableStatus,
	0x0757: PropFollowFocusPositionCurrentValue,
	0x0758: PropFocusBracketShootingStatus,
	0x0759: PropPixelMappingEnableStatus,
	0x075A: PropTimeCodePresetResetEnableStatus,
	0x075B: PropUserBitPresetResetEnableStatus,
	0x075C: PropSensorCleaningEnableStatus,
	0x075D: PropPictureProfileResetEnableStatus,
	0x075E: PropCreativeLookResetEnableStatus,
	0x075F: PropLensVersionNumber,
	0x0760: PropDeviceOverheatingState,
	0x0761: PropMovieIntervalRecCountDownIntervalTime,
	0x0762: PropMovieIntervalRecRecordingDuration,
	0x0763: PropHighResolutionShutterSpeed,
	0x0764: PropBaseLookImportOperationEnableStatus,
	0x0765: PropLensModelName,
	0x0766: PropFocusPositionCurrentValue,
	0x0767: PropFocusDrivingStatus,
	0x0768: PropFlickerScanStatus,
	0x0769: PropFlickerScanEnableStatus,
	0x076B: PropFTPServerSettingOperationEnableStatus,
	0x076C: PropFTPTransferSettingSaveOperationEnableStatus,
	0x076D: PropFTPTransferSettingReadOperationEnableStatus,
	0x076E: PropFTPTransferSettingSaveReadState,
	0x0770: PropCameraShakeStatus,
	0x0771: PropUpdateBodyStatus,
	0x0773: PropMediaSLOT1WritingState,
	0x0774: PropMediaSLOT2WritingState,
	0x0776: PropMediaSLOT1RecordingAvailableType,
	0x0777: PropMediaSLOT2RecordingAvailableType,
	0x0778: PropMediaSLOT3RecordingAvailableType,
	0x0779: PropCameraOperatingMode,
	0x077A: PropPlaybackViewMode,
	0x0781: PropMediaSLOT3Status,
	0x0783: PropMediaSLOT3RemainingTime,
	0x0786: PropMonitoringDeliveringStatus,
	0x0787: PropMonitoringIsDelivering,
	0x0788: PropMonitoringSettingVersion,
	0x0789: PropMonitoringDeliveryTypeSupportInfo,
	0x078B: PropCameraErrorCautionStatus,
	0x078C: PropSystemErrorCautionStatus,
	0x078D: PropCameraButtonFunctionStatus,
	0x078E: PropFlickerLessShootingStatus,
	0x078F: PropContinuousShootingSpotBoostStatus,
	0x0790: PropContinuousShootingSpotBoostEnableStatus,
	0x0791: PropTimeShiftShootingStatus,
	0x0792: PropZoomDrivingStatus,
	0x0793: PropShootingSelfTimerStatus,
	0x0794: PropCreateNewFolderEnableStatus,
	0x0795: PropForcedFileNumberResetEnableStatus,
	0x0796: PropDefaultAFFreeSizeAndPositionSetting,
	0x0797: PropTrackingOnAndAFOnEnableStatus,
	0x0798: PropProgramShiftStatus,
	0x0799: PropMeteredManualLevel,
	0x079B: PropSecondBatteryRemain,
	0x079C: PropSecondBatteryLevel,
	0x079D: PropTotalBatteryRemain,
	0x079E: PropTotalBatteryLevel,
	0x028D: PropCameraLeverFunction,
	0x028E: PropShootingTimingPreNotificationMode,
	0x028F: PropMicrophoneDirectivity,
	0x0290: PropProductShowcaseSet,
	0x0291: PropAmountOfDefocusSetting,
	0x0292: PropCinematicVlogSetting,
	0x0293: PropCinematicVlogLook,
	0x0294: PropCinematicVlogMood,
	0x0295: PropCinematicVlogAFTransitionSpeed,
	0x07A2: PropMonitoringTransportProtocol,
	0x07A3: PropMonitoringAvailableFormat,
	0x07A4: PropMonitoringFormatSupportInformation,
	0x02BB: PropDeSqueezeDisplayRatio,
	0x02BC: PropZoomPositionSetting,
	0x07A1: PropZoomPositionCurrentValue,
	0x0F07: PropPriv0F07,
	0x0F08: PropPriv0F08,
	0x0F09: PropPriv0F09,
	0x0F0A: PropPriv0F0A,
	0x0F0C: PropPriv0F0C,
	0x0296: PropMonitoringOutputDisplaySDI,
	0x07AC: PropCameraSystemErrorInfo,
	0x0297: PropAFAreaPositionAFC,
	0x0298: PropAFAreaPositionAFS,
	0x07AD: PropFaceEyeDetectionAFStatus,
	0x0299: PropAutoFocusHold,
	0x029A: PropPushAFModeSetting,
	0x029B: PropTouchFunctionInMF,
	0x029C: PropPushAutoFocus,
	0x029D: PropPushAGC,
	0x029E: PropPushAutoIris,
	0x02A3: PropNDFilterPreset3Value,
	0x02A2: PropNDFilterPreset2Value,
	0x02A1: PropNDFilterPreset1Value,
	0x02A0: PropNDFilterPresetSelect,
	0x029F: PropPushAutoNDFilter,
	0x02AB: PropWhiteBalanceOffsetColorTemp,
	0x02AA: PropWhiteBalanceOffsetSetting,
	0x02A9: PropWhiteBalanceOffsetTintATW,
	0x02A8: PropWhiteBalanceOffsetColorTempATW,
	0x02A7: PropWhiteBalanceBGain,
	0x02A6: PropWhiteBalanceRGain,
	0x02A5: PropWhiteBalancePresetColorTemperature,
	0x02A4: PropWhiteBalanceSwitch,
	0x02B4: PropPaintLookDetailLevel,
	0x02B3: PropPaintLookDetailSetting,
	0x02B2: PropPaintLookKneeSlope,
	0x02B1: PropPaintLookKneePoint,
	0x02B0: PropPaintLookAutoKnee,
	0x02AF: PropPaintLookKneeSetting,
	0x02AE: PropPaintLookBBlack,
	0x02AD: PropPaintLookRBlack,
	0x02AC: PropPaintLookMasterBlack,
	0x07AE: PropUploadDatasetVersion,
	0x02B5: PropUserBaseLookOutput,
	0x07AF: PropMonitorLUTSettingOutputDestAssign,
	0x02B6: PropMonitorLUTSetting1,
	0x02B7: PropMonitorLUTSetting2,
	0x02B8: PropMonitorLUTSetting3,
	0x07B0: PropMaximumNumberOfBytes,
	0x02B9: PropSQModeSetting,
	0x02BA: PropMovieQualityFullAutoMode,
	0x02BD: PropFileSettingsCameraId,
	0x02BE: PropFileSettingsReelNumber,
	0x02BF: PropFileSettingsCameraPosition,
	0x02C0: PropImageStabilizationFramingStabilizer,
	0x02C1: PropExposureStep,
	0x07A0: PropTeleWideLeverValueCapability,
	0x02C2: PropEnlargeScreenSetting,
	0x07A5: PropMediaSLOT1ContentsInfoListEnableStatus,
	0x07A6: PropMediaSLOT2ContentsInfoListEnableStatus,
	0x07A7: PropMediaSLOT1ContentsInfoListRegenerateUpdateTime,
	0x07A8: PropMediaSLOT2ContentsInfoListRegenerateUpdateTime,
	0x07A9: PropMediaSLOT1ContentsInfoListUpdateTime,
	0x07AA: PropMediaSLOT2ContentsInfoListUpdateTime,
	0x07AB: PropPostViewTransferResourceStatus,
	0x02C3: PropSimulRecSetting,
	0x02C4: PropSimulRecSettingMovieRecButton,
	0x02C9: PropShutterSelectMode,
	0x02C5: PropOSDImageMode,
	0x07B1: PropFirmwareUpdateCommandVersion,
	0x02C6: PropDebugMode,
	0x0F0B: PropPriv0F0B,
	0x02C7: PropReserved18,
	0x02C8: PropReserved19,
	0x0F0D: PropPriv0F0D,
	0x07B5: PropSetPresetPTZFBinaryVersion,
	0x07B6: PropPanPositionStatus,
	0x07B7: PropTiltPositionStatus,
	0x07B8: PropPanPositionCurrentValue,
	0x07B9: PropTiltPositionCurrentValue,
	0x02CB: PropPanTiltAccelerationRampCurve,
	0x02CC: PropPanLimitMode,
	0x02CD: PropPanLimitRangeMinimum,
	0x02CE: PropPanLimitRangeMaximum,
	0x02CF: PropTiltLimitMode,
	0x02D0: PropTiltLimitRangeMinimum,
	0x02D1: PropTiltLimitRangeMaximum,
	0x02D2: PropPresetPTZFSlotNumber,
	0x07BA: PropCameraPowerStatus,
	0x0313: PropTargetStreamingDestinationSelect,
	0x07CD: PropStreamStatus,
	0x02D3: PropIRRemoteSetting,
	0x02D4: PropIPSetupProtocolSetting,
	0x07BB: PropRecordablePowerSources,
	0x07CF: PropStreamSettingListOperationStatus,
	0x02E4: PropPaintLookMultiMatrixAreaIndication,
	0x02D5: PropIrisCloseSetting,
	0x07BC: PropDisplayedMenuStatus,
	0x02D6: PropLanguageSetting,
	0x07BD: PropPlaybackContentsRecordingDateTime,
	0x07BE: PropPlaybackContentsName,
	0x07BF: PropPlaybackContentsNumber,
	0x07C0: PropPlaybackContentsTotalNumber,
	0x07C1: PropPlaybackContentsRecordingResolution,
	0x07C2: PropPlaybackContentsRecordingFrameRate,
	0x07C3: PropPlaybackContentsRecordingFileFormat,
	0x07C4: PropPlaybackContentsGammaType,
	0x07C5: PropBaseLookNameofPlayback,
	0x07C6: PropBaseLookAppliedofPlayback,
	0x02D7: PropPaintLookUserMatrixSetting,
	0x02D8: PropPaintLookUserMatrixLevel,
	0x02D9: PropPaintLookUserMatrixPhase,
	0x02DA: PropPaintLookUserMatrixRG,
	0x02DB: PropPaintLookUserMatrixRB,
	0x02DC: PropPaintLookUserMatrixGR,
	0x02DD: PropPaintLookUserMatrixGB,
	0x02DE: PropPaintLookUserMatrixBR,
	0x02DF: PropPaintLookUserMatrixBG,
	0x02E0: PropPaintLookMultiMatrixSetting,
	0x02E1: PropPaintLookMultiMatrixAxis,
	0x02E2: PropPaintLookMultiMatrixHue,
	0x02E3: PropPaintLookMultiMatrixSaturation,
	0x07C7: PropMonitoringOutputDisplaySettingDestAssign,
	0x02FC: PropMonitoringOutputDisplaySetting1,
	0x02FD: PropMonitoringOutputDisplaySetting2,
	0x07C8: PropFocusModeStatus,
	0x07C9: PropFocusOperationWithInt16EnableStatus,
	0x02E5: PropAudioInputCH1LevelControl,
	0x02E6: PropAudioInputCH2LevelControl,
	0x02E7: PropAudioInputCH3LevelControl,
	0x02E8: PropAudioInputCH4LevelControl,
	0x02E9: PropAudioInputCH1Level,
	0x02EA: PropAudioInputCH2Level,
	0x02EB: PropAudioInputCH3Level,
	0x02EC: PropAudioInputCH4Level,
	0x02ED: PropAudioInputCH1InputSelect,
	0x02EE: PropAudioInputCH2InputSelect,
	0x02EF: PropAudioInputCH3InputSelect,
	0x02F0: PropAudioInputCH4InputSelect,
	0x02F1: PropAudioInputCH1WindFilter,
	0x02F2: PropAudioInputCH2WindFilter,
	0x02F3: PropAudioInputCH3WindFilter,
	0x02F4: PropAudioInputCH4WindFilter,
	0x02F7: PropRemoteKeyThumbnailButton,
	0x02F8: PropRemoteKeySLOTSelectButton,
	0x02F9: PropVideoRecordingFormatBitrateSetting,
	0x07CA: PropValidRecordingVideoFormat,
	0x02FB: PropMonitoringOutputFormat,
	0x02FE: PropFocusSpeedDirectSync,
	0x02F5: PropAudioInput1TypeSelect,
	0x02F6: PropAudioInput2TypeSelect,
	0x02FA: PropVideoRecordingFormatQuality,
	0x02FF: PropLiveViewImageQualityByNumericalValue,
	0x0300: PropTallyLampControlRed,
	0x0301: PropTallyLampControlGreen,
	0x0302: PropTallyLampControlYellow,
	0x0305: PropMovieRecordingResolutionForRTSP,
	0x0307: PropMovieRecordingFrameRateRTSPSetting,
	0x0308: PropPictureCacheRecSetting,
	0x0309: PropPictureCacheRecSizeAndTime,
	0x030A: PropMovieIntervalRecFrames,
	0x030B: PropImagerScanMode,
	0x0306: PropMovieRecordingResolutionForRAW,
	0x07CB: PropLensSerialNumber,
	0x030C: PropShootingEnableSettingLicense,
	0x030D: PropGridLineDisplayPlayback,
	0x030E: PropGridLineType,
	0x07D9: PropCustomGridLineFileCommandVersion,
	0x07CC: PropMaximumSizeOfImageIDString,
	0x07CE: PropStreamButtonEnableStatus,
	0x07B3: PropAutoRecognitionTargetCandidates,
	0x02CA: PropAutoRecognitionTargetSetting,
	0x07D3: PropDeleteContentOperationEnableStatusSLOT1,
	0x07D4: PropDeleteContentOperationEnableStatusSLOT2,
	0x030F: PropDifferentSetForSQMovie,
	0x0310: PropManualInputForNDFilterValue,
	0x0311: PropLogShootingMode,
	0x0312: PropLogShootingModeColorGamut,
	0x0315: PropVideoStreamSelect,
	0x0314: PropStreamDisplayName,
	0x0316: PropVideoStreamResolution,
	0x0317: PropVideoStreamMaxBitRate,
	0x0318: PropVideoStreamAdaptiveRateControl,
	0x0319: PropVideoStreamCodec,
	0x031A: PropStreamLatency,
	0x031B: PropStreamTTL,
	0x031C: PropStreamCipherType,
	0x031D: PropStreamModeSetting,
	0x031E: PropVideoStreamResolutionMethod,
	0x031F: PropVideoStreamMovieRecPermission,
	0x0320: PropVideoStreamBitRateCompressionMode,
	0x0321: PropVideoStreamBitRateVBRMode,
	0x0322: PropAudioStreamCodecType,
	0x0323: PropAudioStreamSamplingFrequency,
	0x0324: PropAudioStreamBitDepth,
	0x0325: PropAudioStreamChannel,
	0x0326: PropHomeMenuSetting,
	0x0327: PropCallSetting,
	0x0328: PropNDFilterPositionSetting,
	0x07D5: PropSceneFileCommandVersion,
	0x07D6: PropSceneFileUploadOperationEnableStatus,
	0x07D7: PropSceneFileDownloadOperationEnableStatus,
	0x07D8: PropSceneFileIndexesAvailableForDownload,
	0x07D1: PropEframingType,
	0x07D2: PropEframingCommandVersion,
	0x0329: PropEframingAutoFraming,
	0x032A: PropEframingTrackingStartMode,
	0x032B: PropEframingProductionEffect,
	0x032C: PropEframingSpeedPTZ,
	0x0601: PropPriv0601,
	0x0331: PropTopOfTheGroupShootingMarkSetting,
	0x0603: PropPriv0603,
	0x0332: PropCompRAWShootingNR,
	0x0333: PropCompRAWShootingNRFileCompressionType,
	0x0334: PropCompRAWShootingNRNumberOfSheets,
	0x07DA: PropElapsedBulbExposureTime,
	0x07DB: PropRemainingBulbExposureTime,
	0x07DC: PropRemainingNoiseReductionTime,
	0x0F01: PropPriv0F01,
	0x0F02: PropPriv0F02,
	0x032D: PropDigitalExtenderMagnificationSetting,
	0x032E: PropMovieRecReviewPlayingState,
	0x011F: PropNearFar,
	0x0120: PropReserved7,
	0x0121: PropAFAreaPosition,
	0x0126: PropZoomOperation,
	0x0134: PropZoomAndFocusPositionSave,
	0x0135: PropZoomAndFocusPositionLoad,
	0x015B: PropColortempStep,
	0x015E: PropWhiteBalanceTintStep,
	0x015F: PropFocusOperation,
	0x0162: PropShutterECSNumberStep,
	0x0193: PropRemoteTouchOperation,
	0x027C: PropZoomAndFocusPresetZoomOnlySet,
	0x0509: PropCustomWBCaptureStandby,
	0x050A: PropCustomWBCaptureStandbyCancel,
	0x050B: PropCustomWBCapture,
	0x0303: PropZoomOperationWithInt16,
	0x0304: PropFocusOperationWithInt16,
	0x032F: PropHighResolutionShutterSpeedAdjust,
	0x0330: PropHighResolutionShutterSpeedAdjustInIntegralMultiples,
	0x034B: PropMovieAngleOfViewPriority,
	0x034C: PropWindNoiseReductForExternalMic,
	0x034D: PropNoiseCutFilter,
	0x034E: PropNoiseCutFilterForExternalMic,
	0x07DD: PropDispModeCandidateStill,
	0x0339: PropDispModeSettingStill,
	0x033A: PropDispModeStill,
	0x07DE: PropDispModeCandidateMovie,
	0x033B: PropDispModeSettingMovie,
	0x033C: PropDispModeMovie,
	0x0335: PropCompRAWShootingHDR,
	0x0338: PropCompRAWShootingHDRDRSetting,
	0x0336: PropCompRAWShootingHDRFileCompressionType,
	0x0337: PropCompRAWShootingHDRNumberOfSheets,
	0x07DF: PropControlGeneralSettingFileEnableStatus,
	0x033D: PropPeakingDisplay,
	0x033E: PropPeakingLevel,
	0x033F: PropPeakingColor,
	0x0340: PropZebraDisplay,
	0x0341: PropZebraLevel,
	0x0342: PropZebraLevelTypeCustom,
	0x0343: PropZebraLevelStandardCustom,
	0x0344: PropZebraLevelRangeCustom,
	0x0345: PropZebraLevelLowerLimitCustom,
	0x0346: PropMarkerDisplay,
	0x0347: PropCenterMarkerDisplay,
	0x0348: PropAspectMarkerRatioMovie,
	0x0349: PropSafetyZoneDisplay,
	0x034A: PropGuideframeDisplay,
	0x034F: PropDualGain,
	0x0350: PropImagerMode,
	0x0351: PropDisplayQualityForFinderOnly,
	0x0352: PropDisplayQualityForMonitorOnly,
	0x0F06: PropPriv0F06,
	0x0602: PropPriv0602,
	0x0F03: PropPriv0F03,
	0x0F04: PropPriv0F04,
	0x0F05: PropPriv0F05,
}
