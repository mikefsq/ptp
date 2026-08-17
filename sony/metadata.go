package sony

import (
	"fmt"

	"github.com/mikefsq/ptp"
)

// String metadata embedded in the files the camera writes.
//
// The set is small and it is NOT the same as Fujifilm's. Sony has no Artist and
// no Comment field; the photographer credit is SetPhotographer, and there is no
// free-text comment at all. Naming a Sony accessor "Artist" for symmetry would
// be inventing a field the body does not have.
//
// These six are every property Sony's own header marks CrDataType_STR, so this
// is the complete string surface rather than a selection:
//
//	SetPhotographer                 0xD1CE  the photographer credit
//	SetCopyright                    0xD1CF  the copyright line
//	FileSettingsTitleNameSettings   0xD1DC  title baked into the filename
//	RecordingSettingFileName        0xD1CA  recording-setting file name
//	CreateNewFolderEnableStatus     0xD1DB  (status, not settable text)
//	body serial number                      read-only, and reported in DeviceInfo
//
// UNVERIFIED, like the rest of the vendor surface.

// WriteCopyrightInfo values (CrWriteCopyrightInfo). Note it starts at 1, not 0 —
// zero is not "off", it is not a value the property takes.
const (
	CopyrightInfoOff uint64 = 0x01
	CopyrightInfoOn  uint64 = 0x02
)

// Photographer reads the photographer credit written into each file.
func (c *Camera) Photographer() (string, error) {
	return c.propString(PropSetPhotographer)
}

// SetPhotographer sets the photographer credit.
//
// It is only embedded when WriteCopyrightInfo is on — setting the text alone
// changes nothing in the files, which is the obvious way to lose an hour.
// SetCopyrightInfo turns it on.
func (c *Camera) SetPhotographer(v string) error {
	return c.SetPropertyString(PropSetPhotographer, v)
}

// Copyright reads the copyright line written into each file.
func (c *Camera) Copyright() (string, error) {
	return c.propString(PropSetCopyright)
}

// SetCopyright sets the copyright line. Gated by WriteCopyrightInfo, as
// SetPhotographer is.
func (c *Camera) SetCopyright(v string) error {
	return c.SetPropertyString(PropSetCopyright, v)
}

// CopyrightInfo reports whether the camera embeds the photographer and
// copyright strings in the files it writes.
func (c *Camera) CopyrightInfo() (bool, error) {
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return false, err
	}
	d := FindProp(props, PropWriteCopyrightInfo)
	if d == nil {
		return false, fmt.Errorf("sony: this body does not report %s", PropName(PropWriteCopyrightInfo))
	}
	return d.Current == CopyrightInfoOn, nil
}

// SetCopyrightInfo turns embedding of the copyright strings on or off.
//
// This is the switch that makes Photographer and Copyright take effect. It is a
// UInt8 property whose off value is 1, not 0.
func (c *Camera) SetCopyrightInfo(on bool) error {
	v := CopyrightInfoOff
	if on {
		v = CopyrightInfoOn
	}
	return c.SetProperty(PropWriteCopyrightInfo, ptp.TypeUint8, v)
}

// TitleName reads the title the camera bakes into filenames.
func (c *Camera) TitleName() (string, error) {
	return c.propString(PropFileSettingsTitleNameSettings)
}

// SetTitleName sets the title baked into filenames.
func (c *Camera) SetTitleName(v string) error {
	return c.SetPropertyString(PropFileSettingsTitleNameSettings, v)
}

// RecordingSettingFileName reads the recording-setting file name.
func (c *Camera) RecordingSettingFileName() (string, error) {
	return c.propString(PropRecordingSettingFileName)
}

// SetRecordingSettingFileName sets the recording-setting file name.
func (c *Camera) SetRecordingSettingFileName(v string) error {
	return c.SetPropertyString(PropRecordingSettingFileName, v)
}

// Metadata is every string field the camera reports, read in one round trip.
//
// Prefer this to the individual accessors when reading more than one: each of
// those costs a full property snapshot.
type Metadata struct {
	Photographer string
	Copyright    string
	TitleName    string
	RecordingSet string

	// Embedding reports whether WriteCopyrightInfo is on. The strings above are
	// stored regardless; this is whether they reach the files.
	Embedding    bool
	HasEmbedding bool
}

// ReadMetadata pulls the string fields out of a property snapshot.
func ReadMetadata(props []DeviceProperty) Metadata {
	var m Metadata
	m.Photographer, _ = PropString(props, PropSetPhotographer)
	m.Copyright, _ = PropString(props, PropSetCopyright)
	m.TitleName, _ = PropString(props, PropFileSettingsTitleNameSettings)
	m.RecordingSet, _ = PropString(props, PropRecordingSettingFileName)
	if d := FindProp(props, PropWriteCopyrightInfo); d != nil {
		m.Embedding, m.HasEmbedding = d.Current == CopyrightInfoOn, true
	}
	return m
}

// Metadata fetches a snapshot and decodes the string fields.
func (c *Camera) Metadata() (Metadata, error) {
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return Metadata{}, err
	}
	return ReadMetadata(props), nil
}

func (m Metadata) String() string {
	embed := "embedding unknown"
	if m.HasEmbedding {
		embed = "embedding off"
		if m.Embedding {
			embed = "embedding on"
		}
	}
	return fmt.Sprintf("photographer=%q copyright=%q title=%q (%s)",
		m.Photographer, m.Copyright, m.TitleName, embed)
}

// propString reads one string property, reporting clearly when the body does
// not expose it as a string rather than returning a silent empty value.
func (c *Camera) propString(p Prop) (string, error) {
	props, err := c.GetAllDevicePropData()
	if err != nil {
		return "", err
	}
	d := FindProp(props, p)
	if d == nil {
		return "", fmt.Errorf("sony: this body does not report %s", PropName(p))
	}
	if d.Type != ptp.TypeString {
		return "", fmt.Errorf("sony: %s is not a string property (type 0x%04X)",
			PropName(p), uint16(d.Type))
	}
	return d.CurrentStr, nil
}
