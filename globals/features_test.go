package globals

import (
	"testing"
)

// func TestFeat(t *testing.T) {
// 	scv := "370.00"
// 	var f Feat
// 	var ff Features = NewFeatures()
// 	ff = ff.Set(SupportsAirPlayAudio)
// 	f.Set(SupportsAirPlayAudio)
// 	if f != Feat(ff.ToUint64()) {
// 		t.Fail()
// 	}
// 	ff = ff.Set(AudioRedundant)
// 	f.Set(AudioRedundant)
// 	if f != Feat(ff.ToUint64()) {
// 		t.Fail()
// 	}
// 	if !f.Contains(SupportsAirPlayAudio) {
// 		t.Error("f, Expected Feat.Contains() to report true")
// 	}

// 	ff = ff.Set(HasUnifiedAdvertiserInfo)
// 	f.Set(HasUnifiedAdvertiserInfo)
// 	if f != Feat(ff.ToUint64()) {
// 		t.Fail()
// 	}
// 	f.Unset(HasUnifiedAdvertiserInfo)
// 	if !f.Contains(HasUnifiedAdvertiserInfo) {
// 		t.Error("f, Expected Feat.Contains() to report false, DUE to unset")
// 	}
// 	f.Set(HasUnifiedAdvertiserInfo)
// 	f.Set(HasUnifiedAdvertiserInfo)
// 	if f != Feat(ff.ToUint64()) {
// 		t.Error("f, Expected to be equal, DUE to f.Set() called x2")
// 	}

// 	ff = ff.Set(SupportsBufferedAudio)
// 	f.Set(SupportsBufferedAudio)
// 	if f != Feat(ff.ToUint64()) {
// 		t.Fail()
// 	}
// 	ff = ff.Set(SupportsPTP)
// 	f.Set(SupportsPTP)
// 	if f != Feat(ff.ToUint64()) {
// 		t.Fail()
// 	}
// 	ff = ff.Set(SupportsUnifiedPairSetupAndMFi)
// 	f.Set(SupportsUnifiedPairSetupAndMFi)
// 	if f != Feat(ff.ToUint64()) {
// 		t.Fail()
// 	}
// 	if !(strings.Compare(f.ToRecord(), ff.ToRecord()) == 0) {
// 		t.Errorf("f: %s\nff: %s\n", f.ToRecord(), ff.ToRecord())
// 	}
// 	if !f.SupportsExtendedWHA(&scv) {
// 		t.Fail()
// 	}

// }

func TestFeat(t *testing.T) {
	var f Feat
	var ff Feat = MetadataFeatures_3 | SupportsAirPlayScreenRotation | FPSAPv2pt5_AES_GCM
	f.Set(MetadataFeatures_3 | SupportsAirPlayScreenRotation | FPSAPv2pt5_AES_GCM)
	if !f.Contains(MetadataFeatures_3) {
		t.Errorf("%d", (SupportsAPSync | SupportsAirPlayScreen))
	}
	if !f.Contains(FPSAPv2pt5_AES_GCM) {
		t.Fail()
	}
	if ff != f {
		t.Fail()
	}
	t.Logf(f.ToRecord(), ff.ToRecord())
	ff |= SupportsSystemPairing
	t.Logf(ff.ToRecord())
}
func FeatAirplayDevice() Feat {
	return SupportsAirPlayAudio | AudioRedundant | HasUnifiedAdvertiserInfo |
		SupportsBufferedAudio |
		SupportsUnifiedMediaControl |
		SupportsHKPairingAndAccessControl | SupportsHKPeerManagement |
		// features = features.Set(SupportsUnifiedMediaControl).Set(SupportsSystemPairing).Set(SupportsCoreUtilsPairingAndEncryption).Set(SupportsHKPairingAndAccessControl)
		Authentication_4 | SupportsPTP |
		AudioFormats_0 | AudioFormats_1 | AudioFormats_2
}

func TestAirplay(t *testing.T) {
	t.Logf("%s", FeatAirplayDevice().ToRecord())
}
