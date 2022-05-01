package mkiso

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
)

var (
	flagVolumeLabel      string
	flagLogicalBlockSize int64
	flagOverwrite        bool
	flagVolumeIdentifier string
)

func init() {
	flag.StringVar(&flagVolumeLabel, "volume-label", "", "Volume label for image")
	flag.Int64Var(&flagLogicalBlockSize, "logical-block-size", 2048, "Logical block size for iso9660 {2048,4096,8192}")
	flag.BoolVar(&flagOverwrite, "overwrite", false, "Overwrite existing file")
	flag.StringVar(&flagVolumeIdentifier, "volume-identifier", "CIDATA", "Volume identifier for image {cidata,CIDATA,<custom>}")
}

type ciFile struct {
	name     string
	realPath string
	content  []byte
}

// return the ciFile's size rounded up to the nearest block size
func (ci *ciFile) logicalSize() int64 {
	clen := int64(len(ci.content))
	if clen%flagLogicalBlockSize != 0 {
		clen = (clen/flagLogicalBlockSize + 1) * flagLogicalBlockSize
	}
	return clen
}

// loads files into byte buffers and leaves them in a channel for later.
// returns the total bytes rounded up to the nearest logical-block-size
func loadFiles(files []string) (int64, []ciFile, error) {
	ciFiles := make([]ciFile, len(files))

	for i, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			return 0, nil, err
		}
		ciFiles[i].name = path.Base(f)
		ciFiles[i].realPath = f
		ciFiles[i].content = content
	}

	var rsize int64
	for _, cif := range ciFiles {
		rsize += cif.logicalSize()
	}
	return rsize, ciFiles, nil
}

func CreateIso(diskImg string, files []string) error {
	if diskImg == "" {
		log.Fatal("must have a valid path for diskImg")
	}
	blockBytes, ciFiles, err := loadFiles(files)
	// disk size is
	//     the system area (32768)
	//   + the single volume descriptor (2048)
	//   + the required volume descriptor terminator (2048)
	//   + blocks (ceil(size of data / 2048)*2048)
	var diskSize int64 = blockBytes + 32768 + 2048 + 2048
	mydisk, err := diskfs.Create(diskImg, diskSize, diskfs.Raw)
	if err != nil {
		return err
	}

	// the following line is required for an ISO, which may have logical block sizes
	// only of 2048, 4096, 8192
	mydisk.LogicalBlocksize = flagLogicalBlockSize
	fspec := disk.FilesystemSpec{Partition: 0, FSType: filesystem.TypeISO9660, VolumeLabel: flagVolumeLabel}
	fs, err := mydisk.CreateFilesystem(fspec)
	if err != nil {
		return err
	}
	for _, cif := range ciFiles {
		rw, err := fs.OpenFile(cif.name, os.O_CREATE|os.O_RDWR)
		if err != nil {
			return err
		}
		_, err = rw.Write(cif.content)
		if err != nil {
			return err
		}
	}
	iso, ok := fs.(*iso9660.FileSystem)
	if !ok {
		return fmt.Errorf("not an iso9660 filesystem")
	}
	err = iso.Finalize(iso9660.FinalizeOptions{
		RockRidge:        true,
		VolumeIdentifier: "cidata",
	})
	if err != nil {
		return err
	}
	return nil
}
