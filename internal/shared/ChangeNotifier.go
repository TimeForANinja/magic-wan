package shared

type ChangeNotifier interface {
	OnManualInterfaceAdd() error
	OnWGInterfaceAdd() error
	OnWGInterfaceRemove() error
	OnIPChange() error
}
