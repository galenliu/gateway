//THis File is AUTO-GENERATED
package characteristic

const (
	ServiceLabelNamespaceArabicNumerals int = 1
	ServiceLabelNamespaceDots           int = 0
)
const TypeServiceLabelNamespace = "CD"

type ServiceLabelNamespace struct {
	*Int
}

func NewServiceLabelNamespace() *ServiceLabelNamespace {

	char := NewInt(TypeServiceLabelNamespace)
	char.Format = FormatUint8
	char.Perms = []string{PermRead}

	char.SetValue(0)
	return &ServiceLabelNamespace{char}
}
