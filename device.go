package main

import "goplay2/globals"

// func airplayDevice() globals.Features {
// 	var features = globals.NewFeatures().Set(globals.SupportsAirPlayAudio).Set(globals.AudioRedundant)
// 	features = features.Set(globals.HasUnifiedAdvertiserInfo).Set(globals.SupportsBufferedAudio)
// 	features = features.Set(globals.SupportsUnifiedMediaControl)
// 	features = features.Set(globals.SupportsHKPairingAndAccessControl).Set(globals.SupportsHKPeerManagement)
// 	//features = features.Set(globals.SupportsUnifiedMediaControl).Set(globals.SupportsSystemPairing).Set(globals.SupportsCoreUtilsPairingAndEncryption).Set(globals.SupportsHKPairingAndAccessControl)
// 	features = features.Set(globals.Authentication_4)
// 	features = features.Set(globals.SupportsPTP)
// 	features = features.Set(globals.AudioFormats_0).Set(globals.AudioFormats_1).Set(globals.AudioFormats_2)

// 	return features
// }

func FeatAirplayDevice() globals.Feat {
	return globals.SupportsAirPlayAudio | globals.AudioRedundant | globals.HasUnifiedAdvertiserInfo |
		globals.SupportsBufferedAudio |
		// globals.SupportsUnifiedMediaControl |
		globals.SupportsHKPairingAndAccessControl | globals.SupportsHKPeerManagement |
		// features = features.Set(globals.SupportsUnifiedMediaControl).Set(globals.SupportsSystemPairing).Set(globals.SupportsCoreUtilsPairingAndEncryption).Set(globals.SupportsHKPairingAndAccessControl)
		globals.Authentication_4 | globals.SupportsPTP |
		globals.AudioFormats_0 | globals.AudioFormats_1 | globals.AudioFormats_2
}

// SupportsAirPlayAudio (bit 9)
// AudioRedundant (bit 11)
// HasUnifiedAdvertiserInfo (bit 30)
// SupportsBufferedAudio (bit 40)
// SupportsPTP (bit 41)
// SupportsUnifiedPairSetupAndMFi (bit 51)
// The respective features bitmask is 0x8030040000a00 and will be declared as features=0x40000a00,0x80300.

func airplayDeviceMin() globals.Feat {
	// var flags = globals.SupportsAirPlayAudio
	return globals.SupportsAirPlayAudio | globals.AudioFormats_1 | globals.AudioFormats_2 |
		globals.SupportsPTP | globals.HasUnifiedAdvertiserInfo | globals.SupportsBufferedAudio |
		globals.SupportsUnifiedPairSetupAndMFi
}
