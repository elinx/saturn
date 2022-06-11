package epub

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"path"
	"strings"

	cssparser "github.com/elinx/saturn/pkg/css_parser"
	log "github.com/sirupsen/logrus"
)

const (
	ContainerFilename = "META-INF/container.xml" // filename of container.xml
)

// This the manifest item id of the epub content
type ManifestId string
type HRef string

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
			ID        ManifestId `xml:"id,attr"`
			Href      HRef       `xml:"href,attr"`
			MediaType string     `xml:"media-type,attr"`
		} `xml:"item"`
	} `xml:"manifest"`
	Spine struct {
		TocID ManifestId `xml:"toc,attr"`
		Items []struct {
			IDref ManifestId `xml:"idref,attr"`
			Link  string     `xml:"linear,attr"`
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
		Src HRef `xml:"src,attr"`
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
				Src HRef `xml:"src,attr"`
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
	Styles       []*cssparser.Rule
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
	if err := epub.parseTableOfContent(); err != nil {
		return err
	}
	if err := epub.parseCssFiles(); err != nil {
		return err
	}
	return nil
}

func (epub *Epub) Close() {
	epub.readercloser.Close()
}

// getContentByFilePath return file content by full filepath(relative to rootfile)
func (epub *Epub) getContentByFilePath(filepath string) (string, error) {
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

// GetContentByHref return file content by href(setted in spine)
func (epub *Epub) GetContentByHref(href HRef) (string, error) {
	filepath := epub.GetFullPath(href)
	return epub.getContentByFilePath(filepath)
}

// GetContentByManifestId return file content by manifest id
func (epub *Epub) GetContentByManifestId(id ManifestId) (string, error) {
	filepath := epub.getManifestFilePathById(id)
	log.Infof("filepath: %s", filepath)
	return epub.getContentByFilePath(filepath)
}

// GetFullPath return full filepath(relative to rootfile) by href(setted in spine)
func (epub *Epub) GetFullPath(href HRef) string {
	namespace := epub.Container.Rootfiles[0].FullPath[:len(epub.Container.Rootfiles[0].FullPath)-len(path.Base(epub.Container.Rootfiles[0].FullPath))]
	return path.Join(namespace, string(href))
}

func (epub *Epub) getManifestFilePathById(id ManifestId) string {
	log.Printf("id: %s\n", id)
	for i, v := range epub.Rootfile.Manifest.Items {
		if v.ID == id {
			return epub.GetFullPath(epub.Rootfile.Manifest.Items[i].Href)
		}
	}
	return ""
}

func (epub *Epub) HrefToManifestId(href HRef) ManifestId {
	for _, v := range epub.Rootfile.Manifest.Items {
		if v.Href == href {
			return v.ID
		}
	}
	return ""
}

func (epub *Epub) ManifestIdToHref(id ManifestId) HRef {
	for _, v := range epub.Rootfile.Manifest.Items {
		if v.ID == id {
			return v.Href
		}
	}
	return ""
}

type SpineContent struct {
	Orders   []ManifestId
	Contents map[ManifestId]string
}

// GetSpinContent return all spine content in `href: html content` format
func (epub *Epub) GetSpinContent() (*SpineContent, error) {
	content := &SpineContent{
		Orders:   make([]ManifestId, 0),
		Contents: make(map[ManifestId]string),
	}
	for _, v := range epub.Rootfile.Spine.Items {
		if c, err := epub.GetContentByManifestId(v.IDref); err != nil {
			return nil, err
		} else {
			content.Orders = append(content.Orders, v.IDref)
			content.Contents[v.IDref] = c
		}
	}
	return content, nil
}

func (epub *Epub) GetNextSection(id ManifestId) (string, ManifestId, error) {
	for i, v := range epub.Rootfile.Spine.Items {
		if v.IDref == id {
			if i+1 < len(epub.Rootfile.Spine.Items) {
				nextId := epub.Rootfile.Spine.Items[i+1].IDref
				content, err := epub.getContentByFilePath(epub.getManifestFilePathById(epub.Rootfile.Spine.Items[i+1].IDref))
				return content, nextId, err
			}
			return "", "", fmt.Errorf("end of epub")
		}
	}
	return "", "", fmt.Errorf("id not found")
}

func (epub *Epub) GetPrevSection(id ManifestId) (string, ManifestId, error) {
	for i, v := range epub.Rootfile.Spine.Items {
		if v.IDref == id {
			if i-1 >= 0 {
				prevId := epub.Rootfile.Spine.Items[i-1].IDref
				content, err := epub.getContentByFilePath(epub.getManifestFilePathById(epub.Rootfile.Spine.Items[i-1].IDref))
				return content, prevId, err
			}
			return "", "", fmt.Errorf("start of epub")
		}
	}
	return "", "", fmt.Errorf("id not found")
}

// func (epub *Epub) GetChapterIndex(id ManifestId) int {
// 	for i, v := range epub.Toc.NavMap.NavPoints {
// 		if v.IDref == id {
// 			return i
// 		}
// 	}
// 	return -1
// }

// func (epub *Epub) GetNextChapter(id ManifestId) (string, error) {
// 	index := epub.GetChapterIndex(id)
// 	if index == -1 {
// 		return "", fmt.Errorf("chapter not found")
// 	}
// 	return epub.GetChapterByIndex(index + 1)
// }

// func (epub *Epub) GetPrevChapter(id ManifestId) (string, error) {
// 	index := epub.GetChapterIndex(id)
// 	if index == -1 {
// 		return "", fmt.Errorf("chapter not found")
// 	}
// 	return epub.GetChapterByIndex(index - 1)
// }

func (epub *Epub) parseTableOfContent() error {
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

type TableOfContent struct {
	Orders  []string // Titles in order
	Content map[string]struct {
		Level int
		ID    ManifestId
	}
}

func (epub *Epub) GetTableOfContent() *TableOfContent {
	content := &TableOfContent{
		Orders: make([]string, 0),
		Content: make(map[string]struct {
			Level int
			ID    ManifestId
		}),
	}
	for _, v := range epub.Toc.NavMap.NavPoints {
		content.Orders = append(content.Orders, v.NavLable.Text)
		content.Content[v.NavLable.Text] = struct {
			Level int
			ID    ManifestId
		}{0, epub.HrefToManifestId(v.Content.Src)}
		for _, v := range v.NavPoints {
			content.Orders = append(content.Orders, v.NavLable.Text)
			content.Content[v.NavLable.Text] = struct {
				Level int
				ID    ManifestId
			}{1, epub.HrefToManifestId(v.Content.Src)}
		}
	}
	return content
}

func (epub *Epub) getCssFiles() []string {
	var cssFiles []string
	for _, item := range epub.Rootfile.Manifest.Items {
		if strings.HasSuffix(string(item.Href), ".css") {
			cssFiles = append(cssFiles, epub.GetFullPath(item.Href))
		}
	}
	return cssFiles
}

func (epub *Epub) parseCssFiles() error {
	for _, filename := range epub.getCssFiles() {
		if content, err := epub.getContentByFilePath(filename); err == nil {
			if rules, err := cssparser.NewParser().Parse(content); err != nil {
				return err
			} else {
				epub.Styles = append(epub.Styles, rules...)
			}
		}
	}
	return nil
}
