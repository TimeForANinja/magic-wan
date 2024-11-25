package pkg

import "magic-wan/pkg/osUtil"

// IMPROVMENT: if we don't even use any other Service we can remove half the osUtil code
var (
	FrrService = &osUtil.Service{
		Name: "frr",
	}
	MagicWanService = &osUtil.Service{
		Name: osUtil.WanServiceName,
	}
)
