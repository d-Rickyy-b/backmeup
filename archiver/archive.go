package archiver

import (
	"archive/tar"
	"backmeup/config"
	"github.com/cheggaaa/pb/v3"
	"github.com/klauspost/compress/gzip"
	"github.com/klauspost/compress/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type BackupFileMetadata struct {
	Path           string
	BackupBasePath string
}

func getPathInArchive(filePath string, backupBasePath string, unit config.Unit) string {
	// Remove the base Path from the file Path within the archiver, if option is set
	pathInArchive := filePath

	if !unit.UseAbsolutePaths {
		parentBasePath := filepath.Dir(backupBasePath)
		pathInArchive = strings.ReplaceAll(filePath, parentBasePath, "")

		pathInArchive = strings.TrimPrefix(pathInArchive, "\\")
	}

	return pathInArchive
}

func WriteArchive(backupArchivePath string, filesToBackup []BackupFileMetadata, unit config.Unit) {
	archiveFile, err := os.Create(backupArchivePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer archiveFile.Close()

	switch unit.ArchiveType {
	case "tar.gz":
		writeTar(archiveFile, filesToBackup, unit)
	case "zip":
		writeZip(archiveFile, filesToBackup, unit)
	default:
		log.Fatalf("Can't handle archiver type '%s'", unit.ArchiveType)
	}
}

func writeTar(archiveFile *os.File, filesToBackup []BackupFileMetadata, unit config.Unit) {
	// set up the gzip and tar writer
	gw := gzip.NewWriter(archiveFile)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Init progress bar
	bar := pb.StartNew(len(filesToBackup))

	for i := range filesToBackup {
		fileMetadata := filesToBackup[i]
		filePath := fileMetadata.Path

		pathInArchive := getPathInArchive(filePath, fileMetadata.BackupBasePath, unit)

		if err := addFileToTar(tw, filePath, pathInArchive); err != nil {
			log.Fatalln(err)
		}

		bar.Increment()
	}

	bar.Finish()
}

func writeZip(archiveFile *os.File, filesToBackup []BackupFileMetadata, unit config.Unit) {
	zw := zip.NewWriter(archiveFile)
	defer zw.Close()

	bar := pb.StartNew(len(filesToBackup))

	for i := range filesToBackup {
		fileMetadata := filesToBackup[i]
		filePath := fileMetadata.Path

		pathInArchive := getPathInArchive(filePath, fileMetadata.BackupBasePath, unit)

		if err := addFileToZip(zw, filePath, pathInArchive); err != nil {
			log.Fatalln(err)
		}

		bar.Increment()
	}

	bar.Finish()
}

func addFileToTar(tw *tar.Writer, path string, pathInArchive string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if stat, err := file.Stat(); err == nil {
		// now lets create the header as needed for this file within the tarball
		header := new(tar.Header)
		header.Format = tar.FormatGNU
		header.Name = pathInArchive
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()
		// write the header to the tarball archiver
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		// copy the file data to the tarball
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
	}

	return nil
}

func addFileToZip(zw *zip.Writer, path string, pathInArchive string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if stat, err := file.Stat(); err == nil {
		header, headerErr := zip.FileInfoHeader(stat)
		if headerErr != nil {
			return headerErr
		}

		header.Name = pathInArchive
		// write the header to the zip archiver
		writer, headerErr := zw.CreateHeader(header)
		if headerErr != nil {
			return err
		}
		// copy the file data to the zip
		if _, err := io.Copy(writer, file); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}
