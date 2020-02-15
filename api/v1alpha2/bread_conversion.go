package v1alpha2

import (
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func (src *Bread) ConvertTo(dstRaw conversion.Hub) error {
	return nil
}

func (src *Bread) ConvertFrom(srcRaw conversion.Hub) error {
	return nil
}
