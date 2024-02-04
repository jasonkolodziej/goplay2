package globals

import (
	"fmt"
	"math/big"
	"strconv"
)

type Features struct {
	Value big.Int
}

type Feat int64

const (
	// video supported
	SupportsAirPlayVideoV1 Feat = 1 << iota
	// photo supported
	SupportsAirPlayPhoto
	//video protected with FairPlay DRM
	VideoFairPlay            // = 2
	VideoVolumeControl       // = 3
	VideoHTTPLiveStreams     // = 4
	SupportsAirPlaySlideshow //= 5
	//! Skip
	_
	// Mirroring supported
	SupportsAirPlayScreen         // = 7
	SupportsAirPlayScreenRotation // = 8
	SupportsAirPlayAudio          // = 9
	//! skip
	_
	// audio packet redundancy supported
	AudioRedundant // = 11
	// FairPlay secure auth supported
	FPSAPv2pt5_AES_GCM // = 12
	// photo preloading supported
	PhotoCaching // = 13
	// Authentication type 4. FairPlay authentication
	Authentication_4 // = 14
	// bit 0 of MetadataFeatures. Text.
	MetadataFeatures_0 // = 15
	// bit 1 of MetadataFeatures. Artwork.
	MetadataFeatures_1 // = 16
	// bit 2 of MetadataFeatures. Progress.
	MetadataFeatures_2 // = 17
	// support for audio format 1
	AudioFormats_0 // = 18
	// 	support for audio format 2.
	//! This bit must be set for AirPlay 2 connection to work
	AudioFormats_1 // = 19
	// support for audio format 3.
	//! This bit must be set for AirPlay 2 connection to work
	AudioFormats_2 // = 20
	// support for audio format 4
	AudioFormats_3 // = 21
	//! Skip
	_
	// Authentication type 1. RSA Authentication
	Authentication_1 // = 23
	//! Skip 2
	_
	_
	//! Conflict: HasUnifiedAdvertiserInfo
	Authentication_8      // = 26
	SupportsLegacyPairing // = 27
	//! Skip 3
	_
	_
	//! Conflict: ROAP
	// RAOP is supported on this port. With this bit set your don't need the AirTunes service
	HasUnifiedAdvertiserInfo // = 30
	//* Skip
	_
	// NOTE: Don’t read key from pk record it is known
	IsCarPlay      // = 32
	SupportsVolume = IsCarPlay

	SupportsAirPlayVideoPlayQueue Feat = 1 << (iota - 1) // = 33
	SupportsAirPlayFromCloud                             // = 34
	//! Conflict: was skipped
	SupportsTLS_PSK // = 35
	_
	_
	//! Conflict: SupportsCoreUtilsPairingAndEncryption <-- SupportsUnifiedMediaControl
	//* SupportsHKPairingAndAccessControl,
	//* SupportsSystemPairing and SupportsTransientPairing
	//*		 implies SupportsCoreUtilsPairingAndEncryption
	SupportsCoreUtilsPairingAndEncryption Feat = 1<<(iota-1) + SupportsHKPairingAndAccessControl |
		SupportsSystemPairing | SupportsTransientPairing // = 38
	//* Skip
	_
	//! Bit needed for device to show as supporting multi-room audio
	SupportsBufferedAudio Feat = 1 << (iota - 1) // = 40
	//! Bit needed for device to show as supporting multi-room audio
	SupportsPTP              // = 41
	SupportsScreenMultiCodec //  = 42

	SupportsSystemPairing Feat = 1<<(iota-1) + SupportsTransientPairing // = 43
	//! Conflict: Empty
	IsAPValeriaScreenSender Feat = 1 << (iota - 1) // = 44
	//* Skip
	_
	SupportsHKPairingAndAccessControl // = 46
	//! Conflict: Empty
	SupportsHKPeerManagement // = 47
	//! SupportsTransientPairing was SupportsCoreUtilsPairingAndEncryption
	//* SupportsSystemPairing implies SupportsTransientPairing
	SupportsTransientPairing // = 48
	//! Empty
	SupportsAirPlayVideoV2 // = 49
	//* bit 4 of MetadataFeatures. binary plist.
	MetadataFeatures_3 // = 50
	// Authentication type 8. MFi authentication
	SupportsUnifiedPairSetupAndMFi  // = 51
	SupportsSetPeersExtendedMessage // = 52

	//* Rest not described
	_
	SupportsAPSync // = 54
	SupportsWoL1   // = 55
	SupportsWoL2   // = 56
	_
	SupportsHangdogRemoteControl       // = 58
	SupportsAudioStreamConnectionSetup // = 59
	SupportsAudioMediaDataControl      // = 60
	SupportsRFC2198Redundancy          // = 61
	SupportsUnifiedMediaControl        = SupportsCoreUtilsPairingAndEncryption
)

// SupportsUnifiedMediaControl
// func (flag Features) Set(i int) Features {
// 	flag.Value.SetBit(&flag.Value, i, 1)
// 	return flag
// }

// func (flag Features) UnSet(i int) Features {
// 	flag.Value.SetBit(&flag.Value, i, 0)
// 	return flag
// }

// func (flag Features) ToRecord() string {
// 	return fmt.Sprintf("0x%x,0x%x", flag.Value.Int64()&0xffffffff, flag.Value.Int64()>>32&0xffffffff)
// }

// func (flag Features) ToUint64() uint64 {
// 	return flag.Value.Uint64()
// }

// func NewFeatures() Features {
// 	return Features{Value: *big.NewInt(0)}
// }

func (f *Feat) Set(i Feat) Feat {
	*f |= i // (1 << i)
	return *f
}

func (f *Feat) Unset(i Feat) Feat {
	(*f) &^= i // (0 << i)
	return *f
}

func (f *Feat) Contains(i Feat) bool {
	return ((*f) & i) != 0
}
func (f Feat) ToUint64() uint64 {
	return uint64(f)
}
func (f Feat) ToRecord() string {
	return fmt.Sprintf("0x%x,0x%x", f&0xffffffff, f>>32&0xffffffff)
}

func (f *Feat) SupportsExtendedWHA(srcvers *string) bool {
	//? srcvers >= 366 && (41 || forceAirPlay2NTP) && 40
	n, _ := strconv.ParseFloat(*srcvers, 32)
	return n >= 366.00 && (f.Contains(SupportsPTP)) && f.Contains(SupportsBufferedAudio)
}
