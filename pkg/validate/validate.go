package validate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/docker/docker/client"
	coreImage "github.com/jparrill/decker/pkg/core/image"
	coreReg "github.com/jparrill/decker/pkg/core/registry"

	imageapi "github.com/openshift/api/image/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	ReleaseImageStreamFile   = "release-manifests/image-references"
	ReleaseImageMetadataFile = "release-manifests/0000_50_installer_coreos-bootimages.yaml"
)

type OCPImage struct {
	coreImage.ContainerImage
	SrcRegistry *coreReg.Registry
	DstRegistry *coreReg.Registry
}

type OCPVersion struct {
	Name        string `json:"name"`
	PullSpec    string `json:"pullSpec"`
	DownloadURL string `json:"downloadURL"`
}

type RegistryClientProvider struct{}

func NewValidateOCPImage(url, auth, filePath string, dCLi *client.Client) *OCPImage {
	ci, err := coreImage.NewContainerImage(url, auth, filePath, dCLi)
	if err != nil {
		panic(err)
	}

	return &OCPImage{
		ContainerImage: *ci,
	}
}

func (p *RegistryClientProvider) Lookup(ctx context.Context, image string, pullSecret string) (releaseImage *ReleaseImage, err error) {
	fileContents, err := ExtractImageFiles(ctx, image, pullSecret, ReleaseImageStreamFile, ReleaseImageMetadataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to extract release metadata: %w", err)
	}

	if _, ok := fileContents[ReleaseImageStreamFile]; !ok {
		return nil, fmt.Errorf("release image references file not found in release image %s", image)
	}
	imageStream, err := DeserializeImageStream(fileContents[ReleaseImageStreamFile])
	if err != nil {
		return nil, err
	}

	if _, ok := fileContents[ReleaseImageMetadataFile]; !ok {
		return nil, fmt.Errorf("release image metadata file not found in release image %s", image)
	}
	coreOSMeta, err := DeserializeImageMetadata(fileContents[ReleaseImageMetadataFile])
	if err != nil {
		return nil, err
	}

	return &ReleaseImage{
		ImageStream:    imageStream,
		StreamMetadata: coreOSMeta,
	}, nil
}

func DeserializeImageStream(data []byte) (*imageapi.ImageStream, error) {
	var imageStream imageapi.ImageStream
	if err := json.Unmarshal(data, &imageStream); err != nil {
		return nil, fmt.Errorf("couldn't read image stream data as a serialized ImageStream: %w\nraw data:\n%s", err, string(data))
	}
	return &imageStream, nil
}

func DeserializeImageMetadata(data []byte) (*CoreOSStreamMetadata, error) {
	var coreOSMetaCM corev1.ConfigMap
	if err := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 100).Decode(&coreOSMetaCM); err != nil {
		return nil, fmt.Errorf("couldn't read image lookup data as serialized ConfigMap: %w\nraw data:\n%s", err, string(data))
	}
	streamData, hasStreamData := coreOSMetaCM.Data["stream"]
	if !hasStreamData {
		return nil, fmt.Errorf("coreos stream metadata configmap is missing the 'stream' key")
	}
	var coreOSMeta CoreOSStreamMetadata
	if err := json.Unmarshal([]byte(streamData), &coreOSMeta); err != nil {
		return nil, fmt.Errorf("couldn't decode stream metadata data: %w\n%s", err, streamData)
	}
	return &coreOSMeta, nil
}

// ExtractImageFiles extracts a list of files from a registry image given the image reference, pull secret and the
// list of files to extract. It returns a map with file contents or an error.
func ExtractImageFiles(ctx context.Context, imageRef string, pullSecret string, files ...string) (map[string][]byte, error) {
	//layers, fromBlobs, err := getMetadata(ctx, imageRef, pullSecret)
	sys := &types.SystemContext{
		AuthFilePath: pullSecret,
	}
	printImageMetadata(imageRef, sys)
	//if err != nil {
	//	return nil, err
	//}

	fileContents := map[string][]byte{}
	for _, file := range files {
		fileContents[file] = nil
	}
	if len(fileContents) == 0 {
		return fileContents, nil
	}

	// Iterate over layers in reverse order to find the most recent version of files
	//for i := len(layers) - 1; i >= 0; i-- {
	//layer := layers[i]
	//err := func() error {
	//r, err := fromBlobs.Open(ctx, layer.Digest)
	//if err != nil {
	//return fmt.Errorf("unable to access the source layer %s: %v", layer.Digest, err)
	//}
	//defer r.Close()
	//rc, err := dockerarchive.DecompressStream(r)
	//if err != nil {
	//return err
	//}
	//defer rc.Close()
	//tr := tar.NewReader(rc)
	//for {
	//hdr, err := tr.Next()
	//if err != nil {
	//if err == io.EOF {
	//break
	//}
	//return err
	//}
	//if hdr.Typeflag == tar.TypeReg {
	//value, needFile := fileContents[hdr.Name]
	//if !needFile {
	//continue
	//}
	//// If value already assigned, the content was found in an earlier layer
	//if value != nil {
	//continue
	//}
	//out := &bytes.Buffer{}
	//if _, err := io.Copy(out, tr); err != nil {
	//return err
	//}
	//fileContents[hdr.Name] = out.Bytes()
	//}
	//if allFound(fileContents) {
	//break
	//}
	//}
	//return nil
	//}()
	//if err != nil {
	//return nil, err
	//}
	//if allFound(fileContents) {
	//break
	//}
	//}
	return fileContents, nil
}

//func getMetadata(imageref, pullSecret string) {
//	// Convierte la referencia en una estructura ImageReference
//	ref, err := alltransports.ParseImageName(imageref)
//	if err != nil {
//		fmt.Printf("Error al analizar la referencia de la imagen: %v\n", err)
//		return
//	}
//
//	// Obtiene los metadatos de la imagen
//	sys := docker.SystemContext{}
//	manifest, _, err := ref.Transport().Get(ref.Context(), sys)
//	if err != nil {
//		fmt.Printf("Error al obtener los metadatos de la imagen: %v\n", err)
//		return
//	}
//
//	// Imprime los metadatos de la imagen
//	fmt.Printf("Metadatos de la imagen:\n%+v\n", manifest)
//}

func printImageMetadata(imageRef string, sys *types.SystemContext) {
	ref, err := alltransports.ParseImageName("docker://" + imageRef)
	if err != nil {
		fmt.Printf("Error parsing image reference: %v\n", err)
		return
	}

	imgSrc, err := ref.NewImageSource(context.Background(), sys)
	if err != nil {
		fmt.Printf("Error getting image source: %v\n", err)
		return
	}
	defer imgSrc.Close()

	configBlob, _, err := imgSrc.GetBlob(context.Background(), types.BlobInfo{URLs: imgSrc.Reference().PolicyConfigurationNamespaces()}, nil)
	if err != nil {
		fmt.Printf("Error getting config blob: %v\n", err)
		return
	}

	fmt.Printf("Image metadata:\n%+v\n", configBlob)
}
