package handlers

import "goplay2/rtsp"

/*
OnAuthSetup

	In case MFi authentication is supported (HasUnifiedAdvertiserInfo),
	it can be raised either with a auth-setup challenge or included in the pairing process
	if supported (SupportsUnifiedPairSetupAndMFi). This section describes the first option.
	Note that even if it is a server authentication (so that clients ensure MFi authenticity),
	with Airplay2 devices this step cannot be ignored from a client implementation point of view.
	Meaning, even if the authentication/signature is not checked on client side, the request has to be done,
	otherwise the server will deny further requests
	The challenge process is the following:
	1. Generation of Curve25119 key pairs on both client and server
	2. Client send its public key to server
	3. Server append its public key with client's one, this will be the message to sign.
		It then gets the signature from Apple authentication IC. RSA-1024 is used, with SHA-1 hash algorithm.
	4. Signature is encrypted with AES-128 in Counter mode with:
	5. Key = 16 first bytes of SHA1(<7:AES-KEY><32:Curve25119 Shared key>)
	6. IV = 16 first bytes of SHA1(<6:AES-IV><32:Curve25119 Shared key>)
	7. Server respond with it's Curve25119 public key, the encrypted signature and the certificate.
*/
func (r *Rstp) OnAuthSetup(req *rtsp.Request) (*rtsp.Response, error) {
	return nil, nil
}
