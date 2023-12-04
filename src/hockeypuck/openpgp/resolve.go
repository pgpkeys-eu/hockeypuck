/*
   Hockeypuck - OpenPGP key server
   Copyright (C) 2012-2014  Casey Marshall

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, version 3.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package openpgp

import (
	"crypto/md5"
	"encoding/hex"
	log "hockeypuck/logrus"

	"github.com/pkg/errors"
)

func ValidSelfSigned(key *PrimaryKey, selfSignedOnly bool) error {
	// Process direct signatures first
	ss, others := key.SigInfo()
	var certs []*Signature
	for _, cert := range ss.Revocations {
		if cert.Error == nil {
			certs = append(certs, cert.Signature)
		} else {
			log.Debugf("Dropped direct revocation sig because %s", cert.Error.Error())
		}
	}
	for _, cert := range ss.Certifications {
		if cert.Error == nil {
			certs = append(certs, cert.Signature)
		} else {
			log.Debugf("Dropped direct certification sig because %s", cert.Error.Error())
		}
	}
	key.Signatures = certs
	if !selfSignedOnly {
		key.Signatures = append(key.Signatures, others...)
	}
	var userIDs []*UserID
	var subKeys []*SubKey
	for _, uid := range key.UserIDs {
		ss, others := uid.SigInfo(key)
		var certs []*Signature
		for _, cert := range ss.Revocations {
			if cert.Error == nil {
				certs = append(certs, cert.Signature)
			} else {
				log.Debugf("Dropped revocation sig on uid '%s' because %s", uid.Keywords, cert.Error.Error())
			}
		}
		for _, cert := range ss.Certifications {
			if cert.Error == nil {
				certs = append(certs, cert.Signature)
			} else {
				log.Debugf("Dropped certification sig on uid '%s' because %s", uid.Keywords, cert.Error.Error())
			}
		}
		if len(certs) > 0 {
			uid.Signatures = certs
			if !selfSignedOnly {
				uid.Signatures = append(uid.Signatures, others...)
			}
			userIDs = append(userIDs, uid)
		} else {
			log.Debugf("Dropped uid '%s' because no valid self-sigs", uid.Keywords)
		}
	}
	for _, subKey := range key.SubKeys {
		ss, others := subKey.SigInfo(key)
		var certs []*Signature
		for _, cert := range ss.Revocations {
			if cert.Error == nil {
				certs = append(certs, cert.Signature)
			} else {
				log.Debugf("Dropped revocation sig on subkey %s because %s", subKey.KeyID(), cert.Error.Error())
			}
		}
		for _, cert := range ss.Certifications {
			if cert.Error == nil {
				certs = append(certs, cert.Signature)
			} else {
				log.Debugf("Dropped certification sig on subkey %s because %s", subKey.KeyID(), cert.Error.Error())
			}
		}
		if len(certs) > 0 {
			subKey.Signatures = certs
			if !selfSignedOnly {
				subKey.Signatures = append(subKey.Signatures, others...)
			}
			subKeys = append(subKeys, subKey)
		} else {
			log.Debugf("Dropped subkey %s because no valid self-sigs", subKey.KeyID())
		}
	}
	key.UserIDs = userIDs
	key.SubKeys = subKeys
	return key.updateMD5()
}

func CollectDuplicates(key *PrimaryKey) error {
	err := dedup(key, func(primary, _ packetNode) {
		primary.packet().Count++
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return key.updateMD5()
}

func Merge(dst, src *PrimaryKey) error {
	dst.UserIDs = append(dst.UserIDs, src.UserIDs...)
	dst.SubKeys = append(dst.SubKeys, src.SubKeys...)
	dst.Signatures = append(dst.Signatures, src.Signatures...)

	err := dedup(dst, func(primary, duplicate packetNode) {
		primaryPacket := primary.packet()
		duplicatePacket := duplicate.packet()
		if duplicatePacket.Count > primaryPacket.Count {
			primaryPacket.Count = duplicatePacket.Count
		}
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return ValidSelfSigned(dst, false)
}

func MergeRevocationSig(dst *PrimaryKey, src *Signature) error {
	dst.Signatures = append(dst.Signatures, src)

	err := dedup(dst, func(primary, duplicate packetNode) {
		primaryPacket := primary.packet()
		duplicatePacket := duplicate.packet()
		if duplicatePacket.Count > primaryPacket.Count {
			primaryPacket.Count = duplicatePacket.Count
		}
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return ValidSelfSigned(dst, false)
}

func hexmd5(b []byte) string {
	d := md5.Sum(b)
	return hex.EncodeToString(d[:])
}

func dedup(root packetNode, handleDuplicate func(primary, duplicate packetNode)) error {
	nodes := map[string]packetNode{}

	for _, node := range root.contents() {
		uuid := node.uuid() + "_" + hexmd5(node.packet().Packet)
		primary, ok := nodes[uuid]
		if ok {
			err := primary.removeDuplicate(root, node)
			if err != nil {
				return errors.WithStack(err)
			}

			err = dedup(primary, nil)
			if err != nil {
				return errors.WithStack(err)
			}

			if handleDuplicate != nil {
				handleDuplicate(primary, node)
			}
		} else {
			nodes[uuid] = node
		}
	}
	return nil
}
