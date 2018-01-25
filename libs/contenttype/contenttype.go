package contenttype

func GetExtensionByContentType(contentType string) string {
	return contentTypeToExt[contentType]
}

func GetContentTypeByExtension(extension string) string {
	return extToContentType[extension]
}
