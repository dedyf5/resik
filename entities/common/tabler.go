// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package common

import (
	uuidPkg "github.com/dedyf5/resik/pkg/uuid"
)

type UserPublicIDProvider interface {
	AllUserPublicIDs() []uuidPkg.UUIDV7
}

func UniqueAllUserPublicIDs[
	T any,
	PT interface {
		*T
		UserPublicIDProvider
	},
](items []T) []uuidPkg.UUIDV7 {
	keys := make(map[uuidPkg.UUIDV7]struct{}, len(items)*2)
	var list []uuidPkg.UUIDV7

	for i := range items {
		var pt PT = &items[i]

		for _, id := range pt.AllUserPublicIDs() {
			if _, exists := keys[id]; !exists {
				keys[id] = struct{}{}
				list = append(list, id)
			}
		}
	}

	return list
}
