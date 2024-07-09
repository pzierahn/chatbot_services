package utils

import pb "github.com/pzierahn/chatbot_services/services/proto"

func GetDocumentTitle(meta *pb.DocumentMetadata) string {
	switch meta.Data.(type) {
	case *pb.DocumentMetadata_Web:
		return meta.GetWeb().Title
	case *pb.DocumentMetadata_File:
		return meta.GetFile().Filename
	default:
		return ""
	}
}
