package engine

import "fmt"

// Filter = gaya poles warna. Satu nilai = satu tanggung jawab (SRP).
// Request bukan filter poles — hanya titip ide, tidak di-FFmpeg.
type Filter string

const (
	FilterNatural   Filter = "natural"
	FilterDramatis  Filter = "dramatis"
	FilterCinematic Filter = "cinematic"
	FilterOriginal  Filter = "original"
	FilterRequest   Filter = "request"
)

// AllFilters: untuk TUI menu 1-5 (4 poles + 1 request).
func AllFilters() []Filter {
	return []Filter{FilterNatural, FilterDramatis, FilterCinematic, FilterOriginal, FilterRequest}
}

// Label: teks ramah untuk TUI.
func (f Filter) Label() string {
	switch f {
	case FilterNatural:
		return "Natural — halus, terang dikit"
	case FilterDramatis:
		return "Dramatis — pekat, kontras vivid ✨"
	case FilterCinematic:
		return "Cinematic — teal-orange, film bioskop"
	case FilterOriginal:
		return "Original — cuma resize, tanpa poles"
	case FilterRequest:
		return "Request — titip ide filter impian"
	default:
		return string(f)
	}
}

// IsValid: cegah input liar.
func (f Filter) IsValid() bool {
	switch f {
	case FilterNatural, FilterDramatis, FilterCinematic, FilterOriginal, FilterRequest:
		return true
	}
	return false
}

// ParseFilter: dari string input (aman, tidak terima FFmpeg bebas).
func ParseFilter(s string) (Filter, error) {
	f := Filter(s)
	if !f.IsValid() {
		return "", fmt.Errorf("filter tidak dikenal: %s — pilih natural/dramatis/cinematic/original/request", s)
	}
	return f, nil
}

// filterVF: rangkaian video filter FFmpeg per gaya (tidak bikin dari 0, manfaatkan FFmpeg).
// Output: bagian setelah scale+pad, dipasang di BuildArgs.
func filterVF(f Filter) string {
	switch f {
	case FilterNatural:
		// Halus — yang sekarang, biar tidak kaget.
		return "unsharp=5:5:0.8:3:3:0.4,eq=brightness=0.06:contrast=1.05:saturation=1.1"
	case FilterDramatis:
		// Pekat, vivid, vignette — kayak CapCut Dramatic.
		return "unsharp=5:5:1.2:3:3:0.6,eq=brightness=0.02:contrast=1.25:saturation=1.30,vignette=angle=PI/4"
	case FilterCinematic:
		// Teal-orange + curves + vignette — film bioskop.
		return "eq=brightness=0.02:contrast=1.15:saturation=1.05,colorbalance=rs=0.10:bs=-0.08,curves=r='0/0 0.5/0.48 1/1':g='0/0 0.5/0.52 1/1':b='0/0 0.5/0.55 1/1',vignette=angle=PI/4,unsharp=5:5:0.6:3:3:0.3"
	case FilterOriginal:
		// Tanpa poles.
		return ""
	default:
		return ""
	}
}
