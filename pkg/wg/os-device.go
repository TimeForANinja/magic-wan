package wg

import (
	"magic-wan/pkg/osUtil"
)

func BaseConfigureInterface(ifcName string, selfIP string) error {
	// IMPROVMENT: expect net.IPNet instead of string
	err := osUtil.SetInterfaceAddress(ifcName, selfIP+"/31")
	if err != nil {
		return err
	}

	err = osUtil.SetInterfaceUp(ifcName)
	return err
}
