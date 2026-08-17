package ptp

import "testing"

// SendObjectInfo answers with three handles: storage, parent, and the handle
// the object will take. A camera that returns fewer has not told us where the
// object went, so sending its bytes would write them against an unknown handle.
func TestSendObjectRequiresThreeParams(t *testing.T) {
	f := &fakeTransport{in: [][]byte{
		mkContainer(ContainerResponse, uint16(RespOK), 1, []uint32{1}, nil),
	}}
	s := NewSession(f)

	_, _, _, err := s.SendObject(1, 0, &ObjectInfo{Filename: "x.jpg"}, []byte{1, 2, 3})
	if err == nil {
		t.Fatal("expected an error when SendObjectInfo returns too few parameters")
	}
}

// Upload depends on the encoder and the parser agreeing exactly. A field-order
// slip would corrupt every uploaded file's metadata, and the camera would
// accept it — nothing on the wire says which field is which.
//
// This is the only test that exercises EncodeObjectInfo: the parser's own
// round-trip test builds its bytes by hand, so it would not catch an encoder
// that emitted them in the wrong order.
func TestObjectInfoEncodeDecodeRoundTrip(t *testing.T) {
	want := &ObjectInfo{
		StorageID: 0x00010001, ObjectFormat: FormatEXIFJPEG, ProtectionStatus: 0,
		CompressedSize: 3080192, ThumbFormat: FormatEXIFJPEG, ThumbCompressedSize: 1124,
		ThumbPixWidth: 160, ThumbPixHeight: 120,
		ImagePixWidth: 4912, ImagePixHeight: 3264, ImageBitDepth: 24,
		ParentObject: 0x00080000, AssociationType: 0, AssociationDesc: 0,
		SequenceNumber: 449, Filename: "_DSC0449.JPG",
		CaptureDate: "20200831T142530", ModificationDate: "20200831T142530", Keywords: "",
	}
	got, err := ParseObjectInfo(EncodeObjectInfo(want))
	if err != nil {
		t.Fatalf("round trip: %v", err)
	}
	if *got != *want {
		t.Errorf("round trip changed the struct:\n got %+v\nwant %+v", *got, *want)
	}
}
