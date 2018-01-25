package contenttype

import "testing"

func TestGetExtensionByContentType(t *testing.T) {
	ext := GetExtensionByContentType("image/jpeg")
	if ext != "jpeg" {
		t.Error("Unmatched (image/jpeg)")
	}
}

func TestGetContentTypeByExtension(t *testing.T) {
	contentType := GetContentTypeByExtension("jpg")
	if contentType != "image/jpeg" {
		t.Error("Unmatched (jpg)")
	}

	contentType = GetContentTypeByExtension("jpeg")
	if contentType != "image/jpeg" {
		t.Error("Unmatched (jpeg)")
	}
}
