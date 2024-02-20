package documents

import pb "github.com/pzierahn/chatbot_services/proto"

type Webpage struct {
	Url string `json:"url"`
}

type File struct {
	Path string `json:"path,omitempty"`
}

type DocumentMeta struct {
	Webpage *Webpage `json:"webpage,omitempty"`
	File    *File    `json:"file,omitempty"`
}

func (meta *DocumentMeta) IsWebpage() bool {
	return meta.Webpage != nil
}

func (meta *DocumentMeta) IsFile() bool {
	return meta.File != nil
}

func metaFromProto(meta *pb.DocumentMetadata) *DocumentMeta {

	if meta == nil {
		return nil
	}

	if meta.GetWeb() != nil {
		return &DocumentMeta{
			Webpage: &Webpage{
				Url: meta.GetWeb().Url,
			},
		}
	}

	if meta.GetFile() != nil {
		return &DocumentMeta{
			File: &File{
				Path: meta.GetFile().Path,
			},
		}
	}

	return nil
}

func metaToProto(meta DocumentMeta) *pb.DocumentMetadata {

	if meta.IsWebpage() {
		return &pb.DocumentMetadata{
			Data: &pb.DocumentMetadata_Web{
				Web: &pb.Webpage{
					Url: meta.Webpage.Url,
				},
			},
		}
	}

	if meta.IsFile() {
		return &pb.DocumentMetadata{
			Data: &pb.DocumentMetadata_File{
				File: &pb.File{
					Path: meta.File.Path,
				},
			},
		}
	}

	return nil
}
