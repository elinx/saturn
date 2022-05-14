package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"

	log "github.com/sirupsen/logrus"
)

const (
	ContainerFilename = "META-INF/container.xml" // filename of container.xml
)

type Container struct {
	Container xml.Name `xml:"container"`
	XMLNS     string   `xml:"xmlns,attr"`
	Version   string   `xml:"version,attr"`
	Rootfiles []struct {
		Rootfile  xml.Name `xml:"rootfile"`
		FullPath  string   `xml:"full-path,attr"`
		MediaType string   `xml:"media-type,attr"`
	} `xml:"rootfiles>rootfile"`
}

type Rootfile struct {
	Metadata struct {
		Rights      string `xml:"rights"`
		ISBN        string `xml:"identifier"`
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Creator     string `xml:"creator"`
		Date        string `xml:"date"`
		Publisher   string `xml:"publisher"`
		Language    string `xml:"language"`
		Format      string `xml:"format"`
	} `xml:"metadata"`
	Manifest struct {
		Items []struct {
			ID        string `xml:"id,attr"`
			Href      string `xml:"href,attr"`
			MediaType string `xml:"media-type,attr"`
		} `xml:"item"`
	} `xml:"manifest"`
	Spine struct {
		TocID string `xml:"toc,attr"`
		Items []struct {
			IDref string `xml:"idref,attr"`
			Link  string `xml:"linear,attr"`
		} `xml:"itemref"`
	} `xml:"spine"`
	Guide struct {
		Items []struct {
			Type  string `xml:"type,attr"`
			Href  string `xml:"href,attr"`
			Title string `xml:"title"`
		} `xml:"reference"`
	} `xml:"guide"`
}

type navPoint struct {
	ID        string `xml:"id,attr"`
	PlayOrder string `xml:"playOrder,attr"`
	NavLable  struct {
		Text string `xml:",chardata"`
	} `xml:"navLabel>text"`
	Content struct {
		Src string `xml:"src,attr"`
	} `xml:"content"`
}

type Toc struct {
	Head struct {
		Metas []struct {
			Name  string `xml:"name,attr"`
			Value string `xml:"content,attr"`
		} `xml:"meta"`
	} `xml:"head"`
	DocTitle struct {
		Text string `xml:",chardata"`
	} `xml:"docTitle"`
	NavMap struct {
		//TODO(elinx): better way to express nested navPoint
		NavPoints []struct {
			ID        string `xml:"id,attr"`
			PlayOrder string `xml:"playOrder,attr"`
			NavLable  struct {
				Text string `xml:",chardata"`
			} `xml:"navLabel>text"`
			Content struct {
				Src string `xml:"src,attr"`
			} `xml:"content"`
			NavPoints []navPoint `xml:"navPoint"`
		} `xml:"navPoint"`
	} `xml:"navMap"`
}

type Epub struct {
	Filename     string
	Container    Container
	Rootfile     Rootfile
	Toc          Toc
	SeqFiles     []*zip.File
	Files        map[string]*zip.File
	readercloser *zip.ReadCloser
}

func NewEpub(filename string) *Epub {
	return &Epub{
		Filename: filename,
		Files:    make(map[string]*zip.File),
	}
}

func (epub *Epub) Open() error {
	reader, err := zip.OpenReader(epub.Filename)
	if err != nil {
		return err
	}
	epub.readercloser = reader
	epub.SeqFiles = reader.File

	for i, f := range reader.File {
		epub.Files[f.Name] = reader.File[i]
	}

	if f, found := epub.Files[ContainerFilename]; !found {
		return fmt.Errorf("%s not found", ContainerFilename)
	} else {
		if rc, err := f.Open(); err != nil {
			return err
		} else {
			if err := xml.NewDecoder(rc).Decode(&epub.Container); err != nil {
				return err
			}
			log.Println(epub.Container)
		}
	}
	if f, found := epub.Files[epub.Container.Rootfiles[0].FullPath]; !found {
		return fmt.Errorf("%s not found", epub.Container.Rootfiles[0].FullPath)
	} else {
		if rc, err := f.Open(); err != nil {
			return err
		} else {
			if err := xml.NewDecoder(rc).Decode(&epub.Rootfile); err != nil {
				return err
			}
		}
	}
	if err := epub.getTableOfContent(); err != nil {
		return err
	}
	return nil
}

func (epub *Epub) Close() {
	epub.readercloser.Close()
}

func (epub *Epub) GetContentByFilePath(filepath string) (string, error) {
	if f, found := epub.Files[filepath]; !found {
		return "", fmt.Errorf("%s not found", filepath)
	} else {
		if rc, err := f.Open(); err != nil {
			return "", err
		} else {
			if content, err := io.ReadAll(rc); err != nil {
				return "", err
			} else {
				return string(content), nil
			}
		}
	}
}

func (epub *Epub) getFullPath(id string) string {
	namespace := epub.Container.Rootfiles[0].FullPath[:len(epub.Container.Rootfiles[0].FullPath)-len(path.Base(epub.Container.Rootfiles[0].FullPath))]
	return path.Join(namespace, id)
}

func (epub *Epub) getManifestFilePathById(id string) string {
	log.Printf("id: %s\n", id)
	for i, v := range epub.Rootfile.Manifest.Items {
		if v.ID == id {
			return epub.getFullPath(epub.Rootfile.Manifest.Items[i].Href)
		}
	}
	return ""
}

func (epub *Epub) getChapterNameByIndex(index int) string {
	ref := epub.Rootfile.Spine.Items[index].IDref
	return epub.getManifestFilePathById(ref)
}

func (epub *Epub) GetChapterByIndex(index int) (string, error) {
	if index < 0 || index >= len(epub.Rootfile.Spine.Items) {
		return "", fmt.Errorf("index out of range")
	}
	return epub.GetContentByFilePath(epub.getChapterNameByIndex(index))
}

func (epub *Epub) getTableOfContent() error {
	filepath := epub.getManifestFilePathById(epub.Rootfile.Spine.TocID)
	if f, found := epub.Files[filepath]; !found {
		return fmt.Errorf("%s not found", filepath)
	} else {
		if rc, err := f.Open(); err != nil {
			return err
		} else {
			if err := xml.NewDecoder(rc).Decode(&epub.Toc); err != nil {
				return err
			}
		}
	}
	return nil
}
