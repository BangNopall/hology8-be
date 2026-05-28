package enums

type Status string

type Phase string

type WinnerPlace string

const (
	Waiting      Status      = "Waiting"
	Rejected     Status      = "Rejected"
	Verified     Status      = "Verified"
	Elimination  Phase       = "Elimination"
	Eliminated   Phase       = "Eliminated"
	Final        Phase       = "Final"
	Disqualified Phase       = "Disqualified"
	First        WinnerPlace = "First"
	Second       WinnerPlace = "Second"
	Third        WinnerPlace = "Third"
	Default      WinnerPlace = "Default"
)

func IsValidStatus(s string) bool {
	switch Status(s) {
	case Waiting, Rejected, Verified:
		return true
	default:
		return false
	}
}

func IsValidPhase(p string) bool {
	switch Phase(p) {
	case Elimination, Final, Disqualified, Eliminated:
		return true
	default:
		return false
	}
}

func IsValidWinnerPlace(wp string) bool {
	switch WinnerPlace(wp) {
	case First, Second, Third, Default:
		return true
	default:
		return false
	}
}
