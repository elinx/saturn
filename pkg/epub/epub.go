package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"
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

type Epub struct {
	Filename     string
	Container    Container
	Rootfile     Rootfile
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

func (epub *Epub) OpenFile() error {
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
			fmt.Println(epub.Container)
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
	return nil
}

func (epub *Epub) Close() {
	epub.readercloser.Close()
}

func (epub *Epub) GetChapter(id string) (string, error) {
	if f, found := epub.Files[id]; !found {
		return "", fmt.Errorf("%s not found", id)
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

func (epub *Epub) getChapterName(index int) string {
	ref := epub.Rootfile.Spine.Items[index].IDref
	for i, v := range epub.Rootfile.Manifest.Items {
		if v.ID == ref {
			return epub.getFullPath(epub.Rootfile.Manifest.Items[i].Href)
		}
	}
	return ""
}

func (epub *Epub) GetChapterByIndex(index int) (string, error) {
	if index < 0 || index >= len(epub.Rootfile.Spine.Items) {
		return "", fmt.Errorf("index out of range")
	}
	return epub.GetChapter(epub.getChapterName(index))
}
