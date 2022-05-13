package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
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
	Package  xml.Name `xml:"package"`
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
}

var container Container
var rootfile Rootfile

func OpenFile(path string) error {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, f := range reader.File {
		if f.Name == "META-INF/container.xml" {
			if rc, err := f.Open(); err != nil {
				return err
			} else {
				if err := xml.NewDecoder(rc).Decode(&container); err != nil {
					return err
				}
				fmt.Println(container)
			}
		}
	}
	if len(container.Rootfiles) == 0 {
		return fmt.Errorf("no rootfile found")
	}

	// only handle one book for now
	idx := 0
	for _, f := range reader.File {
		if f.Name == container.Rootfiles[idx].FullPath {
			if rc, err := f.Open(); err != nil {
				return err
			} else {
				if err := xml.NewDecoder(rc).Decode(&rootfile); err != nil {
					return err
				}
				fmt.Println(rootfile)
			}
		}
	}
	return nil
}
