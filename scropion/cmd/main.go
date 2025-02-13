package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

var exifTagNames = map[uint16]string{
	// TIFF/Basic Tags
	0x0100: "ImageWidth",
	0x0101: "ImageHeight",
	0x0102: "BitsPerSample",
	0x0103: "Compression",
	0x0106: "PhotometricInterpretation",
	0x010E: "ImageDescription",
	0x010F: "Make",
	0x0110: "Model",
	0x0112: "Orientation",
	0x0115: "SamplesPerPixel",
	0x011A: "XResolution",
	0x011B: "YResolution",
	0x0128: "ResolutionUnit",
	0x0131: "Software",
	0x0132: "DateTime",
	0x013B: "Artist",
	0x013E: "WhitePoint",
	0x013F: "PrimaryChromaticities",
	0x0211: "YCbCrCoefficients",
	0x0213: "YCbCrPositioning",
	0x0214: "ReferenceBlackWhite",
	0x8298: "Copyright",

	// EXIF Specific Tags
	0x829A: "ExposureTime",
	0x829D: "FNumber",
	0x8822: "ExposureProgram",
	0x8827: "ISOSpeedRatings",
	0x9000: "ExifVersion",
	0x9003: "DateTimeOriginal",
	0x9004: "DateTimeDigitized",
	0x9201: "ShutterSpeedValue",
	0x9202: "ApertureValue",
	0x9203: "BrightnessValue",
	0x9204: "ExposureBiasValue",
	0x9205: "MaxApertureValue",
	0x9206: "SubjectDistance",
	0x9207: "MeteringMode",
	0x9208: "LightSource",
	0x9209: "Flash",
	0x920A: "FocalLength",
	0x927C: "MakerNote",
	0x9286: "UserComment",
	0xA000: "FlashpixVersion",
	0xA001: "ColorSpace",
	0xA002: "PixelXDimension",
	0xA003: "PixelYDimension",
	0xA004: "RelatedSoundFile",
	0xA005: "InteroperabilityOffset",
	0xA20E: "FocalPlaneXResolution",
	0xA20F: "FocalPlaneYResolution",
	0xA210: "FocalPlaneResolutionUnit",
	0xA217: "SensingMethod",
	0xA300: "FileSource",
	0xA301: "SceneType",
}

func readUint16(r io.Reader, order binary.ByteOrder) (uint16, error) {
	var v uint16
	err := binary.Read(r, order, &v)

	return v, err
}

func readUint32(r io.Reader, order binary.ByteOrder) (uint32, error) {
	var v uint32
	err := binary.Read(r, order, &v)

	return v, err
}

func main() {

	// check args
	if len(os.Args) < 2 {
		log.Fatal("Usage: scorpion FILE1 [FILE2...]")
		return
	}

	filePath := os.Args[1]
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Cannot read image %v", err)
	}

	// ensure image is a jpeg
	if len(data) < 4 || data[0] != 0xFF || data[1] != 0xD8 {
		log.Fatal("Invalid Image Format !")
	}

	// looking for app1 marker that contain exif data (0xFFE1)
	// it may be multiple segments we pick the forst one for simplicity
	offset := 2
	var exifData []byte

	for offset < len(data)-2 {
		if data[offset] == 0xFF && data[offset+1] == 0xE1 {
			segSize := int(data[offset+2])<<8 + int(data[offset+3])
			start := offset + 4
			end := offset + segSize + 2
			if end > len(data) {
				break
			}

			// check for EXIF/0/0
			if bytes.HasPrefix(data[start:end], []byte("Exif\x00\x00")) {
				// The rest after "Exif\0\0" is the TIFF data
				exifData = data[start+6 : end]
				break
			}
		}
		offset++
	}

	if exifData == nil {
		log.Fatal("No Exif Data Found !")
		return
	}

	// exifData should now copntain TIFF header and IDFs
	// the first two bytes are either "II" (Intel : Little Endian) or "MM" (Motorola : Big Endian)
	byteOrder := exifData[0:2]
	var order binary.ByteOrder
	if bytes.Equal(byteOrder, []byte("II")) {
		order = binary.LittleEndian
	} else if bytes.Equal(byteOrder, []byte("MM")) {
		order = binary.BigEndian
	} else {
		log.Fatal("Invalid TIFF Byte Order in EXIF !")
		return
	}

	// the next 2 bytes should be 0x2A00 (42) to indecate a valid TIFF
	tiffCheck := order.Uint16(exifData[2:4])
	if tiffCheck != 42 {
		log.Fatal("Invalid TIFF indicator !")
	}

	// Next 4 bytes is the offset to the 0th IFD from start of the tiff header
	firstIFDOffset := order.Uint32(exifData[4:8])
	if firstIFDOffset >= uint32(len(exifData)) {
		log.Fatal("Invalid IFD Offset !")
		return
	}

	// Parse IFD entries
	// move reader to that offset
	r := bytes.NewReader(exifData[firstIFDOffset:])

	// NUmber of directory entries
	numEntries, err := readUint16(r, order)
	if err != nil {
		log.Fatal("Cannot Read Number of IFD Entries")
		return
	}

	exifValues := make(map[string]string)

	for i := 0; i < int(numEntries); i++ {
		entry := make([]byte, 12)
		_, err := io.ReadFull(r, entry)
		if err != nil {
			log.Fatalf("Error Reading IFD Entry %v\n", err)
			return
		}

		tagID := order.Uint16(entry[0:2])
		dataType := order.Uint16(entry[2:4])
		count := order.Uint32(entry[4:8])
		valueOffset := order.Uint32(entry[8:12])

		tagName, knownTag := exifTagNames[tagID]
		if !knownTag {
			// unkown tag !
			continue
		}

		// if the values fits within 4 bytes is stored right here (for some datatypes)
		// otherwise the offset point to the actual data
		var valueBytes []byte

		typeSize := uint32(1)
		switch dataType {
		case 1, 2, 7: // Byte, ASCII, UNDEFIENDED => size per value = 1
			typeSize = 1
		case 3: // Short => size per value =2
			typeSize = 2
		case 4: // Long => size per value = 4
			typeSize = 4
		case 5: // Rational => size per value = 8
			typeSize = 8
		}

		valueTotalSize := count * typeSize
		if valueTotalSize <= 4 {
			// we treat valueOffset as the actual data in this scenario
			buf := make([]byte, 4)
			order.PutUint32(buf, valueOffset)
			valueBytes = buf[:valueTotalSize]
		} else {
			// otherwise we must read it from exifData on that offset
			// The ExifData start at exifData[0]
			valueStart := int(valueOffset)
			valueEnd := valueStart + int(valueTotalSize)
			if valueEnd > len(exifData) {
				log.Println("Invalid Offse for IFD value data.")
				continue
			}
			valueBytes = exifData[valueStart:valueEnd]
		}

		// decode value according to data type
		switch dataType {
		case 2:
			// ASCII - typically null terminated
			// convert to string trimming any null
			s := string(valueBytes)
			if idx := bytes.IndexByte(valueBytes, 0); idx >= 0 {
				s = s[:idx]
			}
			exifValues[tagName] = s
		case 1, 7:
			// Bytes or undefiened
			// store hex or raw for minimal example
			exifValues[tagName] = fmt.Sprintf("%v", valueBytes)
		case 3:
			// short
			// doesnt handle multple, this assume count = 1
			if len(valueBytes) >= 2 {
				val := order.Uint16(valueBytes)
				exifValues[tagName] = fmt.Sprintf("%v", val)
			}
		case 4:
			// long
			if len(valueBytes) >= 4 {
				val := order.Uint32(valueBytes)
				exifValues[tagName] = fmt.Sprintf("%v", val)
			}
		case 5:
			// Rational = two longs, numerator/denomenator
			// this simplified for count = 1
			if len(valueBytes) >= 8 {
				num := order.Uint32(valueBytes[0:4])
				den := order.Uint32(valueBytes[4:8])
				if den != 0 {
					exifValues[tagName] = fmt.Sprintf("%v", float64(num)/float64(den))
				}
			}
		default:
			// not handled here !
		}

	}

	fmt.Println("Extraxted Exif Data !")
	for k, v := range exifValues {
		fmt.Printf("   %s : %s\n", k, v)
	}

}
